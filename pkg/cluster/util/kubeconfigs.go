package util

import (
	"crypto/x509"
	"fmt"
	"path/filepath"

	"github.com/criticalstack/crit/internal/config"
	"github.com/criticalstack/crit/pkg/kubeconfig"
	"github.com/criticalstack/crit/pkg/kubernetes/pki"
)

const (
	AdminCommonName   = "kubernetes-admin"
	AdminFilename     = "admin.conf"
	AdminOrganization = "system:masters"
)

func WriteAdminConfig(cfg *config.ControlPlaneConfiguration, ca *pki.CertificateAuthority) error {
	return writeConfig(cfg, ca, AdminCommonName, AdminFilename, []string{AdminOrganization})
}

const (
	ControllerManagerCommonName = "system:kube-controller-manager"
	ControllerManagerFilename   = "controller-manager.conf"
)

func WriteControllerManagerConfig(cfg *config.ControlPlaneConfiguration, ca *pki.CertificateAuthority) error {
	return writeConfig(cfg, ca, ControllerManagerCommonName, ControllerManagerFilename, nil)
}

const (
	SchedulerCommonName = "system:kube-scheduler"
	SchedulerFilename   = "scheduler.conf"
)

func WriteSchedulerConfig(cfg *config.ControlPlaneConfiguration, ca *pki.CertificateAuthority) error {
	return writeConfig(cfg, ca, SchedulerCommonName, SchedulerFilename, nil)
}

const (
	KubeletCommonNamePrefix = "system:node:"
	KubeletFilename         = "kubelet.conf"
	KubeletOrganization     = "system:nodes"
)

func WriteKubeletConfig(cfg *config.ControlPlaneConfiguration, ca *pki.CertificateAuthority) error {
	return writeConfig(cfg, ca, KubeletCommonNamePrefix+cfg.NodeConfiguration.Hostname, KubeletFilename, []string{KubeletOrganization})
}

func writeConfig(cfg *config.ControlPlaneConfiguration, ca *pki.CertificateAuthority, cn, filename string, orgs []string) error {
	kp, err := ca.NewSignedKeyPair(filename, &pki.Config{
		CommonName:   cn,
		Usages:       []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
		Organization: orgs,
	})
	if err != nil {
		return err
	}
	return kubeconfig.WriteToFile(kubeconfig.NewForClient(
		fmt.Sprintf("https://%s", cfg.ControlPlaneEndpoint),
		cfg.ClusterName,
		cn,
		pki.EncodeCertPEM(ca.Cert),
		pki.EncodeCertPEM(kp.Cert),
		pki.MustEncodePrivateKeyPem(kp.Key),
	), filepath.Join(cfg.NodeConfiguration.KubeDir, filename))
}
