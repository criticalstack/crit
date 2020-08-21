# Configuring Control Plane Components

## Control plane endpoint

The control plane endpoint is the address (IP or DNS), along with optional port, that represents the control plane. It it is effectively the API server address, however, it is internally used for a few other purposes, such as:

* Discovering other services using the host (e.g. [bootstrap-server](bootstrap-server.md))
* Adding to the SAN for generated cluster certificates

It is specified in the config file like so:

```yaml
apiVersion: crit.sh/v1alpha2
kind: ControlPlaneConfiguration
controlPlaneEndpoint: "example.com:6443"
```

## Disable/Enable Kubernetes Feature Gates

Setting feature gates will be important if you need specific features that are not available by default or maybe to enable a feature that wasn't enabled by default for a particular version of Kubernetes.

For example, CSI-related features were only enabled by default starting with version 1.17, so for older versions of Kubernetes you will need to turn them on manually for the control plane:

```yaml
apiVersion: crit.sh/v1alpha2
kind: ControlPlaneConfiguration
...
kubeAPIServer:
  featureGates:
    CSINodeInfo: true
    CSIDriverRegistry: true
    CSIBlockVolume: true
    VolumeSnapshotDataSource: true
node:
  kubelet:
    featureGates:
      CSINodeInfo: true
      CSIDriverRegistry: true
      CSIBlockVolume: true
```

and for the workers:

```yaml
apiVersion: crit.sh/v1alpha2
kind: WorkerConfiguration
...
node:
  kubelet:
    featureGates:
      CSINodeInfo: true
      CSIDriverRegistry: true
      CSIBlockVolume: true
```

The `kubeAPIServer`, `kubeControllerManager`, `kubeScheduler`, and `kubelet` all have feature gates that can be configured. More info is available in the [Kubernetes docs](https://kubernetes.io/docs/reference/command-line-tools-reference/feature-gates/).


## Configuring Pod/Service Subnets

```yaml
apiVersion: crit.sh/v1alpha2
kind: ControlPlaneConfiguration
podSubnet: "10.153.0.0/16"
serviceSubnet: "10.154.0.0/16"
```
