package kubeconfigs

import (
	"fmt"

	"github.com/spf13/cobra"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/criticalstack/crit/internal/cinder/cluster"
)

var opts struct {
	Name       string
	Kubeconfig string
}

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:           "kubeconfigs",
		Aliases:       []string{"kc"},
		Short:         "Get kubeconfig from cinder cluster",
		Args:          cobra.NoArgs,
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(cmd *cobra.Command, args []string) error {
			kc, err := cluster.GetKubeConfig(opts.Name)
			if err != nil {
				return err
			}
			data, err := clientcmd.Write(*kc)
			if err != nil {
				return err
			}
			fmt.Printf("%s\n", data)
			return nil
		},
	}
	cmd.Flags().StringVar(&opts.Name, "name", "cinder", "cluster name")
	cmd.Flags().StringVar(&opts.Kubeconfig, "kubeconfig", "", "sets kubeconfig path instead of $KUBECONFIG or $HOME/.kube/config")
	return cmd
}
