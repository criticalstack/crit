package v1alpha1

import (
	"fmt"
	"net/url"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clientsetscheme "k8s.io/client-go/kubernetes/scheme"
	kubeletconfigv1beta1 "k8s.io/kubelet/config/v1beta1"
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
	ClusterName string `json:"clusterName,omitempty"`
	// ControlPlaneEndpoint is the IP address or DNS name that represents the
	// control plane.
	// +optional
	ControlPlaneEndpoint string `json:"controlPlaneEndpoint,omitempty"`

	// APIServerURL is the kube-apiserver URL.
	// +optional
	APIServerURL   string `json:"apiServerURL,omitempty"`
	CoreDNSVersion string `json:"coreDNSVersion,omitempty"`
	Verbosity      int    `json:"verbosity,omitempty"`

	CritBootstrapServerConfiguration   CritBootstrapServerConfiguration   `json:"critBootstrapServer,omitempty"`
	EtcdConfiguration                  EtcdConfiguration                  `json:"etcd,omitempty"`
	KubeAPIServerConfiguration         KubeAPIServerConfiguration         `json:"kubeAPIServer,omitempty"`
	KubeControllerManagerConfiguration KubeControllerManagerConfiguration `json:"kubeControllerManager,omitempty"`
	KubeSchedulerConfiguration         KubeSchedulerConfiguration         `json:"kubeScheduler,omitempty"`
	NodeConfiguration                  NodeConfiguration                  `json:"node,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type WorkerConfiguration struct {
	metav1.TypeMeta `json:",inline"`

	ClusterName          string            `json:"clusterName,omitempty"`
	ControlPlaneEndpoint string            `json:"controlPlaneEndpoint,omitempty"`
	APIServerURL         string            `json:"apiServerURL,omitempty"`
	BootstrapServerURL   string            `json:"bootstrapServerURL,omitempty"`
	BootstrapToken       string            `json:"bootstrapToken,omitempty"`
	CACert               string            `json:"caCert,omitempty"`
	NodeConfiguration    NodeConfiguration `json:"node,omitempty"`
}

type NodeConfiguration struct {
	KubernetesVersion    string                                     `json:"kubernetesVersion,omitempty"`
	Hostname             string                                     `json:"hostname,omitempty"`
	KubeDir              string                                     `json:"kubeDir,omitempty"`
	HostIPv4             string                                     `json:"hostIPv4,omitempty"`
	KubeProxyMode        string                                     `json:"kubeProxyMode,omitempty"`
	DNSDomain            string                                     `json:"dnsDomain,omitempty"`
	PodSubnet            string                                     `json:"podSubnet,omitempty"`
	ServiceSubnet        string                                     `json:"serviceSubnet,omitempty"`
	CloudProvider        string                                     `json:"cloudProvider,omitempty"`
	ContainerRuntime     ContainerRuntime                           `json:"containerRuntime,omitempty"`
	Taints               []corev1.Taint                             `json:"taints,omitempty"`
	KubeletConfiguration *kubeletconfigv1beta1.KubeletConfiguration `json:"kubelet,omitempty"`
	KubeletExtraArgs     map[string]string                          `json:"kubeletExtraArgs,omitempty"`
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

type ContainerRuntime string

const (
	Containerd ContainerRuntime = "containerd"
	Docker     ContainerRuntime = "docker"
	CRIO       ContainerRuntime = "crio"
)

func (cr ContainerRuntime) CRISocket() string {
	switch cr {
	case Containerd:
		return "unix:///var/run/containerd/containerd.sock"
	case Docker:
		return "unix:///var/run/docker.sock"
	case CRIO:
		return "unix:///var/run/crio/crio.sock"
	default:
		panic("unknown container runtime")
	}
}

// HostPathMount contains elements describing volumes that are mounted from the
// host.
type HostPathMount struct {
	// Name of the volume inside the pod template.
	Name string `json:"name,omitempty"`
	// HostPath is the path in the host that will be mounted inside
	// the pod.
	HostPath string `json:"hostPath,omitempty"`
	// MountPath is the path inside the pod where hostPath will be mounted.
	MountPath string `json:"mountPath,omitempty"`
	// ReadOnly controls write access to the volume
	ReadOnly bool `json:"readOnly,omitempty"`

	HostPathType corev1.HostPathType `json:"hostPathType,omitempty"`
}

type CritBootstrapServerConfiguration struct {
	// TODO(chrism): maybe add option to disable
	ImageRegistry string            `json:"imageRegistry,omitempty"`
	Version       string            `json:"version,omitempty"`
	BindPort      int               `json:"bindPort,omitempty"`
	CloudProvider string            `json:"cloudProvider,omitempty"`
	ExtraArgs     map[string]string `json:"extraArgs,omitempty"`
}

type KubeAPIServerConfiguration struct {
	BindPort     int               `json:"bindPort,omitempty"`
	ExtraArgs    map[string]string `json:"extraArgs,omitempty"`
	ExtraVolumes []HostPathMount   `json:"extraVolumes,omitempty"`
	FeatureGates map[string]bool   `json:"featureGates,omitempty"`
	ExtraSANs    []string          `json:"extraSans,omitempty"`
	ExtraLabels  map[string]string `json:"extraLabels,omitempty"`
}

type KubeControllerManagerConfiguration struct {
	ExtraArgs    map[string]string `json:"extraArgs,omitempty"`
	ExtraVolumes []HostPathMount   `json:"extraVolumes,omitempty"`
	FeatureGates map[string]bool   `json:"featureGates,omitempty"`
	ExtraLabels  map[string]string `json:"extraLabels,omitempty"`
}

type KubeSchedulerConfiguration struct {
	ExtraArgs    map[string]string `json:"extraArgs,omitempty"`
	ExtraVolumes []HostPathMount   `json:"extraVolumes,omitempty"`
	FeatureGates map[string]bool   `json:"featureGates,omitempty"`
	ExtraLabels  map[string]string `json:"extraLabels,omitempty"`
}

// APIEndpoint represents a reachable Kubernetes API endpoint.
type APIEndpoint struct {
	// The hostname on which the API server is serving.
	Host string `json:"host"`

	// The port on which the API server is serving.
	Port int `json:"port"`
}

// IsZero returns true if host and the port are zero values.
func (v APIEndpoint) IsZero() bool {
	return v.Host == "" && v.Port == 0
}

// String returns a formatted version HOST:PORT of this APIEndpoint.
func (v APIEndpoint) String() string {
	return fmt.Sprintf("%s:%d", v.Host, v.Port)
}
