package app

import (
	"github.com/spf13/cobra"
	"go.uber.org/zap/zapcore"

	e2dapp "github.com/criticalstack/e2d/cmd/e2d/app"

	"github.com/criticalstack/crit/cmd/crit/app/certs"
	"github.com/criticalstack/crit/cmd/crit/app/config"
	"github.com/criticalstack/crit/cmd/crit/app/create"
	"github.com/criticalstack/crit/cmd/crit/app/generate"
	"github.com/criticalstack/crit/cmd/crit/app/template"
	"github.com/criticalstack/crit/cmd/crit/app/up"
	"github.com/criticalstack/crit/cmd/crit/app/version"
	"github.com/criticalstack/crit/pkg/log"
)

var global struct {
	Verbosity int
}

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "crit",
		Short: "bootstrap Critical Stack clusters",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if global.Verbosity > 0 {
				log.SetLevel(zapcore.DebugLevel)
			}
		},
	}

	// add e2d app as subcommand
	cmd.AddCommand(
		e2dapp.NewRootCmd(),
	)

	cmd.AddCommand(
		certs.NewCommand(),
		config.NewCommand(),
		create.NewCommand(),
		generate.NewCommand(),
		template.NewCommand(),
		up.NewCommand(),
		version.NewCommand(),
	)

	cmd.PersistentFlags().CountVarP(&global.Verbosity, "verbose", "v", "log output verbosity")
	return cmd
}
