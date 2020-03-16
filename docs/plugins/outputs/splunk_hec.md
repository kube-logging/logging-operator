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
| data_type | string | No |  event | The type of data that will be sent to Sumo Logic, either event or metric <br> |
| hec_host | string | Yes | - | You can specify SplunkHec host by this parameter.<br> |
| hec_port | int | No |  8088 | The port number for the Hec token or the Hec load balancer. <br> |
| protocol | string | No |  https | This is the protocol to use for calling the Hec API. Available values are: http, https. <br> |
| hec_token | *secret.Secret | Yes | - | Identifier for the Hec token.<br>[Secret](./secret.md)<br> |
| metrics_from_event | bool | No | - | When data_type is set to "metric", the ingest API will treat every key-value pair in the input event as a metric name-value pair. Set metrics_from_event to false to disable this behavior and use metric_name_key and metric_value_key to define metrics. (Default:true)<br> |
| metrics_name_key | string | No |  true | Field name that contains the metric name. This parameter only works in conjunction with the metrics_from_event paramter. When this prameter is set, the metrics_from_event parameter is automatically set to false. <br> |
| metrics_value_key | string | No | - | Field name that contains the metric value, this parameter is required when metric_name_key is configured.<br> |
| coerce_to_utf8 | bool | No |  true | Indicates whether to allow non-UTF-8 characters in user logs. If set to true, any non-UTF-8 character is replaced by the string specified in non_utf8_replacement_string. If set to false, the Ingest API errors out any non-UTF-8 characters. .<br> |
| non_utf8_replacement_string | string | No |  ' ' | If coerce_to_utf8 is set to true, any non-UTF-8 character is replaced by the string you specify in this parameter. .<br> |
| index | string | No | - | Identifier for the Splunk index to be used for indexing events. If this parameter is not set, the indexer is chosen by HEC. Cannot set both index and index_key parameters at the same time.<br> |
| index_key | string | No | - | The field name that contains the Splunk index name. Cannot set both index and index_key parameters at the same time.<br> |
| host | string | No | - | The host location for events. Cannot set both host and host_key parameters at the same time. (Default:hostname)<br> |
| host_key | string | No | - | Key for the host location. Cannot set both host and host_key parameters at the same time.<br> |
| source | string | No | - | The source field for events. If this parameter is not set, the source will be decided by HEC. Cannot set both source and source_key parameters at the same time.<br> |
| source_key | string | No | - | Field name to contain source. Cannot set both source and source_key parameters at the same time.<br> |
| sourcetype | string | No | - | The sourcetype field for events. When not set, the sourcetype is decided by HEC. Cannot set both source and source_key parameters at the same time.<br> |
| sourcetype_key | string | No | - | Field name that contains the sourcetype. Cannot set both source and source_key parameters at the same time.<br> |
| keep_keys | bool | No | - | By default, all the fields used by the *_key parameters are removed from the original input events. To change this behavior, set this parameter to true. This parameter is set to false by default. When set to true, all fields defined in index_key, host_key, source_key, sourcetype_key, metric_name_key, and metric_value_key are saved in the original event.<br> |
| fields | map[string]string | No | - | In this case, parameters inside <fields> are used as indexed fields and removed from the original input events<br> |
| format | *Format | No | - | [Format](./format.md)<br> |
| buffer | *Buffer | No | - | [Buffer](./buffer.md)<br> |
