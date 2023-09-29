## Logging Route

> Warning: Experimental feature

A _Logging Route_ is responsible to define a global rule that instructs `FluentbitAgent` resources that belongs to the same
`Logging` resource to route logs to different target `Logging` aggregators (fluentd or syslog-ng).

The routed logs are filtered based on the `watchNamespaces` and `watchNamespaceSelector` fields of the target `Logging` resources,
which were originally used to limit which _flow_ and _output_ resources are process by the `Logging` resource.

This also means, that the logs routed by _clusterflows_ will be limited to the above namespace list, as the aggregator
will not have any other log in its possession.

### Spec


Example:
```
apiVersion: logging.banzaicloud.io/v1beta1
kind: LoggingRoute
metadata:
  name: tenants
spec:
  source: ops
  targets:
    matchExpressions:
    - key: tenant
      operator: Exists
```

The above logging route configuration means that the `FluentbitAgent` resource in the `ops` _logging_ will route logs
to _logging_ aggregators that has the `tenant` label set.

### Status

The status of the `LoggingRoute` resource is populated with the targets and their namespaces. In case there is an issue
the `problems` field highlights issues that blocks a tenant from receiving any messages, while notices are only informational
messages.

### Example with Logging resources and status

Tenants used by different development teams, where only the teams' own logs should be available. Let's suppose every team has an
ops and an app namespace with tenant labels:
```
apiVersion: logging.banzaicloud.io/v1beta1
kind: Logging
metadata:
  labels:
    tenant: team-a
  name: team-a
spec:
  controlNamespace: team-a-ops
  fluentd: {}
  watchNamespaceSelector:
    matchLabels:
      tenant: team-a
---
apiVersion: logging.banzaicloud.io/v1beta1
kind: Logging
metadata:
  labels:
    tenant: team-b
  name: team-b
spec:
  controlNamespace: team-b-ops
  syslogNG: {}
  watchNamespaceSelector:
    matchLabels:
      tenant: team-b
```
> Note: 
> - these logging resources do not have a corresponding `FluentbitAgent` resource defined as log collection
will be handled by the ops tenant. 
> - the usage of `loggingRef` is not required here.

Ops tenant where all logs should be available:
```
apiVersion: logging.banzaicloud.io/v1beta1
kind: Logging
metadata:
  labels:
    tenant: ops
  name: ops
spec:
  controlNamespace: ops
  fluentd: {}
  loggingRef: ops
---
apiVersion: logging.banzaicloud.io/v1beta1
kind: FluentbitAgent
metadata:
  name: ops
spec:
  loggingRef: ops
  forwardOptions:
    Workers: 0
  syslogng_output:
    Workers: 0
```
> Note: 
> - `loggingRef`s are required here
> - `Workers: 0` is a workaround so that the processing of all tenants (outputs) don't block if one or more tenant is unavailable

And finally the logging route with a populated status:
```
apiVersion: logging.banzaicloud.io/v1beta1
kind: LoggingRoute
metadata:
  name: tenants
spec:
  source: ops
  targets:
    matchExpressions:
    - key: tenant
      operator: Exists
status:
  notices:
  - tenant ops receives logs from ALL namespaces
  noticesCount: 1
  tenants:
  - name: team-a
    namespaces:
    - team-a-ops
    - team-a-app
  - name: b
    namespaces:
    - team-b-ops
    - team-b-app
  - name: ops
```

> Note: there is a notice that the ops tenant receives logs for all namespaces, which is exactly what we
> want here, but for a team or a customer level tenant it can easily be a misconfiguration issue, hence the notice.
