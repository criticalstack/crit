package cluster

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"

	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"

	"github.com/criticalstack/crit/pkg/cluster/bootstrap"
	"github.com/criticalstack/crit/pkg/kubeconfig"
	"github.com/criticalstack/crit/pkg/kubernetes/pki"
)

const krustletConfTmpl = `KUBECONFIG=/var/lib/krustlet/krustlet.conf
KRUSTLET_ADDRESS={{ .IPv4 }}
KRUSTLET_BOOTSTRAP_FILE={{ .BootstrapConfig }}
KRUSTLET_CERT_FILE=/var/lib/kubelet/pki/kubelet.crt
KRUSTLET_PRIVATE_KEY_FILE=/var/lib/kubelet/pki/kubelet.key
KRUSTLET_DATA_DIR=/var/lib/krustlet
KRUSTLET_NODE_IP={{ .IPv4 }}
KRUSTLET_PORT={{ .Port }}
KRUSTLET_NODE_NAME={{ .Name }}
KRUSTLET_INSECURE_REGISTRIES={{ .InsecureRegistries }}
KRUSTLET_ALLOW_LOCAL_MODULES=true`

func BootstrapKrustlet(cluster, name string, port int, node *Node, registries map[string]string) error {
	id, secret := pki.GenerateBootstrapToken()
	token := fmt.Sprintf("%s.%s", id, secret)
	if err := node.Command("crit", "create", "token", token).Run(); err != nil {
		return err
	}
	caCertData, err := node.ReadFile("/etc/kubernetes/pki/ca.crt")
	if err != nil {
		return err
	}
	config := kubeconfig.New(
		fmt.Sprintf("https://%s:6443", node.IP()),
		cluster,
		bootstrap.TokenUser,
		caCertData,
	)
	config.AuthInfos[bootstrap.TokenUser] = &clientcmdapi.AuthInfo{
		Token: token,
	}
	data, err := clientcmd.Write(*config)
	if err != nil {
		return err
	}
	conf := fmt.Sprintf("/var/lib/krustlet/bootstrap-krustlet-%s.conf", name)
	if err := node.WriteFile(conf, data, 0644); err != nil {
		return err
	}

	t, err := template.New("").Parse(krustletConfTmpl)
	if err != nil {
		return err
	}
	var b bytes.Buffer
	if err := t.Execute(&b, map[string]string{
		"IPv4":               node.IP(),
		"Name":               fmt.Sprintf("%s-%s", cluster, name),
		"Port":               fmt.Sprintf("%d", port),
		"BootstrapConfig":    conf,
		"InsecureRegistries": parseInsecureRegistries(registries),
	}); err != nil {
		return err
	}
	if err := node.WriteFile(fmt.Sprintf("/var/lib/krustlet/config-%s.env", name), b.Bytes(), 0644); err != nil {
		return err
	}
	if err := node.Command("systemctl", "enable", "--now", fmt.Sprintf("krustlet-%s", name)).Run(); err != nil {
		return err
	}
	return nil
}

func parseInsecureRegistries(registries map[string]string) string {
	insecureRegistries := make([]string, 0)
	for ref, registry := range registries {
		if strings.HasPrefix(registry, "http://") {
			insecureRegistries = append(insecureRegistries, ref)
		}
	}
	return strings.Join(insecureRegistries, ",")
}
