package config

import (
	externalconfig "github.com/criticalstack/crit/internal/cinder/config/v1alpha1"
)

// The external configuration types being used are aliased to be used internal
// to the project without needing to update import paths.
type (
	ClusterConfiguration = externalconfig.ClusterConfiguration
	Mount                = externalconfig.Mount
	File                 = externalconfig.File
	Encoding             = externalconfig.Encoding
	PortMapping          = externalconfig.PortMapping
	PortMappingProtocol  = externalconfig.PortMappingProtocol
)

var (
	SchemeGroupVersion = externalconfig.SchemeGroupVersion

	Base64     = externalconfig.Base64
	Gzip       = externalconfig.Gzip
	GzipBase64 = externalconfig.GzipBase64
	HostPath   = externalconfig.HostPath
)
