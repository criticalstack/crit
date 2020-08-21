# aws

## Dependencies

 * Terraform >= 0.12
 * Go >= 1.12
 * kubectl

## Getting started

Before getting started make sure that access to AWS is setup either via something like `~/.aws/credentials` or by having the required environment variables set: `AWS_ACCESS_KEY_ID`, `AWS_SECRET_ACCESS_KEY` and `AWS_SECRET_TOKEN` (if it applies).

First, a configuration file is necessary for [Terraform](https://www.terraform.io/docs/index.html) to deploy a Kubernete cluster using crit. The [genconf](https://github.com/criticalstack/crit/tree/master/tools/genconf) tool will be used to help generate the required `terraform.tfvars.json` file and it should be built automatically by the Makefile target (so long as Go is setup and configured correctly). The first time running either of the Makefile targets that run terraform, genconf will detect that there is no configuration file and will provide interactive prompts to create the new `terraform.tfvars.json` file. The aws credentials are crucial to have before this step as genconf helps to contextualize the setup using calls into the target account (currently to discover the available VPCs).

Creating a brand new cluster should be as easy as running:

```bash
make apply
```

This will deploy the nodes as set in the configuration and will wait for the controlplane to become available. Once it is available it will download the admin kubeconfig and set it as your current context. If everything worked then you should be able to run kubectl commands:

```bash
kubectl cluster-info
```

If needing to make changes to existing infrastructure, such as changing the size of the control plane or worker pool, simply run:

```bash
make update
```

This will prompt `Keep current settings? (Y/n)` which will allow changing the values before re-deploying the infrastructure. Go ahead and try changing from a single-node controlplane to a HA cluster (3 or 5 nodes) and this process should seemlessly deploy and configure the new controlplane!

While developing crit, the binary can be rebuilt and calls to `make update` will ensure the new binary is installed on the new instances.
