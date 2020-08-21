#!/bin/bash -eu

verlte() {
    [  "$1" = "$(echo -e "$1\n$2" | sort -V | head -n1)" ]
}

VALUES_FILE=/cinder/cilium-values.yaml
if $(verlte $(uname -r) 4.17); then
    VALUES_FILE=/cinder/cilium-compat-values.yaml
fi

helm install cilium /cinder/charts/cilium.tgz --namespace kube-system \
    --set config.ipam=kubernetes \
    -f $VALUES_FILE
