package cluster

import (
	"context"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/criticalstack/crit/internal/cinder/cluster"
	"github.com/criticalstack/crit/internal/cinder/config"
	"github.com/criticalstack/crit/internal/cinder/config/constants"
	"github.com/criticalstack/crit/internal/cinder/feature"
	"github.com/criticalstack/crit/internal/cinder/utils"
	critconfig "github.com/criticalstack/crit/internal/config"
	"github.com/criticalstack/crit/pkg/kubernetes/pki"
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
		Args:          cobra.NoArgs,
		Use:           "node",
		Short:         "Creates a new cinder worker",
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			nodes, err := cluster.ListNodes(opts.Name)
			if len(nodes) == 0 {
				return errors.Errorf("cannot find cluster %q", opts.Name)
			}
			fmt.Printf("Modifying cluster %q ...\n", opts.Name)

			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
			defer cancel()

			var node *cluster.Node
			if err := utils.NewStep("Creating worker node", ":fire:", opts.Verbose, func() (err error) {
				cfg := &config.ClusterConfiguration{}
				if opts.Config != "" {
					cfg, err = config.LoadFromFile(opts.Config)
					if err != nil {
						return err
					}
				}
				if err := feature.MutableGates.SetFromMap(cfg.FeatureGates); err != nil {
					return err
				}
				if cfg.WorkerConfiguration == nil {
					cfg.WorkerConfiguration = &critconfig.WorkerConfiguration{}
				}
				cluster.SetWorkerConfigurationDefaults(cfg.WorkerConfiguration)
				cn, err := cluster.GetControlPlaneNode(opts.Name)
				if err != nil {
					return err
				}
				cfg.WorkerConfiguration.ControlPlaneEndpoint.Host = cn.IP()
				cfg.WorkerConfiguration.ControlPlaneEndpoint.Port = 6443
				id, secret := pki.GenerateBootstrapToken()
				cfg.WorkerConfiguration.BootstrapToken = fmt.Sprintf("%s.%s", id, secret)
				if err := cn.Command("crit", "create", "token", cfg.WorkerConfiguration.BootstrapToken).Run(); err != nil {
					return err
				}
				data, err := cn.ReadFile("/etc/kubernetes/pki/ca.crt")
				if err != nil {
					return err
				}
				cfg.Files = append(cfg.Files, config.File{
					Path:        "/etc/kubernetes/pki/ca.crt",
					Owner:       "root:root",
					Permissions: "0644",
					Encoding:    config.Base64,
					Content:     base64.StdEncoding.EncodeToString(data),
				})
				if feature.Gates.Enabled(feature.LocalRegistry) {
					cfg.RegistryMirrors[fmt.Sprintf("%s:%d", cfg.LocalRegistryName, cfg.LocalRegistryPort)] = fmt.Sprintf("http://%s:%d", cfg.LocalRegistryName, cfg.LocalRegistryPort)
				}
				if len(cfg.RegistryMirrors) > 0 {
					patch, err := cluster.GetRegistryMirrors(cfg.RegistryMirrors)
					if err != nil {
						return err
					}
					if err := cluster.AppendOrPatchContainerd(cfg, string(patch)); err != nil {
						return err
					}
				}
				node, err = cluster.CreateWorkerNode(ctx, &cluster.WorkerConfig{
					ClusterName:          opts.Name,
					ContainerName:        opts.Name + "-worker",
					Image:                opts.Image,
					Verbose:              opts.Verbose,
					ClusterConfiguration: cfg,
				})
				if err != nil {
					return err
				}
				return nil
			}); err != nil {
				return err
			}
			fmt.Printf("Worker node %q added.\n", node.IP())
			return nil
		},
	}
	cmd.Flags().StringVar(&opts.Name, "name", "cinder", "cluster name")
	cmd.Flags().StringVar(&opts.Image, "image", constants.DefaultNodeImage, "node image")
	cmd.Flags().StringVarP(&opts.Config, "config", "c", "", "cinder configuration file")
	cmd.Flags().StringVar(&opts.Kubeconfig, "kubeconfig", "", "sets kubeconfig path instead of $KUBECONFIG or $HOME/.kube/config")
	cmd.Flags().BoolVarP(&opts.Verbose, "verbose", "v", false, "show verbose output")
	return cmd
}
