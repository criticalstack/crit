package kubeconfig

import (
	"crypto/x509"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/criticalstack/crit/pkg/kubeconfig"
	"github.com/criticalstack/crit/pkg/kubernetes/pki"
)

var opts struct {
	Name         string
	Server       string
	CertName     string
	CertDir      string
	CommonName   string
	Organization string
}

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:           "kubeconfig [filename]",
		Short:         "generates a kubeconfig",
		Args:          cobra.ExactArgs(1),
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			path := args[0]
			if !filepath.IsAbs(path) {
				path, err = filepath.Abs(path)
				if err != nil {
					return err
				}
			}
			ca, err := pki.LoadCertificateAuthority(opts.CertDir, opts.CertName)
			if err != nil {
				return err
			}
			kp, err := ca.NewSignedKeyPair(filepath.Base(path), &pki.Config{
				CommonName:   opts.CommonName,
				Usages:       []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
				Organization: []string{opts.Organization},
			})
			if err != nil {
				return err
			}
			k := kubeconfig.NewForClient(
				opts.Server,
				opts.Name,
				opts.CommonName,
				pki.EncodeCertPEM(ca.Cert),
				pki.EncodeCertPEM(kp.Cert),
				pki.MustEncodePrivateKeyPem(kp.Key),
			)
			return kubeconfig.WriteToFile(k, args[0])
		},
	}

	cmd.Flags().StringVar(&opts.Name, "name", "crit", "")
	cmd.Flags().StringVar(&opts.Server, "server", "", "")
	cmd.Flags().StringVar(&opts.CertName, "cert-name", "ca", "")
	cmd.Flags().StringVar(&opts.CertDir, "cert-dir", ".", "")
	cmd.Flags().StringVar(&opts.CommonName, "CN", "kubernetes-admin", "")
	cmd.Flags().StringVar(&opts.Organization, "O", "system:masters", "")
	return cmd
}
