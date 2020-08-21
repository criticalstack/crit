# Encrypting Kubernetes Secrets


## EncryptionProviderConfig

To encrypt secrets within the cluster you must create an `EncryptionConfiguration` manifest and pass it to the API server.

```sh
touch /etc/kubernetes/encryption-config.yaml
chmod 600 /etc/kubernetes/encryption-config.yaml
cat <<-EOT > /etc/kubernetes/encryption-config.yaml
apiVersion: apiserver.config.k8s.io/v1
kind: EncryptionConfiguration
resources:
  - resources:
    - secrets
    providers:
    - aescbc:
        keys:
        - name: key1
          secret: $(cat /etc/kubernetes/pki/etcd/ca.key | md5sum | cut -f 1 -d ' ' | head -c -1 | base64)
    - identity: {}
EOT
```

This `EncryptionConfiguration` uses the `aescbc` provider for encrypting secrets. Details on other providers, including third-party key management systems, can be found in the [Kubernetes official documentation](https://kubernetes.io/docs/tasks/administer-cluster/encrypt-data/#providers).

```yaml
apiVersion: crit.sh/v1alpha2
kind: ControlPlaneConfiguration
kubeAPIServer:
  extraVolumes:
  - name: encryption-config
    hostPath: /etc/kubernetes/encryption-config.yaml
    mountPath: /etc/kubernetes/encryption-config.yaml
    readOnly: true
```

Once the API server is available, [verify that new secrets are encrypted](https://kubernetes.io/docs/tasks/administer-cluster/encrypt-data/#verifying-that-data-is-encrypted).
