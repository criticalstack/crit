// Package cluster contains the functions for bootstrapping a Kubernetes
// cluster node.
package cluster

import (
	"context"
	"path/filepath"
	"time"

	"github.com/pkg/errors"
	"go.uber.org/zap"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/criticalstack/crit/internal/config"
	"github.com/criticalstack/crit/internal/feature"
	"github.com/criticalstack/crit/pkg/log"
)

type RuntimeConfig struct {
	KubeletTimeout time.Duration
	Verbose        bool
}

type Cluster struct {
	kubeConfigFile string
	rc             *RuntimeConfig
	fns            []interface{}
}

func New(kubeConfigFile string, rc *RuntimeConfig) *Cluster {
	return &Cluster{
		kubeConfigFile: kubeConfigFile,
		rc:             rc,
		fns:            make([]interface{}, 0),
	}
}

func (c *Cluster) Add(fns ...interface{}) {
	c.fns = append(c.fns, fns...)
}

func (c *Cluster) Config() *rest.Config {
	config, err := clientcmd.BuildConfigFromFlags("", c.kubeConfigFile)
	if err != nil {
		log.Debug("Cluster.Config", zap.Error(err))
	}
	return config
}

func (c *Cluster) Client() *kubernetes.Clientset {
	client, err := kubernetes.NewForConfig(c.Config())
	if err != nil {
		log.Debug("Cluster.Client", zap.Error(err))
	}
	return client
}

type (
	controlPlaneFunc = func(context.Context, *config.ControlPlaneConfiguration) error
	workerFunc       = func(context.Context, *config.WorkerConfiguration) error
	nodeFunc         = func(context.Context, *config.NodeConfiguration) error
)

// RunControlPlane creates a new control plane node.
func RunControlPlane(ctx context.Context, rc *RuntimeConfig, cfg *config.ControlPlaneConfiguration) error {
	c := New(filepath.Join(cfg.NodeConfiguration.KubeDir, "admin.conf"), rc)

	// set crit feature gates
	if err := feature.MutableGates.SetFromMap(cfg.FeatureGates); err != nil {
		return err
	}
	c.Add(
		c.ControlPlanePreCheck,
		c.CreateOrDownloadCerts,
		c.CreateNodeCerts,
		c.StopKubelet,
		c.WriteKubeConfigs,
		c.WriteKubeletConfigs,
		c.StartKubelet,
		c.WriteKubeManifests,
		c.WaitClusterAvailable,
	)
	if feature.Gates.Enabled(feature.BootstrapServer) {
		c.Add(c.WriteBootstrapServerManifest)
	}
	c.Add(
		c.DeployCoreDNS,
		c.DeployKubeProxy,
		c.EnableCSRApprover,
		c.MarkControlPlane,
		c.UploadInfo,
	)
	if feature.Gates.Enabled(feature.AuthProxyCA) {
		c.Add(c.UploadAuthProxyCA)
	}
	if feature.Gates.Enabled(feature.UploadETCDSecrets) {
		c.Add(c.UploadETCDSecrets)
	}
	for _, fn := range c.fns {
		switch fn := fn.(type) {
		case controlPlaneFunc:
			if err := fn(ctx, cfg); err != nil {
				return err
			}
		case nodeFunc:
			if err := fn(ctx, &cfg.NodeConfiguration); err != nil {
				return err
			}
		default:
			panic(errors.Errorf("invalid cluster workflow function: %T", fn))
		}
	}
	return nil
}

// RunWorkerNode creates a new worker node.
func RunWorkerNode(ctx context.Context, rc *RuntimeConfig, cfg *config.WorkerConfiguration) error {
	c := New(filepath.Join(cfg.NodeConfiguration.KubeDir, "kubelet.conf"), rc)

	// set crit feature gates
	if err := feature.MutableGates.SetFromMap(cfg.FeatureGates); err != nil {
		return err
	}
	c.Add(
		c.WorkerPreCheck,
		c.StopKubelet,
		c.WriteBootstrapKubeletConfig,
		c.WriteKubeletConfigs,
		c.StartKubelet,
	)
	for _, fn := range c.fns {
		switch fn := fn.(type) {
		case workerFunc:
			if err := fn(ctx, cfg); err != nil {
				return err
			}
		case nodeFunc:
			if err := fn(ctx, &cfg.NodeConfiguration); err != nil {
				return err
			}
		default:
			panic(errors.Errorf("invalid cluster workflow function: %T", fn))
		}
	}
	return nil
}
