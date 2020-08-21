package v1alpha1

import (
	"k8s.io/apimachinery/pkg/runtime"

	"github.com/criticalstack/crit/internal/cinder/config/constants"
)

func addDefaultingFuncs(scheme *runtime.Scheme) error {
	return RegisterDefaults(scheme)
}

func SetDefaults_ClusterConfiguration(obj *ClusterConfiguration) {
	if obj.FeatureGates == nil {
		obj.FeatureGates = make(map[string]bool)
	}
	if obj.LocalRegistryName == "" {
		obj.LocalRegistryName = constants.DefaultLocalRegistryName
	}
	if obj.LocalRegistryPort == 0 {
		obj.LocalRegistryPort = constants.DefaultLocalRegistryPort
	}
	if obj.RegistryMirrors == nil {
		obj.RegistryMirrors = make(map[string]string)
	}
	for _, pm := range obj.ExtraPortMappings {
		if pm.ListenAddress == "" {
			pm.ListenAddress = "127.0.0.1"
		}
		if pm.Protocol == "" {
			pm.Protocol = PortMappingProtocolTCP
		}
	}
}
