# Logging Operator Syslog-ng quickstart

This document is for the preview version of the Logging Operator with syslog-ng as the forward agent.

## What is new?
There are 4 new CRDs to manage Logging Operator with syslog-ng:
- SyslogNGFlow
- SyslogNGClusterFlow
- SyslogNGOutput
- SyslogNGClusterOutput

These resources works identical to their non syslog-ng pairs.

## Examples

These are examples to demonstrate the differences comapred to the standard custom resources.

### Routing
With syslog-ng there is no restriction on routing. You can define metadata and log content as well. The syntax is slightly different thought.

The following example filters for specific Pod labels
```yaml
  match:
    and:
    - regexp:
        value: json.kubernetes.labels.app.kubernetes.io/instance
        pattern: one-eye-log-generator
        type: string
    - regexp:
        value:  json.kubernetes.labels.app.kubernetes.io/name
        pattern: log-generator
        type: string
```

> Note: You need to use the `json.` prefix before your field naming. Field names are dotted json path. Exampla: {"kubernetes": {"namespace_name": "default"}} represented as `json.kubernetes.namespace_name`.

Example kubernetes message:
```json
{"user":"-","ts":"2022-08-02T07:41:34.090763Z","time":"02/Aug/2022:07:41:34 +0000","stream":"stdout","source":"/var/log/log-generator","size":"1628","remote":"85.151.230.190","referer":"-","path":"/index.html","method":"POST","logtag":"F","kubernetes":{"pod_name":"one-eye-log-generator-57988cbd65-gpgc4","pod_id":"010d4598-2e34-4165-accd-2b77e4fc4bb6","namespace_name":"default","labels":{"pod-template-hash":"57988cbd65","app.kubernetes.io/name":"log-generator","app.kubernetes.io/instance":"one-eye-log-generator"},"host":"ip-192-168-6-51.eu-west-1.compute.internal","docker_id":"89e7cf414b5e6ff1fbee977d62cbc96794d2debd6a52803857d5dbad57d4f772","container_name":"log-generator","container_image":"033498657557.dkr.ecr.us-east-2.amazonaws.com/banzaicloud/log-generator:0.3.20","container_hash":"033498657557.dkr.ecr.us-east-2.amazonaws.com/banzaicloud/log-generator@sha256:b031138718194a17fdac2964bacf9543f96b037a65cd50138a5754ddb7897bb5"},"http_x_forwarded_for":"\"-\"","host":"-","code":"403","cluster":"xxxxx","agent":"Mozilla/5.0 (iPhone; CPU iPhone OS 10_3_3 like Mac OS X) AppleWebKit/603.3.8 (KHTML, like Gecko) Version/9.0 Mobile/13B143 Safari/601.1"}
```

Options for match statement content as well. The syntax is slightly different thought.

```yaml
  match:
    and: ... // Logical AND between expressions
    or: ... // // Logical OR between expressions
    regexp: ... // Regular expression for matching
    not: ... // Logical NOT before an expression
```

```yaml
    pattern // expression or string value to compare to
    template // Optional template expression to evaluate against
    value // Specify a field name of the record to match against the value of.
    flags // https://www.syslog-ng.com/technical-documents/doc/syslog-ng-open-source-edition/3.37/administration-guide/81#TOPIC-1829224
    type // Match type can be string for exact match or leave empty fo Regexp  https://www.syslog-ng.com/technical-documents/doc/syslog-ng-open-source-edition/3.37/administration-guide/81#TOPIC-1829223
```

## Filters

## Outputs