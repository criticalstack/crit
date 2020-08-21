package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clientsetscheme "k8s.io/client-go/kubernetes/scheme"

	critconfig "github.com/criticalstack/crit/internal/config"
)

func init() {
	_ = AddToScheme(clientsetscheme.Scheme)
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type ClusterConfiguration struct {
	metav1.TypeMeta `json:",inline"`
	// Files specifies extra files to be passed to user_data upon creation.
	// +optional
	Files []File `json:"files,omitempty"`
	// PreCritCommands specifies extra commands to run before crit runs
	// +optional
	PreCritCommands []string `json:"preCritCommands,omitempty"`
	// PostCritCommands specifies extra commands to run after crit runs
	// +optional
	PostCritCommands []string `json:"postCritCommands,omitempty"`
	// +optional
	FeatureGates map[string]bool `json:"featureGates,omitempty"`
	// +optional
	ControlPlaneConfiguration *critconfig.ControlPlaneConfiguration `json:"controlPlaneConfiguration,omitempty"`
	// +optional
	WorkerConfiguration *critconfig.WorkerConfiguration `json:"workerConfiguration,omitempty"`
	// +optional
	ExtraMounts []*Mount `json:"extraMounts"`
	// +optional
	ExtraPortMappings []*PortMapping `json:"extraPortMappings"`
	// Name of local container registry. Used for DNS resolution.
	// Default: "cinderegg"
	// +optional
	LocalRegistryName string `json:"localRegistryName"`
	// Port of local container registry.
	// Default: 5000
	// +optional
	LocalRegistryPort int `json:"localRegistryPort"`
	// +optional
	RegistryMirrors map[string]string `json:"registryMirrors"`

	// TODO(chrism):
	// add maybe other cloudinit stuff
	// machine-api/machine-api-provider-docker versions
	// cilium/cni options
}

type Mount struct {
	HostPath      string   `json:"hostPath"`
	ContainerPath string   `json:"containerPath"`
	ReadOnly      bool     `json:"readOnly"`
	Attrs         []string `json:"attrs"`
}

// Encoding specifies the cloud-init file encoding.
// +kubebuilder:validation:Enum=base64;gzip;gzip+base64
type Encoding string

const (
	// Base64 implies the contents of the file are encoded as base64.
	Base64 Encoding = "base64"
	// Gzip implies the contents of the file are encoded with gzip.
	Gzip Encoding = "gzip"
	// GzipBase64 implies the contents of the file are first base64 encoded and then gzip encoded.
	GzipBase64 Encoding = "gzip+base64"
	// HostPath implies the contents is a file path that corresponds to an
	// actual file on the host.
	HostPath Encoding = "hostpath"
)

// File defines the input for generating write_files in cloud-init.
type File struct {
	// Path specifies the full path on disk where to store the file.
	Path string `json:"path"`

	// Owner specifies the ownership of the file, e.g. "root:root".
	// +optional
	Owner string `json:"owner,omitempty"`

	// Permissions specifies the permissions to assign to the file, e.g. "0640".
	// +optional
	Permissions string `json:"permissions,omitempty"`

	// Encoding specifies the encoding of the file contents.
	// +optional
	Encoding Encoding `json:"encoding,omitempty"`

	// Content is the actual content of the file.
	Content string `json:"content"`
}

type PortMapping struct {
	ContainerPort int32 `json:"containerPort,omitempty"`
	// +optional
	HostPort int32 `json:"hostPort,omitempty"`
	// +optional
	ListenAddress string `json:"listenAddress,omitempty"`
	// +optional
	Protocol PortMappingProtocol `json:"protocol,omitempty"`
}

// +kubebuilder:validation:Enum=TCP;UDP
type PortMappingProtocol string

const (
	// PortMappingProtocolTCP specifies TCP protocol
	PortMappingProtocolTCP PortMappingProtocol = "TCP"
	// PortMappingProtocolUDP specifies UDP protocol
	PortMappingProtocolUDP PortMappingProtocol = "UDP"
)
