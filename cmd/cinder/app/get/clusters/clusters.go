package clusters

import (
	"fmt"

	"github.com/spf13/cobra"
	"sigs.k8s.io/kind/pkg/cluster"
	kindcmd "sigs.k8s.io/kind/pkg/cmd"
)

var opts struct {
	Name       string
	Kubeconfig string
}

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:           "clusters",
		Short:         "Get running cluster",
		Args:          cobra.NoArgs,
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(cmd *cobra.Command, args []string) error {
			logger := kindcmd.NewLogger()
			provider := cluster.NewProvider(
				cluster.ProviderWithLogger(logger),
			)
			clusters, err := provider.List()
			if err != nil {
				return err
			}
			for _, cluster := range clusters {
				fmt.Println(cluster)
			}
			return nil
		},
	}
	cmd.Flags().StringVar(&opts.Name, "name", "cinder", "cluster name")
	cmd.Flags().StringVar(&opts.Kubeconfig, "kubeconfig", "", "sets kubeconfig path instead of $KUBECONFIG or $HOME/.kube/config")
	return cmd
}
