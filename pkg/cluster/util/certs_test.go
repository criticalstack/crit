package util

import (
	"os"
	"path/filepath"
	"testing"

	kubeletconfigv1beta1 "k8s.io/kubelet/config/v1beta1"

	"github.com/criticalstack/crit/internal/config"
	computil "github.com/criticalstack/crit/pkg/cluster/components/util"
	"github.com/criticalstack/crit/pkg/config/constants"
	"github.com/criticalstack/crit/pkg/kubernetes/pki"
)

func contains(ss []string, match string) bool {
	for _, s := range ss {
		if s == match {
			return true
		}
	}
	return false
}

func TestWriteAPIServerCertAndKey(t *testing.T) {
	if err := os.MkdirAll("testdata", 0755); err != nil {
		t.Fatal(err)
	}
	if err := WriteClusterCA(filepath.Join("testdata", "pki")); err != nil {
		t.Fatal(err)
	}
	cfg := &config.ControlPlaneConfiguration{
		ControlPlaneEndpoint: computil.APIEndpoint{
			Host: "example.com",
			Port: 6443,
		},
		ServiceSubnet: constants.DefaultServiceSubnet,
		NodeConfiguration: config.NodeConfiguration{
			KubeDir:  "testdata",
			HostIPv4: "127.0.0.1",
			KubeletConfiguration: &kubeletconfigv1beta1.KubeletConfiguration{
				ClusterDomain: constants.DefaultClusterDomain,
			},
		},
	}

	err := WriteAPIServerCertAndKey(cfg)
	if err != nil {
		t.Fatal(err)
	}

	kp, err := pki.LoadKeyPair("testdata/pki", "apiserver")
	if err != nil {
		t.Fatal(err)
	}

	if !contains(kp.Cert.DNSNames, "example.com") {
		t.Fatalf("controlPlaneEndpoint was not added to SAN: %v", kp.Cert.DNSNames)
	}
}
