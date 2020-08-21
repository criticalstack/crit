package init

import (
	"os"

	"github.com/spf13/cobra"

	clusterutil "github.com/criticalstack/crit/pkg/cluster/util"
)

var opts struct {
	CertDir string
}

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:           "init",
		Short:         "initialize a new CA",
		Args:          cobra.NoArgs,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := os.MkdirAll(opts.CertDir, 0755); err != nil && !os.IsExist(err) {
				return err
			}
			// add etcd certs as well
			return clusterutil.WriteClusterCA(opts.CertDir)
		},
	}

	cmd.Flags().StringVar(&opts.CertDir, "cert-dir", "", "")
	return cmd
}
