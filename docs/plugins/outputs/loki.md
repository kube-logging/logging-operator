### Loki
#### Fluentd output plugin to ship logs to a Loki server.

| Variable Name | Type | Required | Default | Description |
|---|---|---|---|---|
| url | string | No | https://logs-us-west1.grafana.net | The url of the Loki server to send logs to. <br> |
| username | *secret.Secret | No | - | Specify a username if the Loki server requires authentication.<br>[Secret](./secret.md)<br> |
| password | *secret.Secret | No | - | Specify password if the Loki server requires authentication.<br>[Secret](./secret.md)<br> |
| tenant | string | No | - | Loki is a multi-tenant log storage platform and all requests sent must include a tenant.<br> |
| extra_labels | bool | No |  nil | Set of labels to include with every Loki stream.<br> |
| buffer | *Buffer | No | - | [Buffer](./buffer.md)<br> |
