package v1alpha1

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/conversion"
	kubeproxyconfigv1alpha1 "k8s.io/kube-proxy/config/v1alpha1"

	"github.com/criticalstack/crit/pkg/config/constants"
	"github.com/criticalstack/crit/pkg/config/v1alpha2"
	netutil "github.com/criticalstack/crit/pkg/util/net"
)

func Convert_v1alpha1_ControlPlaneConfiguration_To_v1alpha2_ControlPlaneConfiguration(in *ControlPlaneConfiguration, out *v1alpha2.ControlPlaneConfiguration, s conversion.Scope) error {
	SetDefaults_ControlPlaneConfiguration(in)
	v1alpha2.SetDefaults_ControlPlaneConfiguration(out)

	switch {
	case in.APIServerURL != "":
		host, port, err := netutil.SplitHostPort(strings.TrimPrefix(in.APIServerURL, "https://"))
		if err != nil {
			return err
		}
		out.ControlPlaneEndpoint.Host = host
		out.ControlPlaneEndpoint.Port = port
	case in.ControlPlaneEndpoint != "":
		out.ControlPlaneEndpoint.Host = in.ControlPlaneEndpoint
	}

	// the bootstrap server was originally enabled by default
	if out.FeatureGates == nil {
		out.FeatureGates = make(map[string]bool)
	}
	out.FeatureGates["BootstrapServer"] = true

	out.KubeAPIServerConfiguration.HealthcheckProxyVersion = constants.DefaultHealthcheckProxyVersion
	out.KubeAPIServerConfiguration.HealthcheckProxyBindPort = constants.DefaultHealthcheckProxyBindPort
	switch in.NodeConfiguration.KubeProxyMode {
	case "disabled":
		out.KubeProxyConfiguration.Disabled = true
	case "iptables", "ipvs":
		out.KubeProxyConfiguration.Config.Mode = kubeproxyconfigv1alpha1.ProxyMode(in.NodeConfiguration.KubeProxyMode)
	default:
		return errors.Errorf("unknown KubeProxyMode: %s", in.NodeConfiguration.KubeProxyMode)
	}
	out.PodSubnet = in.NodeConfiguration.PodSubnet
	out.ServiceSubnet = in.NodeConfiguration.ServiceSubnet
	return autoConvert_v1alpha1_ControlPlaneConfiguration_To_v1alpha2_ControlPlaneConfiguration(in, out, s)
}

func Convert_v1alpha1_WorkerConfiguration_To_v1alpha2_WorkerConfiguration(in *WorkerConfiguration, out *v1alpha2.WorkerConfiguration, s conversion.Scope) error {
	SetDefaults_WorkerConfiguration(in)
	v1alpha2.SetDefaults_WorkerConfiguration(out)

	switch {
	case in.APIServerURL != "":
		host, port, err := netutil.SplitHostPort(strings.TrimPrefix(in.APIServerURL, "https://"))
		if err != nil {
			return err
		}
		out.ControlPlaneEndpoint.Host = host
		out.ControlPlaneEndpoint.Port = port
	case in.ControlPlaneEndpoint != "":
		out.ControlPlaneEndpoint.Host = in.ControlPlaneEndpoint
	}
	return autoConvert_v1alpha1_WorkerConfiguration_To_v1alpha2_WorkerConfiguration(in, out, s)
}

func Convert_v1alpha2_ControlPlaneConfiguration_To_v1alpha1_ControlPlaneConfiguration(in *v1alpha2.ControlPlaneConfiguration, out *ControlPlaneConfiguration, s conversion.Scope) error {
	v1alpha2.SetDefaults_ControlPlaneConfiguration(in)
	SetDefaults_ControlPlaneConfiguration(out)

	if !in.ControlPlaneEndpoint.IsZero() {
		out.APIServerURL = fmt.Sprintf("https://%s", in.ControlPlaneEndpoint)
	}
	out.ControlPlaneEndpoint = in.ControlPlaneEndpoint.Host
	out.NodeConfiguration.KubeProxyMode = string(in.KubeProxyConfiguration.Config.Mode)
	if in.KubeProxyConfiguration.Disabled {
		out.NodeConfiguration.KubeProxyMode = "disabled"
	}
	out.NodeConfiguration.PodSubnet = in.PodSubnet
	out.NodeConfiguration.ServiceSubnet = in.ServiceSubnet
	return autoConvert_v1alpha2_ControlPlaneConfiguration_To_v1alpha1_ControlPlaneConfiguration(in, out, s)
}

func Convert_v1alpha2_WorkerConfiguration_To_v1alpha1_WorkerConfiguration(in *v1alpha2.WorkerConfiguration, out *WorkerConfiguration, s conversion.Scope) error {
	v1alpha2.SetDefaults_WorkerConfiguration(in)
	SetDefaults_WorkerConfiguration(out)

	if !in.ControlPlaneEndpoint.IsZero() {
		out.APIServerURL = fmt.Sprintf("https://%s", in.ControlPlaneEndpoint)
	}
	out.ControlPlaneEndpoint = in.ControlPlaneEndpoint.Host
	return autoConvert_v1alpha2_WorkerConfiguration_To_v1alpha1_WorkerConfiguration(in, out, s)
}

func Convert_v1alpha2_CritBootstrapServerConfiguration_To_v1alpha1_CritBootstrapServerConfiguration(in *v1alpha2.CritBootstrapServerConfiguration, out *CritBootstrapServerConfiguration, s conversion.Scope) error {
	// Enabled is ignored as the bootstrap server could not be disabled in
	// v1alpha1
	return autoConvert_v1alpha2_CritBootstrapServerConfiguration_To_v1alpha1_CritBootstrapServerConfiguration(in, out, s)
}

func Convert_v1alpha1_CritBootstrapServerConfiguration_To_v1alpha2_CritBootstrapServerConfiguration(in *CritBootstrapServerConfiguration, out *v1alpha2.CritBootstrapServerConfiguration, s conversion.Scope) error {
	return autoConvert_v1alpha1_CritBootstrapServerConfiguration_To_v1alpha2_CritBootstrapServerConfiguration(in, out, s)
}

func Convert_v1alpha2_KubeAPIServerConfiguration_To_v1alpha1_KubeAPIServerConfiguration(in *v1alpha2.KubeAPIServerConfiguration, out *KubeAPIServerConfiguration, s conversion.Scope) error {
	return autoConvert_v1alpha2_KubeAPIServerConfiguration_To_v1alpha1_KubeAPIServerConfiguration(in, out, s)
}

func Convert_v1alpha1_NodeConfiguration_To_v1alpha2_NodeConfiguration(in *NodeConfiguration, out *v1alpha2.NodeConfiguration, s conversion.Scope) error {
	// KubeProxyMode moved to ControlPlaneConfiguration and conversion is
	// handled there
	return autoConvert_v1alpha1_NodeConfiguration_To_v1alpha2_NodeConfiguration(in, out, s)
}
