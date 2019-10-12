# Loki output plugin 
## Overview
Fluentd output plugin to ship logs to a Loki server.
More info at https://github.com/banzaicloud/fluent-plugin-kubernetes-loki
>Example: [Store Nginx Access Logs in Grafana Loki with Logging Operator](../../../docs/example-loki-nginx.md)

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
| extra_labels | bool | No |  nil | Set of labels to include with every Loki stream.<br> |
| buffer | *Buffer | No | - | [Buffer](./buffer.md)<br> |
