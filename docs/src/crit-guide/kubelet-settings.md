# Kubelet Settings

## Disable swap for Linux-based Operating Systems

Swap cannot be enabled for the kubelet to work (see [here](https://github.com/kubernetes/kubernetes/issues/53533)). This is a helpful drop-in to ensure that swap is disabled on a system:

```toml
[Unit]
After=local-fs.target

[Service]
ExecStart=/sbin/swapoff -a

[Install]
WantedBy=multi-user.target
```

## Reserving Resources

Reserving some resources for the system to use is often times very helpful to ensure that resource hungry pods don't kill the system by causing it to run out of memory.

```yaml
...
node:
  kubelet:
    kubeReserved:
      cpu: 128m
      memory: 64Mi
    kubeReservedCgroup: /podruntime.slice
    kubeletCgroups: /podruntime.slice
    systemReserved:
      cpu: 128m
      memory: 192Mi
    systemReservedCgroup: /system.slice
```

```toml
# /etc/systemd/system/kubelet.service.d/10-cgroup.conf
# Sets the cgroup for the kubelet service
[Service]
CPUAccounting=true
MemoryAccounting=true
Slice=podruntime.slice
```

```toml
# /etc/systemd/system/containers.slice
# Creates a cgroup for kubelet
[Unit]
Description=Grouping resources slice for containers
Documentation=man:systemd.special(7)
DefaultDependencies=no
Before=slices.target
Requires=-.slice
After=-.slice
```

```toml
# /etc/systemd/system/podruntime.slice
# Creates a cgroup for kubelet
[Unit]
Description=Limited resources slice for Kubelet service
Documentation=man:systemd.special(7)
DefaultDependencies=no
Before=slices.target
Requires=-.slice
After=-.slice
```
