package cluster

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"

	"k8s.io/apimachinery/pkg/util/wait"
	kindnodes "sigs.k8s.io/kind/pkg/cluster/nodes"
	"sigs.k8s.io/kind/pkg/cluster/nodeutils"
	kindexec "sigs.k8s.io/kind/pkg/exec"

	"github.com/criticalstack/crit/internal/cinder/config"
)

type Node struct {
	kindnodes.Node

	ipGetter sync.Once
	ip       string

	Stdout, Stderr io.Writer
	b              bytes.Buffer
}

func NewNode(node kindnodes.Node) *Node {
	n := &Node{
		Node: node,
	}
	n.Stdout = &n.b
	n.Stderr = &n.b
	return n
}

func (n *Node) CombinedOutput() []byte {
	return n.b.Bytes()
}

func (n *Node) IP() string {
	n.ipGetter.Do(func() {
		_ = wait.PollImmediate(500*time.Millisecond, 2*time.Second, func() (ok bool, err error) {
			n.ip, _, err = n.Node.IP()
			if err != nil {
				return false, nil
			}
			return true, nil
		})
	})
	return n.ip
}

func (n *Node) Command(cmd string, args ...string) kindexec.Cmd {
	return n.Node.Command(cmd, args...).SetStdout(n.Stdout).SetStderr(n.Stderr)
}

func (n *Node) LoadImage(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	return nodeutils.LoadImageArchive(n.Node, f)
}

func (n *Node) ReadFile(path string) ([]byte, error) {
	var b bytes.Buffer
	if err := n.Command("cat", path).SetStdout(&b).Run(); err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

func (n *Node) MkdirAll(path string, perm os.FileMode) error {
	return n.Command("mkdir", "-p", path).Run()
}

func (n *Node) WriteFile(path string, data []byte, perm os.FileMode) error {
	if err := n.Command("sh", "-c", fmt.Sprintf("`cat > %s`", path)).SetStdin(bytes.NewReader(data)).Run(); err != nil {
		return err
	}
	return n.Command("chmod", fmt.Sprintf("%#o", perm), path).Run()
}

func (n *Node) RunCloudInit(cfg *config.ClusterConfiguration) (err error) {
	for _, f := range cfg.Files {
		data, err := config.ReadFile(f)
		if err != nil {
			return err
		}
		mode, err := strconv.ParseUint(f.Permissions, 8, 32)
		if err != nil {
			return err
		}
		if err := n.MkdirAll(filepath.Dir(f.Path), 0755); err != nil {
			return err
		}
		if err := n.WriteFile(f.Path, data, os.FileMode(mode)); err != nil {
			return err
		}
	}
	var b bytes.Buffer
	b.WriteString("#!/bin/bash\n")
	for _, c := range cfg.PreCritCommands {
		b.WriteString(c)
		b.WriteString("\n")
	}
	if err := n.WriteFile("/cinder/scripts/pre-up.sh", b.Bytes(), 0755); err != nil {
		return err
	}
	b.Reset()
	b.WriteString("#!/bin/bash\n")
	for _, c := range cfg.PostCritCommands {
		b.WriteString(c)
		b.WriteString("\n")
	}
	if err := n.WriteFile("/cinder/scripts/post-up.sh", b.Bytes(), 0755); err != nil {
		return err
	}
	if err := n.Command("bash", "/cinder/scripts/pre-up.sh").Run(); err != nil {
		return err
	}
	return n.Command("crit", "up", "-c", "/var/lib/crit/config.yaml").Run()
}

func (n *Node) SystemdReady(ctx context.Context) error {
	return wait.PollImmediate(500*time.Millisecond, 30*time.Second, func() (bool, error) {
		if err := n.Command("systemctl", "is-system-running").Run(); err != nil {
			return false, nil
		}
		return true, nil
	})
}
