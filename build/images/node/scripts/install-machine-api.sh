#!/bin/bash

kubectl create namespace mapi-system
helm install machine-api /cinder/charts/machine-api.tgz --namespace mapi-system \
    --set externalReadyWait=1s
kubectl create namespace mapd-system
helm install machine-api-provider-docker /cinder/charts/machine-api-provider-docker.tgz --namespace mapd-system
