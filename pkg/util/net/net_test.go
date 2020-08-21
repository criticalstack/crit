package net_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	netutil "github.com/criticalstack/crit/pkg/util/net"
)

func TestGetIndexedIP(t *testing.T) {
	ip, err := netutil.GetDNSIP("10.254.0.0/16")
	if err != nil {
		t.Fatal(err)
	}
	if diff := cmp.Diff("10.254.0.10", ip.String()); diff != "" {
		t.Errorf("GetDNSIP() mismatch (-want +got):\n%s", diff)
	}
}
