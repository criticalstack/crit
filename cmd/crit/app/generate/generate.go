package generate

import (
	"github.com/spf13/cobra"

	generatehash "github.com/criticalstack/crit/cmd/crit/app/generate/hash"
	generatekubeconfig "github.com/criticalstack/crit/cmd/crit/app/generate/kubeconfig"
	generatetoken "github.com/criticalstack/crit/cmd/crit/app/generate/token"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "generate",
		Aliases: []string{"gen"},
		Short:   "Utilities for generating values",
	}
	cmd.AddCommand(
		generatehash.NewCommand(),
		generatetoken.NewCommand(),
		generatekubeconfig.NewCommand(),
	)
	return cmd
}
