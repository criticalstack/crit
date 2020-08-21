package kubernetes

import (
	"context"

	v1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// TODO(chris): all functions in this package will need to check if the
// kubernetes.Clientset is nil and return an error if it is.
func GetConfigMap(k *kubernetes.Clientset, ctx context.Context, name string) (*v1.ConfigMap, error) {
	return k.CoreV1().ConfigMaps(metav1.NamespaceSystem).Get(ctx, name, metav1.GetOptions{})
}

func UpdateConfigMap(k *kubernetes.Clientset, ctx context.Context, cm *v1.ConfigMap) error {
	if _, err := k.CoreV1().ConfigMaps(cm.Namespace).Create(ctx, cm, metav1.CreateOptions{}); err != nil {
		if !apierrors.IsAlreadyExists(err) {
			return err
		}
		if _, err := k.CoreV1().ConfigMaps(cm.Namespace).Update(ctx, cm, metav1.UpdateOptions{}); err != nil {
			return err
		}
	}
	return nil
}
