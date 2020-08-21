#!/bin/bash -eu

kubectl create namespace local-path-storage
helm install local-path-storage /cinder/charts/local-path-provisioner.tgz \
    --namespace local-path-storage \
    --set nameOverride=local-path-storage \
    --set storageClass.defaultClass=true
