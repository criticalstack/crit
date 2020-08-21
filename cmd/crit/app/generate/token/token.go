package token

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/criticalstack/crit/pkg/kubernetes/pki"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:           "token [token]",
		Short:         "generates a bootstrap token",
		Args:          cobra.MinimumNArgs(1),
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			id, secret := pki.GenerateBootstrapToken()
			fmt.Printf("%s.%s", id, secret)
			return nil
		},
	}
	return cmd
}
