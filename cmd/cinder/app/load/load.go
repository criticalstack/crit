package load

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/kyokomi/emoji"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"sigs.k8s.io/kind/pkg/cluster/nodeutils"
	"sigs.k8s.io/kind/pkg/exec"
	"sigs.k8s.io/kind/pkg/fs"

	"github.com/criticalstack/crit/internal/cinder/cluster"
	"github.com/criticalstack/crit/internal/cinder/utils"
)

var opts struct {
	Name       string
	Image      string
	Config     string
	Kubeconfig string
	Verbose    bool
}

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:           "load",
		Short:         "Load container images from host",
		Args:          cobra.ExactArgs(1),
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(cmd *cobra.Command, args []string) (reterr error) {
			imageName := args[0]
			imageID, err := cluster.ImageID(imageName)
			if err != nil {
				return errors.Errorf("image: %q not present locally", imageName)
			}
			fmt.Printf("Host image found:\n Name:%q\n ID: %q\n", imageName, imageID)

			nodes, err := cluster.ListNodes(opts.Name)
			if err != nil {
				return err
			}

			if len(nodes) == 0 {
				return errors.Errorf("cannot find nodes for a cluster with the name %q", opts.Name)
			}
			fmt.Printf("\nLoading image into %d node(s) ...\n", len(nodes))

			dir, err := fs.TempDir("", "image-tar")
			if err != nil {
				return errors.Wrap(err, "failed to create tempdir")
			}
			defer func() {
				if err := utils.NewStep("Cleaning up", ":broom:", opts.Verbose, func() (err error) {
					return os.RemoveAll(dir)
				}); err != nil {
					reterr = err
				}
			}()

			imageTarPath := filepath.Join(dir, "image.tar")
			if err := utils.NewStep("Preparing host image", ":optical_disk:", opts.Verbose, func() (err error) {
				return exec.Command("docker", "save", "-o", imageTarPath, imageName).Run()
			}); err != nil {
				return err
			}
			for _, node := range nodes {
				id, err := nodeutils.ImageID(node.Node, imageName)
				if err == nil || id == imageID {
					emoji.Printf(" :package: Image already loaded on node %q\n", node.String())
					continue
				}
				if err := utils.NewStep(fmt.Sprintf("Loading image on node %q", node.String()), ":package:", opts.Verbose, func() (err error) {
					return node.LoadImage(imageTarPath)
				}); err != nil {
					return err
				}
			}

			return nil
		},
	}

	cmd.Flags().StringVar(&opts.Name, "name", "cinder", "cluster name")
	cmd.Flags().BoolVarP(&opts.Verbose, "verbose", "v", false, "show verbose output")
	return cmd
}
