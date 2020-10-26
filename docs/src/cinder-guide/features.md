# Features

## Side-loading Images

Kind allows you to [side-load images](https://kind.sigs.k8s.io/docs/user/quick-start/#loading-an-image-into-your-cluster) in your local clusters. Cinder exposes the same functionality via [`cinder load`](../cinder-commands/cinder-load.md):

```sh
cinder load criticalstack/quake-kube:v1.0.5
```

This will make the `criticalstack/quake-kube:v1.0.5` image from the host available in the Cinder node. Any image that is available on the host can be loaded, and Cinder lazily pulls images that are not found on the host.

## Registry Mirrors

Mirrors for container image registries can be setup to effectively "alias" them. The key is the alias, and the value is the full endpoint for the registry:

```yaml
apiVersion: cinder.crit.sh/v1alpha1
kind: ClusterConfiguration
registryMirrors:
  docker.io: "https://docker.io"
```

It can be used to alias registries with different names OR it can be used to specify plain http registries:

```yaml
...
registryMirrors:
  myregistry.dev: "http://myregistry.dev"
```

## Local Registry

An instance of [Distribution](https://github.com/docker/distribution) (aka Docker Registry v2) can be setup for a Cinder cluster by specifying a config file with the `LocalRegistry` feature gate:

```yaml
apiVersion: cinder.crit.sh/v1alpha1
kind: ClusterConfiguration
featureGates:
  LocalRegistry: true
```

This will start a Docker container on the host with the running registry (if not already running). The registry is shared for all Cinder clusters on a host and is available at `localhost:5000` (i.e. this is what you `docker push` to). This registry is then available inside the cluster at `cinderegg:5000`.

Cinder also creates the `local-registry-hosting` ConfigMap so that any tooling that supports Local Registry Hosting, such as [Tilt](https://tilt.dev/), will be able to automatically discover and use the local registry.

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: local-registry-hosting
  namespace: kube-public
data:
  localRegistryHosting.v1: |
    host: "localhost:{{ .LocalRegistryPort }}"
    hostFromContainerRuntime: "{{ .LocalRegistryName }}:{{ .LocalRegistryPort }}"
    hostFromClusterNetwork: "{{ .LocalRegistryName }}:{{ .LocalRegistryPort }}"
    help: "https://docs.crit.sh/cinder-guide/local-registry.html"`
```

More information about this Kubernetes standard can be found [here](https://github.com/kubernetes/enhancements/blob/c5b6b632811c21ababa9e3565766b2d70614feec/keps/sig-cluster-lifecycle/generic/1755-communicating-a-local-registry/README.md#design-details).

## Krustlet

[Krustlet](https://github.com/deislabs/krustlet) is a tool to run WebAssembly workloads natively on Kubernetes by acting like node in your Kubernetes cluster. It can be enabled for a Cinder cluster using the following configuration:


```yaml
apiVersion: cinder.crit.sh/v1alpha1
kind: ClusterConfiguration
featureGates:
  LocalRegistry: true
  Krustlet: true
controlPlaneConfiguration:
  kubeProxy:
    affinity:
      nodeAffinity:
        requiredDuringSchedulingIgnoredDuringExecution:
          nodeSelectorTerms:
            - matchExpressions:
              - key: "kubernetes.io/arch"
                operator: NotIn
                values: ["wasm32-wasi", "wasm32-wascc"]
```

*Note that node affinity is being set for `kube-proxy` to ensure it does not try to schedule a pod on either the WASI or WASCC nodes*

This will start two instances of Krustlet for both [WASI](https://wasi.dev/) and [waSCC](https://wascc.dev/) runtimes:

```sh
$ kubectl get no

NAME           STATUS   ROLES    AGE   VERSION
cinder         Ready    master   2m    v1.18.5
cinder-wascc   Ready    <none>   1m    0.5.0
cinder-wasi    Ready    <none>   1m    0.5.0
```

With these nodes ready, we can build and push images to our local registry and run them on our Cinder cluster. For example, the [Hello World Rust for WASI](https://github.com/deislabs/krustlet/tree/master/demos/wasi/hello-world-rust) can be built using cargo and pushed to our local registry using [wasm-to-oci](https://github.com/engineerd/wasm-to-oci):

```sh
cargo build --target wasm32-wasi --release
wasm-to-oci push --use-http \
    target/wasm32-wasi/release/hello-world-rust.wasm \
    localhost:5000/hello-world-rust:v0.2.0
```

The line in [k8s.yaml](https://github.com/deislabs/krustlet/blob/v0.5.0/demos/wasi/hello-world-rust/k8s.yaml) specifying the image to use will need to be modified:

```
...
spec:
  containers:
    - name: hello-world-wasi-rust
      #image: webassembly.azurecr.io/hello-world-wasi-rust:v0.2.0
      image: cinderegg:5000/hello-world-rust:v0.2.0
...
```

Finally, the manifest can be applied:

```sh
kubectl apply -f k8s.yaml
```

Which will result in the pod being scheduled on the waSCC Krustlet:

```
$ kubectl get po -A

NAMESPACE            NAME                                  READY   STATUS                          RESTARTS   AGE
kube-system          cilium-operator-657978fb5b-frrxj      1/1     Running                         0          8m4s
kube-system          cilium-pqmsc                          1/1     Running                         0          8m4s
kube-system          coredns-pqljz                         1/1     Running                         0          7m57s
kube-system          hello-world-wasi-rust                 0/1     ExitCode:0                      0          1s
kube-system          kube-apiserver-cinder                 1/1     Running                         0          8m18s
kube-system          kube-controller-manager-cinder        1/1     Running                         0          8m18s
kube-system          kube-proxy-85lwd                      1/1     Running                         0          8m4s
kube-system          kube-scheduler-cinder                 1/1     Running                         0          8m18s
local-path-storage   local-path-storage-74cd8967f5-vv2mb   1/1     Running                         0          8m4s
```

And should produce the following log output:

```sh
$ kubectl logs hello-world-wasi-rust

hello from stdout!
hello from stderr!
POD_NAME=hello-world-wasi-rust
FOO=bar
CONFIG_MAP_VAL=cool stuff
Args are: []

Bacon ipsum dolor amet chuck turducken porchetta, tri-tip spare ribs t-bone ham hock. Meatloaf
pork belly leberkas, ham beef pig corned beef boudin ground round meatball alcatra jerky.
Pancetta brisket pastrami, flank pork chop ball tip short loin burgdoggen. Tri-tip kevin
shoulder cow andouille. Prosciutto chislic cupim, short ribs venison jerky beef ribs ham hock
short loin fatback. Bresaola meatloaf capicola pancetta, prosciutto chicken landjaeger andouille
swine kielbasa drumstick cupim tenderloin chuck shank. Flank jowl leberkas turducken ham tongue
beef ribs shankle meatloaf drumstick pork t-bone frankfurter tri-tip.
```
