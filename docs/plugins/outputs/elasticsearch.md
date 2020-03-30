---
title: Elasticsearch
weight: 200
---

# Elasticsearch output plugin for Fluentd
## Overview
More info at https://github.com/uken/fluent-plugin-elasticsearch
>Example Deployment: [Save all logs to ElasticSearch](../../../docs/example-es.md)

 #### Example output configurations
 ```
 spec:
   elasticsearch:
     host: elasticsearch-elasticsearch-cluster.default.svc.cluster.local
     port: 9200
     scheme: https
     ssl_verify: false
     ssl_version: TLSv1_2
     buffer:
       timekey: 1m
       timekey_wait: 30s
       timekey_use_utc: true
 ```

## Configuration
### Elasticsearch
#### Send your logs to Elasticsearch

| Variable Name | Type | Required | Default | Description |
|---|---|---|---|---|
| host | string | No | localhost | You can specify Elasticsearch host by this parameter. <br> |
| port | int | No |  9200 | You can specify Elasticsearch port by this parameter.<br> |
| hosts | string | No | - | You can specify multiple Elasticsearch hosts with separator ",". If you specify hosts option, host and port options are ignored.<br> |
| user | string | No | - | User for HTTP Basic authentication. This plugin will escape required URL encoded characters within %{} placeholders. e.g. %{demo+}<br> |
| password | *secret.Secret | No | - | Password for HTTP Basic authentication.<br>[Secret](./secret.md)<br> |
| path | string | No | - | Path for HTTP Basic authentication.<br> |
| scheme | string | No |  http | Connection scheme <br> |
| ssl_verify | *bool | No | true | Skip ssl verification (default: true)<br> |
| ssl_version | string | No | - | If you want to configure SSL/TLS version, you can specify ssl_version parameter. [SSLv23, TLSv1, TLSv1_1, TLSv1_2]<br> |
| logstash_format | bool | No |  false | Enable Logstash log format.<br> |
| include_timestamp | bool | No |  false | Adds a @timestamp field to the log, following all settings logstash_format does, except without the restrictions on index_name. This allows one to log to an alias in Elasticsearch and utilize the rollover API.<br> |
| logstash_prefix | string | No |  logstash | Set the Logstash prefix.<br> |
| logstash_prefix_separator | string | No |  - | Set the Logstash prefix separator.<br> |
| logstash_dateformat | string | No |  %Y.%m.%d | Set the Logstash date format.<br> |
| type_name | string | No |  fluentd | Set the index type for elasticsearch. This is the fallback if `target_type_key` is missing. <br> |
| pipeline | string | No | - | This param is to set a pipeline id of your elasticsearch to be added into the request, you can configure ingest node.<br> |
| time_key_format | string | No | - | The format of the time stamp field (@timestamp or what you specify with time_key). This parameter only has an effect when logstash_format is true as it only affects the name of the index we write to.<br> |
| time_precision | string | No | - | Should the record not include a time_key, define the degree of sub-second time precision to preserve from the time portion of the routed event.<br> |
| time_key | string | No | - | By default, when inserting records in Logstash format, @timestamp is dynamically created with the time at log ingestion. If you'd like to use a custom time, include an @timestamp with your record.<br> |
| utc_index | *bool | No | true | By default, the records inserted into index logstash-YYMMDD with UTC (Coordinated Universal Time). This option allows to use local time if you describe utc_index to false.(default: true)<br> |
| target_index_key | string | No | - | Tell this plugin to find the index name to write to in the record under this key in preference to other mechanisms. Key can be specified as path to nested record using dot ('.') as a separator. https://github.com/uken/fluent-plugin-elasticsearch#target_index_key<br> |
| target_type_key | string | No |  fluentd | Similar to target_index_key config, find the type name to write to in the record under this key (or nested record). If key not found in record - fallback to type_name.<br> |
| template_name | string | No | - | The name of the template to define. If a template by the name given is already present, it will be left unchanged, unless template_overwrite is set, in which case the template will be updated.<br> |
| template_file | string | No | - | The path to the file containing the template to install.<br> |
| templates | string | No | - | Specify index templates in form of hash. Can contain multiple templates.<br> |
| customize_template | string | No | - | Specify the string and its value to be replaced in form of hash. Can contain multiple key value pair that would be replaced in the specified template_file. This setting only creates template and to add rollover index please check the rollover_index configuration.<br> |
| rollover_index | bool | No |  false | Specify this as true when an index with rollover capability needs to be created. https://github.com/uken/fluent-plugin-elasticsearch#rollover_index<br> |
| index_date_pattern | string | No |  now/d | Specify this to override the index date pattern for creating a rollover index.<br> |
| deflector_alias | string | No | - | Specify the deflector alias which would be assigned to the rollover index created. This is useful in case of using the Elasticsearch rollover API<br> |
| index_prefix | string | No |  logstash | Specify the index prefix for the rollover index to be created.<br> |
| application_name | string | No |  default | Specify the application name for the rollover index to be created.<br> |
| template_overwrite | bool | No |  false | Always update the template, even if it already exists.<br> |
| max_retry_putting_template | string | No |  10 | You can specify times of retry putting template.<br> |
| fail_on_putting_template_retry_exceed | *bool | No | true | Indicates whether to fail when max_retry_putting_template is exceeded. If you have multiple output plugin, you could use this property to do not fail on fluentd statup.(default: true)<br> |
| max_retry_get_es_version | string | No |  15 | You can specify times of retry obtaining Elasticsearch version.<br> |
| request_timeout | string | No |  5s | You can specify HTTP request timeout.<br> |
| reload_connections | *bool | No | true | You can tune how the elasticsearch-transport host reloading feature works.(default: true)<br> |
| reload_on_failure | bool | No |  false | Indicates that the elasticsearch-transport will try to reload the nodes addresses if there is a failure while making the request, this can be useful to quickly remove a dead node from the list of addresses.<br> |
| resurrect_after | string | No |  60s | You can set in the elasticsearch-transport how often dead connections from the elasticsearch-transport's pool will be resurrected.<br> |
| include_tag_key | bool | No |  false | This will add the Fluentd tag in the JSON record.<br> |
| tag_key | string | No |  tag | This will add the Fluentd tag in the JSON record.<br> |
| id_key | string | No | - | https://github.com/uken/fluent-plugin-elasticsearch#id_key<br> |
| routing_key | string | No | - | Similar to parent_key config, will add _routing into elasticsearch command if routing_key is set and the field does exist in input event.<br> |
| remove_keys_on_update | string | No | - | Remove keys on update will not update the configured keys in elasticsearch when a record is being updated. This setting only has any effect if the write operation is update or upsert.<br> |
| remove_keys_on_update_key | string | No | - | This setting allows remove_keys_on_update to be configured with a key in each record, in much the same way as target_index_key works.<br> |
| retry_tag | string | No | - | This setting allows custom routing of messages in response to bulk request failures. The default behavior is to emit failed records using the same tag that was provided.<br> |
| write_operation | string | No |  index | The write_operation can be any of: (index,create,update,upsert)<br> |
| reconnect_on_error | bool | No |  false | Indicates that the plugin should reset connection on any error (reconnect on next send). By default it will reconnect only on "host unreachable exceptions". We recommended to set this true in the presence of elasticsearch shield.<br> |
| with_transporter_log | bool | No |  false | This is debugging purpose option to enable to obtain transporter layer log. <br> |
| content_type | string | No |  application/json | With content_type application/x-ndjson, elasticsearch plugin adds application/x-ndjson as Content-Type in payload. <br> |
| include_index_in_url | bool | No | - | With this option set to true, Fluentd manifests the index name in the request URL (rather than in the request body). You can use this option to enforce an URL-based access control.<br> |
| time_parse_error_tag | string | No | - | With logstash_format true, elasticsearch plugin parses timestamp field for generating index name. If the record has invalid timestamp value, this plugin emits an error event to @ERROR label with time_parse_error_tag configured tag.<br> |
| http_backend | string | No |  excon | With http_backend typhoeus, elasticsearch plugin uses typhoeus faraday http backend. Typhoeus can handle HTTP keepalive. <br> |
| prefer_oj_serializer | bool | No |  false | With default behavior, Elasticsearch client uses Yajl as JSON encoder/decoder. Oj is the alternative high performance JSON encoder/decoder. When this parameter sets as true, Elasticsearch client uses Oj as JSON encoder/decoder. <br> |
| flatten_hashes | bool | No | - | Elasticsearch will complain if you send object and concrete values to the same field. For example, you might have logs that look this, from different places:<br>{"people" => 100} {"people" => {"some" => "thing"}}<br>The second log line will be rejected by the Elasticsearch parser because objects and concrete values can't live in the same field. To combat this, you can enable hash flattening.<br> |
| flatten_hashes_separator | string | No | - | Flatten separator<br> |
| validate_client_version | bool | No |  false | When you use mismatched Elasticsearch server and client libraries, fluent-plugin-elasticsearch cannot send data into Elasticsearch. <br> |
| unrecoverable_error_types | string | No | - | Default unrecoverable_error_types parameter is set up strictly. Because es_rejected_execution_exception is caused by exceeding Elasticsearch's thread pool capacity. Advanced users can increase its capacity, but normal users should follow default behavior.<br>If you want to increase it and forcibly retrying bulk request, please consider to change unrecoverable_error_types parameter from default value.<br>Change default value of thread_pool.bulk.queue_size in elasticsearch.yml)<br> |
| verify_es_version_at_startup | *bool | No | true | Because Elasticsearch plugin should change behavior each of Elasticsearch major versions.<br>For example, Elasticsearch 6 starts to prohibit multiple type_names in one index, and Elasticsearch 7 will handle only _doc type_name in index.<br>If you want to disable to verify Elasticsearch version at start up, set it as false.<br>When using the following configuration, ES plugin intends to communicate into Elasticsearch 6. (default: true)<br> |
| default_elasticsearch_version | string | No |  5 | This parameter changes that ES plugin assumes default Elasticsearch version.<br> |
| custom_headers | string | No |  {} | This parameter adds additional headers to request. Example: {"token":"secret"} <br> |
| log_es_400_reason | bool | No |  false | By default, the error logger won't record the reason for a 400 error from the Elasticsearch API unless you set log_level to debug. However, this results in a lot of log spam, which isn't desirable if all you want is the 400 error reasons. You can set this true to capture the 400 error reasons without all the other debug logs. <br> |
| suppress_doc_wrap | bool | No |  false | By default, record body is wrapped by 'doc'. This behavior can not handle update script requests. You can set this to suppress doc wrapping and allow record body to be untouched. <br> |
| ignore_exceptions | string | No | - | A list of exception that will be ignored - when the exception occurs the chunk will be discarded and the buffer retry mechanism won't be called. It is possible also to specify classes at higher level in the hierarchy. For example<br>`ignore_exceptions ["Elasticsearch::Transport::Transport::ServerError"]`<br>will match all subclasses of ServerError - Elasticsearch::Transport::Transport::Errors::BadRequest, Elasticsearch::Transport::Transport::Errors::ServiceUnavailable, etc.<br> |
| exception_backup | *bool | No | true | Indicates whether to backup chunk when ignore exception occurs. (default: true)<br> |
| bulk_message_request_threshold | string | No |  20MB | Configure bulk_message request splitting threshold size.<br>Default value is 20MB. (20 * 1024 * 1024)<br>If you specify this size as negative number, bulk_message request splitting feature will be disabled. <br> |
| buffer | *Buffer | No | - | [Buffer](./buffer.md)<br> |
