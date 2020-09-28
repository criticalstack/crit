# Configuring Authentication

## Configure the Kubernetes API Server

The Kubernetes API server can be configured with [OpenID Connect](https://openid.net/connect/) to use an existing OpenID Identity Provider. It can only trust a single issuer and until the API server can be [configured with component configs](https://github.com/kubernetes/enhancements/blob/master/keps/sig-cluster-lifecycle/wgs/0014-20180707-componentconfig-api-types-to-staging.md#kube-apiserver-changes) it must be specified in the Crit config as command-line arguments:

```yaml
apiVersion: crit.sh/v1alpha2
kind: ControlPlaneConfiguration
kubeAPIServer:
  extraArgs:
    oidc-issuer-url: "https://accounts.google.com"
    oidc-client-id: critical-stack
    oidc-username-claim: email
    oidc-groups-claim: groups
```

The above configuration will allow the API server to use Google as its identity provider, but with some major limitations:

* Kubernetes does not act as a client for the issuer
* Does not provide a way to manage the lifecycle of OpenID Connect tokens

This can be best understood looking in the Kubernetes authentication documentation for [OpenID Connect Tokens](https://kubernetes.io/docs/reference/access-authn-authz/authentication/#openid-connect-tokens). The process of getting a token happens completely outside of the context of the Kubernetes cluster and is passed in as an argument to `kubectl` commands.


## Using an In-cluster Identity Provider

Given the limitations mentioned above, many run their own identity providers inside of the cluster to provide additional auth features to the cluster. This complicates configuration, however, since the API server will either have to be reconfigured and restarted, or will need to be configured with an issuer that is not yet running.

So what if you want to provide a web interface that leverages this authentication? Given the limitations mentioned above, you would have to write authentication logic for the specific upstream identity provider into your application, and should the upstream identity provider change, so does the authentication logic AND the API server configuration. This is where identity providers, such as [Dex](https://github.com/dexidp/dex), come in. Dex uses OpenID Connect to provide authentication for other applications by acting as a shim between the client app and the upstream provider. When using Dex, the `oidc-issuer-url` argument being specified needs to target the expected address of Dex running the cluster, so something like:

```yaml
oidc-issuer-url: "https://dex.kube-system.svc.cluster.local:5556"
```

It is ok that Dex isn't running yet, the API server will function as normal until the issuer is available.

### The auth-proxy CA

The API server uses the host's root CAs by default, but in the case where an application might not be using a CA signed certificate, like during development or automated testing, Crit generates an additional CA that is already available in the API server certs volume. This helps with the chicken/egg problem of needing to specify a CA file when bootstrapping a new cluster before the application has been deployed. To use this auth-proxy CA, just add this to the API server configuration:

```yaml
oidc-ca-file: /etc/kubernetes/pki/auth-proxy-ca.crt
```

Please note that this assumes that the default Kubernetes directory (`/etc/kubernetes`) is being used. From here there are many options to make use of auth-proxy CA. For example, [cert-manager](https://cert-manager.io/) can be installed and the auth-proxy CA can be setup as a `ClusterIssuer`:

```sh
# install cert-manager
kubectl create namespace cert-manager
kubectl apply --validate=false -f https://github.com/jetstack/cert-manager/releases/download/v0.14.0/cert-manager.yaml

# add auth-proxy-ca secret to be used as ClusterIssuer
kubectl -n cert-manager create secret generic auth-proxy-ca --from-file=tls.crt=/etc/kubernetes/pki/auth-proxy-ca.crt --from-file=tls.key=/etc/kubernetes/pki/auth-proxy-ca.key

# wait for cert-manager-webhook readiness
while [[ $(kubectl -n cert-manager get pods -l app=webhook -o 'jsonpath={..status.conditions[?(@.type=="Ready")].status}') != "True" ]]; do echo "waiting for pod" && sleep 1; done

kubectl apply -f - <<EOT
apiVersion: cert-manager.io/v1alpha2
kind: ClusterIssuer
metadata:
  name: auth-proxy-ca
  namespace: cert-manager
spec:
  ca:
    secretName: auth-proxy-ca
EOT
```

Then applications can create cert-manager certificates for their application to use:

```yaml
apiVersion: cert-manager.io/v1alpha2
kind: Certificate
metadata:
  name: myapp-example
spec:
  secretName: myapp-certs
  duration: 8760h # 365d
  renewBefore: 360h # 15d
  organization:
  -  Internet Widgits Pty Ltd
  isCA: false
  keySize: 2048
  keyAlgorithm: rsa
  keyEncoding: pkcs1
  usages:
    - server auth
    - client auth
  dnsNames:
  - myapp.example.com
  issuerRef:
    name: auth-proxy-ca
    kind: ClusterIssuer
```

Of course, this is just one possible way to approach authentication, and configuration will vary greatly depending upon the needs of the application(s) running on the cluster.
