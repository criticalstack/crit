package util

import (
	"io/ioutil"
	"os"
	"path/filepath"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"

	yamlutil "github.com/criticalstack/crit/pkg/kubernetes/yaml"
)

func WriteKubeComponent(obj runtime.Object, path string) error {
	data, err := yamlutil.MarshalToYaml(obj, corev1.SchemeGroupVersion)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return err
	}
	return ioutil.WriteFile(path, data, 0600)
}
