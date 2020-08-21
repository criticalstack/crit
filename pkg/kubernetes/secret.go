package kubernetes

import (
	"context"

	v1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func UpdateSecret(k *kubernetes.Clientset, ctx context.Context, secret *v1.Secret) error {
	if _, err := k.CoreV1().Secrets(metav1.NamespaceSystem).Create(ctx, secret, metav1.CreateOptions{}); err != nil {
		if !apierrors.IsAlreadyExists(err) {
			return err
		}
		if _, err := k.CoreV1().Secrets(metav1.NamespaceSystem).Update(ctx, secret, metav1.UpdateOptions{}); err != nil {
			return err
		}
	}
	return nil
}
