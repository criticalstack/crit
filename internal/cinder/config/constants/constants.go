package constants

import (
	"fmt"

	"github.com/criticalstack/crit/pkg/config/constants"
)

const (
	DefaultNodeImage  = "criticalstack/cinder:v1"
	DefaultNetwork    = "cinder"
	KubernetesVersion = "1.18.5"

	DefaultMachineAPIVersion               = "1.0.1"
	DefaultMachineAPIProviderDockerVersion = "1.0.2"
	DefaultKubeRBACProxyVersion            = "0.5.0"
	DefaultCiliumVersion                   = "1.8.1"
	DefaultCiliumStartupScriptVersion      = "af2a99046eca96c0138551393b21a5c044c7fe79"
	DefaultLocalPathProvisionerVersion     = "0.0.12"
	DefaultRegistryVersion                 = "2.7.1"

	DefaultLocalRegistryName = "cinderegg"
	DefaultLocalRegistryPort = 5000
)

func GetImages() map[string]string {
	images := map[string]string{
		"kube-apiserver":              fmt.Sprintf("%s:v%s", constants.KubeAPIServerImage, KubernetesVersion),
		"kube-controller-manager":     fmt.Sprintf("%s:v%s", constants.KubeControllerManagerImage, KubernetesVersion),
		"kube-scheduler":              fmt.Sprintf("%s:v%s", constants.KubeSchedulerImage, KubernetesVersion),
		"kube-proxy":                  fmt.Sprintf("%s:v%s", constants.KubeProxyImage, KubernetesVersion),
		"pause":                       fmt.Sprintf("%s:%s", constants.PauseImage, constants.DefaultPauseImageVersion),
		"coredns":                     fmt.Sprintf("%s:%s", constants.CoreDNSImage, constants.DefaultCoreDNSVersion),
		"bootstrap-server":            fmt.Sprintf("%s:v%s", constants.CritBootstrapServerImage, constants.DefaultBootstrapServerVersion),
		"healthcheck-proxy":           fmt.Sprintf("%s:v%s", constants.CritHealthCheckProxyImage, constants.DefaultHealthcheckProxyVersion),
		"machine-api":                 fmt.Sprintf("docker.io/criticalstack/machine-api:v%s", DefaultMachineAPIVersion),
		"machine-api-provider-docker": fmt.Sprintf("docker.io/criticalstack/machine-api-provider-docker:v%s", DefaultMachineAPIProviderDockerVersion),
		"kube-rbac-proxy":             fmt.Sprintf("gcr.io/kubebuilder/kube-rbac-proxy:v%s", DefaultKubeRBACProxyVersion),
		"cilium":                      fmt.Sprintf("docker.io/cilium/cilium:v%s", DefaultCiliumVersion),
		"cilium-operator-generic":     fmt.Sprintf("docker.io/cilium/operator-generic:v%s", DefaultCiliumVersion),
		"cilium-startup-script":       fmt.Sprintf("docker.io/cilium/startup-script:%s", DefaultCiliumStartupScriptVersion),
		"local-path-provisioner":      fmt.Sprintf("docker.io/rancher/local-path-provisioner:v%s", DefaultLocalPathProvisionerVersion),
		"registry":                    fmt.Sprintf("docker.io/library/registry:%s", DefaultRegistryVersion),
	}
	return images
}
