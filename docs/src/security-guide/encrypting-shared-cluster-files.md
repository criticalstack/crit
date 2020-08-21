# Encrypting Shared Cluster Files

The pki shared by all control plane nodes are distributed via etcd/e2d using [e2db](https://github.com/criticalstack/e2d/tree/master/pkg/e2db), an ORM-like abstraction over etcd. These files should be protected using strong encryption, and e2db provides a feature for [encrypting entire tables](https://github.com/criticalstack/e2d/tree/master/pkg/e2db#table-encryption). The one requirement is that the etcd ca key is provided in the crit configuration:

```yaml
apiVersion: crit.sh/v1alpha2
kind: ControlPlaneConfiguration
etcd:
  endpoints:
  - "https://${controlPlaneEndpoint.Host}:2379"
  caFile: /etc/kubernetes/pki/etcd/ca.crt
  caKey: /etc/kubernetes/pki/etcd/ca.key
  certFile: /etc/kubernetes/pki/etcd/client.crt
  keyFile: /etc/kubernetes/pki/etcd/client.key
```

where the important file here is `ca.key`, since it is only one suitable to use as a data encryption key.
