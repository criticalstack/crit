package version

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/criticalstack/crit/internal/buildinfo"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:           "version",
		Short:         "Print the version info",
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(cmd *cobra.Command, args []string) error {
			data, err := json.Marshal(map[string]string{
				"Date":         buildinfo.Date,
				"GitSHA":       buildinfo.GitSHA,
				"GitTreeState": buildinfo.GitTreeState,
				"GoVersion":    buildinfo.GoVersion,
				"Version":      buildinfo.Version,
			})
			if err != nil {
				return err
			}
			fmt.Printf("%s\n", data)
			return nil
		},
	}
	return cmd
}
