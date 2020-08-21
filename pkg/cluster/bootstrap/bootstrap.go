// Package bootstrap contains functions for bootstrapping Kubernetes nodes.
package bootstrap

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/pkg/errors"
	"go.uber.org/zap"
	"k8s.io/apimachinery/pkg/util/wait"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"

	"github.com/criticalstack/crit/internal/config"
	"github.com/criticalstack/crit/pkg/cluster/bootstrap/authorizers/ec2metadata"
	"github.com/criticalstack/crit/pkg/kubeconfig"
	"github.com/criticalstack/crit/pkg/log"
)

const (
	// DefaultClusterName defines the default cluster name
	DefaultClusterName = "crit"

	// TokenUser defines token user
	TokenUser = "tls-bootstrap-token-user"
)

func GetBootstrapKubeletKubeconfig(cfg *config.WorkerConfiguration) (*clientcmdapi.Config, error) {
	caCertData, err := ioutil.ReadFile(cfg.CACert)
	if err != nil {
		return nil, err
	}
	if cfg.BootstrapToken != "" {
		config := kubeconfig.New(
			fmt.Sprintf("https://%s", cfg.ControlPlaneEndpoint),
			DefaultClusterName,
			TokenUser,
			caCertData,
		)
		config.AuthInfos[TokenUser] = &clientcmdapi.AuthInfo{
			Token: cfg.BootstrapToken,
		}
		return config, nil
	}
	rootCAs, err := x509.SystemCertPool()
	if err != nil {
		return nil, err
	}
	if ok := rootCAs.AppendCertsFromPEM(caCertData); !ok {
		log.Warnf("cannot append ca cert")
	}
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs: rootCAs,
			},
		},
		Timeout: 2 * time.Second,
	}
	var p struct {
		Provider AuthorizationType `json:"provider"`
	}
	if err := wait.PollImmediateInfinite(5*time.Second, func() (bool, error) {
		resp, err := client.Get(cfg.BootstrapServerURL + "/authorize")
		if err != nil {
			return false, nil
		}
		defer resp.Body.Close()
		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return false, err
		}
		if err := json.Unmarshal(data, &p); err != nil {
			return false, err
		}
		return true, nil
	}); err != nil {
		return nil, err
	}
	var data []byte
	switch p.Provider {
	case AmazonIdentityDocumentAndSignature:
		body, err := ec2metadata.GetSignedDocument()
		if err != nil {
			return nil, err
		}
		data, err = json.Marshal(&Request{
			Type: AmazonIdentityDocumentAndSignature,
			Body: body,
		})
		if err != nil {
			return nil, err
		}
	default:
		return nil, errors.Errorf("received unknown provider: %v", p.Provider)
	}

	var bootstrapToken string
	if err := wait.PollImmediate(5*time.Second, 5*time.Minute, func() (bool, error) {
		resp, err := client.Post(cfg.BootstrapServerURL+"/authorize", "application/json", bytes.NewReader(data))
		if err != nil {
			log.Warn("cannot authorize", zap.Error(err))
			return false, nil
		}
		defer resp.Body.Close()
		data, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			return false, err
		}
		var r Response
		if err := json.Unmarshal(data, &r); err != nil {
			return false, err
		}
		if r.Error != "" {
			return false, errors.New(r.Error)
		}
		if r.BootstrapToken == "" {
			return false, errors.New("BootstrapToken cannot be empty")
		}
		bootstrapToken = r.BootstrapToken
		return true, nil
	}); err != nil {
		return nil, err
	}
	config := kubeconfig.New(
		fmt.Sprintf("https://%s", cfg.ControlPlaneEndpoint),
		DefaultClusterName,
		TokenUser,
		caCertData,
	)
	config.AuthInfos[TokenUser] = &clientcmdapi.AuthInfo{
		Token: bootstrapToken,
	}
	return config, nil
}
