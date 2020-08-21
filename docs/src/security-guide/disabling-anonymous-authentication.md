# Disabling Anonymous Authentication

The API server defaults to allow anonymous auth, meaning that incoming requests that are not authenticated will be implicitly given a username `system:anonymous` and be part of the `system:unauthenticated` group. While this user may not have permission to anything, problems related to allowing anonymous authentication are still possible, such as vulnerabilities like the ["Billion Laughs" attack](https://github.com/kubernetes/kubernetes/issues/83253).

Disabling anonymous authentication only requires passing an argument to the API server:

```yaml
apiVersion: crit.sh/v1alpha2
kind: ControlPlaneConfiguration
kubeAPIServer:
  extraArgs:
    anonymous-auth: false
```

## API Server Healthchecks

[Liveness probes](https://kubernetes.io/docs/tasks/configure-pod-container/configure-liveness-readiness-startup-probes/) will fail for static pods should anonymous-auth be set to false. Crit addresses this by detecting when `--anonymous-auth` has been disabled and adds a special healthcheck-proxy sidecar to the apiserver static pod. It acts as a reverse proxy with the frontend effectively accepting anonymous traffic and the backend using an authenticated user. The backend connection is established with the built-in `system:basic-info-viewer` user to limit the auth to only being able to look at health and version information.
