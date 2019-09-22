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
| ssl_verify | bool | Yes | true | Skip ssl verification (default: true)<br> |
| ssl_version | string | No | - | If you want to configure SSL/TLS version, you can specify ssl_version parameter. [SSLv23, TLSv1, TLSv1_1, TLSv1_2]<br> |
| logstash_format | bool | No |  false | Enable Logstash log format.<br> |
| include_timestamp | bool | No | - | Adds a @timestamp field to the log, following all settings logstash_format does, except without the restrictions on index_name. This allows one to log to an alias in Elasticsearch and utilize the rollover API.<br> |
| logstash_prefix | string | No |  true | Set the Logstash prefix.<br> |
| logstash_prefix_separator | string | No |  - | Set the Logstash prefix separator.<br> |
| logstash_dateformat | string | No |  %Y.%m.%d | Set the Logstash date format.<br> |
| pipeline | string | No | - | This param is to set a pipeline id of your elasticsearch to be added into the request, you can configure ingest node.<br> |
| time_key_format | string | No | - | The format of the time stamp field (@timestamp or what you specify with time_key). This parameter only has an effect when logstash_format is true as it only affects the name of the index we write to.<br> |
| time_precision | string | No | - | Should the record not include a time_key, define the degree of sub-second time precision to preserve from the time portion of the routed event.<br> |
| time_key | string | No | - | By default, when inserting records in Logstash format, @timestamp is dynamically created with the time at log ingestion. If you'd like to use a custom time, include an @timestamp with your record.<br> |
| utc_index | bool | No |  true | By default, the records inserted into index logstash-YYMMDD with UTC (Coordinated Universal Time). This option allows to use local time if you describe utc_index to false.<br> |
| target_index_key | string | No | - | Tell this plugin to find the index name to write to in the record under this key in preference to other mechanisms. Key can be specified as path to nested record using dot ('.') as a separator. https://github.com/uken/fluent-plugin-elasticsearch#target_index_key<br> |
| target_type_key | string | No |  true | Similar to target_index_key config, find the type name to write to in the record under this key (or nested record). If key not found in record - fallback to type_name.<br> |
| template_name | string | No | - | The name of the template to define. If a template by the name given is already present, it will be left unchanged, unless template_overwrite is set, in which case the template will be updated.<br> |
| template_file | string | No | - | The path to the file containing the template to install.<br> |
| templates | string | No | - | Specify index templates in form of hash. Can contain multiple templates.<br> |
| customize_template | string | No | - | Specify the string and its value to be replaced in form of hash. Can contain multiple key value pair that would be replaced in the specified template_file. This setting only creates template and to add rollover index please check the rollover_index configuration.<br> |
| rollover_index | bool | No |  false | Specify this as true when an index with rollover capability needs to be created. https://github.com/uken/fluent-plugin-elasticsearch#rollover_index<br> |
| index_date_pattern | string | No |  now/d | Specify this to override the index date pattern for creating a rollover index.<br> |
| deflector_alias | string | No | - | Specify the deflector alias which would be assigned to the rollover index created. This is useful in case of using the Elasticsearch rollover API<br> |
| index_prefix | string | No | - | Specify the index prefix for the rollover index to be created.<br> |
| application_name | string | No |  default | Specify the application name for the rollover index to be created.<br> |
| template_overwrite | bool | No |  false | Always update the template, even if it already exists.<br> |
| max_retry_putting_template | string | No |  10 | You can specify times of retry putting template.<br> |
| fail_on_putting_template_retry_exceed | bool | No |  true | Indicates whether to fail when max_retry_putting_template is exceeded. If you have multiple output plugin, you could use this property to do not fail on fluentd statup.<br> |
| max_retry_get_es_version | string | No |  15 | You can specify times of retry obtaining Elasticsearch version.<br> |
| request_timeout | string | No |  5s | You can specify HTTP request timeout.<br> |
| reload_connections | bool | No |  true | You can tune how the elasticsearch-transport host reloading feature works.<br> |
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
| prefer_oj_serializer | bool | No |  fqlse | With default behavior, Elasticsearch client uses Yajl as JSON encoder/decoder. Oj is the alternative high performance JSON encoder/decoder. When this parameter sets as true, Elasticsearch client uses Oj as JSON encoder/decoder. <br> |
| buffer | *Buffer | No | - | [Buffer](./buffer.md)<br> |
