package nodes

import (
	"fmt"

	"github.com/spf13/cobra"
	kindconstants "sigs.k8s.io/kind/pkg/cluster/constants"

	"github.com/criticalstack/crit/internal/cinder/cluster"
)

var opts struct {
	Name       string
	Kubeconfig string
}

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:           "ip",
		Short:         "Get node IP address",
		Args:          cobra.MaximumNArgs(1),
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(cmd *cobra.Command, args []string) error {
			nodes, err := cluster.ListNodes(opts.Name)
			if err != nil {
				return err
			}
			for _, node := range nodes {
				if len(args) == 0 {
					role, err := node.Role()
					if err != nil {
						return err
					}
					if role == kindconstants.ControlPlaneNodeRoleValue {
						fmt.Printf("%s\n", node.IP())
						return nil
					}
					continue
				}
				if node.String() == args[0] {
					fmt.Printf("%s\n", node.IP())
					return nil
				}
			}
			return nil
		},
	}
	cmd.Flags().StringVar(&opts.Name, "name", "cinder", "cluster name")
	return cmd
}
