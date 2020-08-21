#!/bin/bash

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

systemctl daemon-reload

cat << EOF > config.toml
control_plane_host = "${CONTROL_PLANE_HOST}"

[network]
host_ipv4 = "${WORKER_HOST_IP}"

[bootstrap]
bootstrap_token = "abcdef.0123456789abcdef"
ca_cert = "/etc/kubernetes/pki/ca.crt"
EOF

crit worker up -c config.toml
