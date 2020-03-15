# Splunk via Hec output plugin for Fluentd
## Overview
More info at https://github.com/splunk/fluent-plugin-splunk-hec

 #### Example output configurations
 ```
 spec:
   SplunkHec:
     host: splunk.default.svc.cluster.local
     port: 8088
     protocol: http
 ```

## Configuration
### SplunkHecOutput
#### SplunkHecOutput sends your logs to Splunk via Hec

| Variable Name | Type | Required | Default | Description |
|---|---|---|---|---|
| hec_host | string | Yes | - | You can specify SplunkHec host by this parameter.<br> |
| hec_port | int | No | 8088 | The port number for the Hec token or the Hec load balancer. <br> |
| protocol | string | No | https | This is the protocol to use for calling the Hec API. Available values are: http, https. <br> |
| hec_token | *secret.Secret | Yes | - | Identifier for the Hec token.<br>[Secret](./secret.md)<br> |
| metrics_from_event | bool | No | - | When data_type is set to "metric", the ingest API will treat every key-value pair in the input event as a metric name-value pair. Set metrics_from_event to false to disable this behavior and use metric_name_key and metric_value_key to define metrics. (Default:true)<br> |
| metrics_name_key | string | No | - | Field name that contains the metric name. This parameter only works in conjunction with the metrics_from_event paramter. When this prameter is set, the metrics_from_event parameter is automatically set to false.<br> |
| metrics_value_key | string | No | - | Field name that contains the metric value, this parameter is required when metric_name_key is configured.<br> |
