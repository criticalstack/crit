package util

import (
	"io/ioutil"
	"os"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"

	"github.com/criticalstack/crit/internal/config"
	yamlutil "github.com/criticalstack/crit/pkg/kubernetes/yaml"
)

func ReadFile(path string) (data []byte, err error) {
	if path == "-" {
		return ioutil.ReadAll(os.Stdin)
	}
	return ioutil.ReadFile(path)
}

func Unmarshal(data []byte) (runtime.Object, error) {
	return yamlutil.UnmarshalFromYaml(data, runtime.NewMultiGroupVersioner(config.SchemeGroupVersion,
		schema.GroupKind{Group: "crit.sh"},
		schema.GroupKind{Group: "crit.criticalstack.com"},
	))
}

func Marshal(obj runtime.Object) ([]byte, error) {
	return yamlutil.MarshalToYaml(obj, config.SchemeGroupVersion)
}

func LoadFromFile(path string) (runtime.Object, error) {
	data, err := ReadFile(path)
	if err != nil {
		return nil, err
	}
	return Unmarshal(data)
}
