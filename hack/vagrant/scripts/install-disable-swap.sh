#!/bin/bash

# kubelet will not run when swap is enabled
cat << EOF > /etc/systemd/system/disable-swap.service
[Unit]
After=local-fs.target

[Service]
ExecStart=swapoff -a

[Install]
WantedBy=multi-user.target
EOF
systemctl enable --now  disable-swap.service
