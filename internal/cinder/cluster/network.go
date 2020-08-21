package cluster

import (
	"bytes"
	"regexp"
	"strconv"
	"strings"

	"sigs.k8s.io/kind/pkg/exec"
)

func createNetwork(name string) error {
	return exec.Command("docker", "network", "create", "-d=bridge",
		"-o", "com.docker.network.bridge.enable_ip_masquerade=true",
		name).Run()
}

func checkIfNetworkExists(name string) (bool, error) {
	out, err := exec.Output(exec.Command(
		"docker", "network", "ls",
		"--filter=name=^"+regexp.QuoteMeta(name)+"$",
		"--format={{.Name}}",
	))
	return strings.HasPrefix(string(out), name), err
}

func getPublishedPort(name string) (int, error) {
	out, err := exec.Output(exec.Command(
		"docker", "inspect",
		"--format={{(index (index .NetworkSettings.Ports \"6443/tcp\") 0).HostPort}}",
		name,
	))
	if err != nil {
		return 0, err
	}
	out = bytes.TrimSuffix(out, []byte("\n"))
	n, err := strconv.Atoi(string(out))
	if err != nil {
		return 0, err
	}
	return n, nil
}
