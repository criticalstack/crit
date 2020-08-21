#!/bin/bash -e

systemctl start containerd-sync.service

mkdir -p /etc/systemd/system/kubelet.service.d

cat << 'EOF' > /etc/systemd/system/kubelet.service.d/20-crit.conf
[Service]
Environment="KUBELET_KUBECONFIG_ARGS=--bootstrap-kubeconfig=/etc/kubernetes/bootstrap-kubelet.conf --kubeconfig=/etc/kubernetes/kubelet.conf"
Environment="KUBELET_CONFIG_ARGS=--config=/var/lib/kubelet/config.yaml"
# This is a file that "crit control up" and "crit worker up" generates at runtime, populating the KUBELET_CRIT_ARGS variable dynamically
EnvironmentFile=-/var/lib/kubelet/crit-flags.env
# This is a file that the user can use for overrides of the kubelet args as a last resort. Preferably, the user should use
# the .NodeRegistration.KubeletExtraArgs object in the configuration files instead. KUBELET_EXTRA_ARGS should be sourced from this file.
EnvironmentFile=-/etc/default/kubelet
ExecStart=
ExecStart=/usr/bin/kubelet $KUBELET_KUBECONFIG_ARGS $KUBELET_CONFIG_ARGS $KUBELET_CRIT_ARGS $KUBELET_EXTRA_ARGS
EOF

crit e2d pki gencerts \
  --ca-cert=/etc/kubernetes/pki/etcd/ca.crt \
  --ca-key=/etc/kubernetes/pki/etcd/ca.key \
  --hosts=${CONTROL_PLANE_HOST} \
  --output-dir=/etc/kubernetes/pki/etcd

# TODO(chris): install e2d service cli
cat << EOF > /etc/systemd/system/e2d.service
[Unit]
Description=e2d

[Service]
ExecStart=/usr/bin/crit e2d run \
  --host=${CONTROL_PLANE_HOST} \
  --bootstrap-addrs=${PEERS} \
  --data-dir=/var/lib/etcd \
  --ca-cert=/etc/kubernetes/pki/etcd/ca.crt \
  --server-cert=/etc/kubernetes/pki/etcd/server.crt \
  --server-key=/etc/kubernetes/pki/etcd/server.key \
  --peer-cert=/etc/kubernetes/pki/etcd/peer.crt \
  --peer-key=/etc/kubernetes/pki/etcd/peer.key \
  --required-cluster-size=${ETCD_CLUSTER_SIZE}
Restart=on-failure
RestartSec=30

[Install]
WantedBy=multi-user.target
EOF
systemctl enable --now e2d.service

cat << EOF > config.toml
control_plane_host = "${CONTROL_PLANE_HOST}"

[network]
host_ipv4 = "${CONTROL_PLANE_HOST}"

[network.extra_args]
prometheus-serve-addr = ":9090"

[kube_apiserver.extra_args]
audit-policy-file = "/etc/kubernetes/audit-policy.yaml"
audit-log-path = "/var/log/kubernetes/kube-apiserver-audit.log"
audit-log-format = "legacy"
audit-log-maxage = "90"
audit-log-maxbackup = "5"
audit-log-maxsize = "100"

[[kube_apiserver.extra_volumes]]
name = "apiserver-logs"
host_path = "/var/log/kubernetes"
mount_path = "/var/log/kubernetes"
read_only = false
host_path_type = "directory"

[[kube_apiserver.extra_volumes]]
name = "apiserver-audit-config"
host_path = "/etc/kubernetes/audit-policy.yaml"
mount_path = "/etc/kubernetes/audit-policy.yaml"
read_only = true

[etcd]
endpoints = [
  "https://${CONTROL_PLANE_HOST}:2379",
]
ca_file = "/etc/kubernetes/pki/etcd/ca.crt"
cert_file = "/etc/kubernetes/pki/etcd/client.crt"
key_file = "/etc/kubernetes/pki/etcd/client.key"
EOF

# TODO(chris): The config is being explicitly set here to prevent a problem
# with how config via stdin is being detected. Should try to figure out a
# better solution going forward, but this works for now.
crit template -c config.toml audit-policy.yaml > /etc/kubernetes/audit-policy.yaml
crit control up -c config.toml
crit token create abcdef.0123456789abcdef

mkdir -p /root/.kube
cp /etc/kubernetes/admin.conf /root/.kube/config

mkdir -p /vagrant/.kube
cp /etc/kubernetes/admin.conf /vagrant/.kube/config
