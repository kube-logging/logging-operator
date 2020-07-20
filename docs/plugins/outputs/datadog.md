---
title: Datadog
weight: 200
generated_file: true
---

# Datadog output plugin for Fluentd
## Overview
It mainly contains a proper JSON formatter and a socket handler that streams logs directly to Datadog - so no need to use a log shipper if you don't wan't to.
More info at https://github.com/DataDog/fluent-plugin-datadog

## Configuration
### Output Config
| Variable Name | Type | Required | Default | Description |
|---|---|---|---|---|
| api_key | *secret.Secret | Yes |  nil | This parameter is required in order to authenticate your fluent agent. <br>+docLink:"Secret,../secret/"<br> |
| use_json | bool | No |  true | Event format, if true, the event is sent in json format. Othwerwise, in plain text.  <br> |
| include_tag_key | bool | No |  false | Automatically include the Fluentd tag in the record.  <br> |
| tag_key | string | No |  "tag" | Where to store the Fluentd tag. <br> |
| timestamp_key | string | No |  "@timestamp" | Name of the attribute which will contain timestamp of the log event. If nil, timestamp attribute is not added. <br> |
| use_ssl | bool | No |  true | If true, the agent initializes a secure connection to Datadog. In clear TCP otherwise.  <br> |
| no_ssl_validation | bool | No |  false | Disable SSL validation (useful for proxy forwarding)  <br> |
| ssl_port | string | No |  "443" | Port used to send logs over a SSL encrypted connection to Datadog. If use_http is disabled, use 10516 for the US region and 443 for the EU region. <br> |
| max_retries | string | No |  "-1" | The number of retries before the output plugin stops. Set to -1 for unlimited retries <br> |
| max_backoff | string | No |  "30" | The maximum time waited between each retry in seconds <br> |
| use_http | bool | No |  true | Enable HTTP forwarding. If you disable it, make sure to change the port to 10514 or ssl_port to 10516  <br> |
| use_compression | bool | No |  true | Enable log compression for HTTP  <br> |
| compression_level | string | No |  "6" | Set the log compression level for HTTP (1 to 9, 9 being the best ratio) <br> |
| dd_source | string | No |  nil | This tells Datadog what integration it is <br> |
| dd_sourcecategory | string | No |  nil | Multiple value attribute. Can be used to refine the source attribute <br> |
| dd_tags | string | No |  nil | Custom tags with the following format "key1:value1, key2:value2" <br> |
| dd_hostname | string | No |  "hostname -f" | Used by Datadog to identify the host submitting the logs. <br> |
| service | string | No |  nil | Used by Datadog to correlate between logs, traces and metrics. <br> |
| port | string | No |  "80" | Proxy port when logs are not directly forwarded to Datadog and ssl is not used <br> |
| host | string | No |  "http-intake.logs.datadoghq.com" | Proxy endpoint when logs are not directly forwarded to Datadog	 <br> |
| buffer | *Buffer | No | - | [Buffer](../buffer/)<br> |
