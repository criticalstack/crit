#!/bin/bash

setenforce enforcing
mkdir -p /etc/kubernetes
chcon -R -t svirt_sandbox_file_t /etc/kubernetes
mkdir -p /var/lib/etcd
chcon -R -t svirt_sandbox_file_t /var/lib/etcd
