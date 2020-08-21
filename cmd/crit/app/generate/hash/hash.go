package hash

import (
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/criticalstack/crit/pkg/kubernetes/pki"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:           "hash [ca-cert-path]",
		Args:          cobra.ExactArgs(1),
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			data, err := pki.GenerateCertHashFromFile(args[0])
			if err != nil {
				return err
			}
			fmt.Printf("sha256:%s", strings.ToLower(hex.EncodeToString(data)))
			return nil
		},
	}
	return cmd
}
