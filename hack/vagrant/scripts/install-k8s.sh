#!/bin/bash

KUBE_VERSION=1.15.3

export INPUT_ARGUMENTS="${@}"
set -u
while [[ $# -gt 0 ]]; do
  case $1 in
    '--version'|-v)
       shift
       if [[ $# -ne 0 ]]; then
           KUBE_VERSION="${1}"
       else
           echo -e "Please provide the desired version. e.g. --version 1.14.3 or -v 1.*"
           exit 0
       fi
       ;;
    '--help'|-h)
       echo "Usage: install-k8s.sh [-h|--help] [-v|--version <version>]"
       echo "  -h --help       Print this help and exit."
       echo "  -v --version    Set the version of kubelet and kubeadm to"
       echo "                  install (accepts wildcards)."
       exit 0
       ;;
    *) exit 1
       ;;
  esac
  shift
done
set +u

cat << EOF > /etc/yum.repos.d/kubernetes.repo
[kubernetes]
name=Kubernetes
baseurl=https://packages.cloud.google.com/yum/repos/kubernetes-el7-x86_64
enabled=1
gpgcheck=1
repo_gpgcheck=1
gpgkey=https://packages.cloud.google.com/yum/doc/yum-key.gpg https://packages.cloud.google.com/yum/doc/rpm-package-key.gpg
exclude=kube*
EOF

mkdir -p /etc/kubernetes \
  /etc/kubernetes/pki \
  /etc/kubernetes/manifests \
  /var/lib/kubelet \
  /var/lib/etcd

dnf install -y kubelet-$KUBE_VERSION kubectl --disableexcludes=kubernetes

systemctl enable kubelet

# br_netfilter must be loaded for the kubelet
cat << EOF > /etc/modules-load.d/br_netfiler.conf
br_netfilter
EOF

# kubelet also requires ip forwarding
cat << EOF >> /etc/sysctl.conf
net.ipv4.ip_forward = 1
EOF
