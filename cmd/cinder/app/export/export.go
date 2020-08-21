package export

import (
	"github.com/spf13/cobra"

	exportkubeconfig "github.com/criticalstack/crit/cmd/cinder/app/export/kubeconfig"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Args:  cobra.NoArgs,
		Use:   "export",
		Short: "Export from local cluster",
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}
	cmd.AddCommand(
		exportkubeconfig.NewCommand(),
		// TODO(chrism): tail file and tail log command
	)
	return cmd
}
