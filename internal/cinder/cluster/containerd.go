package cluster

import (
	"bytes"
	"encoding/base64"

	"github.com/criticalstack/crit/internal/cinder/config"
)

const defaultContainerdConfig = `# explicitly use v2 config format
version = 2

# set default runtime handler to v2, which has a per-pod shim
[plugins."io.containerd.grpc.v1.cri".containerd]
  default_runtime_name = "runc"
[plugins."io.containerd.grpc.v1.cri".containerd.runtimes.runc]
  runtime_type = "io.containerd.runc.v2"

# Setup a runtime with the magic name ("test-handler") used for Kubernetes
# runtime class tests ...
[plugins."io.containerd.grpc.v1.cri".containerd.runtimes.test-handler]
  runtime_type = "io.containerd.runc.v2"`

func AppendOrPatchContainerd(cfg *config.ClusterConfiguration, patchData string) error {
	var b bytes.Buffer
	for _, f := range cfg.Files {
		if f.Path == "/etc/containerd/config.toml" {
			data, err := config.ReadFile(f)
			if err != nil {
				return err
			}
			b.Write(data)
			b.WriteString("\n")
			b.WriteString(patchData)
			f.Encoding = config.Base64
			f.Content = base64.StdEncoding.EncodeToString(b.Bytes())
			return nil
		}
	}
	b.WriteString(defaultContainerdConfig)
	b.WriteString("\n")
	b.WriteString(patchData)
	cfg.Files = append(cfg.Files, config.File{
		Path:        "/etc/containerd/config.toml",
		Owner:       "root:root",
		Permissions: "0644",
		Encoding:    config.Base64,
		Content:     base64.StdEncoding.EncodeToString(b.Bytes()),
	})
	return nil
}

func shouldRestartContainerd(files []config.File) bool {
	for _, f := range files {
		if f.Path == "/etc/containerd/config.toml" {
			return true
		}
	}
	return false
}
