# Installing a Storage Driver

```sh
helm repo add criticalstack https://charts.cscr.io/criticalstack
kubectl create namespace local-path-storage
helm install local-path-storage criticalstack/local-path-provisioner \
	--namespace local-path-storage \
	--set nameOverride=local-path-storage \
	--set storageClass.defaultClass=true
```

## Install the AWS CSI driver via helm

https://github.com/kubernetes-sigs/aws-ebs-csi-driver

```sh
helm repo add criticalstack https://charts.cscr.io/criticalstack
helm install aws-ebs-csi-driver criticalstack/aws-ebs-csi-driver \
	--set enableVolumeScheduling=true \
	--set enableVolumeResizing=true \
	--set enableVolumeSnapshot=true \
	--version 0.3.0
```


## Setting a Default StorageClass

```yaml
kubectl apply -f - <<EOT
kind: StorageClass
apiVersion: storage.k8s.io/v1
metadata:
  name: ebs-sc
  annotations:
    storageclass.kubernetes.io/is-default-class: "true"
provisioner: ebs.csi.aws.com
volumeBindingMode: WaitForFirstConsumer
parameters:
  csi.storage.k8s.io/fstype: xfs
  type: io1
  iopsPerGB: "50"
  encrypted: "true"
EOT
```
