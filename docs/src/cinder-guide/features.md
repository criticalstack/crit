# Features

## Side-loading Images

## Krustlet

```yaml
apiVersion: cinder.crit.sh/v1alpha1
kind: ClusterConfiguration
featureGates:
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

https://github.com/deislabs/krustlet

## Local Registry

```yaml
apiVersion: cinder.crit.sh/v1alpha1
kind: ClusterConfiguration
featureGates:
  LocalRegistry: true
```

const LocalRegistryHostingConfigMap = `apiVersion: v1
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

https://github.com/kubernetes/enhancements/pull/1757

https://github.com/kubernetes/enhancements/blob/c5b6b632811c21ababa9e3565766b2d70614feec/keps/sig-cluster-lifecycle/generic/1755-communicating-a-local-registry/README.md#design-details

## Registry Mirrors

```yaml
apiVersion: cinder.crit.sh/v1alpha1
kind: ClusterConfiguration
registryMirrors:
  docker.io: "https://docker.io"
```
