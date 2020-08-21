package delete

import (
	"github.com/spf13/cobra"

	deletecluster "github.com/criticalstack/crit/cmd/cinder/app/delete/cluster"
	deletenode "github.com/criticalstack/crit/cmd/cinder/app/delete/node"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Args:  cobra.NoArgs,
		Use:   "delete",
		Short: "Delete cinder resources",
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}
	cmd.AddCommand(
		deletecluster.NewCommand(),
		deletenode.NewCommand(),
	)
	return cmd
}
