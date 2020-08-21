package api

import (
	"context"

	"github.com/criticalstack/crit/internal/cinder/cluster"
	"github.com/criticalstack/crit/internal/cinder/config"
	"github.com/criticalstack/crit/internal/cinder/config/constants"
)

type (
	Config = cluster.Config
	Node   = cluster.Node

	ControlPlaneConfig = cluster.ControlPlaneConfig
	WorkerConfig       = cluster.WorkerConfig

	ClusterConfiguration = config.ClusterConfiguration
	File                 = config.File
	Encoding             = config.Encoding
)

var (
	Base64           = config.Base64
	Gzip             = config.Gzip
	GzipBase64       = config.GzipBase64
	HostPath         = config.HostPath
	Version          = constants.KubernetesVersion
	DefaultNodeImage = constants.DefaultNodeImage
)

func CreateWorkerNode(ctx context.Context, cfg *WorkerConfig) (*Node, error) {
	return cluster.CreateWorkerNode(ctx, cfg)
}

func CreateNode(ctx context.Context, cfg *Config) (*Node, error) {
	return cluster.CreateNode(ctx, cfg)
}

func DeleteNodes(n []*Node) error {
	return cluster.DeleteNodes(n)
}

func ListNodes(clusterName string) ([]*Node, error) {
	return cluster.ListNodes(clusterName)
}
