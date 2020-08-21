package cluster

import (
	"context"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/pkg/errors"
	"go.uber.org/zap"
	"k8s.io/apimachinery/pkg/util/wait"
	runtimeapi "k8s.io/cri-api/pkg/apis/runtime/v1alpha2"

	"github.com/criticalstack/crit/internal/config"
	"github.com/criticalstack/crit/pkg/cluster/components"
	computil "github.com/criticalstack/crit/pkg/cluster/components/util"
	clusterutil "github.com/criticalstack/crit/pkg/cluster/util"
	"github.com/criticalstack/crit/pkg/kubernetes/pki"
	"github.com/criticalstack/crit/pkg/kubernetes/remote"
	"github.com/criticalstack/crit/pkg/log"
	executil "github.com/criticalstack/crit/pkg/util/exec"
	"github.com/criticalstack/crit/pkg/util/systemd"
)

func (c *Cluster) WriteKubeConfigs(ctx context.Context, cfg *config.ControlPlaneConfiguration) error {
	log.Info("kubeconfigs", zap.String("description", "write kubeconfigs to disk"))
	ca, err := pki.LoadCertificateAuthority(filepath.Join(cfg.NodeConfiguration.KubeDir, "pki"), "ca")
	if err != nil {
		return err
	}
	fns := []func(*config.ControlPlaneConfiguration, *pki.CertificateAuthority) error{
		clusterutil.WriteAdminConfig,
		clusterutil.WriteControllerManagerConfig,
		clusterutil.WriteSchedulerConfig,
		clusterutil.WriteKubeletConfig,
	}
	for _, fn := range fns {
		if err := fn(cfg, ca); err != nil {
			return err
		}
	}
	return nil
}

func (c *Cluster) WriteKubeManifests(ctx context.Context, cfg *config.ControlPlaneConfiguration) error {
	log.Info("kubemanifests", zap.String("description", "write kubernetes static pod manifests to disk"))
	p, err := components.NewAPIServerStaticPod(cfg)
	if err != nil {
		return err
	}
	if err := computil.WriteKubeComponent(p, filepath.Join(cfg.NodeConfiguration.KubeDir, "manifests/kube-apiserver.yaml")); err != nil {
		return err
	}
	if err := computil.WriteKubeComponent(components.NewControllerManagerStaticPod(cfg), filepath.Join(cfg.NodeConfiguration.KubeDir, "manifests/kube-controller-manager.yaml")); err != nil {
		return err
	}
	return computil.WriteKubeComponent(components.NewSchedulerStaticPod(cfg), filepath.Join(cfg.NodeConfiguration.KubeDir, "manifests/kube-scheduler.yaml"))
}

func (c *Cluster) WriteBootstrapServerManifest(ctx context.Context, cfg *config.ControlPlaneConfiguration) error {
	log.Info("bootstrap-server-manifest", zap.String("description", "write bootstrap server static pod manifest to disk"))
	return computil.WriteKubeComponent(components.NewBootstrapServerStaticPod(cfg), filepath.Join(cfg.NodeConfiguration.KubeDir, "manifests/crit-bootstrap-server.yaml"))
}

func (c *Cluster) WaitClusterAvailable(ctx context.Context, cfg *config.ControlPlaneConfiguration) error {
	log.Info("cluster-available", zap.String("description", "wait for cluster to become available"))
	ctx, cancel := context.WithTimeout(ctx, 4*time.Minute)
	defer cancel()

	r, err := remote.NewRuntimeServiceClient(ctx, cfg.NodeConfiguration.ContainerRuntime.CRISocket())
	if err != nil {
		return err
	}

	// The apiserver container is only ever queried once, because the
	// presumption is made that it should not have to restart during
	// initial bootstrapping. Should this no longer be the case in the
	// future, this will need to be adapted to account for scenarios where
	// the intended flow of the apiserver allows for 1 or more container
	// restarts.
	var container *runtimeapi.Container
	if err := wait.PollImmediateUntil(500*time.Millisecond, func() (bool, error) {
		var err error
		container, err = r.GetContainerByName(ctx, "kube-apiserver")
		if err != nil {
			return false, nil
		}
		return true, nil
	}, ctx.Done()); err != nil {
		return err
	}
	status, err := r.GetContainerStatus(ctx, container.GetId())
	if err != nil {
		return err
	}

	// While waiting for the apiserver to start, the status of the
	// container is checked and short circuits should the container be
	// unavailable or stopped with error.
	errCh := make(chan error, 1)
	go func() {
		defer close(errCh)

		//st := time.Now()
		errCh <- wait.PollImmediateUntil(500*time.Millisecond, func() (bool, error) {
			status, err := r.GetContainerStatus(ctx, container.GetId())
			if err != nil {
				return false, errors.Wrap(err, "kube-apiserver container exited")
			}

			if status.ExitCode != 0 {
				return false, errors.Errorf("kube-apiserver container exited with code: %d", status.ExitCode)
			}
			return false, nil
		}, ctx.Done())
	}()

	// The apiserver container logs will only stream logs when verbose
	// output is requested from the user. For those not familiar with the
	// output of the apiserver, these logs can be confusing and unhelpful
	// in most situations.
	if c.rc.Verbose {
		go func() {
			stdout := executil.NewPrefixWriter(os.Stdout, "\t")
			defer stdout.Close()

			if err := r.TailLogs(ctx, status.GetLogPath(), stdout); err != nil {
				log.Error("cannot tail log", zap.Error(err))
			}
		}()
	}

	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			// Finally, check for the apiserver to start reporting as healthy.
			status := 0
			c.Client().Discovery().RESTClient().Get().AbsPath("/healthz").Do(ctx).StatusCode(&status)
			if status == http.StatusOK {
				return nil
			}
		case reterr := <-errCh:
			// already printing the logs, so just exit
			if c.rc.Verbose {
				return reterr
			}
			stdout := executil.NewPrefixWriter(os.Stdout, "\t")
			defer stdout.Close()

			// We have to stop the kubelet and range over the container log
			// path directory. This is for cases like cinder, where the
			// systemd service restarts so quickly (every 1s) that the
			// stored log path is no longer available.
			if err := systemd.StopUnit("kubelet.service"); err != nil {
				return err
			}
			dir := filepath.Dir(status.GetLogPath())
			infos, err := ioutil.ReadDir(dir)
			if err != nil {
				return err
			}
			files := make([]string, 0)
			for _, info := range infos {
				files = append(files, filepath.Join(dir, info.Name()))
			}
			sort.Sort(sort.Reverse(sort.StringSlice(files)))
			for _, file := range files {
				if err := r.ReadLogs(ctx, file, stdout); err != nil {
					log.Debug("cannot read log path",
						zap.Error(err),
						zap.String("path", file),
					)
					continue
				}
				return reterr
			}
			return reterr
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}
