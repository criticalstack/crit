package cluster

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/criticalstack/crit/internal/cinder/cluster"
)

var opts struct {
	Name string
}

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Args:          cobra.NoArgs,
		Use:           "cluster",
		Short:         "Deletes a cinder cluster",
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(cmd *cobra.Command, args []string) error {
			n, err := cluster.ListNodes(opts.Name)
			if err != nil {
				return err
			}
			if len(n) == 0 {
				return errors.Errorf("cannot find nodes for a cluster with the name %q", opts.Name)
			}

			fmt.Printf("Deleting cluster %q ...\n", opts.Name)
			if err := cluster.DeleteNodes(n); err != nil {
				return err
			}
			return nil
		},
	}
	cmd.Flags().StringVar(&opts.Name, "name", "cinder", "cluster name")
	return cmd
}
