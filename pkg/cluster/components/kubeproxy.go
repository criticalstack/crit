package components

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clientset "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
	kubeproxyconfigv1alpha1 "k8s.io/kube-proxy/config/v1alpha1"

	"github.com/criticalstack/crit/internal/config"
	"github.com/criticalstack/crit/pkg/kubernetes"
	yamlutil "github.com/criticalstack/crit/pkg/kubernetes/yaml"
)

var (
	KubeProxyName = "kube-proxy"

	KubeProxyServiceAccount = &corev1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Name:      KubeProxyName,
			Namespace: metav1.NamespaceSystem,
		},
	}

	KubeProxyClusterRoleBinding = &rbacv1.ClusterRoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name: "crit:node-proxier",
		},
		RoleRef: rbacv1.RoleRef{
			APIGroup: rbacv1.GroupName,
			Kind:     "ClusterRole",
			Name:     "system:node-proxier",
		},
		Subjects: []rbacv1.Subject{
			{
				Kind:      rbacv1.ServiceAccountKind,
				Name:      KubeProxyName,
				Namespace: metav1.NamespaceSystem,
			},
		},
	}

	KubeProxyRole = &rbacv1.Role{
		ObjectMeta: metav1.ObjectMeta{
			Name:      KubeProxyName,
			Namespace: metav1.NamespaceSystem,
		},
		Rules: []rbacv1.PolicyRule{
			{
				Verbs:         []string{"get"},
				APIGroups:     []string{""},
				Resources:     []string{"configmaps"},
				ResourceNames: []string{KubeProxyName},
			},
		},
	}

	KubeProxyRoleBinding = &rbacv1.RoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name:      KubeProxyName,
			Namespace: metav1.NamespaceSystem,
		},
		RoleRef: rbacv1.RoleRef{
			APIGroup: rbacv1.GroupName,
			Kind:     "Role",
			Name:     "kube-proxy",
		},
		Subjects: []rbacv1.Subject{
			{
				Kind: rbacv1.GroupKind,
				Name: "system:bootstrappers:crit:default-node-token",
			},
		},
	}
)

func ApplyKubeProxyRBAC(client *clientset.Clientset, ctx context.Context) error {
	if err := kubernetes.UpdateServiceAccount(client, ctx, KubeProxyServiceAccount); err != nil {
		return err
	}
	if err := kubernetes.UpdateClusterRoleBinding(client, ctx, KubeProxyClusterRoleBinding); err != nil {
		return err
	}
	if err := kubernetes.UpdateRole(client, ctx, KubeProxyRole); err != nil {
		return err
	}
	return kubernetes.UpdateRoleBinding(client, ctx, KubeProxyRoleBinding)
}

func NewKubeProxyConfigMap(cfg *config.ControlPlaneConfiguration) (*corev1.ConfigMap, error) {
	data, err := yamlutil.MarshalToYaml(cfg.KubeProxyConfiguration.Config, kubeproxyconfigv1alpha1.SchemeGroupVersion)
	if err != nil {
		return nil, err
	}
	kubeconfigData, err := clientcmd.Write(clientcmdapi.Config{
		Clusters: map[string]*clientcmdapi.Cluster{
			"default": {
				Server:               fmt.Sprintf("https://%s", cfg.ControlPlaneEndpoint),
				CertificateAuthority: "/var/run/secrets/kubernetes.io/serviceaccount/ca.crt",
			},
		},
		Contexts: map[string]*clientcmdapi.Context{
			"default": {
				Cluster:   "default",
				AuthInfo:  "default",
				Namespace: "default",
			},
		},
		CurrentContext: "default",
		AuthInfos: map[string]*clientcmdapi.AuthInfo{
			"default": {
				TokenFile: "/var/run/secrets/kubernetes.io/serviceaccount/token",
			},
		},
	})
	if err != nil {
		return nil, err
	}
	cm := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      KubeProxyName,
			Namespace: metav1.NamespaceSystem,
		},
		Data: map[string]string{
			"config.conf":     string(data),
			"kubeconfig.conf": string(kubeconfigData),
		},
	}
	return cm, nil
}
