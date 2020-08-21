package cluster

import (
	"context"
	"crypto/sha512"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/criticalstack/e2d/pkg/e2db"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/criticalstack/crit/internal/config"
	clusterutil "github.com/criticalstack/crit/pkg/cluster/util"
	"github.com/criticalstack/crit/pkg/log"
)

// sharedClusterFiles represents files that must be shared by all nodes in the
// control plane. All other secrets/configuration are derivative of these
// initial files.
var sharedClusterFiles = []string{
	"/etc/kubernetes/pki/ca.crt",
	"/etc/kubernetes/pki/ca.key",
	"/etc/kubernetes/pki/auth-proxy-ca.crt",
	"/etc/kubernetes/pki/auth-proxy-ca.key",
	"/etc/kubernetes/pki/front-proxy-ca.crt",
	"/etc/kubernetes/pki/front-proxy-ca.key",
	"/etc/kubernetes/pki/sa.key",
	"/etc/kubernetes/pki/sa.pub",
}

func (c *Cluster) CreateOrDownloadCerts(ctx context.Context, cfg *config.ControlPlaneConfiguration) error {
	log.Info("cluster-certs", zap.String("description", "download or create cluster certs"))
	t, _ := ctx.Deadline()
	log.Info("waiting for etcd to become available ...",
		zap.String("etcd-address", cfg.EtcdConfiguration.ClientAddr()),
		zap.Duration("timeout", time.Until(t).Round(time.Minute)),
	)
	opts := make([]e2db.TableOption, 0)
	if cfg.EtcdConfiguration.CAKey != "" {
		data, err := ioutil.ReadFile(cfg.EtcdConfiguration.CAKey)
		if err != nil {
			return err
		}
		opts = append(opts, e2db.WithEncryption(sha512.New512_256().Sum(data)))
	} else {
		log.Warn("The etcd CAKey was not specified in the provided configuration. Without it, the shared clusters files cannot be encrypted at rest.")
	}
	db, err := e2db.New(ctx, &e2db.Config{
		ClientAddr: cfg.EtcdConfiguration.ClientAddr(),
		CAFile:     cfg.EtcdConfiguration.CAFile,
		CertFile:   cfg.EtcdConfiguration.CertFile,
		KeyFile:    cfg.EtcdConfiguration.KeyFile,
		Namespace:  "crit",
	})
	if err != nil {
		return err
	}
	defer db.Close()

	log.Info("connected to etcd",
		zap.String("etcd-address", cfg.EtcdConfiguration.ClientAddr()),
	)

	return db.Table(new(ClusterFile), opts...).Tx(func(tx *e2db.Tx) error {
		var files []*ClusterFile
		if err := tx.All(&files); err != nil && errors.Cause(err) != e2db.ErrNoRows {
			return err
		}

		if len(files) > 0 {
			log.Info("existing cluster pki found")
			for _, f := range files {
				if err := f.Write(); err != nil {
					return err
				}
			}
			return nil
		}

		// If this is the first node the certs won't exist, so we must
		// create them. This will only ever happen once for any given
		// cluster.
		log.Info("cluster pki not found in table, generating new pki locally ...")
		fns := []func(string) error{
			clusterutil.WriteClusterCA,
			clusterutil.WriteFrontProxyCA,
			clusterutil.WriteServiceAccountCA,
			clusterutil.WriteAuthProxyCA,
		}
		for _, fn := range fns {
			if err := fn(filepath.Join(cfg.NodeConfiguration.KubeDir, "pki")); err != nil {
				return err
			}
		}
		for _, path := range sharedClusterFiles {
			file, err := newClusterFile(path)
			if err != nil {
				return err
			}
			if err := tx.Insert(file); err != nil {
				return err
			}
			log.Debug("shared cluster file created", zap.String("path", file.Name), zap.Stringer("mode", file.Mode))
		}
		return nil
	})
}

// CreateNodeCerts generates certs specific to the node. This should run after
// the shared cluster certs have been created/downloaded.
func (c *Cluster) CreateNodeCerts(ctx context.Context, cfg *config.ControlPlaneConfiguration) error {
	log.Info("node-certs", zap.String("description", "create node certs"))
	fns := []func(*config.ControlPlaneConfiguration) error{
		clusterutil.WriteAPIServerCertAndKey,
		clusterutil.WriteAPIServerKubeletClientCertAndKey,
		clusterutil.WriteFrontProxyClientCertAndKey,
		clusterutil.WriteAPIServerHealthcheckClientCertAndKey,
	}
	for _, fn := range fns {
		if err := fn(cfg); err != nil {
			return err
		}
	}
	return nil
}

type ClusterFile struct {
	Name string `e2db:"id"`
	Mode os.FileMode
	Data []byte
}

func newClusterFile(path string) (*ClusterFile, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, errors.Wrap(err, path)
	}
	fi, err := os.Stat(path)
	if err != nil {
		return nil, errors.Wrap(err, path)
	}
	return &ClusterFile{path, fi.Mode(), data}, nil
}

func (c *ClusterFile) Write() error {
	if err := os.MkdirAll(filepath.Dir(c.Name), 0600); err != nil {
		return err
	}
	return ioutil.WriteFile(c.Name, c.Data, c.Mode)
}
