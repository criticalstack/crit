# Enabling Pod Security Policies

## What is a Pod Security Policy

Pod Security Policies are in-cluster Kubernetes resources that provides ways of securing pods. The official [Pod Security Policy](https://kubernetes.io/docs/concepts/policy/pod-security-policy/) of the official Kubernetes docs provides a great deal of helpful information and a walkthrough of how to use them, and is highly recommended reading. For the purposes of this documentation, we really just want to focus on getting them running on your Crit cluster.

## Configuration

The APIServer has quite a few [admission plugins enabled by default](https://kubernetes.io/docs/reference/access-authn-authz/admission-controllers/#which-plugins-are-enabled-by-default), however, the `PodSecurityPolicy` plugin must be enabled when configuring the APIServer with the `enable-admission-plugin` option:

```yaml
apiVersion: crit.sh/v1alpha2
kind: ControlPlaneConfiguration
kubeAPIServer:
  extraArgs:
    enable-admission-plugins: PodSecurityPolicy
```

`enable-admission-plugin` can be provided a comma-delimited list of admission plugins to enable. While the order that admission plugins run does matter, it does not matter for this particular option as it simply enables the plugin.

The admission plugin `SecurityContextDeny` must _**NOT**_ be enabled along with `PodSecurityPolicy`. In the case that `PodSecurityPolicy` is enabled, the usage completely [supplants the functionality provided by `SecurityContextDeny`](https://github.com/kubernetes/kubernetes/issues/53797#issuecomment-336153103).
