# Control Plane Sizing

### External Etcd

When dealing with etcd running external to the Kubernetes control plane components, there are not a lot of restrictions on how many control plane nodes one can have. There can be any number of nodes that meet demand and availability needs, and can even be auto-scaled. With that said, however, the performance of Kubernetes is tied heavily to the performance of etcd, so more nodes does not mean more performance.

### Colocated Etcd

Colocation of etcd, or "stacked etcd" (as it's referred to in the Kubernetes documentation), is the practice of installing etcd alongside the Kubernetes control plane components (`kube-apiserver`, `kube-controller-manager`, and `kube-scheduler`). This has some obvious benefits like reducing cost by reducing the virtual machines needed, but introduces a lot of complexity and restrictions.

Etcd's performance goes down the more nodes that are added, because more members are required to vote to commit to the raft log, so there should never be more than 5 voting members in a cluster (unless performing a rolling upgrade). Also, the number of members should always be odd to help protect against the split-brain problem. This means that the control plane can only safely be made up of 1, 3, or 5 nodes.

Etcd also should not be scaled up or down (at least, at this time). The reason is that the etcd cluster is being put at risk each time there is a membership change, so this also means that the control plane size needs to be selected ahead of time and not be altered.

#### General Recommendations

In cloud environments, 3 is a good size to balance resiliency and performance. The reasoning here is that cloud environments provide ways to quickly automate replacing failed members, so losing a node does not put etcd in danger of losing quorum for long until a new node can replace the existing one. As etcd moves towards adding more functionality around the learners member type, this will also open up to having a "hot spare" ready to take the place of the failed member immediately.

For bare-metal, 5 is a good size to ensure that failed nodes have more time to be replaced since a new node might need to be physically allocated.
