#!/bin/bash -eu
# This is a workaround to several cgroup mount problems:
# https://github.com/kubernetes-sigs/kind/issues/1614

set -o pipefail

CURRENT_CGROUP=$(grep systemd /proc/self/cgroup | cut -d: -f3)
CGROUP_SUBSYSTEMS=$(findmnt -lun -o source,target -t cgroup | grep "${CURRENT_CGROUP}" | awk '{print $2}')

mount --make-rprivate /sys/fs/cgroup

echo "${CGROUP_SUBSYSTEMS}" |
while IFS= read -r SUBSYSTEM; do
  mkdir -p "${SUBSYSTEM}/kubelet"
  if [ "${SUBSYSTEM}" == "/sys/fs/cgroup/cpuset" ]; then
    cat "${SUBSYSTEM}/cpuset.cpus" > "${SUBSYSTEM}/kubelet/cpuset.cpus"
    cat "${SUBSYSTEM}/cpuset.mems" > "${SUBSYSTEM}/kubelet/cpuset.mems"
  fi
  mount --bind "${SUBSYSTEM}/kubelet" "${SUBSYSTEM}/kubelet"
done
