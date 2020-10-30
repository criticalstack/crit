package constants

import (
	"fmt"

	"github.com/criticalstack/crit/pkg/config/constants"
)

var (
	DefaultNodeImage = "criticalstack/cinder:latest"
)

const (
	// the Kubernetes version to beused when building the cinder image
	KubernetesVersion = "1.18.10"

	// Kind base image to use as the base for the cinder image
	KindBaseImage = "docker.io/kindest/base:v20200928-02f74589"

	DefaultNetwork           = "cinder"
	DefaultLocalRegistryName = "cinderegg"
	DefaultLocalRegistryPort = 5000

	// container image versions
	CiliumStartupScriptVersion      = "af2a99046eca96c0138551393b21a5c044c7fe79"
	CiliumVersion                   = "1.8.5"
	KubeRBACProxyVersion            = "0.5.0"
	LocalPathProvisionerVersion     = "0.0.12"
	MachineAPIProviderDockerVersion = "1.0.7"
	MachineAPIVersion               = "1.0.6"
	RegistryVersion                 = "2.7.1"
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
		"machine-api":                 fmt.Sprintf("docker.io/criticalstack/machine-api:v%s", MachineAPIVersion),
		"machine-api-provider-docker": fmt.Sprintf("docker.io/criticalstack/machine-api-provider-docker:v%s", MachineAPIProviderDockerVersion),
		"kube-rbac-proxy":             fmt.Sprintf("gcr.io/kubebuilder/kube-rbac-proxy:v%s", KubeRBACProxyVersion),
		"cilium":                      fmt.Sprintf("docker.io/cilium/cilium:v%s", CiliumVersion),
		"cilium-operator-generic":     fmt.Sprintf("docker.io/cilium/operator-generic:v%s", CiliumVersion),
		"cilium-startup-script":       fmt.Sprintf("docker.io/cilium/startup-script:%s", CiliumStartupScriptVersion),
		"local-path-provisioner":      fmt.Sprintf("docker.io/rancher/local-path-provisioner:v%s", LocalPathProvisionerVersion),
		"registry":                    fmt.Sprintf("docker.io/library/registry:%s", RegistryVersion),
	}
	return images
}
