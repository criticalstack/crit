package cluster

import (
	"context"
	"encoding/base64"
	"os"
	"strings"
	"time"

	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kubeletconfigv1beta1 "k8s.io/kubelet/config/v1beta1"
	"k8s.io/utils/pointer"
	kindconstants "sigs.k8s.io/kind/pkg/cluster/constants"

	"github.com/criticalstack/crit/internal/cinder/config"
	"github.com/criticalstack/crit/internal/cinder/config/constants"
	critconfig "github.com/criticalstack/crit/internal/config"
	yamlutil "github.com/criticalstack/crit/pkg/kubernetes/yaml"
)

type WorkerConfig struct {
	ClusterName          string
	ContainerName        string
	Image                string
	Verbose              bool
	ClusterConfiguration *config.ClusterConfiguration
}

func CreateWorkerNode(ctx context.Context, cfg *WorkerConfig) (*Node, error) {
	node, err := CreateNode(ctx, &Config{
		ClusterName:       cfg.ClusterName,
		ContainerName:     cfg.ContainerName,
		Image:             cfg.Image,
		Role:              kindconstants.WorkerNodeRoleValue,
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
	data, err := node.ReadFile("/cinder/kubernetes_version")
	if err != nil {
		return nil, err
	}
	cfg.ClusterConfiguration.WorkerConfiguration.NodeConfiguration.KubernetesVersion = strings.TrimSpace(string(data))
	data, err = yamlutil.MarshalToYaml(cfg.ClusterConfiguration.WorkerConfiguration, critconfig.SchemeGroupVersion)
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
	if err := node.Command("bash", "/cinder/scripts/fix-cgroup-mounts.sh").Run(); err != nil {
		return nil, err
	}
	if err := node.SystemdReady(ctx); err != nil {
		return nil, err
	}
	if shouldRestartContainerd(cfg.ClusterConfiguration.Files) {
		cfg.ClusterConfiguration.PreCritCommands = append([]string{"systemctl restart containerd"}, cfg.ClusterConfiguration.PreCritCommands...)
	}
	if err := node.RunCloudInit(cfg.ClusterConfiguration); err != nil {
		return nil, err
	}
	return node, nil
}

func SetWorkerConfigurationDefaults(cfg *critconfig.WorkerConfiguration) {
	cfg.NodeConfiguration.KubernetesVersion = constants.KubernetesVersion
	if cfg.NodeConfiguration.KubeletConfiguration == nil {
		cfg.NodeConfiguration.KubeletConfiguration = &kubeletconfigv1beta1.KubeletConfiguration{}
	}
	cfg.NodeConfiguration.KubeletConfiguration.FileCheckFrequency = metav1.Duration{Duration: 1 * time.Second}
	cfg.NodeConfiguration.KubeletConfiguration.HTTPCheckFrequency = metav1.Duration{Duration: 1 * time.Second}
	cfg.NodeConfiguration.KubeletConfiguration.NodeStatusReportFrequency = metav1.Duration{Duration: 1 * time.Second}
	cfg.NodeConfiguration.KubeletConfiguration.NodeStatusUpdateFrequency = metav1.Duration{Duration: 1 * time.Second}
	cfg.NodeConfiguration.KubeletConfiguration.SyncFrequency = metav1.Duration{Duration: 1 * time.Second}
	cfg.NodeConfiguration.KubeletConfiguration.SerializeImagePulls = pointer.BoolPtr(false)
	cfg.NodeConfiguration.KubeletConfiguration.CgroupRoot = "/kubelet"
}

func GetControlPlaneNode(clusterName string) (*Node, error) {
	nodes, err := ListNodes(clusterName)
	if err != nil {
		return nil, err
	}
	for _, node := range nodes {
		role, err := node.Role()
		if err != nil {
			return nil, err
		}
		if role != kindconstants.ControlPlaneNodeRoleValue {
			continue
		}
		return node, nil
	}
	return nil, errors.Errorf("cannot find cluster: %q", clusterName)
}
