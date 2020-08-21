package renew

import (
	"path/filepath"

	"github.com/spf13/cobra"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/criticalstack/crit/pkg/config/constants"
	"github.com/criticalstack/crit/pkg/kubeconfig"
	"github.com/criticalstack/crit/pkg/kubernetes/pki"
)

var opts struct {
	CertDir string
	DryRun  bool
}

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:           "renew",
		Short:         "renew cluster certificates",
		Args:          cobra.NoArgs,
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(cmd *cobra.Command, args []string) error {
			certTree := map[string][]string{
				"ca": {
					"apiserver",
					"apiserver-kubelet-client",
					"apiserver-healthcheck-client",
				},
				"front-proxy-ca": {
					"front-proxy-client",
				},
			}
			for caName, certs := range certTree {
				ca, err := pki.LoadCertificateAuthority(opts.CertDir, caName)
				if err != nil {
					return err
				}
				for _, certName := range certs {
					cert, err := pki.ReadCertFromFile(filepath.Join(opts.CertDir, certName+".crt"))
					if err != nil {
						return err
					}
					kp, err := ca.NewSignedKeyPair(certName, &pki.Config{
						CommonName:   cert.Subject.CommonName,
						Organization: cert.Subject.Organization,
						AltNames: pki.AltNames{
							IPs:      cert.IPAddresses,
							DNSNames: cert.DNSNames,
						},
						Usages: cert.ExtKeyUsage,
					})
					if err != nil {
						return err
					}
					if !opts.DryRun {
						if err := kp.WriteFiles(opts.CertDir); err != nil {
							return err
						}
					}
				}
			}
			kubeconfigs := []string{
				"admin.conf",
				"controller-manager.conf",
				"scheduler.conf",
			}
			for _, configName := range kubeconfigs {
				config, err := clientcmd.LoadFromFile(filepath.Join(opts.CertDir, configName))
				if err != nil {
					return err
				}
				cert, err := kubeconfig.LoadClientCertificateFromConfig(config)
				if err != nil {
					return err
				}
				ca, err := pki.LoadCertificateAuthority(opts.CertDir, "ca")
				if err != nil {
					return err
				}
				kp, err := ca.NewSignedKeyPair(configName, &pki.Config{
					CommonName:   cert.Subject.CommonName,
					Organization: cert.Subject.Organization,
					Usages:       cert.ExtKeyUsage,
				})
				if err != nil {
					return err
				}
				config.AuthInfos[config.Contexts[config.CurrentContext].AuthInfo].ClientCertificateData = pki.EncodeCertPEM(kp.Cert)
				config.AuthInfos[config.Contexts[config.CurrentContext].AuthInfo].ClientKeyData = pki.MustEncodePrivateKeyPem(kp.Key)
				if !opts.DryRun {
					if err := kubeconfig.WriteToFile(config, filepath.Join(opts.CertDir, configName)); err != nil {
						return err
					}
				}
			}
			return nil
		},
	}

	cmd.Flags().BoolVar(&opts.DryRun, "dry-run", false, "")
	cmd.Flags().StringVar(&opts.CertDir, "cert-dir", filepath.Join(constants.DefaultKubeDir, "pki"), "")
	return cmd
}
