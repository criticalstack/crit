package components

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	clientsetscheme "k8s.io/client-go/kubernetes/scheme"
	kubeletconfigv1beta1 "k8s.io/kubelet/config/v1beta1"

	"github.com/criticalstack/crit/internal/config"
	computil "github.com/criticalstack/crit/pkg/cluster/components/util"
	yamlutil "github.com/criticalstack/crit/pkg/kubernetes/yaml"
	"github.com/criticalstack/crit/pkg/log"
	processutil "github.com/criticalstack/crit/pkg/util/process"
)

func init() {
	_ = kubeletconfigv1beta1.AddToScheme(clientsetscheme.Scheme)
}

func WriteKubeletConfigFile(kc *kubeletconfigv1beta1.KubeletConfiguration, path string) error {
	data, err := yamlutil.MarshalToYaml(kc, kubeletconfigv1beta1.SchemeGroupVersion)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return err
	}
	return ioutil.WriteFile(path, data, 0644)
}

const (
	KubeletEnvFileName = "crit-flags.env"

	// KubeletEnvFileVariableName specifies the environment variable name crit
	// should write a value to dynamic environment file
	KubeletEnvFileVariableName = "KUBELET_CRIT_ARGS"
)

// WriteKubeletDynamicEnvFile writes an environment file with dynamic flags to
// the kubelet. This is identical to how kubeadm creates a dynamic environment
// file, but with names changed. This will most likely be changed in the future
// to fit better with crit usage.
func WriteKubeletDynamicEnvFile(cfg *config.NodeConfiguration, registerTaintsUsingFlags bool, kubeletDir string) error {
	kubeletFlags := map[string]string{
		"container-runtime":          "remote",
		"container-runtime-endpoint": cfg.ContainerRuntime.CRISocket(),
	}
	if registerTaintsUsingFlags && cfg.Taints != nil && len(cfg.Taints) > 0 {
		ts := []string{}
		for _, t := range cfg.Taints {
			ts = append(ts, t.String())
		}
		kubeletFlags["register-with-taints"] = strings.Join(ts, ",")
	}

	// If systemd-resolved is running then we can presume that /etc/resolv.conf
	// will be pointing to the systemd-resolved stub resolver. Since the
	// systemd-resolved stub resolver runs on lo, the address:port it specifies
	// will not be available inside of containers. To handle this case the
	// systemd resolv.conf is used.
	running, err := processutil.IsCommandRunning("systemd-resolved")
	if err != nil {
		return err
	}
	if running {
		kubeletFlags["resolv-conf"] = "/run/systemd/resolve/resolv.conf"
	}
	nodeName := cfg.Hostname
	if name, ok := cfg.KubeletExtraArgs["hostname-override"]; ok {
		nodeName = name
	}

	// Make sure the node name we're passed will work with Kubelet
	if nodeName != "" && nodeName != cfg.Hostname {
		log.Infof("setting kubelet hostname-override to %q", nodeName)
		kubeletFlags["hostname-override"] = nodeName
	}
	if cfg.KubeletExtraArgs == nil {
		cfg.KubeletExtraArgs = make(map[string]string)
	}

	if cfg.CloudProvider != "" {
		cfg.KubeletExtraArgs["cloud-provider"] = cfg.CloudProvider
	}
	cfg.KubeletExtraArgs["address"] = cfg.HostIPv4
	cfg.KubeletExtraArgs["node-ip"] = cfg.HostIPv4

	argList := computil.BuildArgumentListFromMap(kubeletFlags, cfg.KubeletExtraArgs)
	envFileContent := fmt.Sprintf("%s=%q\n", KubeletEnvFileVariableName, strings.Join(argList, " "))
	if err := os.MkdirAll(kubeletDir, 0700); err != nil {
		return err
	}
	return ioutil.WriteFile(filepath.Join(kubeletDir, KubeletEnvFileName), []byte(envFileContent), 0644)
}
