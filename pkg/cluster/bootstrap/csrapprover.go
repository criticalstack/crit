package bootstrap

import (
	"context"

	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clientset "k8s.io/client-go/kubernetes"

	"github.com/criticalstack/crit/pkg/kubernetes"
)

var NodeBootstrapTokenRBAC = []*rbacv1.ClusterRoleBinding{
	// Allow to post CSR
	{
		ObjectMeta: metav1.ObjectMeta{
			Name: "crit:kubelet-bootstrap",
		},
		RoleRef: rbacv1.RoleRef{
			APIGroup: rbacv1.GroupName,
			Kind:     "ClusterRole",
			Name:     "system:node-bootstrapper",
		},
		Subjects: []rbacv1.Subject{
			{
				Kind: rbacv1.GroupKind,
				Name: "system:bootstrappers:crit:default-node-token",
			},
		},
	},

	// Allow csrapprover controller to auto-approve CSRs
	{
		ObjectMeta: metav1.ObjectMeta{
			Name: "crit:node-autoapprove-bootstrap",
		},
		RoleRef: rbacv1.RoleRef{
			APIGroup: rbacv1.GroupName,
			Kind:     "ClusterRole",
			Name:     "system:certificates.k8s.io:certificatesigningrequests:nodeclient",
		},
		Subjects: []rbacv1.Subject{
			{
				Kind: "Group",
				Name: "system:bootstrappers:crit:default-node-token",
			},
		},
	},

	// Allow csrapprover controller to auto-approve certificate rotation CSRs
	{
		ObjectMeta: metav1.ObjectMeta{
			Name: "crit:node-autoapprove-certificate-rotation",
		},
		RoleRef: rbacv1.RoleRef{
			APIGroup: rbacv1.GroupName,
			Kind:     "ClusterRole",
			Name:     "system:certificates.k8s.io:certificatesigningrequests:selfnodeclient",
		},
		Subjects: []rbacv1.Subject{
			{
				Kind: "Group",
				Name: "system:nodes",
			},
		},
	},
}

func ApplyCSRApproverRBAC(client *clientset.Clientset, ctx context.Context) error {
	for _, crb := range NodeBootstrapTokenRBAC {
		if err := kubernetes.UpdateClusterRoleBinding(client, ctx, crb); err != nil {
			return err
		}
	}
	return nil
}
