package up

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"go.uber.org/zap/zapcore"

	"github.com/criticalstack/crit/internal/config"
	"github.com/criticalstack/crit/pkg/cluster"
	configutil "github.com/criticalstack/crit/pkg/config/util"
	"github.com/criticalstack/crit/pkg/log"
)

var opts struct {
	ConfigFile     string
	Timeout        time.Duration
	KubeletTimeout time.Duration
}

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:           "up",
		Short:         "Bootstraps a new node",
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, cancel := context.WithTimeout(context.Background(), opts.Timeout)
			defer cancel()

			cfg, err := configutil.LoadFromFile(opts.ConfigFile)
			if err != nil {
				return err
			}
			rc := &cluster.RuntimeConfig{
				KubeletTimeout: opts.KubeletTimeout,
			}
			if log.Level() == zapcore.DebugLevel {
				rc.Verbose = true
			}
			switch c := cfg.(type) {
			case *config.ControlPlaneConfiguration:
				return cluster.RunControlPlane(ctx, rc, c)
			case *config.WorkerConfiguration:
				return cluster.RunWorkerNode(ctx, rc, c)
			default:
				return errors.Errorf("received invalid configuration type: %T", cfg)
			}
		},
	}

	cmd.Flags().StringVarP(&opts.ConfigFile, "config", "c", "config.yaml", "config file")
	cmd.Flags().DurationVar(&opts.Timeout, "timeout", 20*time.Minute, "")
	cmd.Flags().DurationVar(&opts.KubeletTimeout, "kubelet-timeout", 15*time.Second, "timeout for Kubelet to become healthy")
	return cmd
}
