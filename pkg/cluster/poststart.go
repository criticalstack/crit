package cluster

import (
	"context"
	"encoding/base64"
	"io/ioutil"
	"path/filepath"

	"go.uber.org/zap"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/criticalstack/crit/internal/config"
	"github.com/criticalstack/crit/pkg/cluster/bootstrap"
	"github.com/criticalstack/crit/pkg/kubernetes"
	nodeutil "github.com/criticalstack/crit/pkg/kubernetes/util/node"
	yamlutil "github.com/criticalstack/crit/pkg/kubernetes/yaml"
	"github.com/criticalstack/crit/pkg/log"
)

func (c *Cluster) EnableCSRApprover(ctx context.Context, cfg *config.ControlPlaneConfiguration) error {
	log.Info("enable-csrapprover", zap.String("description", "adds RBAC that allows csrapprover to bootstrap nodes"))
	return bootstrap.ApplyCSRApproverRBAC(c.Client(), ctx)
}

func (c *Cluster) MarkControlPlane(ctx context.Context, cfg *config.ControlPlaneConfiguration) error {
	log.Info("mark-control-plane", zap.String("description", "add taint to control plane node"))
	return nodeutil.PatchNodeWithContext(ctx, c.Client(), cfg.NodeConfiguration.Hostname, func(n *corev1.Node) {
		nodeutil.AddTaint(n, corev1.Taint{
			Key:    "node-role.kubernetes.io/master",
			Effect: corev1.TaintEffectNoSchedule,
		})
		n.ObjectMeta.Labels["node-role.kubernetes.io/master"] = ""
	})
}

var (
	CritConfigName = "crit-config"

	CritConfigRole = &rbacv1.Role{
		ObjectMeta: metav1.ObjectMeta{
			Name:      CritConfigName,
			Namespace: metav1.NamespaceSystem,
		},
		Rules: []rbacv1.PolicyRule{
			{
				Verbs:         []string{"get"},
				APIGroups:     []string{""},
				Resources:     []string{"configmaps"},
				ResourceNames: []string{CritConfigName},
			},
		},
	}

	CritConfigRoleBinding = &rbacv1.RoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name:      CritConfigName,
			Namespace: metav1.NamespaceSystem,
		},
		RoleRef: rbacv1.RoleRef{
			APIGroup: rbacv1.GroupName,
			Kind:     "Role",
			Name:     CritConfigName,
		},
		Subjects: []rbacv1.Subject{
			{
				Kind: rbacv1.GroupKind,
				Name: "system:nodes",
			},
			{
				Kind: rbacv1.GroupKind,
				Name: "system:bootstrappers:crit:default-node-token",
			},
		},
	}
)

func (c *Cluster) UploadInfo(ctx context.Context, cfg *config.ControlPlaneConfiguration) error {
	log.Info("upload-info", zap.String("description", "upload crit cluster info to ConfigMap"))
	caCertData, err := ioutil.ReadFile(filepath.Join(cfg.NodeConfiguration.KubeDir, "pki/ca.crt"))
	if err != nil {
		return err
	}
	data, err := yamlutil.MarshalToYaml(cfg, config.SchemeGroupVersion)
	if err != nil {
		return err
	}
	client := c.Client()
	if err := kubernetes.UpdateConfigMap(client, ctx, &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      CritConfigName,
			Namespace: metav1.NamespaceSystem,
		},
		Data: map[string]string{
			"ca":     base64.StdEncoding.EncodeToString(caCertData),
			"config": string(data),
		},
	}); err != nil {
		return err
	}
	if err := kubernetes.UpdateRole(client, ctx, CritConfigRole); err != nil {
		return err
	}
	return kubernetes.UpdateRoleBinding(client, ctx, CritConfigRoleBinding)
}

func (c *Cluster) UploadAuthProxyCA(ctx context.Context, cfg *config.ControlPlaneConfiguration) error {
	log.Info("upload-auth-proxy-ca", zap.String("description", "upload self-signed auth-proxy ca"))
	if _, err := c.Client().CoreV1().Namespaces().Create(ctx, &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: "cert-manager",
		},
	}, metav1.CreateOptions{}); err != nil {
		if !apierrors.IsAlreadyExists(err) {
			return err
		}
	}
	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "auth-proxy-ca",
			Namespace: "cert-manager",
		},
		StringData: map[string]string{
			"tls.crt": fromFile("/etc/kubernetes/pki/auth-proxy-ca.crt"),
			"tls.key": fromFile("/etc/kubernetes/pki/auth-proxy-ca.key"),
		},
	}
	if _, err := c.Client().CoreV1().Secrets(secret.Namespace).Create(ctx, secret, metav1.CreateOptions{}); err != nil {
		if !apierrors.IsAlreadyExists(err) {
			return err
		}
	}
	return nil
}

func (c *Cluster) UploadETCDSecrets(ctx context.Context, cfg *config.ControlPlaneConfiguration) error {
	log.Info("upload-etcd-client-secrets", zap.String("description", ""))
	if _, err := c.Client().CoreV1().Namespaces().Create(ctx, &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: "critical-stack",
		},
	}, metav1.CreateOptions{}); err != nil {
		if !apierrors.IsAlreadyExists(err) {
			return err
		}
	}
	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "etcd-secrets",
			Namespace: "critical-stack",
		},
		StringData: map[string]string{
			"etcd-client-ca.crt": fromFile("/etc/kubernetes/pki/etcd/ca.crt"),
			"etcd-client.crt":    fromFile("/etc/kubernetes/pki/etcd/client.crt"),
			"etcd-client.key":    fromFile("/etc/kubernetes/pki/etcd/client.key"),
		},
	}
	if _, err := c.Client().CoreV1().Secrets(secret.Namespace).Create(ctx, secret, metav1.CreateOptions{}); err != nil {
		if !apierrors.IsAlreadyExists(err) {
			return err
		}
	}
	return nil
}

func fromFile(path string) string {
	data, _ := ioutil.ReadFile(path)
	return string(data)
}
