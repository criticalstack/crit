# Audit Policy Logging

```yaml
apiVersion: crit.sh/v1alpha2
kind: ControlPlaneConfiguration
kubeAPIServer:
  extraArgs:
    audit-policy-file: "/etc/kubernetes/audit-policy.yaml"
    audit-log-path: "/var/log/kubernetes/kube-apiserver-audit.log"
    audit-log-maxage: "30"
    audit-log-maxbackup: "10"
    audit-log-maxsize: "100"
  extraVolumes:
  - name: apiserver-logs
    hostPath: /var/log/kubernetes
    mountPath: /var/log/kubernetes
    readOnly: false
    hostPathType: directory
  - name: apiserver-audit-config
    hostPath: /etc/kubernetes/audit-policy.yaml
    mountPath: /etc/kubernetes/audit-policy.yaml
    readOnly: true
```
