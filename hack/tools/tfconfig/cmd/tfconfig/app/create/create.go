package create

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/criticalstack/crit/hack/tools/tfconfig/pkg/hclutil"
	"github.com/criticalstack/crit/hack/tools/tfconfig/pkg/prompt"
)

var opts struct {
	Name string
}

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "create [path]",
		Short:        "Create a variable definitions file.",
		Args:         cobra.ExactArgs(1),
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			path := filepath.Join(args[0], opts.Name)

			if _, err := os.Stat(path); !os.IsNotExist(err) {
				data, err := ioutil.ReadFile(path)
				if err != nil {
					return err
				}
				fmt.Printf("%s\n", data)
				if prompt.Confirm("use existing config") {
					return nil
				}
			}

			vars, err := hclutil.ReadVars(args[0])
			if err != nil {
				return err
			}

			// present custom prompts to the user first
			for _, v := range vars {
				if fn, ok := customPrompts[v.Name]; ok {
					if err := fn(v); err != nil {
						return err
					}
				}
			}

			// use a generic input for the remaining prompts
			for _, v := range vars {
				if _, ok := customPrompts[v.Name]; ok {
					continue
				}
				if err := prompt.Input(v); err != nil {
					return err
				}
			}
			fmt.Printf("\n%s\n", vars)
			if prompt.Confirmf("write this file to %s", path) {
				return hclutil.WriteVars(path, vars)
			}
			return nil
		},
	}
	cmd.Flags().StringVarP(&opts.Name, "name", "n", "terraform.tfvars", "Name of variable definitions file")
	return cmd
}

var customPrompts = map[string]prompt.CustomValuePrompt{
	"control_plane_size": prompt.NewSelectNumberPrompt(1, 3, 5),
	"kubernetes_version": prompt.NewDeferredSelectPrompt(func() []string {
		versions, err := getKubernetesVersions(5)
		if err != nil {
			panic(err)
		}
		return versions
	}),
}
