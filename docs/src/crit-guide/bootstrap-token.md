# Bootstrap Token

Crit supports a worker bootstrap flow using [bootstrap tokens](https://kubernetes.io/docs/reference/access-authn-authz/bootstrap-tokens/) and the cluster CA certificate (e.g. `/etc/kubernetes/pki/ca.crt`):

```yaml
apiVersion: crit.sh/v1alpha2
kind: WorkerConfiguration
bootstrapToken: abcdef.0123456789abcdef
caCert: /etc/kubernetes/pki/ca.crt
controlPlaneEndpoint: mycluster.domain
node:
  cloudProvider: aws
  kubernetesVersion: 1.17.3
```

This method is adapted from the [kubeadm join workflow](https://kubernetes.io/docs/reference/setup-tools/kubeadm/kubeadm-join/#join-workflow), but uses the full CA certificate instead of using CA pinning. It also does not depend upon clients getting a signed configmap, and therefore does not require anonymous auth to be turned on.
