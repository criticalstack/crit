package nodes

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/criticalstack/crit/internal/cinder/cluster"
)

var opts struct {
	Name       string
	Kubeconfig string
}

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:           "nodes",
		Short:         "List cinder cluster nodes",
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
					fmt.Printf("%s %s\n", node, node.IP())
					continue
				}
				if node.String() == args[0] {
					fmt.Printf("%s %s\n", node, node.IP())
					return nil
				}
			}
			fmt.Printf("%d node(s) found.\n", len(nodes))
			return nil
		},
	}
	cmd.Flags().StringVar(&opts.Name, "name", "cinder", "cluster name")
	return cmd
}
