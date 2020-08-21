package config

import (
	"github.com/spf13/cobra"

	configimport "github.com/criticalstack/crit/cmd/crit/app/config/import"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Handle Kubernetes config files",
	}

	cmd.AddCommand(
		configimport.NewCommand(),
	)
	return cmd
}
