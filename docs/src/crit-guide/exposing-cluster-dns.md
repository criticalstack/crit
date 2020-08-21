# Exposing Cluster DNS

## Replace Systemd-resolved With Dnsmasq

Sometimes systemd-resolved, default stub resolver for many linux systems, needs to be replaced with dnsmasq. This dnsmasq systemd drop-in is useful to ensure that systemd-resolved is not running when the dnsmasq service is started:

```
# /etc/systemd/system/dnsmasq.service.d/10-resolved-fix.conf
[Unit]
After=systemd-resolved.service

[Service]
ExecStartPre=/bin/systemctl stop systemd-resolved.service
ExecStartPost=/bin/systemctl start systemd-resolved.service
```

It works by allowing systemd-resolved to start, but stopping it once the dnsmasq service is started. This is helpful because it doesn't require changing any of the systemd-resolved specific settings but allows the dnsmasq service to be enabled/disabled when desired.

## Forwarding Cluster-bound DNS on the Host

A reason why one might want to use something like dnsmasq, instead of systemd-resolved, is to expose the cluster DNS to the host. This would allow resolution of DNS for [service and pod subnets](configuring-control-plane-components.md#configuring-podservice-subnets) from the host that is running the Kubernetes components. It only requires adding this dnsmasq configuration drop-in:

```
# /etc/dnsmasq.d/kube.conf
server=/cluster.local/10.254.0.10
```

This tells dnsmasq to forward any DNS queries it receives that end in the cluster domain, to the Kubernetes cluster dns, CoreDNS. In this case, it is presuming that the default cluster domain (`cluster.local`) and services subnet, have been configured. The address of CoreDNS is chosen automatically based upon the services subnet, so if the services subnet is `10.254.0.0/16` (the default), CoreDNS will be listening at `10.254.0.10`.

## Specifying the resolv.conf

The default resolv.conf from the host is used, `/etc/resolv.conf` or a different conf file can be set in the Kubelet:

```yaml
apiVersion: crit.sh/v1alpha2
kind: ControlPlaneConfiguration
node:
  kubelet:
    resolvConf: /other/resolv.conf
```

Crit attempts to determine if systemd-resolved is running, and dynamically sets `resolvConf` to be `/run/systemd/resolve/resolv.conf`.
