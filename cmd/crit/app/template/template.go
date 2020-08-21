package template

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/criticalstack/crit/pkg/cluster"
	configutil "github.com/criticalstack/crit/pkg/config/util"
)

var opts struct {
	ConfigFile string
}

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:           "template [path]",
		Short:         "Render embedded assets",
		Args:          cobra.MinimumNArgs(1),
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := configutil.LoadFromFile(opts.ConfigFile)
			if err != nil {
				return err
			}
			path := args[0]
			data, err := cluster.Execute(path, cfg)
			if err != nil {
				return err
			}
			fmt.Print(string(data))
			return nil
		},
	}

	cmd.Flags().StringVarP(&opts.ConfigFile, "config", "c", "config.yaml", "config file")
	return cmd
}
