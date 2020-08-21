package config

import (
	clientsetscheme "k8s.io/client-go/kubernetes/scheme"

	configv1alpha1 "github.com/criticalstack/crit/pkg/config/v1alpha1"
	externalconfig "github.com/criticalstack/crit/pkg/config/v1alpha2"
)

// The external configuration types being used are aliased to be used internal
// to the project without needing to update import paths.
type (
	ControlPlaneConfiguration          = externalconfig.ControlPlaneConfiguration
	WorkerConfiguration                = externalconfig.WorkerConfiguration
	NodeConfiguration                  = externalconfig.NodeConfiguration
	EtcdConfiguration                  = externalconfig.EtcdConfiguration
	CritBootstrapServerConfiguration   = externalconfig.CritBootstrapServerConfiguration
	KubeAPIServerConfiguration         = externalconfig.KubeAPIServerConfiguration
	KubeControllerManagerConfiguration = externalconfig.KubeControllerManagerConfiguration
)

var SchemeGroupVersion = externalconfig.SchemeGroupVersion

func init() {
	_ = configv1alpha1.RegisterConversions(clientsetscheme.Scheme)
}
