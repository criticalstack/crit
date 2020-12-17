package components

import (
	"fmt"
	"path/filepath"
	"strings"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"

	"github.com/criticalstack/crit/internal/config"
	computil "github.com/criticalstack/crit/pkg/cluster/components/util"
	"github.com/criticalstack/crit/pkg/kubernetes/util/pointer"
	"github.com/criticalstack/crit/pkg/log"
)

var (
	AdmissionPlugins = []string{
		"NodeRestriction",
		"PodPreset",
	}
	FeatureGates = map[string]bool{
		"TTLAfterFinished": true,
	}
)

func NewAPIServerStaticPod(cfg *config.ControlPlaneConfiguration) (*corev1.Pod, error) {
	certsDir := filepath.Join(cfg.NodeConfiguration.KubeDir, "pki")

	defaultArguments := map[string]string{
		"advertise-address":               cfg.NodeConfiguration.HostIPv4,
		"insecure-port":                   "0",
		"enable-admission-plugins":        "NodeRestriction",
		"service-cluster-ip-range":        cfg.ServiceSubnet,
		"service-account-key-file":        filepath.Join(certsDir, "sa.pub"),
		"client-ca-file":                  filepath.Join(certsDir, "ca.crt"),
		"tls-cert-file":                   filepath.Join(certsDir, "apiserver.crt"),
		"tls-private-key-file":            filepath.Join(certsDir, "apiserver.key"),
		"kubelet-client-certificate":      filepath.Join(certsDir, "apiserver-kubelet-client.crt"),
		"kubelet-client-key":              filepath.Join(certsDir, "apiserver-kubelet-client.key"),
		"enable-bootstrap-token-auth":     "true",
		"secure-port":                     fmt.Sprintf("%d", cfg.KubeAPIServerConfiguration.BindPort),
		"allow-privileged":                "true",
		"kubelet-preferred-address-types": "InternalIP,ExternalIP,Hostname",
		// add options to configure the front proxy.  Without the generated client cert, this will never be useable
		// so add it unconditionally with recommended values
		"requestheader-username-headers":     "X-Remote-User",
		"requestheader-group-headers":        "X-Remote-Group",
		"requestheader-extra-headers-prefix": "X-Remote-Extra-",
		"requestheader-client-ca-file":       filepath.Join(certsDir, "front-proxy-ca.crt"),
		"requestheader-allowed-names":        "front-proxy-client",
		"proxy-client-cert-file":             filepath.Join(certsDir, "front-proxy-client.crt"),
		"proxy-client-key-file":              filepath.Join(certsDir, "front-proxy-client.key"),
		"etcd-servers":                       strings.Join(cfg.EtcdConfiguration.Endpoints, ","),
		"etcd-cafile":                        cfg.EtcdConfiguration.CAFile,
		"etcd-certfile":                      cfg.EtcdConfiguration.CertFile,
		"etcd-keyfile":                       cfg.EtcdConfiguration.KeyFile,
	}

	modes := []string{"Node", "RBAC"}
	if v, ok := cfg.KubeAPIServerConfiguration.ExtraArgs["authorization-mode"]; ok {
		switch v {
		case "ABAC", "Webhook":
			modes = append(modes, v)
		}
	}

	if _, ok := cfg.KubeAPIServerConfiguration.ExtraArgs["secure-port"]; ok {
		delete(cfg.KubeAPIServerConfiguration.ExtraArgs, "secure-port")
		log.Warn(`ignoring apiserver extraArgs "secure-port", use BindPort instead`)
	}
	defaultArguments["authorization-mode"] = strings.Join(modes, ",")
	defaultArguments["enable-admission-plugins"] = strings.Join(AdmissionPlugins, ",")
	for k, v := range cfg.KubeAPIServerConfiguration.FeatureGates {
		FeatureGates[k] = v
	}
	featureGates := make([]string, 0)
	for k, v := range FeatureGates {
		featureGates = append(featureGates, fmt.Sprintf("%s=%t", k, v))
	}
	defaultArguments["feature-gates"] = strings.Join(featureGates, ",")
	defaultArguments["runtime-config"] = "settings.k8s.io/v1alpha1=true"

	if cfg.NodeConfiguration.CloudProvider != "" {
		defaultArguments["cloud-provider"] = cfg.NodeConfiguration.CloudProvider
	}

	command := []string{"kube-apiserver"}
	command = append(command, computil.BuildArgumentListFromMap(defaultArguments, cfg.KubeAPIServerConfiguration.ExtraArgs)...)

	p := &corev1.Pod{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "Pod",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "kube-apiserver",
			Namespace: metav1.NamespaceSystem,
			Labels: map[string]string{
				"component": "kube-apiserver",
				"tier":      "control-plane",
			},
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:            "kube-apiserver",
					Image:           fmt.Sprintf("k8s.gcr.io/kube-apiserver:v%s", cfg.NodeConfiguration.KubernetesVersion),
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
					},
					LivenessProbe: &corev1.Probe{
						Handler: corev1.Handler{
							HTTPGet: &corev1.HTTPGetAction{
								Host:   cfg.NodeConfiguration.HostIPv4,
								Path:   "/healthz",
								Port:   intstr.FromInt(cfg.KubeAPIServerConfiguration.BindPort),
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
					Name: "ca-certs",
					VolumeSource: corev1.VolumeSource{
						HostPath: &corev1.HostPathVolumeSource{
							Path: "/etc/ssl/certs",
							Type: pointer.HostPathTypePtr(corev1.HostPathDirectoryOrCreate),
						},
					},
				},
			},
		},
	}

	// Liveness probes will fail for static pods should anonymous-auth be set
	// to false, so a special healthcheck-proxy sidecar is added to the
	// apiserver static pod. It acts as a reverse proxy with the frontend
	// effectively accepting anonymous traffic and the backend using an
	// authenticated user. The backend connection is established with the
	// built-in system:basic-info-viewer user to limit the auth to only being
	// able to look at health and version information.
	if cfg.KubeAPIServerConfiguration.ExtraArgs["anonymous-auth"] == "false" {
		p.Spec.Containers = append(p.Spec.Containers, corev1.Container{
			Name:  "crit-healthcheck-proxy",
			Image: fmt.Sprintf("docker.io/criticalstack/healthcheck-proxy:v%s", cfg.KubeAPIServerConfiguration.HealthcheckProxyVersion),
			Command: append([]string{"/healthcheck-proxy"}, computil.BuildArgumentListFromMap(map[string]string{
				"client-ca-file":                 filepath.Join(certsDir, "ca.crt"),
				"tls-cert-file":                  filepath.Join(certsDir, "apiserver.crt"),
				"tls-private-key-file":           filepath.Join(certsDir, "apiserver.key"),
				"healthcheck-client-certificate": filepath.Join(certsDir, "apiserver-healthcheck-client.crt"),
				"healthcheck-client-key":         filepath.Join(certsDir, "apiserver-healthcheck-client.key"),
				"secure-port":                    fmt.Sprintf("%d", cfg.KubeAPIServerConfiguration.HealthcheckProxyBindPort),
				"apiserver-port":                 fmt.Sprintf("%d", cfg.KubeAPIServerConfiguration.BindPort),
			}, nil)...),
			ImagePullPolicy: corev1.PullIfNotPresent,
			VolumeMounts: []corev1.VolumeMount{
				{
					Name:      "k8s-certs",
					MountPath: certsDir,
					ReadOnly:  true,
				},
			},
			Resources: corev1.ResourceRequirements{
				Requests: corev1.ResourceList{
					corev1.ResourceName(corev1.ResourceCPU): resource.MustParse("250m"),
				},
			},
			Env: computil.GetProxyEnvVars(),
		})
		p.Spec.Containers[0].LivenessProbe.Handler.HTTPGet.Port = intstr.FromInt(cfg.KubeAPIServerConfiguration.HealthcheckProxyBindPort)
	}

	p.Spec.Volumes = append(p.Spec.Volumes, getCACertsExtraVolumes()...)
	p.Spec.Containers[0].VolumeMounts = append(p.Spec.Containers[0].VolumeMounts, getCACertsExtraVolumeMounts()...)

	if err := appendExtraVolumes(p, cfg.KubeAPIServerConfiguration.ExtraVolumes); err != nil {
		return nil, err
	}

	appendExtraLabels(p, cfg.KubeAPIServerConfiguration.ExtraLabels)

	return p, nil
}
