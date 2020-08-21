package token

import (
	"context"
	"strings"
	"time"

	"github.com/spf13/cobra"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clientset "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/criticalstack/crit/pkg/kubernetes"
	"github.com/criticalstack/crit/pkg/kubernetes/pki"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:           "token [token]",
		Short:         "creates a bootstrap token resource",
		Args:          cobra.MinimumNArgs(1),
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			id, secret := pki.GenerateBootstrapToken()
			if len(args) != 0 {
				parts := strings.SplitN(args[0], ".", 2)
				if len(parts) == 2 {
					id = parts[0]
					secret = parts[1]
				}
			}
			s := &corev1.Secret{
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
					"expiration":                     time.Now().UTC().Add(10 * 365 * 24 * time.Hour).Format("2006-01-02T15:04:05Z"),
				},
			}
			config, err := clientcmd.BuildConfigFromFlags("", "/etc/kubernetes/admin.conf")
			if err != nil {
				return err
			}
			client, err := clientset.NewForConfig(config)
			if err != nil {
				return err
			}
			if err := kubernetes.UpdateSecret(client, context.TODO(), s); err != nil {
				return err
			}
			return nil
		},
	}
	return cmd
}
