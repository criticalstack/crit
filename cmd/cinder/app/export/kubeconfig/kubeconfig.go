package kubeconfig

import (
	"fmt"
	"path/filepath"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"

	"github.com/criticalstack/crit/internal/cinder/cluster"
)

var opts struct {
	Name       string
	Kubeconfig string
}

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:           "kubeconfig",
		Aliases:       []string{"kc"},
		Short:         "Export kubeconfig from cinder cluster and merge with $HOME/.kube/config",
		Args:          cobra.NoArgs,
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if opts.Kubeconfig == "" {
				home, err := homedir.Dir()
				if err != nil {
					return err
				}
				opts.Kubeconfig = filepath.Join(home, ".kube/config")
			}
			if err := cluster.ExportKubeConfig(opts.Name, opts.Kubeconfig); err != nil {
				return err
			}
			fmt.Printf("Set kubectl context to \"kubernetes-admin@%s\"\n", opts.Name)
			return nil
		},
	}
	cmd.Flags().StringVar(&opts.Name, "name", "cinder", "cluster name")
	cmd.Flags().StringVar(&opts.Kubeconfig, "kubeconfig", "", "sets kubeconfig path instead of $KUBECONFIG or $HOME/.kube/config")
	return cmd
}
