// Package kubeconfig contains helper functions for working with Kubernetes
// config files.
package kubeconfig

import (
	"crypto/x509"
	"fmt"

	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
	certutil "k8s.io/client-go/util/cert"
)

func New(server, clusterName, userName string, caCert []byte) *clientcmdapi.Config {
	// Use the cluster and the username as the context name
	contextName := fmt.Sprintf("%s@%s", userName, clusterName)

	return &clientcmdapi.Config{
		Clusters: map[string]*clientcmdapi.Cluster{
			clusterName: {
				Server:                   server,
				CertificateAuthorityData: caCert,
			},
		},
		Contexts: map[string]*clientcmdapi.Context{
			contextName: {
				Cluster:  clusterName,
				AuthInfo: userName,
			},
		},
		AuthInfos:      map[string]*clientcmdapi.AuthInfo{},
		CurrentContext: contextName,
	}
}

func NewForClient(serverURL, clusterName, userName string, ca, cert, key []byte) *clientcmdapi.Config {
	config := New(serverURL, clusterName, userName, ca)
	config.AuthInfos[userName] = &clientcmdapi.AuthInfo{
		ClientKeyData:         key,
		ClientCertificateData: cert,
	}
	config.Contexts[fmt.Sprintf("%s@%s", userName, clusterName)].Namespace = metav1.NamespaceSystem
	return config
}

func WriteToFile(config *clientcmdapi.Config, path string) error {
	return clientcmd.WriteToFile(*config, path)
}

func LoadClientCertificateFromConfig(config *clientcmdapi.Config) (*x509.Certificate, error) {
	ctx, ok := config.Contexts[config.CurrentContext]
	if !ok {
		return nil, errors.New("cannot get config context")
	}
	authInfo, ok := config.AuthInfos[ctx.AuthInfo]
	if !ok || authInfo == nil {
		return nil, errors.New("cannot get config authinfo")
	}
	certs, err := certutil.ParseCertsPEM(authInfo.ClientCertificateData)
	if err != nil {
		return nil, err
	}
	return certs[0], nil
}
