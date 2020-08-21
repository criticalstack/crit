package cluster

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/util/wait"
	"sigs.k8s.io/kind/pkg/cluster"
	kindconstants "sigs.k8s.io/kind/pkg/cluster/constants"
	kindcmd "sigs.k8s.io/kind/pkg/cmd"

	"github.com/criticalstack/crit/internal/cinder/config"
	"github.com/criticalstack/crit/internal/cinder/config/constants"
)

// clusterLabelKey is applied to each "node" docker container for identification
const clusterLabelKey = "io.x-k8s.kind.cluster"

// nodeRoleLabelKey is applied to each "node" docker container for categorization
// of nodes by role
const nodeRoleLabelKey = "io.x-k8s.kind.role"

type Config struct {
	ClusterName       string
	ContainerName     string
	Image             string
	Role              string
	ExtraMounts       []*config.Mount
	ExtraPortMappings []*config.PortMapping
}

func CreateNode(ctx context.Context, cfg *Config) (node *Node, reterr error) {
	exists, err := checkIfNetworkExists(constants.DefaultNetwork)
	if err != nil {
		return nil, err
	}
	if !exists {
		if err := createNetwork(constants.DefaultNetwork); err != nil {
			return nil, err
		}
	}

	args := []string{
		"run",
		"--detach", // run the container detached
		"--tty",    // allocate a tty for entrypoint logs
		// label the node with the cluster ID
		"--label", fmt.Sprintf("%s=%s", clusterLabelKey, cfg.ClusterName),
		"--hostname", cfg.ContainerName, // make hostname match container name
		"--name", cfg.ContainerName, // ... and set the container name
		// label the node with the role ID
		"--label", fmt.Sprintf("%s=%s", nodeRoleLabelKey, cfg.Role),
		// network to attach container to
		"--network", constants.DefaultNetwork,
		// running containers in a container requires privileged
		// NOTE: we could try to replicate this with --cap-add, and use less
		// privileges, but this flag also changes some mounts that are necessary
		// including some ones docker would otherwise do by default.
		// for now this is what we want. in the future we may revisit this.
		"--privileged",
		"--security-opt", "seccomp=unconfined", // also ignore seccomp
		"--security-opt", "apparmor=unconfined", // also ignore apparmor
		// runtime temporary storage
		"--tmpfs", "/tmp", // various things depend on working /tmp
		"--tmpfs", "/run", // systemd wants a writable /run
		// runtime persistent storage
		// this ensures that E.G. pods, logs etc. are not on the container
		// filesystem, which is not only better for performance, but allows
		// running kind in kind for "party tricks"
		// (please don't depend on doing this though!)
		"--volume", "/var",
		// some k8s things want to read /lib/modules
		"--volume", "/lib/modules:/lib/modules:ro",
	}
	providerName := "docker"
	if v, ok := os.LookupEnv("KIND_EXPERIMENTAL_PROVIDER"); ok {
		providerName = v
	}
	if cfg.Role == kindconstants.ControlPlaneNodeRoleValue {
		switch providerName {
		case "docker":
			args = append(args, "-p", "127.0.0.1:0:6443")

			// allow docker to be used inside of the "host" container
			args = append(args, "--volume", "/var/run/docker.sock:/var/run/docker.sock")
		case "podman":
			// Podman expects empty string instead of 0 to assign a random port
			// https://github.com/containers/libpod/blob/master/pkg/spec/ports.go#L68-L69
			args = append(args, "-p", "127.0.0.1::6443")

			// allow podman to be used inside of the "host" container
			args = append(args, "--volume", "/var/run/podman/podman.sock:/var/run/podman/podman.sock")
		}
	}
	if len(cfg.ExtraMounts) > 0 {
		for _, m := range cfg.ExtraMounts {
			if !filepath.IsAbs(m.HostPath) {
				m.HostPath, _ = filepath.Abs(m.HostPath)
			}
			bind := fmt.Sprintf("%s:%s", m.HostPath, m.ContainerPath)
			if len(m.Attrs) > 0 {
				bind = fmt.Sprintf("%s:%s", bind, strings.Join(m.Attrs, ","))
			}
			args = append(args, fmt.Sprintf("--volume=%s", bind))
		}
	}
	if len(cfg.ExtraPortMappings) > 0 {
		for _, pm := range cfg.ExtraPortMappings {
			switch providerName {
			case "docker":
				args = append(args, "-p", fmt.Sprintf("%s:%d:%d/%s", pm.ListenAddress, pm.HostPort, pm.ContainerPort, pm.Protocol))
			case "podman":
				var hostPort string
				if pm.HostPort != 0 {
					hostPort = strconv.Itoa(int(pm.HostPort))
				}
				args = append(args, "-p", fmt.Sprintf("%s:%s:%d/%s", pm.ListenAddress, hostPort, pm.ContainerPort, pm.Protocol))
			}
		}
	}
	args = append(args, cfg.Image)
	cmd := exec.Command(providerName, args...)
	if data, err := cmd.CombinedOutput(); err != nil {
		fmt.Printf("cannot create node: %s\n", data)
		return nil, errors.Wrap(err, "docker run error")
	}
	logger := kindcmd.NewLogger()
	provider := cluster.NewProvider(
		cluster.ProviderWithLogger(logger),
		GetDefault(logger),
	)
	if err := wait.PollImmediateUntil(500*time.Millisecond, func() (bool, error) {
		nodes, err := provider.ListNodes(cfg.ClusterName)
		if err != nil {
			return false, nil
		}
		for _, n := range nodes {
			if n.String() == cfg.ContainerName {
				node = NewNode(n)
				return true, nil
			}
		}
		return false, nil
	}, ctx.Done()); err != nil {
		return nil, err
	}
	return node, nil
}
