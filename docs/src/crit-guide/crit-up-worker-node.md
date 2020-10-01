Running `crit up` with a [worker configuration](bootstrapping-a-worker.md
) perfoms the following steps:
 
| Step      | Description 
| ----------- | ----------- 
|WorkerPreCheck | Validate configuration |
|StopKubelet |  Stop the kubelet using systemd |
|WriteBootstrapKubeletConfig | Write kubelet boostrap kubeconfig |
|WriteKubeletConfigs | Write kubelet settings |
|StartKubelet | Start Kubelet using systemd |

