package util

import (
	"crypto/x509"
	"fmt"
	"net"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/util/validation"
	netutils "k8s.io/utils/net"

	"github.com/criticalstack/crit/internal/config"
	"github.com/criticalstack/crit/pkg/kubernetes/pki"
	"github.com/criticalstack/crit/pkg/log"
)

func WriteAPIServerCertAndKey(cfg *config.ControlPlaneConfiguration) error {
	dir := filepath.Join(cfg.NodeConfiguration.KubeDir, "pki")
	if exists(filepath.Join(dir, "apiserver.key")) {
		log.Warn("apiserver cert/key already exists")
		return nil
	}

	// advertise address
	advertiseAddress := net.ParseIP(cfg.NodeConfiguration.HostIPv4)
	if advertiseAddress == nil {
		return errors.Errorf("error parsing LocalAPIEndpoint AdvertiseAddress %v: is not a valid textual representation of an IP address",
			cfg.NodeConfiguration.HostIPv4)
	}
	_, svcSubnet, err := net.ParseCIDR(cfg.ServiceSubnet)
	if err != nil {
		return err
	}
	internalAPIServerVirtualIP, err := netutils.GetIndexedIP(svcSubnet, 1)
	if err != nil {
		return errors.Wrapf(err, "unable to get first IP address from the given CIDR (%s)", svcSubnet.String())
	}
	hostname, err := os.Hostname()
	if err != nil {
		return err
	}
	certConfig := &pki.Config{
		CommonName: "kube-apiserver",
		Usages:     []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		AltNames: pki.AltNames{
			DNSNames: []string{
				hostname,
				"kubernetes",
				"kubernetes.default",
				"kubernetes.default.svc",
				fmt.Sprintf("kubernetes.default.svc.%s", cfg.NodeConfiguration.KubeletConfiguration.ClusterDomain),
			},
			IPs: []net.IP{
				internalAPIServerVirtualIP,
				advertiseAddress,
			},
		},
	}

	if cfg.ControlPlaneEndpoint.Host != "" {
		if ip := net.ParseIP(cfg.ControlPlaneEndpoint.Host); ip != nil {
			certConfig.AltNames.IPs = append(certConfig.AltNames.IPs, ip)
		} else {
			certConfig.AltNames.DNSNames = append(certConfig.AltNames.DNSNames, cfg.ControlPlaneEndpoint.Host)
		}
	}

	certSANs := []string{
		"localhost",
		"127.0.0.1",
	}

	certSANs = append(certSANs, cfg.KubeAPIServerConfiguration.ExtraSANs...)

	for _, altname := range certSANs {
		if ip := net.ParseIP(altname); ip != nil {
			certConfig.AltNames.IPs = append(certConfig.AltNames.IPs, ip)
		} else if len(validation.IsDNS1123Subdomain(altname)) == 0 {
			certConfig.AltNames.DNSNames = append(certConfig.AltNames.DNSNames, altname)
		} else if len(validation.IsWildcardDNS1123Subdomain(altname)) == 0 {
			certConfig.AltNames.DNSNames = append(certConfig.AltNames.DNSNames, altname)
		} else {
			fmt.Printf(
				"[certificates] WARNING: '%s' was not added to the 'apiserver.crt' SAN, because it is not a valid IP or RFC-1123 compliant DNS entry\n",
				altname,
			)
		}
	}
	ca, err := pki.LoadCertificateAuthority(dir, "ca")
	if err != nil {
		return err
	}
	kp, err := ca.NewSignedKeyPair("apiserver", certConfig)
	if err != nil {
		return err
	}
	return kp.WriteFiles(dir)
}

func WriteAPIServerKubeletClientCertAndKey(cfg *config.ControlPlaneConfiguration) error {
	dir := filepath.Join(cfg.NodeConfiguration.KubeDir, "pki")
	if exists(filepath.Join(dir, "apiserver-kubelet-client.key")) {
		log.Warn("apiserver-kubelet-client cert/key already exists")
		return nil
	}
	ca, err := pki.LoadCertificateAuthority(dir, "ca")
	if err != nil {
		return err
	}
	kp, err := ca.NewSignedKeyPair("apiserver-kubelet-client", &pki.Config{
		CommonName:   "kube-apiserver-kubelet-client",
		Organization: []string{"system:masters"},
		Usages:       []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
	})
	if err != nil {
		return err
	}
	return kp.WriteFiles(dir)
}

func WriteFrontProxyClientCertAndKey(cfg *config.ControlPlaneConfiguration) error {
	dir := filepath.Join(cfg.NodeConfiguration.KubeDir, "pki")
	if exists(filepath.Join(dir, "front-proxy-client.key")) {
		log.Warn("front-proxy-client cert/key already exists")
		return nil
	}
	ca, err := pki.LoadCertificateAuthority(dir, "front-proxy-ca")
	if err != nil {
		return err
	}
	kp, err := ca.NewSignedKeyPair("front-proxy-client", &pki.Config{
		CommonName: "front-proxy-client",
		Usages:     []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
	})
	if err != nil {
		return err
	}
	return kp.WriteFiles(dir)
}

func WriteAPIServerHealthcheckClientCertAndKey(cfg *config.ControlPlaneConfiguration) error {
	dir := filepath.Join(cfg.NodeConfiguration.KubeDir, "pki")
	if exists(filepath.Join(dir, "apiserver-healthcheck-client.key")) {
		log.Warn("apiserver-healthcheck-client cert/key already exists")
		return nil
	}
	ca, err := pki.LoadCertificateAuthority(dir, "ca")
	if err != nil {
		return err
	}
	kp, err := ca.NewSignedKeyPair("apiserver-healthcheck-client", &pki.Config{
		CommonName: "system:basic-info-viewer",
		Usages:     []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
	})
	if err != nil {
		return err
	}
	return kp.WriteFiles(dir)
}

func WriteClusterCA(dir string) error {
	if exists(filepath.Join(dir, "ca.key")) {
		log.Warn("cluster CA already exists")
		return nil
	}

	ca, err := pki.NewCertificateAuthority("ca", &pki.Config{
		CommonName: "kubernetes",
	})
	if err != nil {
		return err
	}
	return ca.WriteFiles(dir)
}

func WriteFrontProxyCA(dir string) error {
	if exists(filepath.Join(dir, "front-proxy-ca.key")) {
		log.Warn("front proxy CA already exists")
		return nil
	}
	ca, err := pki.NewCertificateAuthority("front-proxy-ca", &pki.Config{
		CommonName: "front-proxy-ca",
	})
	if err != nil {
		return err
	}
	return ca.WriteFiles(dir)
}

// WriteAuthProxyCA creates a new CA key/cert pair in the provided directory
// named auth-proxy-ca.{crt,key}. This CA is intended for use with settings
// such as the oidc-ca-file flags, and is generated ahead of time because of
// the chicken/egg problem when requiring the CA file be specified during
// cluster bootstrapping and the application used for oidc will be ultimately
// running on the same cluster.
func WriteAuthProxyCA(dir string) error {
	if exists(filepath.Join(dir, "auth-proxy-ca.key")) {
		log.Warn("auth proxy CA already exists")
		return nil
	}
	ca, err := pki.NewCertificateAuthority("auth-proxy-ca", &pki.Config{
		CommonName: "auth-proxy-ca",
	})
	if err != nil {
		return err
	}
	return ca.WriteFiles(dir)
}

func WriteServiceAccountCA(dir string) error {
	if exists(filepath.Join(dir, "sa.key")) {
		log.Warn("service account CA already exists")
		return nil
	}
	key, err := pki.NewPrivateKey()
	if err != nil {
		return err
	}
	if err := pki.WriteKey(dir, "sa", key); err != nil {
		return err
	}
	return pki.WritePublicKey(dir, "sa", key.Public())
}

func exists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}
