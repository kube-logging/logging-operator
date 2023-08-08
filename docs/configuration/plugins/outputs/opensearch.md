---
title: OpenSearch
weight: 200
generated_file: true
---

# OpenSearch output plugin for Fluentd
## Overview
 More info at https://github.com/fluent/fluent-plugin-opensearch
 >Example Deployment: [Save all logs to OpenSearch](../../../../quickstarts/es-nginx/)

 ## Example output configurations
 ```yaml
 spec:

	opensearch:
	  host: opensearch-cluster.default.svc.cluster.local
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
## OpenSearch

Send your logs to OpenSearch

### host (string, optional) {#opensearch-host}

You can specify OpenSearch host by this parameter.  

Default: localhost

### port (int, optional) {#opensearch-port}

You can specify OpenSearch port by this parameter. 

Default:  9200

### user (string, optional) {#opensearch-user}

User for HTTP Basic authentication. This plugin will escape required URL encoded characters within %{} placeholders. e.g. %{demo+} 

Default: -

### password (*secret.Secret, optional) {#opensearch-password}

Password for HTTP Basic authentication. [Secret](../secret/) 

Default: -

### path (string, optional) {#opensearch-path}

Path for HTTP Basic authentication. 

Default: -

### scheme (string, optional) {#opensearch-scheme}

Connection scheme  

Default:  http

### hosts (string, optional) {#opensearch-hosts}

You can specify multiple OpenSearch hosts with separator ",". If you specify hosts option, host and port options are ignored. 

Default: -

### target_index_key (string, optional) {#opensearch-target_index_key}

Tell this plugin to find the index name to write to in the record under this key in preference to other mechanisms. Key can be specified as path to nested record using dot ('.') as a separator. 

Default: -

### time_key_format (string, optional) {#opensearch-time_key_format}

The format of the time stamp field (@timestamp or what you specify with time_key). This parameter only has an effect when logstash_format is true as it only affects the name of the index we write to. 

Default: -

### time_precision (string, optional) {#opensearch-time_precision}

Should the record not include a time_key, define the degree of sub-second time precision to preserve from the time portion of the routed event. 

Default: -

### include_timestamp (bool, optional) {#opensearch-include_timestamp}

Adds a @timestamp field to the log, following all settings logstash_format does, except without the restrictions on index_name. This allows one to log to an alias in OpenSearch and utilize the rollover API. 

Default:  false

### logstash_format (bool, optional) {#opensearch-logstash_format}

Enable Logstash log format. 

Default:  false

### logstash_prefix (string, optional) {#opensearch-logstash_prefix}

Set the Logstash prefix. 

Default:  logstash

### logstash_prefix_separator (string, optional) {#opensearch-logstash_prefix_separator}

Set the Logstash prefix separator. 

Default:  -

### logstash_dateformat (string, optional) {#opensearch-logstash_dateformat}

Set the Logstash date format. 

Default:  %Y.%m.%d

### utc_index (*bool, optional) {#opensearch-utc_index}

By default, the records inserted into index logstash-YYMMDD with UTC (Coordinated Universal Time). This option allows to use local time if you describe utc_index to false.(default: true) 

Default: true

### suppress_type_name (*bool, optional) {#opensearch-suppress_type_name}

Suppress type name to avoid warnings in OpenSearch 

Default: -

### index_name (string, optional) {#opensearch-index_name}

The index name to write events to  

Default:  fluentd

### id_key (string, optional) {#opensearch-id_key}

Field on your data to identify the data uniquely 

Default: -

### write_operation (string, optional) {#opensearch-write_operation}

The write_operation can be any of: (index,create,update,upsert) 

Default:  index

### parent_key (string, optional) {#opensearch-parent_key}

parent_key 

Default: -

### routing_key (string, optional) {#opensearch-routing_key}

routing_key 

Default: -

### request_timeout (string, optional) {#opensearch-request_timeout}

You can specify HTTP request timeout. 

Default:  5s

### reload_connections (*bool, optional) {#opensearch-reload_connections}

You can tune how the OpenSearch-transport host reloading feature works.(default: true) 

Default: true

### reload_on_failure (bool, optional) {#opensearch-reload_on_failure}

Indicates that the OpenSearch-transport will try to reload the nodes addresses if there is a failure while making the request, this can be useful to quickly remove a dead node from the list of addresses. 

Default:  false

### retry_tag (string, optional) {#opensearch-retry_tag}

This setting allows custom routing of messages in response to bulk request failures. The default behavior is to emit failed records using the same tag that was provided. 

Default: -

### resurrect_after (string, optional) {#opensearch-resurrect_after}

You can set in the OpenSearch-transport how often dead connections from the OpenSearch-transport's pool will be resurrected. 

Default:  60s

### time_key (string, optional) {#opensearch-time_key}

By default, when inserting records in Logstash format, @timestamp is dynamically created with the time at log ingestion. If you'd like to use a custom time, include an @timestamp with your record. 

Default: -

### time_key_exclude_timestamp (bool, optional) {#opensearch-time_key_exclude_timestamp}

time_key_exclude_timestamp  

Default:  false

### ssl_verify (*bool, optional) {#opensearch-ssl_verify}

Skip ssl verification (default: true) 

Default: true

### client_key (*secret.Secret, optional) {#opensearch-client_key}

Client certificate key 

Default: -

### client_cert (*secret.Secret, optional) {#opensearch-client_cert}

Client certificate 

Default: -

### client_key_pass (*secret.Secret, optional) {#opensearch-client_key_pass}

Client key password 

Default: -

### ca_file (*secret.Secret, optional) {#opensearch-ca_file}

CA certificate 

Default: -

### ssl_version (string, optional) {#opensearch-ssl_version}

If you want to configure SSL/TLS version, you can specify ssl_version parameter. [SSLv23, TLSv1, TLSv1_1, TLSv1_2] 

Default: -

### remove_keys_on_update (string, optional) {#opensearch-remove_keys_on_update}

Remove keys on update will not update the configured keys in OpenSearch when a record is being updated. This setting only has any effect if the write operation is update or upsert. 

Default: -

### remove_keys_on_update_key (string, optional) {#opensearch-remove_keys_on_update_key}

This setting allows remove_keys_on_update to be configured with a key in each record, in much the same way as target_index_key works. 

Default: -

### flatten_hashes (bool, optional) {#opensearch-flatten_hashes}

https://github.com/fluent/fluent-plugin-opensearch#hash-flattening 

Default: -

### flatten_hashes_separator (string, optional) {#opensearch-flatten_hashes_separator}

Flatten separator 

Default: -

### template_name (string, optional) {#opensearch-template_name}

The name of the template to define. If a template by the name given is already present, it will be left unchanged, unless template_overwrite is set, in which case the template will be updated. 

Default: -

### template_file (*secret.Secret, optional) {#opensearch-template_file}

The path to the file containing the template to install. [Secret](../secret/) 

Default: -

### template_overwrite (bool, optional) {#opensearch-template_overwrite}

Always update the template, even if it already exists. 

Default:  false

### customize_template (string, optional) {#opensearch-customize_template}

Specify the string and its value to be replaced in form of hash. Can contain multiple key value pair that would be replaced in the specified template_file. This setting only creates template and to add rollover index please check the rollover_index configuration. 

Default: -

### index_date_pattern (*string, optional) {#opensearch-index_date_pattern}

Specify this to override the index date pattern for creating a rollover index. 

Default:  now/d

### index_separator (string, optional) {#opensearch-index_separator}

index_separator  

Default:  -

### application_name (*string, optional) {#opensearch-application_name}

Specify the application name for the rollover index to be created. 

Default:  default

### templates (string, optional) {#opensearch-templates}

Specify index templates in form of hash. Can contain multiple templates. 

Default: -

### max_retry_putting_template (string, optional) {#opensearch-max_retry_putting_template}

You can specify times of retry putting template. 

Default:  10

### fail_on_putting_template_retry_exceed (*bool, optional) {#opensearch-fail_on_putting_template_retry_exceed}

Indicates whether to fail when max_retry_putting_template is exceeded. If you have multiple output plugin, you could use this property to do not fail on fluentd statup.(default: true) 

Default: true

### fail_on_detecting_os_version_retry_exceed (*bool, optional) {#opensearch-fail_on_detecting_os_version_retry_exceed}

fail_on_detecting_os_version_retry_exceed (default: true) 

Default: true

### max_retry_get_os_version (int, optional) {#opensearch-max_retry_get_os_version}

max_retry_get_os_version  

Default:  15

### include_tag_key (bool, optional) {#opensearch-include_tag_key}

This will add the Fluentd tag in the JSON record. 

Default:  false

### tag_key (string, optional) {#opensearch-tag_key}

This will add the Fluentd tag in the JSON record. 

Default:  tag

### time_parse_error_tag (string, optional) {#opensearch-time_parse_error_tag}

With logstash_format true, OpenSearch plugin parses timestamp field for generating index name. If the record has invalid timestamp value, this plugin emits an error event to @ERROR label with time_parse_error_tag configured tag. 

Default: -

### reconnect_on_error (bool, optional) {#opensearch-reconnect_on_error}

Indicates that the plugin should reset connection on any error (reconnect on next send). By default it will reconnect only on "host unreachable exceptions". We recommended to set this true in the presence of OpenSearch shield. 

Default:  false

### pipeline (string, optional) {#opensearch-pipeline}

This param is to set a pipeline id of your OpenSearch to be added into the request, you can configure ingest node. 

Default: -

### with_transporter_log (bool, optional) {#opensearch-with_transporter_log}

This is debugging purpose option to enable to obtain transporter layer log.  

Default:  false

### emit_error_for_missing_id (bool, optional) {#opensearch-emit_error_for_missing_id}

emit_error_for_missing_id  

Default:  false

### sniffer_class_name (string, optional) {#opensearch-sniffer_class_name}

TThe default Sniffer used by the OpenSearch::Transport class works well when Fluentd has a direct connection to all of the OpenSearch servers and can make effective use of the _nodes API. This doesn't work well when Fluentd must connect through a load balancer or proxy. The parameter sniffer_class_name gives you the ability to provide your own Sniffer class to implement whatever connection reload logic you require. In addition, there is a new Fluent::Plugin::OpenSearchSimpleSniffer class which reuses the hosts given in the configuration, which is typically the hostname of the load balancer or proxy. For example, a configuration like this would cause connections to logging-os to reload every 100 operations: https://github.com/fluent/fluent-plugin-opensearch#sniffer-class-name 

Default: -

### selector_class_name (string, optional) {#opensearch-selector_class_name}

selector_class_name 

Default: -

### reload_after (string, optional) {#opensearch-reload_after}

When reload_connections true, this is the integer number of operations after which the plugin will reload the connections. The default value is 10000. 

Default: -

### include_index_in_url (bool, optional) {#opensearch-include_index_in_url}

With this option set to true, Fluentd manifests the index name in the request URL (rather than in the request body). You can use this option to enforce an URL-based access control. 

Default: -

### http_backend (string, optional) {#opensearch-http_backend}

With http_backend typhoeus, opensearch plugin uses typhoeus faraday http backend. Typhoeus can handle HTTP keepalive.  

Default:  excon

### http_backend_excon_nonblock (*bool, optional) {#opensearch-http_backend_excon_nonblock}

http_backend_excon_nonblock (default: true) 

Default: true

### validate_client_version (bool, optional) {#opensearch-validate_client_version}

When you use mismatched OpenSearch server and client libraries, fluent-plugin-opensearch cannot send data into OpenSearch.   

Default:  false

### prefer_oj_serializer (bool, optional) {#opensearch-prefer_oj_serializer}

With default behavior, OpenSearch client uses Yajl as JSON encoder/decoder. Oj is the alternative high performance JSON encoder/decoder. When this parameter sets as true, OpenSearch client uses Oj as JSON encoder/decoder.  

Default:  false

### unrecoverable_error_types (string, optional) {#opensearch-unrecoverable_error_types}

Default unrecoverable_error_types parameter is set up strictly. Because rejected_execution_exception is caused by exceeding OpenSearch's thread pool capacity. Advanced users can increase its capacity, but normal users should follow default behavior. 

Default: -

### unrecoverable_record_types (string, optional) {#opensearch-unrecoverable_record_types}

unrecoverable_record_types 

Default: -

### emit_error_label_event (*bool, optional) {#opensearch-emit_error_label_event}

emit_error_label_event (default: true) 

Default: true

### verify_os_version_at_startup (*bool, optional) {#opensearch-verify_os_version_at_startup}

verify_os_version_at_startup (default: true) 

Default: true

### default_opensearch_version (int, optional) {#opensearch-default_opensearch_version}

max_retry_get_os_version  

Default:  1

### log_os_400_reason (bool, optional) {#opensearch-log_os_400_reason}

log_os_400_reason  

Default:  false

### custom_headers (string, optional) {#opensearch-custom_headers}

This parameter adds additional headers to request. Example: {"token":"secret"}  

Default:  {}

### suppress_doc_wrap (bool, optional) {#opensearch-suppress_doc_wrap}

By default, record body is wrapped by 'doc'. This behavior can not handle update script requests. You can set this to suppress doc wrapping and allow record body to be untouched.  

Default:  false

### ignore_exceptions (string, optional) {#opensearch-ignore_exceptions}

A list of exception that will be ignored - when the exception occurs the chunk will be discarded and the buffer retry mechanism won't be called. It is possible also to specify classes at higher level in the hierarchy. 

Default: -

### exception_backup (*bool, optional) {#opensearch-exception_backup}

Indicates whether to backup chunk when ignore exception occurs. (default: true) 

Default: true

### bulk_message_request_threshold (string, optional) {#opensearch-bulk_message_request_threshold}

Configure bulk_message request splitting threshold size. Default value is 20MB. (20 * 1024 * 1024) If you specify this size as negative number, bulk_message request splitting feature will be disabled.  

Default:  20MB

### compression_level (string, optional) {#opensearch-compression_level}

compression_level 

Default: -

### truncate_caches_interval (string, optional) {#opensearch-truncate_caches_interval}

truncate_caches_interval 

Default: -

### use_legacy_template (*bool, optional) {#opensearch-use_legacy_template}

use_legacy_template (default: true) 

Default: true

### catch_transport_exception_on_retry (*bool, optional) {#opensearch-catch_transport_exception_on_retry}

catch_transport_exception_on_retry (default: true) 

Default: true

### target_index_affinity (bool, optional) {#opensearch-target_index_affinity}

target_index_affinity  

Default:  false

### buffer (*Buffer, optional) {#opensearch-buffer}

Default: -

### slow_flush_log_threshold (string, optional) {#opensearch-slow_flush_log_threshold}

The threshold for chunk flush performance check. Parameter type is float, not time, default: 20.0 (seconds) If chunk flush takes longer time than this threshold, fluentd logs warning message and increases metric fluentd_output_status_slow_flush_count. 

Default: -

### data_stream_enable (*bool, optional) {#opensearch-data_stream_enable}

Use @type opensearch_data_stream 

Default: -

### data_stream_name (string, optional) {#opensearch-data_stream_name}

You can specify Opensearch data stream name by this parameter. This parameter is mandatory for opensearch_data_stream. 

Default: -

### data_stream_template_name (string, optional) {#opensearch-data_stream_template_name}

Specify an existing index template for the data stream. If not present, a new template is created and named after the data stream.  

Default:  data_stream_name


