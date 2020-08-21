package kubernetes

import (
	"context"

	v1 "k8s.io/api/rbac/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func UpdateClusterRoleBinding(k *kubernetes.Clientset, ctx context.Context, crb *v1.ClusterRoleBinding) error {
	if _, err := k.RbacV1().ClusterRoleBindings().Create(ctx, crb, metav1.CreateOptions{}); err != nil {
		if !apierrors.IsAlreadyExists(err) {
			return err
		}
		if _, err := k.RbacV1().ClusterRoleBindings().Update(ctx, crb, metav1.UpdateOptions{}); err != nil {
			return err
		}
	}
	return nil
}

func UpdateRoleBinding(k *kubernetes.Clientset, ctx context.Context, rb *v1.RoleBinding) error {
	if _, err := k.RbacV1().RoleBindings(rb.Namespace).Create(ctx, rb, metav1.CreateOptions{}); err != nil {
		if !apierrors.IsAlreadyExists(err) {
			return err
		}
		if _, err := k.RbacV1().RoleBindings(rb.Namespace).Update(ctx, rb, metav1.UpdateOptions{}); err != nil {
			return err
		}
	}
	return nil
}

func UpdateRole(k *kubernetes.Clientset, ctx context.Context, r *v1.Role) error {
	if _, err := k.RbacV1().Roles(r.Namespace).Create(ctx, r, metav1.CreateOptions{}); err != nil {
		if !apierrors.IsAlreadyExists(err) {
			return err
		}
		if _, err := k.RbacV1().Roles(r.Namespace).Update(ctx, r, metav1.UpdateOptions{}); err != nil {
			return err
		}
	}
	return nil
}
