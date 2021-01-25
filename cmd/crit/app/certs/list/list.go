package list

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/spf13/cobra"

	"github.com/criticalstack/crit/pkg/config/constants"
	"github.com/criticalstack/crit/pkg/kubernetes/pki"
)

var opts struct {
	KubeDir string
}

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:           "list",
		Short:         "list cluster certificates",
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

			w := new(tabwriter.Writer)
			w.Init(os.Stdout, 0, 8, 0, '\t', 0)
			defer w.Flush()

			fmt.Fprintln(w, "Certificate Authorities:")
			fmt.Fprintln(w, "========================")
			fmt.Fprintln(w, strings.Join([]string{
				"Name",
				"CN",
				"Expires",
				"NotAfter",
			}, "\t"))

			for caName := range certTree {
				ca, err := pki.LoadCertificateAuthority(filepath.Join(opts.KubeDir, "pki"), caName)
				if err != nil {
					return err
				}
				fmt.Fprintln(w, strings.Join([]string{
					caName,
					ca.Cert.Subject.CommonName,
					formatDuration(time.Until(ca.Cert.NotAfter)),
					ca.Cert.NotAfter.Format(time.RFC3339),
				}, "\t"))
			}
			fmt.Fprint(w, "\n")

			fmt.Fprintln(w, "Certificates:")
			fmt.Fprintln(w, "=============")
			fmt.Fprintln(w, strings.Join([]string{
				"Name",
				"CN",
				"Expires",
				"NotAfter",
			}, "\t"))
			for _, certs := range certTree {
				for _, certName := range certs {
					cert, err := pki.ReadCertFromFile(filepath.Join(opts.KubeDir, "pki", certName+".crt"))
					if err != nil {
						return err
					}
					fmt.Fprintln(w, strings.Join([]string{
						certName,
						cert.Subject.CommonName,
						formatDuration(time.Until(cert.NotAfter)),
						cert.NotAfter.Format(time.RFC3339),
					}, "\t"))
				}
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&opts.KubeDir, "kube-dir", constants.DefaultKubeDir, "")
	return cmd
}

func formatDuration(d time.Duration) string {
	days := int64(d.Hours() / 24)
	if days > 365 {
		return fmt.Sprintf("%dy", days/365)
	}
	return fmt.Sprintf("%dd", days)

}
