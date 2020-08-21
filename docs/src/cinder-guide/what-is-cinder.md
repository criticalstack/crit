# What is Cinder

Cinder, or Crit-in-Docker, is very similar to [kind](https://kind.sigs.k8s.io/). In fact, it uses many packages from kind under-the-hood along with the base container image that makes it all work. Think of cinder as like a flavor of kind (kind is quite good, to say the least). Just like kind, cinder won't work on all platforms, and right now only supports amd64 architectures running macOS and linux, and requires running Docker.

Cinder bootstraps each node with Crit and installs several helpful additional components, such as the machine-api and machine-api-provider-docker.

## Using Cinder to Develop Crit


```yaml
# dev.yaml
apiVersion: cinder.crit.sh/v1alpha1
kind: ClusterConfiguration
files:
  - path: "/usr/bin/crit"
    owner: "root:root"
    permissions: "0755"
    encoding: hostpath
    content: bin/crit
```

```sh
❯ make crit
❯ cinder create cluster -c dev.yaml
Creating cluster "cinder" ...
 🔥  Generating certificates
 🔥  Creating control-plane node
 🔥  Installing CNI
 🔥  Installing StorageClass
 🔥  Running post-up commands
Set kubectl context to "kubernetes-admin@cinder". Prithee, be careful.
```
