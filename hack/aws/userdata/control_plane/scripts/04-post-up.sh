#!/bin/bash


#####################################################################
# Post-up script
#####################################################################

mkdir -p /root/.kube
cp /etc/kubernetes/admin.conf /root/.kube/config

export KUBECONFIG=/root/.kube/config

helm repo add criticalstack https://charts.cscr.io/criticalstack

helm install cilium criticalstack/cilium \
    --namespace kube-system \
    --set global.prometheus.enabled=true \
    --version 1.7.1

helm install aws-ebs-csi-driver criticalstack/aws-ebs-csi-driver \
    --set enableVolumeScheduling=true \
    --set enableVolumeResizing=true \
    --set enableVolumeSnapshot=true \
    --version 0.3.0

kubectl apply -f /etc/kubernetes/default-storage-class.yaml

# install cert-manager
kubectl create namespace cert-manager
kubectl apply --validate=false -f https://github.com/jetstack/cert-manager/releases/download/v0.14.0/cert-manager.yaml

# add auth-proxy-ca secret to be used as ClusterIssuer
kubectl -n cert-manager create secret generic auth-proxy-ca --from-file=tls.crt=/etc/kubernetes/pki/auth-proxy-ca.crt --from-file=tls.key=/etc/kubernetes/pki/auth-proxy-ca.key

# wait for cert-manager-webhook readiness
while [[ $(kubectl -n cert-manager get pods -l app=webhook -o 'jsonpath={..status.conditions[?(@.type=="Ready")].status}') != "True" ]]; do echo "waiting for pod" && sleep 1; done

# setup the auth-proxy clusterissuer
kubectl apply -f /etc/kubernetes/auth-proxy-ca.yaml
