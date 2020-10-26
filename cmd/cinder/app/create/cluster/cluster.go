package cluster

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"path/filepath"
	"time"

	"github.com/mitchellh/go-homedir"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/criticalstack/crit/internal/cinder/cluster"
	"github.com/criticalstack/crit/internal/cinder/config"
	"github.com/criticalstack/crit/internal/cinder/config/constants"
	"github.com/criticalstack/crit/internal/cinder/feature"
	"github.com/criticalstack/crit/internal/cinder/utils"
	critconfig "github.com/criticalstack/crit/internal/config"
	yamlutil "github.com/criticalstack/crit/pkg/kubernetes/yaml"
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
		Use:           "cluster",
		Short:         "Creates a new cinder cluster",
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			if _, err := cluster.ImageID(opts.Image); err != nil {
				fmt.Printf("Image %q not found locally ...\n", opts.Image)
				if err := utils.NewStep("Downloading base image", ":fire:", opts.Verbose, func() (err error) {
					return cluster.PullImage(opts.Image)
				}); err != nil {
					return err
				}
				fmt.Printf("\n")
			}

			nodes, err := cluster.ListNodes(opts.Name)
			if len(nodes) != 0 {
				return errors.Errorf("node(s) already exist for a cluster with the name %q", opts.Name)
			}
			fmt.Printf("Creating cluster %q ...\n", opts.Name)

			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
			defer cancel()

			files := make([]config.File, 0)
			if err := utils.NewStep("Generating certificates", ":fire:", opts.Verbose, func() error {
				cert, key, err := utils.CreateCA("kubernetes")
				if err != nil {
					return err
				}
				files = append(files, config.File{
					Path:        "/etc/kubernetes/pki/ca.crt",
					Owner:       "root:root",
					Permissions: "0644",
					Encoding:    config.Base64,
					Content:     base64.StdEncoding.EncodeToString(cert),
				})
				files = append(files, config.File{
					Path:        "/etc/kubernetes/pki/ca.key",
					Owner:       "root:root",
					Permissions: "0600",
					Encoding:    config.Base64,
					Content:     base64.StdEncoding.EncodeToString(key),
				})
				cert, key, err = utils.CreateCA("etcd")
				if err != nil {
					return err
				}
				files = append(files, config.File{
					Path:        "/etc/kubernetes/pki/etcd/ca.crt",
					Owner:       "root:root",
					Permissions: "0644",
					Encoding:    config.Base64,
					Content:     base64.StdEncoding.EncodeToString(cert),
				})
				files = append(files, config.File{
					Path:        "/etc/kubernetes/pki/etcd/ca.key",
					Owner:       "root:root",
					Permissions: "0600",
					Encoding:    config.Base64,
					Content:     base64.StdEncoding.EncodeToString(key),
				})
				return nil
			}); err != nil {
				return err
			}

			var node *cluster.Node
			defer func() {
				// TODO(chrism): The abstraction here for printing out verbose
				// information should be improved. It should also prompt the
				// user on whether to print out the full information in the
				// event of a failure to create a cluster.
				if err != nil && node != nil && !opts.Verbose {
					fmt.Print(string(node.CombinedOutput()))
				}
			}()

			cfg := &config.ClusterConfiguration{}
			if err := utils.NewStep("Creating control-plane node", ":fire:", opts.Verbose, func() (err error) {
				if opts.Config != "" {
					cfg, err = config.LoadFromFile(opts.Config)
					if err != nil {
						return err
					}
				}
				if err := feature.MutableGates.SetFromMap(cfg.FeatureGates); err != nil {
					return err
				}
				if cfg.ControlPlaneConfiguration == nil {
					cfg.ControlPlaneConfiguration = &critconfig.ControlPlaneConfiguration{}
				}
				cfg.ControlPlaneConfiguration.ClusterName = opts.Name
				cluster.SetControlPlaneConfigurationDefaults(cfg.ControlPlaneConfiguration)
				cfg.Files = append(cfg.Files, files...)
				data, err := yamlutil.MarshalToYaml(cfg.ControlPlaneConfiguration, critconfig.SchemeGroupVersion)
				if err != nil {
					return err
				}
				cfg.Files = append(cfg.Files, config.File{
					Path:        "/var/lib/crit/config.yaml",
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
				node, err = cluster.CreateControlPlaneNode(ctx, &cluster.ControlPlaneConfig{
					ClusterName:          opts.Name,
					ContainerName:        opts.Name,
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

			if err := utils.NewStep("Installing CNI", ":fire:", opts.Verbose, func() (err error) {
				return node.Command("bash", "/cinder/scripts/install-cni.sh").SetEnv(
					"KUBECONFIG=/etc/kubernetes/admin.conf",
					fmt.Sprintf("CONTROL_PLANE_HOST=%s", node.IP()),
				).Run()
			}); err != nil {
				return err
			}

			if err := utils.NewStep("Installing StorageClass", ":fire:", opts.Verbose, func() (err error) {
				return node.Command("bash", "/cinder/scripts/install-storageclass.sh").SetEnv("KUBECONFIG=/etc/kubernetes/admin.conf").Run()
			}); err != nil {
				return err
			}

			if feature.Gates.Enabled(feature.MachineAPI) {
				if err := utils.NewStep("Installing machine-api", ":fire:", opts.Verbose, func() (err error) {
					return node.Command("bash", "/cinder/scripts/install-machine-api.sh").SetEnv("KUBECONFIG=/etc/kubernetes/admin.conf").Run()
				}); err != nil {
					return err
				}
			}
			if feature.Gates.Enabled(feature.LocalRegistry) {
				if err := utils.NewStep("Installing local registry", ":fire:", opts.Verbose, func() (err error) {
					if !cluster.IsContainerRunning(cfg.LocalRegistryName) {
						if err := cluster.CreateRegistry(cfg.LocalRegistryName, cfg.LocalRegistryPort); err != nil {
							return err
						}
					}
					data, err := cluster.GetLocalRegistryHostingConfigMap(cfg)
					if err != nil {
						return err
					}
					return node.Command("kubectl", "apply", "-f", "-").SetStdin(bytes.NewReader(data)).SetEnv("KUBECONFIG=/etc/kubernetes/admin.conf").Run()
				}); err != nil {
					return err
				}
			}

			if feature.Gates.Enabled(feature.Krustlet) {
				if err := utils.NewStep("Installing Krustlet", ":fire:", opts.Verbose, func() (err error) {
					if err := cluster.BootstrapKrustlet(opts.Name, "wasi", 3000, node, cfg.RegistryMirrors); err != nil {
						return err
					}
					if err := cluster.BootstrapKrustlet(opts.Name, "wascc", 3001, node, cfg.RegistryMirrors); err != nil {
						return err
					}
					return nil
				}); err != nil {
					return err
				}
			}

			if err := utils.NewStep("Running post-up commands", ":fire:", opts.Verbose, func() (err error) {
				return node.Command("bash", "/cinder/scripts/post-up.sh").SetEnv("KUBECONFIG=/etc/kubernetes/admin.conf").Run()
			}); err != nil {
				return err
			}

			if opts.Kubeconfig == "" {
				home, err := homedir.Dir()
				if err != nil {
					return err
				}
				opts.Kubeconfig = filepath.Join(home, ".kube/config")
			}
			if err := cluster.ExportKubeConfig(opts.Name, opts.Kubeconfig); err != nil {
				return err
			}
			fmt.Printf("Set kubectl context to \"kubernetes-admin@%s\". Prithee, be careful.\n", opts.Name)
			return nil
		},
	}
	cmd.Flags().StringVar(&opts.Name, "name", "cinder", "cluster name")
	// TODO(chrism): make this lookup envvar for default image first
	cmd.Flags().StringVar(&opts.Image, "image", constants.DefaultNodeImage, "node image")
	cmd.Flags().StringVarP(&opts.Config, "config", "c", "", "cinder configuration file")
	cmd.Flags().StringVar(&opts.Kubeconfig, "kubeconfig", "", "sets kubeconfig path instead of $KUBECONFIG or $HOME/.kube/config")
	cmd.Flags().BoolVarP(&opts.Verbose, "verbose", "v", false, "show verbose output")
	return cmd
}
