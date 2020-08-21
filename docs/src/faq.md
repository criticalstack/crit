# FAQ

# TODO: modify this a bit

### Can e2d scale up (or down) after cluster initialization?

The short answer is No, because it is unsafe to scale etcd and any solution that scales etcd is increasing the chance of cluster failure. This is a feature that will be supported in the future, but it relies on new features and fixes to etcd. Some context will be necessary to explain why:

A common misconception about etcd is that it is scalable. While etcd is a distributed key/value store, the reason it is distributed is to provide for distributed consensus, *NOT* to scale in/out for performance (or flexibility). In fact, the best performing etcd cluster is when it only has 1 member and the performance goes down as more members are added. In etcd v3.4, a new type of member called learners was introduced. These are members that can receive raft log updates, but are not part of the quorum voting process. This will be an important feature for many reasons, like stability/safety and faster recovery from faults, but will also potentially<sup>[[1]](#faq-fn-1)</sup> enable etcd clusters of arbitrary sizes.

So why not scale within the [recommended cluster sizes](https://github.com/etcd-io/etcd/blob/master/Documentation/faq.md#what-is-maximum-cluster-size) if the only concern is performance? Previously, etcd clusters have been vulnerable to corruption during membership changes due to the way etcd implemented raft. This has only recently been addressed by incredible work from CockroachDB, and it is worth reading about the issue and the solution in this blog post: [Availability and Region Failure: Joint Consensus in CockroachDB](https://www.cockroachlabs.com/blog/joint-consensus-raft/).

The last couple features needed to safely scale have been roadmapped for v3.5 and are highlighted in the [etcd learner design doc](https://github.com/etcd-io/etcd/blob/master/Documentation/learning/design-learner.md#features-in-v35):

> Make learner state only and default: Defaulting a new member state to learner will greatly improve membership reconfiguration safety, because learner does not change the size of quorum. Misconfiguration will always be reversible without losing the quorum.

> Make voting-member promotion fully automatic: Once a learner catches up to leader’s logs, a cluster can automatically promote the learner. etcd requires certain thresholds to be defined by the user, and once the requirements are satisfied, learner promotes itself to a voting member. From a user’s perspective, “member add” command would work the same way as today but with greater safety provided by learner feature.

Since we want to implement this feature as safely and reliably as possible, we are waiting for this confluence of features to become stable before finally implementing scaling into e2d.

<a name="faq-fn-1">[1]</a> Only potentially, because the maximum is currently set to allow only 1 learner. There is a concern that too many learners could have a negative impact on the leader which is discussed briefly [here](https://github.com/etcd-io/etcd/issues/11401). It is also worth noting that other features may also fulfill the same need like some kind of follower replication: [etcd#11357](https://github.com/etcd-io/etcd/issues/11357).
