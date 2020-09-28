# Configuring Control Plane Components

See [here](https://pkg.go.dev/github.com/criticalstack/crit@v1.0.1/pkg/config/v1alpha2#ControlPlaneConfiguration) for a complete list of available configuration options.
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

## Configuring a Cloud Provider

A cloud provider can be specified to integrate with the underlying infrastructure provider. Note, the specified cloud will most likely require authentication/authorization to access their APIs.

Crit supports both [In-tree and out-of-tree](https://kubernetes.io/blog/2019/04/17/the-future-of-cloud-providers-in-kubernetes/#in-tree-out-of-tree-providers) cloud providers. 

### In-tree Cloud Provider

[In-tree](https://github.com/kubernetes/kubernetes/blob/master/pkg/cloudprovider/providers/providers.go#L21-L27) cloud providers can be specified with the following: 

```yaml
apiVersion: crit.sh/v1alpha2
kind: ControlPlaneConfiguration
...
node:
  cloudProvider: aws
```

and for the workers: 

```yaml
apiVersion: crit.sh/v1alpha2
kind: WorkerConfiguration
...
node:
  cloudProvider: aws
```

### Out-of-tree Cloud Provider

Out-of-tree cloud providers can be specified with the following: 

```yaml
apiVersion: crit.sh/v1alpha2
kind: ControlPlaneConfiguration
...
node:
  kubeletExtraArgs: 
    cloud-provider: external
```

and for the workers: 

```yaml
apiVersion: crit.sh/v1alpha2
kind: WorkerConfiguration
...
node:
  kubeletExtraArgs: 
    cloudProvider: external
```

A manifest specific to cloud environment must then be applied to run the external [cloud controller manager](https://kubernetes.io/docs/tasks/administer-cluster/running-cloud-controller/#examples). 


