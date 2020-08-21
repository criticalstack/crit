#!/bin/bash -eu


DEBIAN_FRONTEND=noninteractive clean-install gnupg software-properties-common
curl -sL https://packagecloud.io/criticalstack/public/gpgkey | apt-key add -
apt-add-repository "deb https://packagecloud.io/criticalstack/public/ubuntu/ bionic main"
DEBIAN_FRONTEND=noninteractive clean-install e2d
curl -sL https://packages.cloud.google.com/apt/doc/apt-key.gpg | apt-key add -
apt-add-repository "deb http://apt.kubernetes.io/ kubernetes-xenial main"
DEBIAN_FRONTEND=noninteractive && clean-install kubectl=${KUBERNETES_VERSION}-00
curl -L https://storage.googleapis.com/kubernetes-release/release/v${KUBERNETES_VERSION}/bin/linux/amd64/kubelet -o /usr/bin/kubelet
chmod +x /usr/bin/kubelet
echo "KUBELET_EXTRA_ARGS=--fail-swap-on=false" >> /etc/default/kubelet
curl https://raw.githubusercontent.com/helm/helm/master/scripts/get-helm-3 | bash
curl -LO https://download.docker.com/linux/static/stable/x86_64/docker-19.03.1.tgz
tar zxvf docker-19.03.1.tgz --strip 1 -C /usr/bin docker/docker
rm docker-19.03.1.tgz
curl -LO https://krustlet.blob.core.windows.net/releases/krustlet-v0.3.0-linux-amd64.tar.gz
tar zxvf krustlet-v0.3.0-linux-amd64.tar.gz -C /usr/bin
rm krustlet-v0.3.0-linux-amd64.tar.gz
mkdir -p /var/lib/krustlet
mkdir -p /etc/kubernetes/pki
mkdir -p /kind
mkdir -p /etc/systemd/system/kubelet.service.d
mkdir -p /etc/cni/net.d
mkdir -p /var/lib/crit
mkdir -p /var/log/kubernetes
chmod 700 /var/lib/etcd
