package util

import (
	"testing"

	"github.com/criticalstack/crit/internal/config"
)

const (
	v1alpha1Config = `apiVersion: crit.criticalstack.com/v1alpha1
kind: ControlPlaneConfiguration
controlPlaneEndpoint: "example.com"
`
)

func TestConfigConversion(t *testing.T) {
	obj, err := Unmarshal([]byte(v1alpha1Config))
	if err != nil {
		t.Fatal(err)
	}
	cfg, ok := obj.(*config.ControlPlaneConfiguration)
	if !ok {
		t.Fatalf("expected %T, received %T", &config.ControlPlaneConfiguration{}, obj)
	}
	if cfg.KubeProxyConfiguration.Config.Mode != "iptables" {
		t.Fatalf("KubeProxyConfiguration conversion failed")
	}
}
