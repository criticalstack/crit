package cluster

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/criticalstack/crit/internal/cinder/cluster"
)

var opts struct {
	Name string
}

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:           "node",
		Short:         "Deletes a cinder node",
		Args:          cobra.ExactArgs(1),
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(cmd *cobra.Command, args []string) error {
			nodes, err := cluster.ListNodes(opts.Name)
			if err != nil {
				return err
			}
			if len(nodes) == 0 {
				return fmt.Errorf("cannot find node %q", opts.Name)
			}

			for _, node := range nodes {
				if node.String() == args[0] {
					fmt.Printf("Deleting node %q ...\n", opts.Name)
					if err := cluster.DeleteNodes([]*cluster.Node{node}); err != nil {
						return err
					}
				}
			}
			return nil
		},
	}
	cmd.Flags().StringVar(&opts.Name, "name", "cinder", "cluster name")
	return cmd
}
