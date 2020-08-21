package v1alpha2

import (
	"testing"

	yamlutil "github.com/criticalstack/crit/pkg/kubernetes/yaml"
)

const (
	apiEndpointString = `apiVersion: crit.sh/v1alpha2
kind: ControlPlaneConfiguration
controlPlaneEndpoint: "example.com:6443"`

	apiEndpointObject = `apiVersion: crit.sh/v1alpha2
kind: ControlPlaneConfiguration
controlPlaneEndpoint:
  host: example.com
  port: 6443`
)

func TestAPIEndpointUnmarshal(t *testing.T) {
	obj, err := yamlutil.UnmarshalFromYaml([]byte(apiEndpointString), SchemeGroupVersion)
	if err != nil {
		t.Fatal(err)
	}
	cfg, _ := obj.(*ControlPlaneConfiguration)
	if cfg.ControlPlaneEndpoint.Host != "example.com" {
		t.Fatalf("APIEndpoint unmarshaled incorrectly: %v", cfg.ControlPlaneEndpoint.Host)
	}
	if cfg.ControlPlaneEndpoint.Port != 6443 {
		t.Fatalf("APIEndpoint unmarshaled incorrectly: %v", cfg.ControlPlaneEndpoint.Port)
	}
	obj, err = yamlutil.UnmarshalFromYaml([]byte(apiEndpointObject), SchemeGroupVersion)
	if err != nil {
		t.Fatal(err)
	}
	cfg, _ = obj.(*ControlPlaneConfiguration)
	if cfg.ControlPlaneEndpoint.Host != "example.com" {
		t.Fatalf("APIEndpoint unmarshaled incorrectly: %v", cfg.ControlPlaneEndpoint.Host)
	}
	if cfg.ControlPlaneEndpoint.Port != 6443 {
		t.Fatalf("APIEndpoint unmarshaled incorrectly: %v", cfg.ControlPlaneEndpoint.Port)
	}
}
