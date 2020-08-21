package app

import (
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"

	proxyproto "github.com/pires/go-proxyproto"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

type serverOptions struct {
	CertFile   string
	KeyFile    string
	Provider   string
	Filters    string
	Kubeconfig string
	Port       int
}

func NewRootCmd() *cobra.Command {
	o := &serverOptions{}

	cmd := &cobra.Command{
		Use:          "bootstrap-server",
		Short:        "run bootstrap-server",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if o.CertFile == "" {
				return errors.New("must provide CertFile")
			}
			if o.KeyFile == "" {
				return errors.New("must provide KeyFile")
			}
			s := &http.Server{
				Addr: fmt.Sprintf(":%d", o.Port),
				Handler: newBootstrapRouter(&bootstrapConfig{
					Provider:   o.Provider,
					Filters:    parseFilters(o.Filters),
					Kubeconfig: o.Kubeconfig,
				}),
				ReadTimeout:    10 * time.Second,
				WriteTimeout:   10 * time.Second,
				MaxHeaderBytes: 1 << 20,
			}
			l, err := net.Listen("tcp", s.Addr)
			if err != nil {
				return err
			}
			return s.ServeTLS(&proxyproto.Listener{Listener: l}, o.CertFile, o.KeyFile)
		},
	}

	cmd.Flags().StringVar(&o.Provider, "provider", "", "")
	cmd.Flags().StringVar(&o.Filters, "filters", "", "")
	cmd.Flags().StringVar(&o.CertFile, "cert-file", "", "server certificate")
	cmd.Flags().StringVar(&o.KeyFile, "key-file", "", "server key")
	cmd.Flags().StringVar(&o.Kubeconfig, "kubeconfig", "", "")
	cmd.Flags().IntVar(&o.Port, "port", 8080, "")

	return cmd
}

func parseFilters(s string) map[string]string {
	filters := make(map[string]string)
	pairs := strings.Split(s, ",")
	for _, pair := range pairs {
		parts := strings.SplitN(pair, "=", 2)
		if len(parts) != 2 {
			continue
		}
		filters[parts[0]] = parts[1]
	}
	return filters
}
