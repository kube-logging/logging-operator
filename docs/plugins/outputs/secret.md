# Secret definition

## Define secret

Secrets can be used in logging-operator `Output` definitions.

> Secrets *MUST* be in the *SAME* namespace as the `Output` or `ClusterOutput` custom resource

**Example secret definition**
```yaml
aws_key_id:
  valueFrom:
    secretKeyRef:
      name: <kubernetes-secret-name>
      key: <kubernetes-secret-key>
```

For debug purposes you can define secret values directly. However this is *NOT* recommended in production.
```yaml
aws_key_id:
  value: "secretvalue"
```