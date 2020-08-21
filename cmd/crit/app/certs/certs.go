package certs

import (
	"github.com/spf13/cobra"

	certsinit "github.com/criticalstack/crit/cmd/crit/app/certs/init"
	certslist "github.com/criticalstack/crit/cmd/crit/app/certs/list"
	certsrenew "github.com/criticalstack/crit/cmd/crit/app/certs/renew"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "certs",
		Short: "Handle Kubernetes certificates",
	}
	cmd.AddCommand(
		certsinit.NewCommand(),
		certsrenew.NewCommand(),
		certslist.NewCommand(),
	)
	return cmd
}
