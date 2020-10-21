#!/bin/bash

kubectl create namespace mapi-system
kubectl apply -n mapi-system -f /cinder/manifests/machine-api.yaml
kubectl create namespace mapd-system
kubectl apply -n mapd-system -f /cinder/manifests/machine-api-provider-docker.yaml
