package create

import (
	"github.com/spf13/cobra"

	createtoken "github.com/criticalstack/crit/cmd/crit/app/create/token"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create Kubernetes resources",
	}

	cmd.AddCommand(
		createtoken.NewCommand(),
	)
	return cmd
}
