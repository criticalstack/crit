package cluster

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/pkg/errors"
	"go.uber.org/zap"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/util/wait"
	clientset "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/criticalstack/crit/internal/config"
	"github.com/criticalstack/crit/pkg/cluster/bootstrap"
	"github.com/criticalstack/crit/pkg/cluster/components"
	configutil "github.com/criticalstack/crit/pkg/config/util"
	"github.com/criticalstack/crit/pkg/kubeconfig"
	"github.com/criticalstack/crit/pkg/kubernetes"
	"github.com/criticalstack/crit/pkg/log"
	executil "github.com/criticalstack/crit/pkg/util/exec"
	netutil "github.com/criticalstack/crit/pkg/util/net"
	"github.com/criticalstack/crit/pkg/util/systemd"
)

const (
	DefaultKubeletDir = "/var/lib/kubelet"

	kubeletFailureMessage = `Attempt to start the kubelet.service was not successful. The
kubelet is required to start Kubernetes components, so
Kubernetes will not be available. Check the full logs in the
systemd journal:

	journalctl -xu kubelet.service

`
)

func (c *Cluster) WriteKubeletConfigs(ctx context.Context, cfg *config.NodeConfiguration) error {
	log.Info("write-kubelet-configs", zap.String("description", "write the kubelet configs"))
	if err := components.WriteKubeletDynamicEnvFile(cfg, false, DefaultKubeletDir); err != nil {
		return err
	}
	cfg.KubeletConfiguration.Authentication.X509.ClientCAFile = filepath.Join(cfg.KubeDir, "pki/ca.crt")
	if cfg.KubeletConfiguration.StaticPodPath == "" {
		cfg.KubeletConfiguration.StaticPodPath = filepath.Join(cfg.KubeDir, "manifests")
	}
	cfg.KubeletConfiguration.RotateCertificates = true
	return components.WriteKubeletConfigFile(cfg.KubeletConfiguration, filepath.Join(DefaultKubeletDir, "config.yaml"))
}

func (c *Cluster) StopKubelet(ctx context.Context, cfg *config.NodeConfiguration) error {
	log.Info("stop-kubelet", zap.String("description", "stop kubelet service"))
	return systemd.StopUnit("kubelet.service")
}

func (c *Cluster) StartKubelet(ctx context.Context, cfg *config.NodeConfiguration) error {
	log.Info("start-kubelet", zap.String("description", "start kubelet service"))
	if err := systemd.StartUnit("kubelet.service"); err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(ctx, c.rc.KubeletTimeout)
	defer cancel()

	if err := wait.PollImmediateUntil(500*time.Millisecond, func() (bool, error) {
		status, err := systemd.UnitStatus("kubelet.service")
		if err != nil {
			log.Debug("kubelet unit status", zap.Error(err))
			return false, nil
		}
		if !systemd.UnitIsReady(status) {
			return false, nil
		}
		client := &http.Client{Timeout: 1 * time.Second}
		url := fmt.Sprintf("http://%s:%d/healthz", cfg.KubeletConfiguration.HealthzBindAddress, *cfg.KubeletConfiguration.HealthzPort)
		resp, err := client.Get(url)
		if err != nil {
			log.Debug("kubelet health check", zap.Error(err))
			return false, nil
		}
		defer resp.Body.Close()

		return true, nil
	}, ctx.Done()); err != nil {
		stderr := executil.NewPrefixWriter(os.Stderr, "\t")
		defer stderr.Close()

		stderr.Write([]byte(kubeletFailureMessage))
		return errors.New("Attempt to start the kubelet.service was not successful")
	}
	return nil
}

func (c *Cluster) WriteBootstrapKubeletConfig(ctx context.Context, cfg *config.WorkerConfiguration) error {
	log.Info("write-bootstrap-kubelet-config", zap.String("description", "create bootstrap-kubelet.conf"))
	bootstrapKubeletConf, err := bootstrap.GetBootstrapKubeletKubeconfig(cfg)
	if err != nil {
		return err
	}
	clientConfig, err := clientcmd.NewDefaultClientConfig(*bootstrapKubeletConf, &clientcmd.ConfigOverrides{}).ClientConfig()
	if err != nil {
		return errors.Wrap(err, "failed to create API client configuration from kubeconfig")
	}
	client, err := clientset.NewForConfig(clientConfig)
	if err != nil {
		return err
	}
	var cm *corev1.ConfigMap
	if err := wait.PollImmediateUntil(500*time.Millisecond, func() (ok bool, err error) {
		cm, err = kubernetes.GetConfigMap(client, ctx, CritConfigName)
		if err != nil {
			if apierrors.IsUnauthorized(err) || apierrors.IsNotFound(err) {
				return false, nil
			}
			log.Error("cannot get worker crit-config", zap.Error(err))
			return false, nil
		}
		return true, nil
	}, ctx.Done()); err != nil {
		return err
	}
	obj, err := configutil.Unmarshal([]byte(cm.Data["config"]))
	if err != nil {
		return err
	}
	cCfg, ok := obj.(*config.ControlPlaneConfiguration)
	if !ok {
		return errors.Errorf("expected ControlPlaneConfiguration, received %T", cCfg)
	}
	if cfg.NodeConfiguration.KubeletConfiguration.ClusterDNS == nil {
		dnsIP, err := netutil.GetDNSIP(cCfg.ServiceSubnet)
		if err != nil {
			log.Warn("cannot parse service subnet", zap.Error(err))
			cfg.NodeConfiguration.KubeletConfiguration.ClusterDNS = []string{DefaultClusterDNSIP}
		} else {
			cfg.NodeConfiguration.KubeletConfiguration.ClusterDNS = []string{dnsIP.String()}
		}
	}
	return kubeconfig.WriteToFile(bootstrapKubeletConf, filepath.Join(cfg.NodeConfiguration.KubeDir, "bootstrap-kubelet.conf"))
}
