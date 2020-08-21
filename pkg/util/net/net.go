package net

import (
	"net"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	netutils "k8s.io/utils/net"
)

// DetectHostIPv4 attempts to determine the host IPv4 address by finding the
// first non-loopback device with an assigned IPv4 address.
func DetectHostIPv4() (string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", errors.WithStack(err)
	}
	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() == nil {
				continue
			}
			return ipnet.IP.String(), nil
		}
	}
	return "", errors.New("cannot detect host IPv4 address")
}

// CalcNodeCidrSize determines the size of the subnets used on each node, based
// on the pod subnet provided.  For IPv4, we assume that the pod subnet will
// be /16 and use /24. If the pod subnet cannot be parsed, the IPv4 value will
// be used (/24).
//
// For IPv6, the algorithm will do two three. First, the node CIDR will be set
// to a multiple of 8, using the available bits for easier readability by user.
// Second, the number of nodes will be 512 to 64K to attempt to maximize the
// number of nodes (see NOTE below). Third, pod networks of /113 and larger will
// be rejected, as the amount of bits available is too small.
//
// A special case is when the pod network size is /112, where /120 will be used,
// only allowing 256 nodes and 256 pods.
//
// If the pod network size is /113 or larger, the node CIDR will be set to the same
// size and this will be rejected later in validation.
//
// NOTE: Currently, the design allows a maximum of 64K nodes. This algorithm splits
// the available bits to maximize the number used for nodes, but still have the node
// CIDR be a multiple of eight.
//
// Copied from github.com/kubernetes/kubernetes/cmd/kubeadm/app/phases/controlplane/manifests.go
func CalcNodeCidrSize(podSubnet string) string {
	maskSize := "24"
	if ip, podCidr, err := net.ParseCIDR(podSubnet); err == nil {
		if netutils.IsIPv6(ip) {
			var nodeCidrSize int
			podNetSize, totalBits := podCidr.Mask.Size()
			switch {
			case podNetSize == 112:
				// Special case, allows 256 nodes, 256 pods/node
				nodeCidrSize = 120
			case podNetSize < 112:
				// Use multiple of 8 for node CIDR, with 512 to 64K nodes
				nodeCidrSize = totalBits - ((totalBits-podNetSize-1)/8-1)*8
			default:
				// Not enough bits, will fail later, when validate
				nodeCidrSize = podNetSize
			}
			maskSize = strconv.Itoa(nodeCidrSize)
		}
	}
	return maskSize
}

// GetDNSIP returns a dnsIP, which is 10th IP in svcSubnet CIDR range
//
// Copied from github.com/kubernetes/kubernetes/cmd/kubeadm/app/constants/constants.go
func GetDNSIP(svcSubnet string) (net.IP, error) {
	// Get the service subnet CIDR
	_, svcSubnetCIDR, err := net.ParseCIDR(svcSubnet)
	if err != nil {
		return nil, errors.Wrapf(err, "couldn't parse service subnet CIDR %q", svcSubnet)
	}

	// Selects the 10th IP in service subnet CIDR range as dnsIP
	dnsIP, err := netutils.GetIndexedIP(svcSubnetCIDR, 10)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to get tenth IP address from service subnet CIDR %s", svcSubnetCIDR.String())
	}

	return dnsIP, nil
}

func SplitHostPort(s string) (string, int32, error) {
	if !strings.Contains(s, ":") {
		s = s + ":0"
	}
	host, port, err := net.SplitHostPort(s)
	if err != nil {
		return "", 0, err
	}
	n, err := strconv.ParseInt(port, 10, 32)
	if err != nil {
		return "", 0, err
	}
	return host, int32(n), nil
}
