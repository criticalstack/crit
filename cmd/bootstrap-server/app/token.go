package app

import (
	"context"
	"time"

	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clientset "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/criticalstack/crit/pkg/kubernetes"
	"github.com/criticalstack/crit/pkg/kubernetes/pki"
)

func createNewToken(kubeconfigFile string) (string, error) {
	config, err := clientcmd.LoadFromFile(kubeconfigFile)
	if err != nil {
		return "", errors.Wrapf(err, "failed to load kubeconfig: %#v", kubeconfigFile)
	}
	clientConfig, err := clientcmd.NewDefaultClientConfig(*config, &clientcmd.ConfigOverrides{}).ClientConfig()
	if err != nil {
		return "", errors.Wrap(err, "failed to create API client configuration from kubeconfig")
	}

	client, err := clientset.NewForConfig(clientConfig)
	if err != nil {
		return "", errors.Wrap(err, "failed to create API client")
	}
	id, secret := pki.GenerateBootstrapToken()
	if err := kubernetes.UpdateSecret(client, context.TODO(), &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "bootstrap-token-" + id,
			Namespace: metav1.NamespaceSystem,
		},
		Type: corev1.SecretTypeBootstrapToken,
		StringData: map[string]string{
			"token-id":                       id,
			"token-secret":                   secret,
			"usage-bootstrap-authentication": "true",
			"usage-bootstrap-signing":        "true",
			"auth-extra-groups":              "system:bootstrappers:crit:default-node-token",
			"expiration":                     time.Now().Add(15 * time.Minute).Format("2006-01-02T15:04:05Z"),
		},
	}); err != nil {
		return "", err
	}
	return id + "." + secret, nil
}
