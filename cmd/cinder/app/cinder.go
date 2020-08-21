package app

import (
	"github.com/spf13/cobra"

	"github.com/criticalstack/crit/cmd/cinder/app/create"
	"github.com/criticalstack/crit/cmd/cinder/app/delete"
	"github.com/criticalstack/crit/cmd/cinder/app/export"
	"github.com/criticalstack/crit/cmd/cinder/app/get"
	"github.com/criticalstack/crit/cmd/cinder/app/load"
	"github.com/criticalstack/crit/cmd/crit/app/version"
)

var global struct {
	Verbosity int
}

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cinder",
		Short: "Create local Kubernetes clusters",
		Long:  "Cinder is a tool for creating and managing local Kubernetes clusters using containers as nodes. It builds upon kind, but using Crit and Cilium to configure a Critical Stack cluster locally.",
	}
	cmd.AddCommand(
		create.NewCommand(),
		delete.NewCommand(),
		export.NewCommand(),
		get.NewCommand(),
		load.NewCommand(),
		version.NewCommand(),
	)

	cmd.PersistentFlags().CountVarP(&global.Verbosity, "verbose", "v", "log output verbosity")
	return cmd
}
