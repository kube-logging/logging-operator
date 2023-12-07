## Standalone Fluentd config

The standalone `FluentdConfig` and `SyslogNGConfig` are namespaced resources that allow the configuration of the Fluentd / SyslogNG
aggregator components in the control namespace separately from the Logging resource.

The primary benefit of this behaviour is that it enables a multi-tenant model, where tenant owners are responsible
for operating their own aggregator, while the Logging resource is in control of the central operations team.
For more information about the multi-tenancy model where the collector is capable of routing logs based on namespaces
to individual aggregators and where aggregators are fully isolated, please see [Multi-tenancy](multi-tenancy.md)

Traditional configuration of fluentd within the logging resource (the case with syslog-ng is the same expect the field is `syslogNG`):
```
apiVersion: logging.banzaicloud.io/v1beta1
kind: Logging
metadata:
  name: example
spec:
  controlNamespace: logging
  fluentd:
    scaling:
      replicas: 2
```

The alternative configuration is as follows:
```
apiVersion: logging.banzaicloud.io/v1beta1
kind: Logging
metadata:
  name: example
spec:
  controlNamespace: logging
---
apiVersion: logging.banzaicloud.io/v1beta1
kind: FluentdConfig
metadata:
  name: example
  namespace: logging
spec:
  scaling:
    replicas: 2
```

Note: In case of syslog-ng the name of the standalone config resource is `SyslogNGConfig`.

### Schema and migration

The schema for `FluentdConfig.spec` / `SyslogNGConfig.spec` is the same as it was withing `Logging.spec.fluentd` / `Logging.spec.syslogNG`,
so the migration should be a trivial lift and shift exercise.

### Restrictions and status

There can only be one active `FluentdConfig` or `SyslogNGConfig` for a single `Logging` resource at a time. The controller will make
sure to register the active resource into the `Logging` resource's status under `fluentdConfigName` / `syslogNGConfigName`,
while registering the `Logging` resource name under `logging` in the `FluentdConfig` / `SyslogNGConfig` resource's status.

```
kubectl get logging example -o jsonpath='{.status}' | jq .
{
  "configCheckResults": {
    "ac2d4553": true
  },
  "fluentdConfigName": "example"
}
```

```
kubectl get fluentdconfig example -o jsonpath='{.status}' | jq .
{
  "active": true,
  "logging": "example"
}
```

If there is a conflict, then the controller will add a problem to both resources so that both operations and tenant users can be aware of it.
For example add another `FluentdConfig` resource on top of the existing one:

```
apiVersion: logging.banzaicloud.io/v1beta1
kind: FluentdConfig
metadata:
  name: example2
  namespace: logging
spec: {}
```

The first `FluentdConfig` should be left intact, while the second one should have the following status:
```
kubectl get fluentdconfig example2 -o jsonpath='{.status}' | jq .
{
  "active": false,
  "problems": [
    "logging already has a detached fluentd configuration, remove excess configuration objects"
  ],
  "problemsCount": 1
}
```

The `Logging` resource will also highlight the issue
```
kubectl get logging example -o jsonpath='{.status}' | jq .
{
  "configCheckResults": {
    "ac2d4553": true
  },
  "fluentdConfigName": "example",
  "problems": [
    "multiple fluentd configurations found, couldn't associate it with logging"
  ],
  "problemsCount": 1
}
```

Once the extra `FluentdConfig` resource is removed the `Logging` resource status should return back to normal:
```
kubectl get logging example -o jsonpath='{.status}' | jq .
{
  "configCheckResults": {
    "ac2d4553": true
  },
  "fluentdConfigName": "example"
}
```
