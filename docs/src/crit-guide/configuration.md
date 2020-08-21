# Configuration

Configuration is passed to Crit via yaml and is separated into two different types: `ControlPlaneConfiguration` and `WorkerConfiguration`. The concerns are split between these two configs for any given node, and they each contain a `NodeConfiguration` that specifies node-specific settings like for the [Kubelet], networking, etc.

## Embedded ComponentConfigs

The [ComponentConfigs](https://github.com/kubernetes/enhancements/blob/master/keps/sig-cluster-lifecycle/wgs/0014-20180707-componentconfig-api-types-to-staging.md) are part of an ongoing effort to make configuration of Kubernetes components (API server, kubelet, etc) more dynamic by making configuration directly through Kubernetes API types. Crit will be using these ComponentConfigs when available since they simplify all aspects of taking user configuration and transforming that into Kubernetes component configuration. Kubernetes components are being changed to support direct configuration from file with the ComponentConfig API types, so Crit embeds these to make configuration more straightforward.

Currently, only the [`kube-proxy`](https://github.com/kubernetes/kube-proxy) and [`kubelet`](https://github.com/kubernetes/kubelet) ComponentConfigs are ready to be used, but more are currently being worked on and will be adopted by Crit as other components begin supporting configuration from file.

## Runtime Defaults

Some configuration defaults are set at the time of running [`crit up`](../crit-commands/crit-up.md). These mostly include settings that are based upon the host that is running the command, such as the hostname.

If left unset, the `controlPlaneEndpoint` value will be set to the ipv4 of the host. In the case there are multiple network interfaces, the first non-loopback network interface is used.

The default directory for Kubernetes files is `/etc/kubernetes` and any paths to manifests, certificates, etc are derived from this.

Etcd is also configured presuming that mTLS is used and that the etcd nodes are colocated with the Kubernetes control plane components, effectively making this the default configuration:

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

The CA certificate required for the worker to validate the cluster it's joining is also derived from the default Kubernetes configuration directory:

```yaml
apiVersion: crit.sh/v1alpha2
kind: WorkerConfiguration
caCert: /etc/kubernetes/pki/ca.crt
```
