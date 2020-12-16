package components

import (
	"fmt"
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
)

func NewSchedulerStaticPod(cfg *config.ControlPlaneConfiguration) *corev1.Pod {
	kubeconfigFile := filepath.Join(cfg.NodeConfiguration.KubeDir, "scheduler.conf")

	defaultArguments := map[string]string{
		"bind-address":              "127.0.0.1",
		"leader-elect":              "true",
		"kubeconfig":                kubeconfigFile,
		"authentication-kubeconfig": kubeconfigFile,
		"authorization-kubeconfig":  kubeconfigFile,
	}

	for k, v := range cfg.KubeSchedulerConfiguration.FeatureGates {
		FeatureGates[k] = v
	}
	featureGates := make([]string, 0)
	for k, v := range FeatureGates {
		featureGates = append(featureGates, fmt.Sprintf("%s=%t", k, v))
	}
	defaultArguments["feature-gates"] = strings.Join(featureGates, ",")

	command := []string{"kube-scheduler"}
	command = append(command, computil.BuildArgumentListFromMap(defaultArguments, cfg.KubeSchedulerConfiguration.ExtraArgs)...)
	p := &corev1.Pod{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "Pod",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "kube-scheduler",
			Namespace: metav1.NamespaceSystem,
			Labels: map[string]string{
				"component": "kube-scheduler",
				"tier":      "control-plane",
			},
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:            "kube-scheduler",
					Image:           fmt.Sprintf("k8s.gcr.io/kube-scheduler:v%s", cfg.NodeConfiguration.KubernetesVersion),
					ImagePullPolicy: corev1.PullIfNotPresent,
					Command:         command,
					VolumeMounts: []corev1.VolumeMount{
						{
							Name:      "kubeconfig",
							MountPath: "/etc/kubernetes/scheduler.conf",
							ReadOnly:  true,
						},
					},
					LivenessProbe: &corev1.Probe{
						Handler: corev1.Handler{
							HTTPGet: &corev1.HTTPGetAction{
								Host:   cfg.NodeConfiguration.HostIPv4,
								Path:   "/healthz",
								Port:   intstr.FromInt(10251),
								Scheme: corev1.URISchemeHTTP,
							},
						},
						InitialDelaySeconds: 15,
						TimeoutSeconds:      15,
						FailureThreshold:    8,
					},
					Resources: corev1.ResourceRequirements{
						Requests: corev1.ResourceList{
							corev1.ResourceName(corev1.ResourceCPU): resource.MustParse("100m"),
						},
					},
					Env: computil.GetProxyEnvVars(),
				},
			},
			PriorityClassName: "system-cluster-critical",
			HostNetwork:       true,
			Volumes: []corev1.Volume{
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

	if err := appendExtraVolumes(p, cfg.KubeSchedulerConfiguration.ExtraVolumes); err != nil {
		log.Debug("scheduler extra volumes", zap.Error(err))
	}

	if err := appendExtraLabels(p, cfg.KubeSchedulerConfiguration.ExtraLabels); err != nil {
		log.Info("scheduler extra labels", zap.Error(err))
	}
	return p
}
