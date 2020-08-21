# vagrant

**This is no longer being used**

## Getting started

First, the base vagrant box must be built locally:

```bash
make base
```

The base box will configure the OS and install packages required to run crit. This only needs to be ran **once**, or again whenever changes to the base box are desired.


If the required [dependencies](#dependencies) are installed, start a 1 control/1 worker crit cluster with simply:

```bash
make
```


The default makefile target handles creating initial secrets, destroying any existing VMs, and provisioning new ones. The default configuration starts 1 control and worker node but can be changed within the `Vagrantfile` by tweaking the `$control_count` and `$worker_count` values (this will probably be improved in the future).

The admin kubeconfig is placed in `.kube/config` and is automatically merged into your user's default kubeconfig (`$HOME/.kube/config`). You can remove the kubeconfig by running `make remove-ctx`.


## Dependencies

 * Vagrant
 * VirtualBox

The [bento/fedora](https://github.com/chef/bento) box was used as a base here because it installs the VirtualBox Guest Additions necessary for features such as sync folders. It is not intended to be an integration test of any kind with a particular operating system.

It is recommended to install `vagrant` directly from Hashicorp's [website](https://www.vagrantup.com/downloads.html) instead of through the package manager if the package manager's version is not up to date (common for Debian/Ubuntu based distros). The Arch package can be used directly.

## Troubleshooting
If you get a build error with `cannot find package` errors, make sure you have `GO111MODULE=on` set in your environment:
```bash
export GO111MODULE=on
```

If `make` fails, it may need to be run twice if you have a VM in a state that causes the `destroy` make target to fail.
