package create

import (
	"github.com/spf13/cobra"

	createcluster "github.com/criticalstack/crit/cmd/cinder/app/create/cluster"
	createnode "github.com/criticalstack/crit/cmd/cinder/app/create/node"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Args:  cobra.NoArgs,
		Use:   "create",
		Short: "Create cinder resources",
		Long:  "Create cinder resources such a new cinder clusters or add nodes to existing clusters.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}
	cmd.AddCommand(
		createcluster.NewCommand(),
		createnode.NewCommand(),
	)
	return cmd
}
