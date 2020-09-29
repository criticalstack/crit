# Managing Certificates


## Check Certificate Expiration

You can use the `crit certs list` command to check when certificates expire: 

```sh 
$ crit certs list
Certificate Authorities:
========================
Name		CN		Expires	NotAfter
ca		kubernetes	9y	2030-09-27T01:45:12Z
front-proxy-ca	front-proxy-ca	9y	2030-09-27T16:36:08Z

Certificates:
=============
Name				CN				Expires	NotAfter
apiserver			kube-apiserver			364d	2021-09-29T23:54:16Z
apiserver-kubelet-client	kube-apiserver-kubelet-client	364d	2021-09-29T23:54:16Z
apiserver-healthcheck-client	system:basic-info-viewer	364d	2021-09-29T23:54:16Z
front-proxy-client		front-proxy-client		364d	2021-09-29T23:54:17Z
```


## Rotating Certificates

There are several different solutions pertaining to certificate rotation. The appropriate solution greatly depends on an organization's use case. Some things to consider:
 
* Does certificate rotation need to intergrate with an organization's existing certificate infrastructure?
* Can certificate approval and signing be automated, or does it require a cluster administrator?
* How often do certificates need to be rotated?
* How many clusters need to be supported?


### Rotating with Crit

Certificates can be renewed with [`crit certs renew`](/crit-commands/crit-certs-renew.md). Note, this does not renew the CA.  

### Rotating with the Kubernetes certificates API

Kubernetes provides a [Certificate API](https://kubernetes.io/docs/tasks/tls/managing-tls-in-a-cluster/) that can be used to provision certificates using [certificate signing requests](https://kubernetes.io/docs/reference/access-authn-authz/certificate-signing-requests/). 

#### Kubelet Certificate

The kubelet certificate [can be automatically renewed](https://kubernetes.io/docs/tasks/tls/certificate-rotation/) using the kubernetes api.

#### Advanced Certificate Rotation

Organizations that require an automated certificate rotation solution that integrates with existing certificate infrastructure should consider projects like [cert-manager](https://cert-manager.io/docs/installation/kubernetes/). 

