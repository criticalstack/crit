package configimport

import (
	"fmt"
	"path/filepath"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/criticalstack/crit/pkg/kubeconfig"
)

var opts struct {
	Kubeconfig string
}

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:           "import [kubeconfig]",
		Short:         "import a kubeconfig",
		Args:          cobra.MinimumNArgs(1),
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if opts.Kubeconfig == "" {
				home, err := homedir.Dir()
				if err != nil {
					return err
				}
				opts.Kubeconfig = filepath.Join(home, ".kube/config")
			}
			kc, err := clientcmd.LoadFromFile(args[0])
			if err != nil {
				return err
			}
			if err := kubeconfig.MergeConfigToFile(kc, opts.Kubeconfig); err != nil {
				return err
			}
			fmt.Printf("Set kubectl context to \"%s\"\n", kc.CurrentContext)
			return nil
		},
	}
	cmd.Flags().StringVar(&opts.Kubeconfig, "kubeconfig", "", "sets kubeconfig path instead of $KUBECONFIG or $HOME/.kube/config")
	return cmd
}
