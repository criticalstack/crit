package cluster

import (
	"context"
	"encoding/base64"
	"os"
	"strings"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kubeletconfigv1beta1 "k8s.io/kubelet/config/v1beta1"
	"k8s.io/utils/pointer"
	kindconstants "sigs.k8s.io/kind/pkg/cluster/constants"

	"github.com/criticalstack/crit/internal/cinder/config"
	"github.com/criticalstack/crit/internal/cinder/config/constants"
	"github.com/criticalstack/crit/internal/cinder/feature"
	critconfig "github.com/criticalstack/crit/internal/config"
	yamlutil "github.com/criticalstack/crit/pkg/kubernetes/yaml"
)

var (
	preCritCommands = []string{
		"e2d pki gencerts --ca-cert /etc/kubernetes/pki/etcd/ca.crt --ca-key /etc/kubernetes/pki/etcd/ca.key --output-dir /etc/kubernetes/pki/etcd",
		"systemctl enable --now e2d",
	}
	postCritCommands = []string{
		"kubectl taint nodes --all node-role.kubernetes.io/master-",
	}
)

type ControlPlaneConfig struct {
	ClusterName          string
	ContainerName        string
	Image                string
	Verbose              bool
	ClusterConfiguration *config.ClusterConfiguration
}

func CreateControlPlaneNode(ctx context.Context, cfg *ControlPlaneConfig) (*Node, error) {
	node, err := CreateNode(ctx, &Config{
		ClusterName:       cfg.ClusterName,
		ContainerName:     cfg.ContainerName,
		Image:             cfg.Image,
		Role:              kindconstants.ControlPlaneNodeRoleValue,
		ExtraMounts:       cfg.ClusterConfiguration.ExtraMounts,
		ExtraPortMappings: cfg.ClusterConfiguration.ExtraPortMappings,
	})
	if err != nil {
		return nil, err
	}
	if cfg.Verbose {
		node.Stdout = os.Stdout
		node.Stderr = os.Stderr
	}
	if cfg.ClusterConfiguration.ControlPlaneConfiguration == nil {
		cfg.ClusterConfiguration.ControlPlaneConfiguration = &critconfig.ControlPlaneConfiguration{}
	}
	cfg.ClusterConfiguration.ControlPlaneConfiguration.ClusterName = cfg.ClusterName
	SetControlPlaneConfigurationDefaults(cfg.ClusterConfiguration.ControlPlaneConfiguration)
	data, err := node.ReadFile("/cinder/kubernetes_version")
	if err != nil {
		return nil, err
	}
	cfg.ClusterConfiguration.ControlPlaneConfiguration.NodeConfiguration.KubernetesVersion = strings.TrimSpace(string(data))
	data, err = yamlutil.MarshalToYaml(cfg.ClusterConfiguration.ControlPlaneConfiguration, critconfig.SchemeGroupVersion)
	if err != nil {
		return nil, err
	}
	cfg.ClusterConfiguration.Files = append(cfg.ClusterConfiguration.Files, config.File{
		Path:        "/var/lib/crit/config.yaml",
		Owner:       "root:root",
		Permissions: "0644",
		Encoding:    config.Base64,
		Content:     base64.StdEncoding.EncodeToString(data),
	})
	if feature.Gates.Enabled(feature.FixCgroupMounts) {
		if err := node.Command("bash", "/cinder/scripts/fix-cgroup-mounts.sh").Run(); err != nil {
			return node, err
		}
	}
	if err := node.SystemdReady(ctx); err != nil {
		return node, err
	}
	cfg.ClusterConfiguration.PreCritCommands = append(preCritCommands, cfg.ClusterConfiguration.PreCritCommands...)
	if shouldRestartContainerd(cfg.ClusterConfiguration.Files) {
		cfg.ClusterConfiguration.PreCritCommands = append([]string{"systemctl restart containerd"}, cfg.ClusterConfiguration.PreCritCommands...)
	}
	cfg.ClusterConfiguration.PostCritCommands = append(postCritCommands, cfg.ClusterConfiguration.PostCritCommands...)
	if err := node.RunCloudInit(cfg.ClusterConfiguration); err != nil {
		return node, err
	}
	return node, nil
}

func SetControlPlaneConfigurationDefaults(cfg *critconfig.ControlPlaneConfiguration) {
	// this will be overriden by the kubernetes version file in the image
	cfg.NodeConfiguration.KubernetesVersion = constants.KubernetesVersion
	if cfg.KubeAPIServerConfiguration.ExtraArgs == nil {
		cfg.KubeAPIServerConfiguration.ExtraArgs = make(map[string]string)
	}
	if _, ok := cfg.KubeAPIServerConfiguration.ExtraArgs["enable-admission-plugins"]; !ok {
		cfg.KubeAPIServerConfiguration.ExtraArgs["enable-admission-plugins"] = "NodeRestriction"
	}
	if cfg.KubeControllerManagerConfiguration.ExtraArgs == nil {
		cfg.KubeControllerManagerConfiguration.ExtraArgs = make(map[string]string)
	}
	cfg.KubeControllerManagerConfiguration.ExtraArgs["enable-hostpath-provisioner"] = "true"
	if cfg.NodeConfiguration.KubeletConfiguration == nil {
		cfg.NodeConfiguration.KubeletConfiguration = &kubeletconfigv1beta1.KubeletConfiguration{}
	}
	cfg.NodeConfiguration.KubeletConfiguration.FileCheckFrequency = metav1.Duration{Duration: 1 * time.Second}
	cfg.NodeConfiguration.KubeletConfiguration.HTTPCheckFrequency = metav1.Duration{Duration: 1 * time.Second}
	cfg.NodeConfiguration.KubeletConfiguration.NodeStatusReportFrequency = metav1.Duration{Duration: 1 * time.Second}
	cfg.NodeConfiguration.KubeletConfiguration.NodeStatusUpdateFrequency = metav1.Duration{Duration: 1 * time.Second}
	cfg.NodeConfiguration.KubeletConfiguration.SyncFrequency = metav1.Duration{Duration: 1 * time.Second}
	cfg.NodeConfiguration.KubeletConfiguration.SerializeImagePulls = pointer.BoolPtr(false)
	if feature.Gates.Enabled(feature.FixCgroupMounts) {
		cfg.NodeConfiguration.KubeletConfiguration.CgroupRoot = "/kubelet"
	}
}
