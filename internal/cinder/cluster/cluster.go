package cluster

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/pkg/errors"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
	"sigs.k8s.io/kind/pkg/cluster"
	kindcmd "sigs.k8s.io/kind/pkg/cmd"
	kindexec "sigs.k8s.io/kind/pkg/exec"
	"sigs.k8s.io/kind/pkg/log"

	"github.com/criticalstack/crit/pkg/kubeconfig"
)

func ListNodes(clusterName string) ([]*Node, error) {
	logger := kindcmd.NewLogger()
	provider := cluster.NewProvider(
		cluster.ProviderWithLogger(logger),
		GetDefault(logger),
	)
	nodeList, err := provider.ListNodes(clusterName)
	if err != nil {
		return nil, err
	}
	nodes := make([]*Node, 0)
	for _, node := range nodeList {
		nodes = append(nodes, NewNode(node))
	}
	return nodes, nil
}

func GetKubeConfig(clusterName string) (*clientcmdapi.Config, error) {
	nodes, err := ListNodes(clusterName)
	if err != nil {
		return nil, err
	}
	for _, node := range nodes {
		if node.String() == clusterName {
			data, err := node.ReadFile("/etc/kubernetes/admin.conf")
			if err != nil {
				return nil, err
			}
			kc, err := clientcmd.Load(data)
			if err != nil {
				return nil, err
			}
			n, err := getPublishedPort(clusterName)
			if err != nil {
				return nil, errors.Wrapf(err, "cannot get published port for cluster %q", clusterName)
			}
			for _, cluster := range kc.Clusters {
				cluster.Server = fmt.Sprintf("https://127.0.0.1:%d", n)
			}
			return kc, nil
		}
	}
	return nil, errors.Errorf("cannot find cluster %q", clusterName)
}

func ExportKubeConfig(clusterName string, kubeconfigPath string) error {
	kc, err := GetKubeConfig(clusterName)
	if err != nil {
		return err
	}
	return kubeconfig.MergeConfigToFile(kc, kubeconfigPath)
}

// DeleteNodes is part of the providers.Provider interface
func DeleteNodes(n []*Node) error {
	if len(n) == 0 {
		return nil
	}
	providerName := "docker"
	if v, ok := os.LookupEnv("KIND_EXPERIMENTAL_PROVIDER"); ok {
		providerName = v
	}
	args := make([]string, 0, len(n)+3) // allocate once
	args = append(args,
		"rm",
		"-f", // force the container to be delete now
		"-v", // delete volumes
	)
	for _, node := range n {
		args = append(args, node.String())
	}
	if err := exec.Command(providerName, args...).Run(); err != nil {
		return errors.Wrap(err, "failed to delete nodes")
	}
	return nil
}

// GetDefault selected the default runtime from the environment override
func GetDefault(logger log.Logger) cluster.ProviderOption {
	switch p := os.Getenv("KIND_EXPERIMENTAL_PROVIDER"); p {
	case "":
		return nil
	case "podman":
		logger.Warn("using podman due to KIND_EXPERIMENTAL_PROVIDER")
		return cluster.ProviderWithPodman()
	case "docker":
		logger.Warn("using docker due to KIND_EXPERIMENTAL_PROVIDER")
		return cluster.ProviderWithDocker()
	default:
		logger.Warnf("ignoring unknown value %q for KIND_EXPERIMENTAL_PROVIDER", p)
		return nil
	}
}

// ImageID return the Id of the container image
func ImageID(containerNameOrID string) (string, error) {
	providerName := "docker"
	if v, ok := os.LookupEnv("KIND_EXPERIMENTAL_PROVIDER"); ok {
		providerName = v
	}
	cmd := kindexec.Command(providerName, "image", "inspect",
		"-f", "{{ .Id }}",
		containerNameOrID, // ... against the container
	)
	lines, err := kindexec.CombinedOutputLines(cmd)
	if err != nil {
		return "", err
	}
	if len(lines) != 1 {
		return "", errors.Errorf("Docker image ID should only be one line, got %d lines", len(lines))
	}
	return lines[0], nil
}

func PullImage(image string) error {
	providerName := "docker"
	if v, ok := os.LookupEnv("KIND_EXPERIMENTAL_PROVIDER"); ok {
		providerName = v
	}
	args := []string{
		"pull",
		image,
	}
	if data, err := exec.Command(providerName, args...).CombinedOutput(); err != nil {
		fmt.Printf("cannot pull image: %s\n", data)
		return errors.Wrapf(err, "failed to pull image: %v", image)
	}
	return nil
}
