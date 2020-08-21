package nodes

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/criticalstack/crit/internal/cinder/config/constants"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:           "images",
		Short:         "List all containers images used by cinder",
		Args:          cobra.MaximumNArgs(1),
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(cmd *cobra.Command, args []string) error {
			for _, ref := range constants.GetImages() {
				fmt.Printf("%s\n", ref)
			}
			return nil
		},
	}
	return cmd
}
