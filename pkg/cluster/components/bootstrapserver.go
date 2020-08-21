package components

import (
	"fmt"
	"path/filepath"
	"strconv"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"

	"github.com/criticalstack/crit/internal/config"
	computil "github.com/criticalstack/crit/pkg/cluster/components/util"
	"github.com/criticalstack/crit/pkg/kubernetes/util/pointer"
)

func NewBootstrapServerStaticPod(cfg *config.ControlPlaneConfiguration) *corev1.Pod {
	kubeconfigFile := filepath.Join(cfg.NodeConfiguration.KubeDir, "admin.conf")
	certsDir := filepath.Join(cfg.NodeConfiguration.KubeDir, "pki")
	serverPort := 8080

	defaultArguments := map[string]string{
		"cert-file":  filepath.Join(certsDir, "apiserver.crt"),
		"key-file":   filepath.Join(certsDir, "apiserver.key"),
		"kubeconfig": kubeconfigFile,
		"provider":   cfg.CritBootstrapServerConfiguration.CloudProvider,
	}

	if portStr, ok := cfg.CritBootstrapServerConfiguration.ExtraArgs["port"]; ok {
		if port, err := strconv.Atoi(portStr); err == nil {
			serverPort = port
		}
	}

	command := []string{"/bootstrap-server"}
	command = append(command, computil.BuildArgumentListFromMap(defaultArguments, cfg.CritBootstrapServerConfiguration.ExtraArgs)...)

	p := &corev1.Pod{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "Pod",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "bootstrap-server",
			Namespace: metav1.NamespaceSystem,
			Labels: map[string]string{
				"component": "bootstrap-server",
				"tier":      "control-plane",
			},
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:            "bootstrap-server",
					Image:           fmt.Sprintf("docker.io/criticalstack/bootstrap-server:v%s", cfg.CritBootstrapServerConfiguration.Version),
					ImagePullPolicy: corev1.PullIfNotPresent,
					Command:         command,
					VolumeMounts: []corev1.VolumeMount{
						{
							Name:      "k8s-certs",
							MountPath: certsDir,
							ReadOnly:  true,
						},
						{
							Name:      "kubeconfig",
							MountPath: kubeconfigFile,
							ReadOnly:  true,
						},
					},
					LivenessProbe: &corev1.Probe{
						Handler: corev1.Handler{
							HTTPGet: &corev1.HTTPGetAction{
								Host:   cfg.NodeConfiguration.HostIPv4,
								Path:   "/healthz",
								Port:   intstr.FromInt(serverPort),
								Scheme: corev1.URISchemeHTTPS,
							},
						},
						InitialDelaySeconds: 15,
						TimeoutSeconds:      15,
						FailureThreshold:    8,
					},
					Resources: corev1.ResourceRequirements{
						Requests: corev1.ResourceList{
							corev1.ResourceName(corev1.ResourceCPU): resource.MustParse("250m"),
						},
					},
					Env: computil.GetProxyEnvVars(),
				},
			},
			PriorityClassName: "system-cluster-critical",
			HostNetwork:       true,
			Volumes: []corev1.Volume{
				{
					Name: "k8s-certs",
					VolumeSource: corev1.VolumeSource{
						HostPath: &corev1.HostPathVolumeSource{
							Path: certsDir,
							Type: pointer.HostPathTypePtr(corev1.HostPathDirectoryOrCreate),
						},
					},
				},
				{
					Name: "kubeconfig",
					VolumeSource: corev1.VolumeSource{
						HostPath: &corev1.HostPathVolumeSource{
							Path: kubeconfigFile,
							Type: pointer.HostPathTypePtr(corev1.HostPathFileOrCreate),
						},
					},
				},
			},
		},
	}

	return p
}
