package v1alpha2

import (
	"net/url"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clientsetscheme "k8s.io/client-go/kubernetes/scheme"
	kubeproxyconfigv1alpha1 "k8s.io/kube-proxy/config/v1alpha1"
	kubeletconfigv1beta1 "k8s.io/kubelet/config/v1beta1"

	computil "github.com/criticalstack/crit/pkg/cluster/components/util"
	"github.com/criticalstack/crit/pkg/config/constants"
)

func init() {
	_ = AddToScheme(clientsetscheme.Scheme)
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type ControlPlaneConfiguration struct {
	metav1.TypeMeta `json:",inline"`
	// ClusterName
	// Default: "crit"
	// +optional
	ClusterName string `json:"clusterName"`
	// ControlPlaneEndpoint is the IP address or DNS name that represents the
	// control plane, along with optional port. The host portion is
	// automatically added to the cluster CA SANs.
	// +optional
	ControlPlaneEndpoint computil.APIEndpoint `json:"controlPlaneEndpoint,omitempty"`
	// PodSubnet is the CIDR range for allocating private IP addresses for
	// pods.
	// Default: "10.253.0.0/16"
	PodSubnet string `json:"podSubnet,omitempty"`
	// ServiceSubnet is the CIDR range for allocating private IP addresses for
	// services.
	// Default: "10.254.0.0/16"
	// +optional
	ServiceSubnet string `json:"serviceSubnet,omitempty"`
	// CoreDNSVersion is the version given to the CoreDNS template.
	// Default: "1.6.9"
	// +optional
	CoreDNSVersion string `json:"coreDNSVersion,omitempty"`
	// FeatureGates is a map of feature names to bools that enable or disable
	// alpha/experimental or optional features.
	// +optional
	FeatureGates map[string]bool `json:"featureGates,omitempty"`
	// EtcdConfiguration provides configuration for the client etcd connection
	// used by the apiserver.
	// +optional
	EtcdConfiguration EtcdConfiguration `json:"etcd"`
	// KubeAPIServerConfiguration provides configuration for the kube-apiserver
	// static pod.
	// +optional
	KubeAPIServerConfiguration KubeAPIServerConfiguration `json:"kubeAPIServer"`
	// KubeControllerManagerConfiguration provides configuration for the
	// kube-controller-manager static pod.
	// +optional
	KubeControllerManagerConfiguration KubeControllerManagerConfiguration `json:"kubeControllerManager"`
	// KubeSchedulerConfiguration provides configuration for the kube-scheduler
	// static pod.
	// +optional
	KubeSchedulerConfiguration KubeSchedulerConfiguration `json:"kubeScheduler"`
	// KubeProxyConfiguration provides configuration for the kube-proxy
	// daemonset.
	// +optional
	KubeProxyConfiguration KubeProxyConfiguration `json:"kubeProxy"`
	// CritBootstrapServerConfiguration provides configuration for the
	// crit-bootstrap-server static pod.
	// +optional
	CritBootstrapServerConfiguration CritBootstrapServerConfiguration `json:"critBootstrapServer"`
	// NodeConfiguration provides configuration for the particular node being
	// bootstrapped. This includes host-specific information, such as hostname
	// or IP address, as well as, kubelet configuration.
	NodeConfiguration NodeConfiguration `json:"node"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type WorkerConfiguration struct {
	metav1.TypeMeta `json:",inline"`
	// ClusterName
	// Default: "crit"
	// +optional
	ClusterName string `json:"clusterName"`
	// ControlPlaneEndpoint is the IP address or DNS name that represents the
	// control plane, along with optional port. The host portion is
	// automatically added to the cluster CA SANs.
	// +optional
	ControlPlaneEndpoint computil.APIEndpoint `json:"controlPlaneEndpoint,omitempty"`
	// FeatureGates is a map of feature names to bools that enable or disable
	// alpha/experimental or optional features.
	// +optional
	FeatureGates map[string]bool `json:"featureGates"`
	// BootstrapServerURL is the full URL to the crit-bootstrap-server static
	// pod. This should only be specified for a cluster using the
	// bootstrap-server, otherwise a bootstrap token should be provided.
	// +optional
	BootstrapServerURL string `json:"bootstrapServerURL,omitempty"`
	// BootstrapToken is the Kubernetes bootstrap auth token used for
	// bootstrapping a worker to a control plane.
	// The token format is described here:
	//   https://kubernetes.io/docs/reference/access-authn-authz/bootstrap-tokens/#token-format
	// +optional
	BootstrapToken string `json:"bootstrapToken,omitempty"`
	// CACert is the full file path of the cluster CA certificate. This must be
	// provided during bootstrapping because it is used to verify that the
	// control plane being joined by the worker.
	CACert string `json:"caCert,omitempty"`
	// NodeConfiguration provides configuration for the particular node being
	// bootstrapped. This includes host-specific information, such as hostname
	// or IP address, as well as, kubelet configuration.
	NodeConfiguration NodeConfiguration `json:"node"`
}

type NodeConfiguration struct {
	// KubernetesVersion is the version of Kubernetes for this node.
	KubernetesVersion string `json:"kubernetesVersion,omitempty"`
	// Hostname is the hostname for this node. This defaults to the hostname
	// provided by the host. It is unlikely that this will need to be changed.
	// +optional
	Hostname string `json:"hostname,omitempty"`
	// KubeDir is the base directory for important Kubernetes configuration
	// files (manifests, configuration files, pki, etc). It is unlikely this
	// will ever need to be changed.
	// Default: "/etc/kubernetes
	// +optional
	KubeDir string `json:"kubeDir,omitempty"`
	// HostIPv4 is the IPv4 address of the host for the node being
	// bootstrapped. If this is not provided the first non-loopback network
	// adapter address is used.
	// +optional
	HostIPv4 string `json:"hostIPv4,omitempty"`
	// CloudProvider is used to configured in-tree cloud providers.
	// +optional
	CloudProvider string `json:"cloudProvider,omitempty"`
	// ContainerRuntime is the container runtime being used by the Kubelet.
	// Default: "containerd"
	// +optional
	ContainerRuntime constants.ContainerRuntime `json:"containerRuntime,omitempty"`
	// Taints is any taints to be applied to the node after initial
	// bootstrapping.
	// +optional
	Taints []corev1.Taint `json:"taints,omitempty"`
	// KubeletConfiguration is the component config for the kubelet. There are
	// quite a few defaults that are being set that can be found in defaults.go
	// of this package.
	// +optional
	KubeletConfiguration *kubeletconfigv1beta1.KubeletConfiguration `json:"kubelet,omitempty"`
	// KubeletExtraArgs is a map of arguments to provide to the kubelet binary.
	// This is useful for settings that are not available in the component
	// config. It should not be used to set deprecated flags that have been
	// moved into the component config.
	// +optional
	KubeletExtraArgs map[string]string `json:"kubeletExtraArgs,omitempty"`
}

type EtcdConfiguration struct {
	Endpoints []string `json:"endpoints,omitempty"`
	CAFile    string   `json:"caFile,omitempty"`
	CertFile  string   `json:"certFile,omitempty"`
	KeyFile   string   `json:"keyFile,omitempty"`

	// CAKey is the etcd CA private key. It is only used to encrypt e2db
	// tables, so any file containing data to be used as a secret can be
	// provided here to enable e2db table encryption for shared cluster files.
	CAKey string `json:"caKey,omitempty"`
}

func (ec *EtcdConfiguration) ClientAddr() string {
	for _, ep := range ec.Endpoints {
		u, _ := url.Parse(ep)
		return u.Host
	}
	return "127.0.0.1:2379"
}

type CritBootstrapServerConfiguration struct {
	Version       string            `json:"version,omitempty"`
	BindPort      int               `json:"bindPort,omitempty"`
	CloudProvider string            `json:"cloudProvider,omitempty"`
	ExtraArgs     map[string]string `json:"extraArgs,omitempty"`
}

type KubeAPIServerConfiguration struct {
	BindPort                 int                      `json:"bindPort,omitempty"`
	ExtraArgs                map[string]string        `json:"extraArgs,omitempty"`
	ExtraVolumes             []computil.HostPathMount `json:"extraVolumes,omitempty"`
	FeatureGates             map[string]bool          `json:"featureGates,omitempty"`
	ExtraSANs                []string                 `json:"extraSans,omitempty"`
	HealthcheckProxyVersion  string                   `json:"healthcheckProxyVersion,omitempty"`
	HealthcheckProxyBindPort int                      `json:"healthcheckProxyBindPort,omitempty"`
	ExtraLabels              map[string]string        `json:"extraLabels,omitempty"`
}

type KubeControllerManagerConfiguration struct {
	ExtraArgs    map[string]string        `json:"extraArgs,omitempty"`
	ExtraVolumes []computil.HostPathMount `json:"extraVolumes,omitempty"`
	FeatureGates map[string]bool          `json:"featureGates,omitempty"`
	ExtraLabels  map[string]string        `json:"extraLabels,omitempty"`
}

type KubeSchedulerConfiguration struct {
	ExtraArgs    map[string]string        `json:"extraArgs,omitempty"`
	ExtraVolumes []computil.HostPathMount `json:"extraVolumes,omitempty"`
	FeatureGates map[string]bool          `json:"featureGates,omitempty"`
	ExtraLabels  map[string]string        `json:"extraLabels,omitempty"`
}

type KubeProxyConfiguration struct {
	// NOTE(chrism): KubeProxyConfiguration defines fields using types from
	// component-base. These contain float values and the package
	// controller-tools, used by controller-gen, did not provide a way to
	// encode float types until commit b45abdb. It requires using the crd flag
	// `allowDangerousTypes` for controller-gen to work with embedding
	// KubeProxyConfiguration. It is possible that this should be re-evaluated
	// in the future.
	//
	// https://github.com/kubernetes-sigs/controller-tools/issues/245
	Config   *kubeproxyconfigv1alpha1.KubeProxyConfiguration `json:"config,omitempty"`
	Disabled bool                                            `json:"disabled"`
	Affinity *corev1.Affinity                                `json:"affinity,omitempty"`
}
