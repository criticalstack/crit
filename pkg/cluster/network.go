package cluster

import (
	"bytes"
	"context"
	"io/ioutil"
	"path/filepath"
	"strings"
	"text/template"

	"go.uber.org/zap"
	"sigs.k8s.io/yaml"

	"github.com/criticalstack/crit/internal/config"
	"github.com/criticalstack/crit/pkg/cluster/components"
	"github.com/criticalstack/crit/pkg/kubernetes"
	"github.com/criticalstack/crit/pkg/kubernetes/dynamic"
	"github.com/criticalstack/crit/pkg/log"
)

func (c *Cluster) DeployCoreDNS(ctx context.Context, cfg *config.ControlPlaneConfiguration) error {
	log.Info("CoreDNS", zap.String("description", "deploy CoreDNS"))
	data, err := Execute("coredns.yaml", cfg)
	if err != nil {
		return err
	}
	return dynamic.Apply(ctx, c.Config(), data)
}

func (c *Cluster) DeployKubeProxy(ctx context.Context, cfg *config.ControlPlaneConfiguration) error {
	if cfg.KubeProxyConfiguration.Disabled {
		log.Debug("kube-proxy disabled")
		return nil
	}
	log.Info("kube-proxy", zap.String("description", "deploy kube-proxy"))
	cm, err := components.NewKubeProxyConfigMap(cfg)
	if err != nil {
		return err
	}
	if err := kubernetes.UpdateConfigMap(c.Client(), ctx, cm); err != nil {
		return err
	}
	if err := components.ApplyKubeProxyRBAC(c.Client(), ctx); err != nil {
		return err
	}
	data, err := Execute("kube-proxy.yaml", cfg)
	if err != nil {
		return err
	}
	return dynamic.Apply(ctx, c.Config(), data)
}

// Execute applies a template from embedded files in this package.
func Execute(path string, v interface{}) ([]byte, error) {
	f, err := Files.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	data, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}
	filename := filepath.Base(path)
	name := strings.TrimSuffix(filename, filepath.Ext(filename))
	t, err := template.New(name).Funcs(template.FuncMap{
		"indent": func(spaces int, v string) string {
			pad := strings.Repeat(" ", spaces)
			return pad + strings.Replace(v, "\n", "\n"+pad, -1)
		},
		"toYaml": func(v interface{}) string {
			data, err := yaml.Marshal(v)
			if err != nil {
				log.Debug("cannot parse template", zap.Error(err))
				return ""
			}
			return strings.TrimSuffix(string(data), "\n")
		},
	}).Parse(string(data))
	if err != nil {
		return nil, err
	}
	t = t.Option("missingkey=error")

	var buf bytes.Buffer
	if err := t.Execute(&buf, v); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
