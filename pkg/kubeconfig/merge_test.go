package kubeconfig

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

const (
	kubeConfigFile = `clusters:
- cluster:
    certificate-authority-data: YmFzZTY0
    server: https://172.17.0.2:6443
  name: kind-kind
contexts:
- context:
    cluster: kind-kind
    namespace: kube-system
    user: kind-kind
  name: kind-kind
users:
- name: kind-kind
  user:
    client-certificate-data: YmFzZTY0
    client-key-data: YmFzZTY0`

	kubeConfigFile2 = `apiVersion: v1
clusters:
- cluster:
    certificate-authority-data: YmFzZTY0
    server: https://blerg.com:6443
  name: cinder
- cluster:
    certificate-authority-data: YmFzZTY0
    server: https://172.17.0.2:6443
  name: kind-kind
contexts:
- context:
    cluster: cinder
    user: admin
  name: admin@cinder
- context:
    cluster: kind-kind
    namespace: kube-system
    user: kind-kind
  name: kind-kind
current-context: admin@cinder
kind: Config
preferences: {}
users:
- name: kind-kind
  user:
    client-certificate-data: YmFzZTY0
    client-key-data: YmFzZTY0`
)

func init() {
	if err := os.MkdirAll("testdata", 0755); err != nil {
		panic(err)
	}
}

func TestMergeConfig(t *testing.T) {
	t.Skip()
	if err := ioutil.WriteFile("testdata/kubeconfig1", []byte(kubeConfigFile), 0644); err != nil {
		t.Fatal(err)
	}

	cfg := New("https://blerg.com:6443", "cinder", "admin", []byte("base64"))

	if err := MergeConfigToFile(cfg, "testdata/kubeconfig1"); err != nil {
		t.Fatal(err)
	}

	data, err := ioutil.ReadFile("testdata/kubeconfig1")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("%s\n", data)

}

func TestRemoveConfig(t *testing.T) {
	t.Skip()
	if err := ioutil.WriteFile("testdata/kubeconfig2", []byte(kubeConfigFile2), 0644); err != nil {
		t.Fatal(err)
	}

	if err := RemoveConfigFromFile("cinder", "testdata/kubeconfig2"); err != nil {
		t.Fatal(err)
	}

	data, err := ioutil.ReadFile("testdata/kubeconfig2")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("%s\n", data)

}
