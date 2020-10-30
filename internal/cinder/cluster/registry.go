package cluster

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"text/template"

	"github.com/pkg/errors"
	kindexec "sigs.k8s.io/kind/pkg/exec"

	"github.com/criticalstack/crit/internal/cinder/config"
	"github.com/criticalstack/crit/internal/cinder/config/constants"
)

func CreateRegistry(name string, port int) error {
	args := []string{
		"run",
		"--detach", // run the container detached
		"--restart", "always",
		"--name", name,
		"-p", fmt.Sprintf("%d:5000", port),
		"--network", constants.DefaultNetwork,
		fmt.Sprintf("docker.io/library/registry:%s", constants.RegistryVersion),
	}
	providerName := "docker"
	if v, ok := os.LookupEnv("KIND_EXPERIMENTAL_PROVIDER"); ok {
		providerName = v
	}
	cmd := exec.Command(providerName, args...)
	if data, err := cmd.CombinedOutput(); err != nil {
		fmt.Printf("cannot create registry: %s\n", data)
		return errors.Wrap(err, "docker run error")
	}
	return nil
}

func IsContainerRunning(name string) bool {
	_, err := kindexec.Output(kindexec.Command(
		"docker", "inspect", "-f",
		"{{ .State.Running }}", name,
	))
	return err == nil
}

const LocalRegistryHostingConfigMap = `apiVersion: v1
kind: ConfigMap
metadata:
  name: local-registry-hosting
  namespace: kube-public
data:
  localRegistryHosting.v1: |
    host: "localhost:{{ .LocalRegistryPort }}"
    hostFromContainerRuntime: "{{ .LocalRegistryName }}:{{ .LocalRegistryPort }}"
    hostFromClusterNetwork: "{{ .LocalRegistryName }}:{{ .LocalRegistryPort }}"
    help: "https://docs.crit.sh/cinder-guide/local-registry.html"`

func GetLocalRegistryHostingConfigMap(cfg *config.ClusterConfiguration) ([]byte, error) {
	t, err := template.New("").Parse(LocalRegistryHostingConfigMap)
	if err != nil {
		return nil, err
	}
	var out bytes.Buffer
	if err := t.Execute(&out, cfg); err != nil {
		return nil, errors.Wrap(err, "failed to generate Local Registry Hosting Config Map template")
	}
	return out.Bytes(), nil
}

const registryMirrorTmpl = `{{ range $key, $value := . }}
[plugins."io.containerd.grpc.v1.cri".registry.mirrors."{{ $key }}"]
  endpoint = ["{{ $value }}"]
{{ end }}`

func GetRegistryMirrors(mirrors map[string]string) ([]byte, error) {
	t, err := template.New("").Parse(registryMirrorTmpl)
	if err != nil {
		return nil, err
	}
	var out bytes.Buffer
	if err := t.Execute(&out, mirrors); err != nil {
		return nil, errors.Wrap(err, "failed to generate registry mirrors template")
	}
	return out.Bytes(), nil
}
