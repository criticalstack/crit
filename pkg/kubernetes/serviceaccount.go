package kubernetes

import (
	"context"

	v1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func UpdateServiceAccount(k *kubernetes.Clientset, ctx context.Context, sa *v1.ServiceAccount) error {
	if _, err := k.CoreV1().ServiceAccounts(metav1.NamespaceSystem).Create(ctx, sa, metav1.CreateOptions{}); err != nil {
		if !apierrors.IsAlreadyExists(err) {
			return err
		}
		if _, err := k.CoreV1().ServiceAccounts(metav1.NamespaceSystem).Update(ctx, sa, metav1.UpdateOptions{}); err != nil {
			return err
		}
	}
	return nil
}
