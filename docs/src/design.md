# Design

### Decoupled from Etcd Management

The Kubernetes control plane requires etcd for storage, however, bootstrapping and managing etcd is not a responsibility of Crit. This decreases code complexity and results in more maintainable code. Rather than handle all aspects of installing and managing Kubernetes, Crit is designed to be one tool in the toolbox, specific to bootstrapping Kubernetes components.

Safely handling etcd in a cloud environment is not as easy as it may seem, so we have a separate project, [e2d](https://github.com/criticalstack/e2d), designed to bootstrap etcd and manage cluster membership.

### Lazy Cluster Initialization

Crit leverages the unique features of etcd to handle how brand new clusters are bootstrapped. With other tooling, this is often accomplished by handling cluster initialization separately from all subsequent nodes joining the cluster (even if done so implicitly). The complexity for handling this initial case can be difficult to automate in distributed systems. Instead, the distributed locking capabilities of etcd are used to synchronize nodes and initialize a cluster automatically. All nodes race to acquire the distributed lock, and should the cluster not exist (signified by the presence of shared cluster files), a new cluster is initialized by the node that was first to acquire the lock, otherwise the node joins the cluster.

This ends up being really cool when working with projects like [cluster-api](https://cluster-api.sigs.k8s.io/), since all control plane nodes can be initialized simultaneously, greatly reducing the time to create a HA cluster (especially a 5 node control plane).

### Node Roles

Nodes are distinguished as having only one of two roles, either control plane or worker. All the same configurations for clusters are possible, such as colocating etcd on the control plane, but Crit is only concerned with how it needs to bootstrap the two basic node roles.


### Cluster Upgrades

There are a several important considerations for upgrading a cluster. Crit itself is only a bootstrapper, in that it takes on the daunting task of ensuring that the cluster components are all configured, but afterwards, there is not much left for it to do. However, the most important aspects of the philosophy behind Crit and e2d is to ensure that colocated control planes can:

1. Have all nodes deployed simultaneously, and crit/e2d will ensure that they are bootstrapped regardless of the order they come up.
2. It can safely perform a rolling upgrade.


### Built for Automation

