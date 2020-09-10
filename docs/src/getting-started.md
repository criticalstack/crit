# Getting Started

## Quick Start

A local Critical Stack cluster can be setup using `cinder` with one easy command:

```shell
$ cinder create cluster
```

This quickly creates a ready-to-use Kubernetes cluster running completely within a single Docker container.

Cinder, or Crit-in-Docker, can be useful for developing on Critical Stack clusters locally, or simply to learn more about Crit. You can read more about requirements, configuration, etc over in the [Cinder Guide](cinder-guide/overview.md).

## Running in Production

Setting up a production Kubernetes cluster requires quite a bit of planning and configuration. For one, there are many considerations that influence the way a cluster should be configured. When starting a new cluster or setting up a standard cluster configuration, one should consider the following:

* Where will it be running? (e.g. AWS, GCP, bare-metal, etc)
* What level of resiliency is required?
  * This is about how the cluster can deal with faults and depending upon factors like colocation of etcd, how it fails can become more complicated.
* What will provide out-of-band storage for cluster secrets?
  * This applies mostly to the initial cluster secrets, the Kubernetes and Etcd CA cert/key pairs.
* What kind of applications will run on the cluster?
* What cost-based factors are there?
* What discovery mechanisms are available for new nodes?
* Are there specific performance requirements that affect the infrastructure being used?

The [Crit Guide](crit-guide/overview.md), and the accompanying [Security Guide](security-guide/overview.md), exists to help answer these questions and provide general guidance for setting up a typical Kubernetes cluster to meet various use-cases.

In particular, a few good places to start planning your Kubernetes cluster:

* [System Requirements](crit-guide/system-requirements.md)
* [Running Etcd](crit-guide/running-etcd.md)
* [Control Plane Sizing](crit-guide/control-plane-sizing.md)
* [Generating Certificates](crit-guide/generating-certificates.md)
