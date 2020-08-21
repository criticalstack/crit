package get

import (
	"github.com/spf13/cobra"

	getclusters "github.com/criticalstack/crit/cmd/cinder/app/get/clusters"
	getimages "github.com/criticalstack/crit/cmd/cinder/app/get/images"
	getip "github.com/criticalstack/crit/cmd/cinder/app/get/ip"
	getkubeconfigs "github.com/criticalstack/crit/cmd/cinder/app/get/kubeconfigs"
	getnodes "github.com/criticalstack/crit/cmd/cinder/app/get/nodes"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Args:  cobra.NoArgs,
		Use:   "get",
		Short: "Get cinder resources",
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}
	cmd.AddCommand(
		getclusters.NewCommand(),
		getkubeconfigs.NewCommand(),
		getnodes.NewCommand(),
		getip.NewCommand(),
		getimages.NewCommand(),
	)
	return cmd
}
