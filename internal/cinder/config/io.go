package config

import (
	"encoding/base64"
	"io/ioutil"

	"github.com/pkg/errors"

	"github.com/criticalstack/crit/internal/cinder/utils"
	configutil "github.com/criticalstack/crit/pkg/config/util"
	yamlutil "github.com/criticalstack/crit/pkg/kubernetes/yaml"
)

func LoadFromFile(path string) (*ClusterConfiguration, error) {
	data, err := configutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	v, err := yamlutil.UnmarshalFromYaml(data, SchemeGroupVersion)
	if err != nil {
		return nil, err
	}
	cfg, ok := v.(*ClusterConfiguration)
	if !ok {
		return nil, errors.Errorf("expected %T, received %T", cfg, v)
	}
	return cfg, nil
}

func ReadFile(f File) ([]byte, error) {
	switch f.Encoding {
	case Base64:
		return base64.StdEncoding.DecodeString(f.Content)
	case Gzip:
		return utils.Gunzip([]byte(f.Content))
	case GzipBase64:
		data, err := utils.Gunzip([]byte(f.Content))
		if err != nil {
			return nil, err
		}
		return base64.StdEncoding.DecodeString(string(data))
	case HostPath:
		return ioutil.ReadFile(f.Content)
	default:
		return []byte(f.Content), nil
	}
}
