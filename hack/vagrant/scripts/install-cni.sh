#!/bin/bash


dnf install -y containernetworking-cni containernetworking-plugins

mkdir -p /opt/cni \
  /etc/cni/net.d

ln -s /usr/libexec/cni /opt/cni/bin

cat << EOF > /etc/cni/net.d/10-bridge.conf
{
    "cniVersion": "0.3.1",
    "name": "bridge",
    "type": "bridge",
    "bridge": "cnio0",
    "isGateway": true,
    "ipMasq": true,
    "ipam": {
        "type": "host-local",
        "ranges": [
          [{"subnet": ""}]
        ],
        "routes": [{"dst": "0.0.0.0/0"}]
    }
}
EOF

cat << EOF > /etc/cni/net.d/99-loopback.conf
{
    "cniVersion": "0.3.1",
    "type": "loopback"
}
EOF

systemctl restart containerd
