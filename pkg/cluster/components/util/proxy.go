package util

import (
	"os"
	"strings"

	corev1 "k8s.io/api/core/v1"
)

// GetProxyEnvVars builds a list of environment variables to use in the control
// plane containers in order to use the right proxy
func GetProxyEnvVars() []corev1.EnvVar {
	envs := []corev1.EnvVar{}
	for _, env := range os.Environ() {
		pos := strings.Index(env, "=")
		if pos == -1 {
			// malformed environment variable, skip it.
			continue
		}
		name := env[:pos]
		value := env[pos+1:]
		if strings.HasSuffix(strings.ToLower(name), "_proxy") && value != "" {
			envVar := corev1.EnvVar{Name: name, Value: value}
			envs = append(envs, envVar)
		}
	}
	return envs
}
