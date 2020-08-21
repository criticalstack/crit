# System Requirements

Exact system requirements will be dependent upon a lot of factors, however, for the most part, any relatively modern linux operating system will fit the bill.

 * `Linux kernel` >= 4.9.17
 * systemd
 * iptables (optional)

_Newer versions of the kernel will enable using cilium's kube-proxy replacement feature, which will replace the need to deploy kube-proxy (and therefore not need iptables also)._

### Dependencies

 * kubelet >= 1.14.x
 * containerd >= 1.2.6
 * CNI >= 0.7.5

### References

 * https://kubernetes.io/docs/concepts/extend-kubernetes/compute-storage-net/network-plugins/#cni
 * https://kubernetes.io/docs/concepts/extend-kubernetes/compute-storage-net/network-plugins/#support-hostport
 * https://docs.cilium.io/en/v1.6/gettingstarted/cni-chaining-portmap/#portmap-hostport
