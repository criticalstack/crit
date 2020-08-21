package util

import (
	corev1 "k8s.io/api/core/v1"
)

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
