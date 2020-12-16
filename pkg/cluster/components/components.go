// Package components contains functions for configuring and creating
// Kubernetes components.
package components

import (
	"os"

	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"

	computil "github.com/criticalstack/crit/pkg/cluster/components/util"
	"github.com/criticalstack/crit/pkg/kubernetes/util/pointer"
)

func appendExtraVolumes(p *corev1.Pod, volumes []computil.HostPathMount) (err error) {
	for _, v := range volumes {
		hostPathType := pointer.HostPathTypePtr(v.HostPathType)
		if v.HostPathType == corev1.HostPathUnset {
			hostPathType, err = pointer.DetectHostPathType(v.HostPath)
			if err != nil {
				return errors.Wrap(err, "cannot determine hostPath type, must be provided")
			}
		}

		p.Spec.Volumes = append(p.Spec.Volumes, corev1.Volume{
			Name: v.Name,
			VolumeSource: corev1.VolumeSource{
				HostPath: &corev1.HostPathVolumeSource{
					Path: v.HostPath,
					Type: hostPathType,
				},
			},
		})
		p.Spec.Containers[0].VolumeMounts = append(p.Spec.Containers[0].VolumeMounts, corev1.VolumeMount{
			Name:      v.Name,
			MountPath: v.MountPath,
			ReadOnly:  v.ReadOnly,
		})
	}
	return nil
}

// caCertsExtraVolumePaths specifies the paths that can be conditionally mounted into the apiserver and controller-manager containers
// as /etc/ssl/certs might be or contain a symlink to them. It's a variable since it may be changed in unit testing. This var MUST
// NOT be changed in normal codepaths during runtime.
var caCertsExtraVolumePaths = map[string]string{
	"etcd-pki":                        "/etc/pki",
	"usr-share-ca-certificates":       "/usr/share/ca-certificates",
	"usr-local-share-ca-certificates": "/usr/local/share/ca-certificates",
	"etc-ca-certificates":             "/etc/ca-certificates",
}

func getCACertsExtraVolumeMounts() []corev1.VolumeMount {
	mounts := make([]corev1.VolumeMount, 0)
	for name, path := range caCertsExtraVolumePaths {
		if _, err := os.Stat(path); err != nil {
			continue
		}
		mounts = append(mounts, corev1.VolumeMount{
			Name:      name,
			MountPath: path,
			ReadOnly:  true,
		})
	}
	return mounts
}

func getCACertsExtraVolumes() []corev1.Volume {
	volumes := make([]corev1.Volume, 0)
	for name, path := range caCertsExtraVolumePaths {
		if _, err := os.Stat(path); err != nil {
			continue
		}
		volumes = append(volumes, corev1.Volume{
			Name: name,
			VolumeSource: corev1.VolumeSource{
				HostPath: &corev1.HostPathVolumeSource{
					Path: path,
					Type: pointer.HostPathTypePtr(corev1.HostPathDirectoryOrCreate),
				},
			},
		})
	}
	return volumes
}

func appendExtraLabels(p *corev1.Pod, labels map[string]string) (err error) {
	for k, v := range labels {
		p.ObjectMeta.Labels[k] = v
	}
	return nil
}
