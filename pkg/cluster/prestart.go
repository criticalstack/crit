package cluster

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"go.uber.org/zap"
	"k8s.io/apimachinery/pkg/util/version"

	"github.com/criticalstack/crit/internal/config"
	"github.com/criticalstack/crit/pkg/log"
	executil "github.com/criticalstack/crit/pkg/util/exec"
	fmtutil "github.com/criticalstack/crit/pkg/util/fmt"
	netutil "github.com/criticalstack/crit/pkg/util/net"
)

const (
	// DefaultClusterDNSIP defines default DNS IP
	DefaultClusterDNSIP            = "10.254.0.10"
	DefaultBootstrapServerBindPort = 8080
	DefaultKubeAPIServerPort       = 6443
)

func setNodeRuntimeDefaults(cfg *config.NodeConfiguration) {
	if cfg.HostIPv4 == "" {
		cfg.HostIPv4, _ = netutil.DetectHostIPv4()
	}
	if cfg.Hostname == "" {
		cfg.Hostname, _ = os.Hostname()
	}
}

var (
	// MinKubeVersion indicates the lowest version that crit will provision
	MinKubeVersion = version.MustParseSemantic("v1.14.0")

	// MaxKubeVersion indicates the highest version that crit will provision.
	// However, versions higher than this will only produce a warning.
	MaxKubeVersion = version.MustParseSemantic("v1.19.0")
)

func validateNodeConfiguration(cfg *config.NodeConfiguration) (errs []error) {
	v, err := version.ParseSemantic(cfg.KubernetesVersion)
	if err != nil {
		errs = append(errs, errors.Errorf("invalid KubernetesVersion: %#v", cfg.KubernetesVersion))
	}
	if v != nil && v.Minor() < MinKubeVersion.Minor() {
		errs = append(errs, errors.Errorf("invalid KubernetesVersion: %#v", cfg.KubernetesVersion))
	}
	if v != nil && v.Minor() > MaxKubeVersion.Minor() {
		log.Warn("The KubernetesVersion is newer than expected. Older versions of crit will be capable of bootstrapping newer versions of Kubernetes, but may produce undesired behavior", zap.String("KubernetesVersion", cfg.KubernetesVersion))
	}
	return
}

func (c *Cluster) ControlPlanePreCheck(ctx context.Context, cfg *config.ControlPlaneConfiguration) error {
	log.Info("precheck-control-plane", zap.String("description", "perform host system configuration checks"))
	setControlPlaneRuntimeDefaults(cfg)
	errs := validateControlPlaneConfiguration(cfg)
	if len(errs) > 0 {
		stderr := executil.NewPrefixWriter(os.Stderr, "\t")
		defer stderr.Close()

		stderr.Write([]byte(fmtutil.FormatErrors(errs)))
		return errors.New("failed precheck")
	}
	return nil
}

func setControlPlaneRuntimeDefaults(cfg *config.ControlPlaneConfiguration) {
	// some defaults are derived from NodeConfiguration defaults, so this must
	// run first
	setNodeRuntimeDefaults(&cfg.NodeConfiguration)

	if cfg.ControlPlaneEndpoint.Host == "" {
		log.Warn("ControlPlaneEndpoint is being set implicitly to the host IPv4. It is recommended to use a Load Balancer or DNS for this value to ensure that cluster services, like kube-proxy, will always be able to connect to the control plane.")
		cfg.ControlPlaneEndpoint.Host = cfg.NodeConfiguration.HostIPv4
	}
	if cfg.ControlPlaneEndpoint.Port == 0 {
		cfg.ControlPlaneEndpoint.Port = int32(cfg.KubeAPIServerConfiguration.BindPort)
	}
	if len(cfg.EtcdConfiguration.Endpoints) == 0 {
		cfg.EtcdConfiguration.Endpoints = append(cfg.EtcdConfiguration.Endpoints, fmt.Sprintf("https://%s:2379", cfg.ControlPlaneEndpoint.Host))
	}
	if strings.HasPrefix(cfg.EtcdConfiguration.Endpoints[0], "https://") {
		if cfg.EtcdConfiguration.CAFile == "" {
			cfg.EtcdConfiguration.CAFile = filepath.Join(cfg.NodeConfiguration.KubeDir, "pki/etcd/ca.crt")
		}
		if cfg.EtcdConfiguration.CertFile == "" {
			cfg.EtcdConfiguration.CertFile = filepath.Join(cfg.NodeConfiguration.KubeDir, "pki/etcd/client.crt")
		}
		if cfg.EtcdConfiguration.KeyFile == "" {
			cfg.EtcdConfiguration.KeyFile = filepath.Join(cfg.NodeConfiguration.KubeDir, "pki/etcd/client.key")
		}
		if cfg.EtcdConfiguration.CAKey == "" {
			cfg.EtcdConfiguration.CAKey = filepath.Join(cfg.NodeConfiguration.KubeDir, "pki/etcd/ca.key")
		}
	}
	if cfg.KubeProxyConfiguration.Config.ClusterCIDR == "" {
		cfg.KubeProxyConfiguration.Config.ClusterCIDR = cfg.PodSubnet
	}
	if cfg.NodeConfiguration.KubeletConfiguration.ClusterDNS == nil {
		dnsIP, err := netutil.GetDNSIP(cfg.ServiceSubnet)
		if err != nil {
			log.Warn("cannot parse service subnet", zap.Error(err))
			cfg.NodeConfiguration.KubeletConfiguration.ClusterDNS = []string{DefaultClusterDNSIP}
		} else {
			cfg.NodeConfiguration.KubeletConfiguration.ClusterDNS = []string{dnsIP.String()}
		}
	}
}

func validateControlPlaneConfiguration(cfg *config.ControlPlaneConfiguration) (errs []error) {
	errs = append(errs, validateNodeConfiguration(&cfg.NodeConfiguration)...)

	for _, ep := range cfg.EtcdConfiguration.Endpoints {
		if !strings.HasPrefix(ep, "http") {
			errs = append(errs, errors.Errorf("must specify scheme (http,https) for etcd endpoint: %#v", ep))
			continue
		}
		_, err := url.Parse(ep)
		if err != nil {
			errs = append(errs, errors.Errorf("invalid etcd endpoint url: %#v", ep))
		}
	}
	switch strings.ToLower(string(cfg.KubeProxyConfiguration.Config.Mode)) {
	case "iptables", "ipvs":
	default:
		errs = append(errs, errors.Errorf("invalid KubeProxyConfiguration Mode: %#v", cfg.KubeProxyConfiguration.Config.Mode))
	}
	return
}

func (c *Cluster) WorkerPreCheck(ctx context.Context, cfg *config.WorkerConfiguration) error {
	log.Info("precheck-worker", zap.String("description", "perform host system configuration checks"))
	setWorkerRuntimeDefaults(cfg)
	errs := validateWorkerConfiguration(cfg)
	if len(errs) > 0 {
		stderr := executil.NewPrefixWriter(os.Stderr, "\t")
		defer stderr.Close()

		stderr.Write([]byte(fmtutil.FormatErrors(errs)))
		return errors.New("failed precheck")
	}
	return nil
}

func setWorkerRuntimeDefaults(cfg *config.WorkerConfiguration) {
	// some defaults are derived from NodeConfiguration defaults, so this must
	// run first
	setNodeRuntimeDefaults(&cfg.NodeConfiguration)

	if cfg.BootstrapServerURL == "" && cfg.ControlPlaneEndpoint.Host != "" {
		cfg.BootstrapServerURL = fmt.Sprintf("https://%s:%d", cfg.ControlPlaneEndpoint.Host, DefaultBootstrapServerBindPort)
	}
	if cfg.ControlPlaneEndpoint.Port == 0 {
		cfg.ControlPlaneEndpoint.Port = DefaultKubeAPIServerPort
		log.Warn("ControlPlaneEndpoint not provided with port, defaulting to 6443", zap.Stringer("control_plane_endpoint", cfg.ControlPlaneEndpoint))
	}
	if cfg.CACert == "" {
		cfg.CACert = filepath.Join(cfg.NodeConfiguration.KubeDir, "pki/ca.crt")
	}
}

func validateWorkerConfiguration(cfg *config.WorkerConfiguration) (errs []error) {
	errs = append(errs, validateNodeConfiguration(&cfg.NodeConfiguration)...)

	if cfg.ControlPlaneEndpoint.IsZero() {
		errs = append(errs, errors.New("must provide ControlPlaneEndpoint for WorkerConfiguration"))
	}
	if cfg.BootstrapServerURL == "" && cfg.BootstrapToken == "" {
		errs = append(errs, errors.New("must provide either BootstrapServerURL or BootstrapToken for WorkerConfiguration"))
	}
	return
}
