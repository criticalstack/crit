package components

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"go.uber.org/zap"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"

	"github.com/criticalstack/crit/internal/config"
	computil "github.com/criticalstack/crit/pkg/cluster/components/util"
	"github.com/criticalstack/crit/pkg/kubernetes/util/pointer"
	"github.com/criticalstack/crit/pkg/log"
	netutil "github.com/criticalstack/crit/pkg/util/net"
)

func NewControllerManagerStaticPod(cfg *config.ControlPlaneConfiguration) *corev1.Pod {
	kubeconfigFile := filepath.Join(cfg.NodeConfiguration.KubeDir, "controller-manager.conf")
	certsDir := filepath.Join(cfg.NodeConfiguration.KubeDir, "pki")
	caFile := filepath.Join(certsDir, "ca.crt")

	defaultArguments := map[string]string{
		"bind-address":                     "127.0.0.1",
		"leader-elect":                     "true",
		"kubeconfig":                       kubeconfigFile,
		"authentication-kubeconfig":        kubeconfigFile,
		"authorization-kubeconfig":         kubeconfigFile,
		"client-ca-file":                   caFile,
		"requestheader-client-ca-file":     filepath.Join(certsDir, "front-proxy-ca.crt"),
		"root-ca-file":                     caFile,
		"service-account-private-key-file": filepath.Join(certsDir, "sa.key"),
		"cluster-signing-cert-file":        caFile,
		"cluster-signing-key-file":         filepath.Join(certsDir, "ca.key"),
		"use-service-account-credentials":  "true",
		"controllers":                      "*,bootstrapsigner,tokencleaner",
	}

	// Let the controller-manager allocate Node CIDRs for the Pod network.
	// Each node will get a subspace of the address CIDR provided with --pod-network-cidr.
	if cfg.PodSubnet != "" {
		// TODO(Arvinderpal): Needs to be fixed once PR #73977 lands. Should be a list of maskSizes.
		maskSize := netutil.CalcNodeCidrSize(cfg.PodSubnet)
		defaultArguments["allocate-node-cidrs"] = "true"
		defaultArguments["cluster-cidr"] = cfg.PodSubnet
		defaultArguments["node-cidr-mask-size"] = maskSize
		if cfg.ServiceSubnet != "" {
			defaultArguments["service-cluster-ip-range"] = cfg.ServiceSubnet
		}
	}

	for k, v := range cfg.KubeControllerManagerConfiguration.FeatureGates {
		FeatureGates[k] = v
	}
	featureGates := make([]string, 0)
	for k, v := range FeatureGates {
		featureGates = append(featureGates, fmt.Sprintf("%s=%t", k, v))
	}
	defaultArguments["feature-gates"] = strings.Join(featureGates, ",")

	if cfg.NodeConfiguration.CloudProvider != "" {
		defaultArguments["cloud-provider"] = cfg.NodeConfiguration.CloudProvider
		defaultArguments["allocate-node-cidrs"] = "true"
		// The allocate-node-cidrs flag enables two internal controllers:
		//
		//   1. A controller that assigns a CIDR block to every node, enabling
		//   pod network traffic across nodes.
		//   2. A controller that syncs these routes with the cloud provider
		//   routing tables (requires cloud-provider flag).
		//
		// The second controller (2) is activated assuming that internal
		// cluster traffic routing should be deferred to the cloud provider, but
		// our CNI (Cilium) preempts that. Additionally, the cloud route tables
		// are functionally inconsistent and limited (e.g. AWS only supports 50
		// endpoints).
		//
		// To prevent the controller-manager from trying to set the cloud route
		// tables, configure-cloud-routes is set `false`, ensuring the second
		// controller (2) is never started.
		defaultArguments["configure-cloud-routes"] = "false"
	}

	command := []string{"kube-controller-manager"}
	command = append(command, computil.BuildArgumentListFromMap(defaultArguments, cfg.KubeControllerManagerConfiguration.ExtraArgs)...)

	p := &corev1.Pod{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "Pod",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "kube-controller-manager",
			Namespace: metav1.NamespaceSystem,
			Labels: map[string]string{
				"component": "kube-controller-manager",
				"tier":      "control-plane",
			},
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:            "kube-controller-manager",
					Image:           fmt.Sprintf("k8s.gcr.io/kube-controller-manager:v%s", cfg.NodeConfiguration.KubernetesVersion),
					ImagePullPolicy: corev1.PullIfNotPresent,
					Command:         command,
					VolumeMounts: []corev1.VolumeMount{
						{
							Name:      "k8s-certs",
							MountPath: certsDir,
							ReadOnly:  true,
						},
						{
							Name:      "ca-certs",
							MountPath: "/etc/ssl/certs",
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
								Port:   intstr.FromInt(10252),
								Scheme: corev1.URISchemeHTTP,
							},
						},
						InitialDelaySeconds: 15,
						TimeoutSeconds:      15,
						FailureThreshold:    8,
					},
					Resources: corev1.ResourceRequirements{
						Requests: corev1.ResourceList{
							corev1.ResourceName(corev1.ResourceCPU): resource.MustParse("200m"),
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
					Name: "ca-certs",
					VolumeSource: corev1.VolumeSource{
						HostPath: &corev1.HostPathVolumeSource{
							Path: "/etc/ssl/certs",
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

	// Flexvolume dir must NOT be readonly as it is used for third-party
	// plugins to integrate with their storage backends via unix domain socket.
	if stat, err := os.Stat("/usr/libexec/kubernetes/kubelet-plugins/volume/exec"); err == nil && stat.IsDir() {
		p.Spec.Volumes = append(p.Spec.Volumes, corev1.Volume{
			Name: "flexvolume-dir",
			VolumeSource: corev1.VolumeSource{
				HostPath: &corev1.HostPathVolumeSource{
					Path: "/usr/libexec/kubernetes/kubelet-plugins/volume/exec",
					Type: pointer.HostPathTypePtr(corev1.HostPathDirectoryOrCreate),
				},
			},
		})
		p.Spec.Containers[0].VolumeMounts = append(p.Spec.Containers[0].VolumeMounts, corev1.VolumeMount{
			Name:      "flexvolume-dir",
			MountPath: "/usr/libexec/kubernetes/kubelet-plugins/volume/exec",
			ReadOnly:  false,
		})
	}

	p.Spec.Volumes = append(p.Spec.Volumes, getCACertsExtraVolumes()...)
	p.Spec.Containers[0].VolumeMounts = append(p.Spec.Containers[0].VolumeMounts, getCACertsExtraVolumeMounts()...)

	if err := appendExtraVolumes(p, cfg.KubeControllerManagerConfiguration.ExtraVolumes); err != nil {
		log.Debug("controller manager extra volumes", zap.Error(err))
	}

	appendExtraLabels(p, cfg.KubeControllerManagerConfiguration.ExtraLabels)

	return p
}
