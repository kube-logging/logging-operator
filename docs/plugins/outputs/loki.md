# Loki output plugin 
## Overview
Fluentd output plugin to ship logs to a Loki server.
More info at https://github.com/banzaicloud/fluent-plugin-kubernetes-loki
>Example: [Store Nginx Access Logs in Grafana Loki with Logging Operator](../../examples/example-loki-nginx.md)

 #### Example output configurations
 ```
 spec:
   loki:
     url: http://loki:3100
     buffer:
       timekey: 1m
       timekey_wait: 30s
       timekey_use_utc: true
 ```

## Configuration
### Output Config
| Variable Name | Type | Required | Default | Description |
|---|---|---|---|---|
| url | string | No | https://logs-us-west1.grafana.net | The url of the Loki server to send logs to. <br> |
| username | *secret.Secret | No | - | Specify a username if the Loki server requires authentication.<br>[Secret](./secret.md)<br> |
| password | *secret.Secret | No | - | Specify password if the Loki server requires authentication.<br>[Secret](./secret.md)<br> |
| tenant | string | No | - | Loki is a multi-tenant log storage platform and all requests sent must include a tenant.<br> |
| labels | Label | No | - | Set of labels to include with every Loki stream.<br> |
| extra_labels | map[string]string | No | - | Set of extra labels to include with every Loki stream.<br> |
| line_format | string | No | json | Format to use when flattening the record to a log line: json, key_value (default: key_value)<br> |
| extract_kubernetes_labels | bool | No |  false | Extract kubernetes labels as loki labels <br> |
| remove_keys | []string | No |  [] | Comma separated list of needless record keys to remove <br> |
| drop_single_key | bool | No |  false | If a record only has 1 key, then just set the log line to the value and discard the key. <br> |
| configure_kubernetes_labels | bool | No | - | Configure Kubernetes metadata in a Prometheus like format<br> |
| buffer | *Buffer | No | - | [Buffer](./buffer.md)<br> |
