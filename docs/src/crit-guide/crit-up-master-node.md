
Running `crit up` with a [control plane configuration](configuring-control-plane-components.md
) perfoms the following steps:

| Step      | Description 
| ----------- | ------------------------------------------------------------------------------------ 
|ControlPlanePreCheck | Validate configuration
|CreateOrDownloadCerts | Generate CAs; if already present, don't overwrite
|CreateNodeCerts | Generate certificates for kubernetes components; if already present, dont overwrite 
|StopKubelet | Stop the kubelet using systemd
|WriteKubeConfigs | Generate control plane kubeconfigs and the admin kubeconfig
|WriteKubeletConfigs | Write kubelet settings
|StartKubelet |        Start Kubelet using systemd
|WriteKubeManifests | Write static pod manifests for control plane
|WaitClusterAvailable | Wait for the control plane to be available
|WriteBootstrapServerManifest [optional] | Write the crit boostrap server pod manifest
|DeployCoreDNS | Deploy CoreDNS after cluster is available 
|DeployKubeProxy | Deploy KubeProxy 
|EnableCSRApprover | Add RBAC to allow csrapprover to boostrap nodes 
|MarkControlPlane | Add taint to control plane node
|UploadInfo | Upload crit config map that holds info regarding the cluster




