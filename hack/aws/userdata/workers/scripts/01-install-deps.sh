#!/bin/bash


export DEBIAN_FRONTEND=noninteractive
export AWS_DEFAULT_REGION="$(curl -s http://169.254.169.254/latest/meta-data/placement/availability-zone | sed 's/[a-z]$//')"
apt-get update && apt-get install -y vim jq bash-completion bind9utils netcat etcd-client awscli


#####################################################################
# Install Kubernetes and e2d
#####################################################################

# the hostname MUST be set correctly so that cert SANs match for authz
hostnamectl set-hostname $(curl http://169.254.169.254/latest/meta-data/hostname)

curl -sL https://packages.cloud.google.com/apt/doc/apt-key.gpg | apt-key add -
apt-add-repository "deb http://apt.kubernetes.io/ kubernetes-xenial main"

apt-get install -y kubelet=${kubernetes_version}-00 kubectl kubernetes-cni containerd
apt-mark hold kubelet
systemctl stop kubelet.service

kubectl completion bash > /usr/share/bash-completion/completions/kubectl

#####################################################################
# Install helm
#####################################################################

curl https://raw.githubusercontent.com/helm/helm/master/scripts/get-helm-3 | bash
