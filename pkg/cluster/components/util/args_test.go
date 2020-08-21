package util

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestBuildArgumentListFromMap(t *testing.T) {
	defaultArguments := map[string]string{
		"advertise-address":        "127.0.0.1",
		"insecure-port":            "0",
		"enable-admission-plugins": "NodeRestriction",
	}
	extraArguments := map[string]string{
		"enable-admission-plugins": "PodPreset",
		"service-cluster-ip-range": "10.254.0.0/16",
	}
	expected := []string{
		"--advertise-address=127.0.0.1",
		"--enable-admission-plugins=PodPreset",
		"--insecure-port=0",
		"--service-cluster-ip-range=10.254.0.0/16",
	}
	args := BuildArgumentListFromMap(defaultArguments, extraArguments)

	if diff := cmp.Diff(expected, args); diff != "" {
		t.Errorf("wrong result: (-want +got)\n%s", diff)
	}
}
