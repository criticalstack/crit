# Configuration

## Adding Files

```yaml
apiVersion: cinder.crit.sh/v1alpha1
kind: ClusterConfiguration
files:
  - path: "/etc/kubernetes/auth-proxy-ca.yaml"
    owner: "root:root"
    permissions: "0644"
    content: |
      apiVersion: cert-manager.io/v1alpha2
      kind: ClusterIssuer
      metadata:
        name: auth-proxy-ca
        namespace: cert-manager
      spec:
        ca:
          secretName: auth-proxy-ca
```

### HostPath

```yaml
apiVersion: cinder.crit.sh/v1alpha1
kind: ClusterConfiguration
  kubeAPIServer:
    extraArgs:
      audit-policy-file: "/etc/kubernetes/audit-policy.yaml"
      audit-log-path: "/var/log/kubernetes/kube-apiserver-audit.log"
      audit-log-maxage: "30"
      audit-log-maxbackup: "10"
      audit-log-maxsize: "100"
    extraVolumes:
    - name: apiserver-logs
      hostPath: /var/log/kubernetes
      mountPath: /var/log/kubernetes
      readOnly: false
      hostPathType: Directory
    - name: apiserver-audit-config
      hostPath: /etc/kubernetes/audit-policy.yaml
      mountPath: /etc/kubernetes/audit-policy.yaml
      readOnly: true
files:
  - path: "/etc/kubernetes/audit-policy.yaml"
    owner: "root:root"
    permissions: "0644"
    encoding: hostpath
    content: audit-policy.yaml
```

## Running Additional Commands

```yaml
apiVersion: cinder.crit.sh/v1alpha1
kind: ClusterConfiguration
preCritCommands:
  - crit version
postCritCommands:
  - |
    helm repo add jetstack https://charts.jetstack.io
    helm install cert-manager jetstack/cert-manager \
      --namespace cert-manager \
      --version v0.15.1 \
      --set tolerations[0].effect=NoSchedule \
      --set tolerations[0].key="node.kubernetes.io/not-ready" \
      --set tolerations[0].operator=Exists \
      --set installCRDs=true
    kubectl rollout status -n cert-manager deployment/cert-manager-webhook -w
```


## Updating the Containerd Configuration

```yaml
apiVersion: cinder.crit.sh/v1alpha1
kind: ClusterConfiguration
files:
  - path: "/etc/containerd/config.toml"
    owner: "root:root"
    permissions: "0644"
    content: |
      # explicitly use v2 config format
      version = 2

      # set default runtime handler to v2, which has a per-pod shim
      [plugins."io.containerd.grpc.v1.cri".containerd]
        default_runtime_name = "runc"
      [plugins."io.containerd.grpc.v1.cri".containerd.runtimes.runc]
        runtime_type = "io.containerd.runc.v2"

      # Setup a runtime with the magic name ("test-handler") used for Kubernetes
      # runtime class tests ...
      [plugins."io.containerd.grpc.v1.cri".containerd.runtimes.test-handler]
        runtime_type = "io.containerd.runc.v2"

      [plugins."io.containerd.grpc.v1.cri".registry.mirrors]
        [plugins."io.containerd.grpc.v1.cri".registry.mirrors."docker.io"]
          endpoint = ["https://docker.io"]
```

## Adding Volume Mounts

```yaml
apiVersion: cinder.crit.sh/v1alpha1
kind: ClusterConfiguration
extraMounts:
  - hostPath: templates
    containerPath: /cinder/templates
    readOnly: true
```

## Forwarding Ports to the Host

```yaml
apiVersion: cinder.crit.sh/v1alpha1
kind: ClusterConfiguration
extraPortMappings:
  - containerPort: 2379
    hostPort: 2379
```
