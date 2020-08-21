package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/spf13/cobra"
)

var opts struct {
	TLSCertFile       string
	TLSKeyFile        string
	CAFile            string
	CertFile          string
	KeyFile           string
	Port              int
	APIServerBindPort int
}

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:           "healthcheck-proxy",
		Short:         "crit healthcheck proxy sidecar",
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(cmd *cobra.Command, args []string) error {
			caCert, err := ioutil.ReadFile(opts.CAFile)
			if err != nil {
				return err
			}
			cert, err := ioutil.ReadFile(opts.CertFile)
			if err != nil {
				return err
			}
			key, err := ioutil.ReadFile(opts.KeyFile)
			if err != nil {
				return err
			}
			caPool := x509.NewCertPool()
			caPool.AppendCertsFromPEM(caCert)
			tlsCert, err := tls.X509KeyPair(cert, key)
			if err != nil {
				return err
			}
			u, err := url.Parse(fmt.Sprintf("https://localhost:%d", opts.APIServerBindPort))
			if err != nil {
				return err
			}
			e := echo.New()
			e.Use(middleware.Logger())
			e.Use(middleware.Recover())
			e.Use(middleware.ProxyWithConfig(middleware.ProxyConfig{
				Balancer: middleware.NewRoundRobinBalancer([]*middleware.ProxyTarget{{URL: u}}),
				Transport: &http.Transport{
					TLSClientConfig: &tls.Config{
						RootCAs:      caPool,
						Certificates: []tls.Certificate{tlsCert},
					},
				},
			}))
			e.Use(middleware.BodyLimit("2M"))
			e.Use(middleware.Secure())
			s := &http.Server{
				Addr:           fmt.Sprintf("0.0.0.0:%d", opts.Port),
				Handler:        e,
				ReadTimeout:    10 * time.Second,
				WriteTimeout:   10 * time.Second,
				MaxHeaderBytes: 1 << 20,
			}
			l, err := net.Listen("tcp", s.Addr)
			if err != nil {
				return err
			}
			return s.ServeTLS(l, opts.TLSCertFile, opts.TLSKeyFile)
		},
	}
	cmd.Flags().StringVar(&opts.CAFile, "client-ca-file", "/etc/kubernetes/pki/ca.crt", "")
	cmd.Flags().StringVar(&opts.CertFile, "healthcheck-client-certificate", "/etc/kubernetes/pki/apiserver-healthcheck-client.crt", "")
	cmd.Flags().StringVar(&opts.KeyFile, "healthcheck-client-key", "/etc/kubernetes/pki/apiserver-healthcheck-client.key", "")
	cmd.Flags().StringVar(&opts.TLSCertFile, "tls-cert-file", "/etc/kubernetes/pki/apiserver.crt", "")
	cmd.Flags().StringVar(&opts.TLSKeyFile, "tls-private-key-file", "/etc/kubernetes/pki/apiserver.key", "")
	cmd.Flags().IntVar(&opts.Port, "secure-port", 6444, "")
	cmd.Flags().IntVar(&opts.APIServerBindPort, "apiserver-port", 6443, "")
	return cmd
}

func main() {
	if err := NewCommand().Execute(); err != nil {
		log.Fatal(err)
	}
}
