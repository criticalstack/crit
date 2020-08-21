# Bootstrap Server

_Experimental_

The bootstrap protocol used by Kubernetes/kubeadm relies on operations that imply manual work to be performed, in particular, the bootstrap token creation and how that is distributed to new worker nodes. Crit introduces a new bootstrap protocol that tries to work better in environments that are completely automated.

A bootstrap-server static pod is created alongside the Kubernetes components that run on each control plane node. This provides a service to new nodes before they have joined the cluster that allows them to be authorized and given a bootstrap token. This also has the benefit of making the bootstrap token expiration very small, limited the window greatly that it can be used.

### Configuration

Here is an example of using Amazon Instance Identity Document w/ signature verification while also limiting the accounts bootstrap tokens will be issued for:

```yaml
apiVersion: crit.sh/v1alpha2
kind: ControlPlaneConfiguration
critBootstrapServer:
  cloudProvider: aws
  extraArgs:
    filters: account-id=${account_id}
```

Override bootstrap-server default port:

```yaml
apiVersion: crit.sh/v1alpha2
kind: ControlPlaneConfiguration
critBootstrapServer:
  extraArgs:
    port: 8080
```

### Authorizers

#### AWS

The AWS authorizer uses [Instance Identity Documents](https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/instance-identity-documents.html) and RSA SHA 256 signature verification to confirm the identity of new nodes requesting bootstrap tokens.
