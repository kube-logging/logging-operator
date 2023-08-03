---
title: Elasticsearch
weight: 200
generated_file: true
---

# Elasticsearch output plugin for Fluentd
## Overview
 More info at https://github.com/uken/fluent-plugin-elasticsearch
 >Example Deployment: [Save all logs to ElasticSearch](../../../../quickstarts/es-nginx/)

 ## Example output configurations
 ```yaml
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
## Elasticsearch

Send your logs to Elasticsearch

### host (string, optional) {#elasticsearch-host}

You can specify Elasticsearch host by this parameter.  

Default: localhost

### port (int, optional) {#elasticsearch-port}

You can specify Elasticsearch port by this parameter. 

Default:  9200

### hosts (string, optional) {#elasticsearch-hosts}

You can specify multiple Elasticsearch hosts with separator ",". If you specify hosts option, host and port options are ignored. 

Default: -

### user (string, optional) {#elasticsearch-user}

User for HTTP Basic authentication. This plugin will escape required URL encoded characters within %{} placeholders. e.g. %{demo+} 

Default: -

### password (*secret.Secret, optional) {#elasticsearch-password}

Password for HTTP Basic authentication. [Secret](../secret/) 

Default: -

### path (string, optional) {#elasticsearch-path}

Path for HTTP Basic authentication. 

Default: -

### scheme (string, optional) {#elasticsearch-scheme}

Connection scheme  

Default:  http

### ssl_verify (*bool, optional) {#elasticsearch-ssl_verify}

Skip ssl verification (default: true) 

Default: true

### ssl_version (string, optional) {#elasticsearch-ssl_version}

If you want to configure SSL/TLS version, you can specify ssl_version parameter. [SSLv23, TLSv1, TLSv1_1, TLSv1_2] 

Default: -

### ssl_max_version (string, optional) {#elasticsearch-ssl_max_version}

Specify min/max SSL/TLS version 

Default: -

### ssl_min_version (string, optional) {#elasticsearch-ssl_min_version}

Default: -

### ca_file (*secret.Secret, optional) {#elasticsearch-ca_file}

CA certificate 

Default: -

### client_cert (*secret.Secret, optional) {#elasticsearch-client_cert}

Client certificate 

Default: -

### client_key (*secret.Secret, optional) {#elasticsearch-client_key}

Client certificate key 

Default: -

### client_key_pass (*secret.Secret, optional) {#elasticsearch-client_key_pass}

Client key password 

Default: -

### logstash_format (bool, optional) {#elasticsearch-logstash_format}

Enable Logstash log format. 

Default:  false

### include_timestamp (bool, optional) {#elasticsearch-include_timestamp}

Adds a @timestamp field to the log, following all settings logstash_format does, except without the restrictions on index_name. This allows one to log to an alias in Elasticsearch and utilize the rollover API. 

Default:  false

### logstash_prefix (string, optional) {#elasticsearch-logstash_prefix}

Set the Logstash prefix. 

Default:  logstash

### logstash_prefix_separator (string, optional) {#elasticsearch-logstash_prefix_separator}

Set the Logstash prefix separator. 

Default:  -

### logstash_dateformat (string, optional) {#elasticsearch-logstash_dateformat}

Set the Logstash date format. 

Default:  %Y.%m.%d

### index_name (string, optional) {#elasticsearch-index_name}

The index name to write events to  

Default:  fluentd

### type_name (string, optional) {#elasticsearch-type_name}

Set the index type for elasticsearch. This is the fallback if `target_type_key` is missing.  

Default:  fluentd

### pipeline (string, optional) {#elasticsearch-pipeline}

This param is to set a pipeline id of your elasticsearch to be added into the request, you can configure ingest node. 

Default: -

### time_key_format (string, optional) {#elasticsearch-time_key_format}

The format of the time stamp field (@timestamp or what you specify with time_key). This parameter only has an effect when logstash_format is true as it only affects the name of the index we write to. 

Default: -

### time_precision (string, optional) {#elasticsearch-time_precision}

Should the record not include a time_key, define the degree of sub-second time precision to preserve from the time portion of the routed event. 

Default: -

### time_key (string, optional) {#elasticsearch-time_key}

By default, when inserting records in Logstash format, @timestamp is dynamically created with the time at log ingestion. If you'd like to use a custom time, include an @timestamp with your record. 

Default: -

### utc_index (*bool, optional) {#elasticsearch-utc_index}

By default, the records inserted into index logstash-YYMMDD with UTC (Coordinated Universal Time). This option allows to use local time if you describe utc_index to false.(default: true) 

Default: true

### suppress_type_name (*bool, optional) {#elasticsearch-suppress_type_name}

Suppress type name to avoid warnings in Elasticsearch 7.x 

Default: -

### target_index_key (string, optional) {#elasticsearch-target_index_key}

Tell this plugin to find the index name to write to in the record under this key in preference to other mechanisms. Key can be specified as path to nested record using dot ('.') as a separator. https://github.com/uken/fluent-plugin-elasticsearch#target_index_key 

Default: -

### target_type_key (string, optional) {#elasticsearch-target_type_key}

Similar to target_index_key config, find the type name to write to in the record under this key (or nested record). If key not found in record - fallback to type_name. 

Default:  fluentd

### template_name (string, optional) {#elasticsearch-template_name}

The name of the template to define. If a template by the name given is already present, it will be left unchanged, unless template_overwrite is set, in which case the template will be updated. 

Default: -

### template_file (*secret.Secret, optional) {#elasticsearch-template_file}

The path to the file containing the template to install. [Secret](../secret/) 

Default: -

### templates (string, optional) {#elasticsearch-templates}

Specify index templates in form of hash. Can contain multiple templates. 

Default: -

### customize_template (string, optional) {#elasticsearch-customize_template}

Specify the string and its value to be replaced in form of hash. Can contain multiple key value pair that would be replaced in the specified template_file. This setting only creates template and to add rollover index please check the rollover_index configuration. 

Default: -

### rollover_index (bool, optional) {#elasticsearch-rollover_index}

Specify this as true when an index with rollover capability needs to be created. https://github.com/uken/fluent-plugin-elasticsearch#rollover_index 

Default:  false

### index_date_pattern (*string, optional) {#elasticsearch-index_date_pattern}

Specify this to override the index date pattern for creating a rollover index. 

Default:  now/d

### deflector_alias (string, optional) {#elasticsearch-deflector_alias}

Specify the deflector alias which would be assigned to the rollover index created. This is useful in case of using the Elasticsearch rollover API 

Default: -

### index_prefix (string, optional) {#elasticsearch-index_prefix}

Specify the index prefix for the rollover index to be created. 

Default:  logstash

### application_name (*string, optional) {#elasticsearch-application_name}

Specify the application name for the rollover index to be created. 

Default:  default

### template_overwrite (bool, optional) {#elasticsearch-template_overwrite}

Always update the template, even if it already exists. 

Default:  false

### max_retry_putting_template (string, optional) {#elasticsearch-max_retry_putting_template}

You can specify times of retry putting template. 

Default:  10

### fail_on_putting_template_retry_exceed (*bool, optional) {#elasticsearch-fail_on_putting_template_retry_exceed}

Indicates whether to fail when max_retry_putting_template is exceeded. If you have multiple output plugin, you could use this property to do not fail on fluentd statup.(default: true) 

Default: true

### fail_on_detecting_es_version_retry_exceed (*bool, optional) {#elasticsearch-fail_on_detecting_es_version_retry_exceed}

fail_on_detecting_es_version_retry_exceed (default: true) 

Default: true

### max_retry_get_es_version (string, optional) {#elasticsearch-max_retry_get_es_version}

You can specify times of retry obtaining Elasticsearch version. 

Default:  15

### request_timeout (string, optional) {#elasticsearch-request_timeout}

You can specify HTTP request timeout. 

Default:  5s

### reload_connections (*bool, optional) {#elasticsearch-reload_connections}

You can tune how the elasticsearch-transport host reloading feature works.(default: true) 

Default: true

### reload_on_failure (bool, optional) {#elasticsearch-reload_on_failure}

Indicates that the elasticsearch-transport will try to reload the nodes addresses if there is a failure while making the request, this can be useful to quickly remove a dead node from the list of addresses. 

Default:  false

### reload_after (string, optional) {#elasticsearch-reload_after}

When reload_connections true, this is the integer number of operations after which the plugin will reload the connections. The default value is 10000. 

Default: -

### resurrect_after (string, optional) {#elasticsearch-resurrect_after}

You can set in the elasticsearch-transport how often dead connections from the elasticsearch-transport's pool will be resurrected. 

Default:  60s

### include_tag_key (bool, optional) {#elasticsearch-include_tag_key}

This will add the Fluentd tag in the JSON record. 

Default:  false

### tag_key (string, optional) {#elasticsearch-tag_key}

This will add the Fluentd tag in the JSON record. 

Default:  tag

### id_key (string, optional) {#elasticsearch-id_key}

https://github.com/uken/fluent-plugin-elasticsearch#id_key 

Default: -

### routing_key (string, optional) {#elasticsearch-routing_key}

Similar to parent_key config, will add _routing into elasticsearch command if routing_key is set and the field does exist in input event. 

Default: -

### remove_keys (string, optional) {#elasticsearch-remove_keys}

https://github.com/uken/fluent-plugin-elasticsearch#remove_keys 

Default: -

### remove_keys_on_update (string, optional) {#elasticsearch-remove_keys_on_update}

Remove keys on update will not update the configured keys in elasticsearch when a record is being updated. This setting only has any effect if the write operation is update or upsert. 

Default: -

### remove_keys_on_update_key (string, optional) {#elasticsearch-remove_keys_on_update_key}

This setting allows remove_keys_on_update to be configured with a key in each record, in much the same way as target_index_key works. 

Default: -

### retry_tag (string, optional) {#elasticsearch-retry_tag}

This setting allows custom routing of messages in response to bulk request failures. The default behavior is to emit failed records using the same tag that was provided. 

Default: -

### write_operation (string, optional) {#elasticsearch-write_operation}

The write_operation can be any of: (index,create,update,upsert) 

Default:  index

### reconnect_on_error (bool, optional) {#elasticsearch-reconnect_on_error}

Indicates that the plugin should reset connection on any error (reconnect on next send). By default it will reconnect only on "host unreachable exceptions". We recommended to set this true in the presence of elasticsearch shield. 

Default:  false

### with_transporter_log (bool, optional) {#elasticsearch-with_transporter_log}

This is debugging purpose option to enable to obtain transporter layer log.  

Default:  false

### content_type (string, optional) {#elasticsearch-content_type}

With content_type application/x-ndjson, elasticsearch plugin adds application/x-ndjson as Content-Profile in payload.  

Default:  application/json

### include_index_in_url (bool, optional) {#elasticsearch-include_index_in_url}

With this option set to true, Fluentd manifests the index name in the request URL (rather than in the request body). You can use this option to enforce an URL-based access control. 

Default: -

### time_parse_error_tag (string, optional) {#elasticsearch-time_parse_error_tag}

With logstash_format true, elasticsearch plugin parses timestamp field for generating index name. If the record has invalid timestamp value, this plugin emits an error event to @ERROR label with time_parse_error_tag configured tag. 

Default: -

### http_backend (string, optional) {#elasticsearch-http_backend}

With http_backend typhoeus, elasticsearch plugin uses typhoeus faraday http backend. Typhoeus can handle HTTP keepalive.  

Default:  excon

### prefer_oj_serializer (bool, optional) {#elasticsearch-prefer_oj_serializer}

With default behavior, Elasticsearch client uses Yajl as JSON encoder/decoder. Oj is the alternative high performance JSON encoder/decoder. When this parameter sets as true, Elasticsearch client uses Oj as JSON encoder/decoder.  

Default:  false

### flatten_hashes (bool, optional) {#elasticsearch-flatten_hashes}

Elasticsearch will complain if you send object and concrete values to the same field. For example, you might have logs that look this, from different places: {"people" => 100} {"people" => {"some" => "thing"}} The second log line will be rejected by the Elasticsearch parser because objects and concrete values can't live in the same field. To combat this, you can enable hash flattening. 

Default: -

### flatten_hashes_separator (string, optional) {#elasticsearch-flatten_hashes_separator}

Flatten separator 

Default: -

### validate_client_version (bool, optional) {#elasticsearch-validate_client_version}

When you use mismatched Elasticsearch server and client libraries, fluent-plugin-elasticsearch cannot send data into Elasticsearch.  

Default:  false

### unrecoverable_error_types (string, optional) {#elasticsearch-unrecoverable_error_types}

Default unrecoverable_error_types parameter is set up strictly. Because es_rejected_execution_exception is caused by exceeding Elasticsearch's thread pool capacity. Advanced users can increase its capacity, but normal users should follow default behavior. If you want to increase it and forcibly retrying bulk request, please consider to change unrecoverable_error_types parameter from default value. Change default value of thread_pool.bulk.queue_size in elasticsearch.yml) 

Default: -

### verify_es_version_at_startup (*bool, optional) {#elasticsearch-verify_es_version_at_startup}

Because Elasticsearch plugin should change behavior each of Elasticsearch major versions. For example, Elasticsearch 6 starts to prohibit multiple type_names in one index, and Elasticsearch 7 will handle only _doc type_name in index. If you want to disable to verify Elasticsearch version at start up, set it as false. When using the following configuration, ES plugin intends to communicate into Elasticsearch 6. (default: true) 

Default: true

### default_elasticsearch_version (string, optional) {#elasticsearch-default_elasticsearch_version}

This parameter changes that ES plugin assumes default Elasticsearch version. 

Default:  5

### custom_headers (string, optional) {#elasticsearch-custom_headers}

This parameter adds additional headers to request. Example: {"token":"secret"}  

Default:  {}

### api_key (*secret.Secret, optional) {#elasticsearch-api_key}

api_key parameter adds authentication header. 

Default: -

### log_es_400_reason (bool, optional) {#elasticsearch-log_es_400_reason}

By default, the error logger won't record the reason for a 400 error from the Elasticsearch API unless you set log_level to debug. However, this results in a lot of log spam, which isn't desirable if all you want is the 400 error reasons. You can set this true to capture the 400 error reasons without all the other debug logs.  

Default:  false

### suppress_doc_wrap (bool, optional) {#elasticsearch-suppress_doc_wrap}

By default, record body is wrapped by 'doc'. This behavior can not handle update script requests. You can set this to suppress doc wrapping and allow record body to be untouched.  

Default:  false

### ignore_exceptions (string, optional) {#elasticsearch-ignore_exceptions}

A list of exception that will be ignored - when the exception occurs the chunk will be discarded and the buffer retry mechanism won't be called. It is possible also to specify classes at higher level in the hierarchy. For example `ignore_exceptions ["Elasticsearch::Transport::Transport::ServerError"]` will match all subclasses of ServerError - Elasticsearch::Transport::Transport::Errors::BadRequest, Elasticsearch::Transport::Transport::Errors::ServiceUnavailable, etc. 

Default: -

### exception_backup (*bool, optional) {#elasticsearch-exception_backup}

Indicates whether to backup chunk when ignore exception occurs. (default: true) 

Default: true

### bulk_message_request_threshold (string, optional) {#elasticsearch-bulk_message_request_threshold}

Configure bulk_message request splitting threshold size. Default value is 20MB. (20 * 1024 * 1024) If you specify this size as negative number, bulk_message request splitting feature will be disabled.  

Default:  20MB

### sniffer_class_name (string, optional) {#elasticsearch-sniffer_class_name}

The default Sniffer used by the Elasticsearch::Transport class works well when Fluentd has a direct connection to all of the Elasticsearch servers and can make effective use of the _nodes API. This doesn't work well when Fluentd must connect through a load balancer or proxy. The parameter sniffer_class_name gives you the ability to provide your own Sniffer class to implement whatever connection reload logic you require. In addition, there is a new Fluent::Plugin::ElasticsearchSimpleSniffer class which reuses the hosts given in the configuration, which is typically the hostname of the load balancer or proxy. https://github.com/uken/fluent-plugin-elasticsearch#sniffer-class-name 

Default: -

### buffer (*Buffer, optional) {#elasticsearch-buffer}

[Buffer](../buffer/) 

Default: -

### slow_flush_log_threshold (string, optional) {#elasticsearch-slow_flush_log_threshold}

The threshold for chunk flush performance check. Parameter type is float, not time, default: 20.0 (seconds) If chunk flush takes longer time than this threshold, fluentd logs warning message and increases metric fluentd_output_status_slow_flush_count. 

Default: -

### enable_ilm (bool, optional) {#elasticsearch-enable_ilm}

Enable Index Lifecycle Management (ILM). 

Default: -

### ilm_policy_id (string, optional) {#elasticsearch-ilm_policy_id}

Specify ILM policy id. 

Default: -

### ilm_policy (string, optional) {#elasticsearch-ilm_policy}

Specify ILM policy contents as Hash. 

Default: -

### ilm_policy_overwrite (bool, optional) {#elasticsearch-ilm_policy_overwrite}

Specify whether overwriting ilm policy or not. 

Default: -

### data_stream_enable (*bool, optional) {#elasticsearch-data_stream_enable}

Use @type elasticsearch_data_stream 

Default: -

### data_stream_name (string, optional) {#elasticsearch-data_stream_name}

You can specify Elasticsearch data stream name by this parameter. This parameter is mandatory for elasticsearch_data_stream. There are some limitations about naming rule. For more details https://www.elastic.co/guide/en/elasticsearch/reference/master/indices-create-data-stream.html#indices-create-data-stream-api-path-params 

Default: -

### data_stream_template_name (string, optional) {#elasticsearch-data_stream_template_name}

Specify an existing index template for the data stream. If not present, a new template is created and named after the data stream.  Further details here https://github.com/uken/fluent-plugin-elasticsearch#configuration---elasticsearch-output-data-stream 

Default:  data_stream_name

### data_stream_ilm_name (string, optional) {#elasticsearch-data_stream_ilm_name}

Specify an existing ILM policy to be applied to the data stream. If not present, either the specified template's or a new ILM default policy is applied.  Further details here https://github.com/uken/fluent-plugin-elasticsearch#configuration---elasticsearch-output-data-stream 

Default:  data_stream_name

### data_stream_ilm_policy (string, optional) {#elasticsearch-data_stream_ilm_policy}

Specify data stream ILM policy contents as Hash. 

Default: -

### data_stream_ilm_policy_overwrite (bool, optional) {#elasticsearch-data_stream_ilm_policy_overwrite}

Specify whether overwriting data stream ilm policy or not. 

Default: -


