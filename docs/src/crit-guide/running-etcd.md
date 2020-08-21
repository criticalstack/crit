# Running Etcd

Crit requires a connection to etcd to coordinate the bootstrapping process. The etcd cluster does not have to be colocated on the node. For bootstrapping and managing etcd, we prefer using our own [e2d](https://github.com/criticalstack/e2d) tool. It embeds etcd and combines it with the [hashicorp/memberlist](https://github.com/hashicorp/memberlist) gossip network to manage etcd membership.
