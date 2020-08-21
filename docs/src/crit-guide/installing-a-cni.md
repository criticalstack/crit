# Installing a CNI

The CNI can be installed at any point after a node bootstrapping (i.e. after `crit up` finishes successfully). For example, when we install cilium via helm it looks something like this:

```sh
helm repo add cilium https://helm.cilium.io/
helm install cilium cilium/cilium --namespace kube-system \
    --version 1.8.2
```
