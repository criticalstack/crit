package kubeconfig

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/pkg/errors"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
	"k8s.io/client-go/util/homedir"

	"github.com/criticalstack/crit/pkg/util/lockfile"
)

// MergeConfig merges the provided Config into the current users local
// kubeconfig.
func MergeConfig(cfg *clientcmdapi.Config) error {
	return MergeConfigToFile(cfg, fmt.Sprintf("%s/.kube/config", homedir.HomeDir()))
}

func MergeConfigToFile(cfg *clientcmdapi.Config, path string) error {
	if err := validateConfig(cfg); err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	mu := lockfile.New(path + ".lock")
	if err := mu.Lock(ctx); err != nil {
		return err
	}
	defer mu.Unlock()

	existing, err := clientcmd.LoadFromFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return clientcmd.WriteToFile(*cfg, path)
		}
		return err
	}
	if err := merge(existing, cfg); err != nil {
		return err
	}
	return clientcmd.WriteToFile(*existing, path)
}

func merge(a, b *clientcmdapi.Config) error {
	for k, v := range b.Clusters {
		a.Clusters[k] = v
	}
	for k, v := range b.AuthInfos {
		a.AuthInfos[k] = v
	}
	for k, v := range b.Contexts {
		a.Contexts[k] = v
	}
	a.CurrentContext = b.CurrentContext
	return nil
}

func RemoveConfig(clusterName string) error {
	return RemoveConfigFromFile(clusterName, fmt.Sprintf("%s/.kube/config", homedir.HomeDir()))
}

func RemoveConfigFromFile(clusterName, path string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	mu := lockfile.New(path + ".lock")
	if err := mu.Lock(ctx); err != nil {
		return err
	}
	defer mu.Unlock()

	existing, err := clientcmd.LoadFromFile(path)
	if err != nil {
		return err
	}
	if updated := remove(existing, clusterName); !updated {
		return nil
	}
	return clientcmd.WriteToFile(*existing, path)
}

func remove(cfg *clientcmdapi.Config, name string) (updated bool) {
	if _, ok := cfg.Clusters[name]; ok {
		delete(cfg.Clusters, name)
		updated = true
	}
	if _, ok := cfg.AuthInfos[name]; ok {
		delete(cfg.AuthInfos, name)
		updated = true
	}
	if _, ok := cfg.Contexts[name]; ok {
		delete(cfg.Contexts, name)
		updated = true
	}
	if cfg.CurrentContext == name {
		cfg.CurrentContext = ""
		updated = true
	}
	return
}

func validateConfig(cfg *clientcmdapi.Config) error {
	switch {
	case len(cfg.Clusters) != 1:
		return errors.Errorf("received Config with %d clusters, expect only 1", len(cfg.Clusters))
	case len(cfg.AuthInfos) != 1:
		return errors.Errorf("received Config with %d authinfos, expect only 1", len(cfg.AuthInfos))
	case len(cfg.Contexts) != 1:
		return errors.Errorf("received Config with %d contexts, expect only 1", len(cfg.Contexts))
	default:
		return nil
	}
}
