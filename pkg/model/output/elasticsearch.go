package output

import (
	"github.com/banzaicloud/logging-operator/pkg/model/secret"
	"github.com/banzaicloud/logging-operator/pkg/model/types"
)

// +kubebuilder:object:generate=true
// +docName:"Elasticsearch"
// Send your logs to Elasticsearch
type ElasticsearchOutput struct {
	// You can specify Elasticsearch host by this parameter. (default:localhost)
	Host string `json:"host,omitempty"`
	// You can specify Elasticsearch port by this parameter.(default: 9200)
	Port int `json:"port,omitempty"`
	// You can specify multiple Elasticsearch hosts with separator ",". If you specify hosts option, host and port options are ignored.
	Hosts string `json:"hosts,omitempty"`
	// User for HTTP Basic authentication. This plugin will escape required URL encoded characters within %{} placeholders. e.g. %{demo+}
	User string `json:"user,omitempty"`
	// Password for HTTP Basic authentication.
	// +docLink:"Secret,./secret.md"
	Password *secret.Secret `json:"password,omitempty"`
	// Path for HTTP Basic authentication.
	Path string `json:"path,omitempty"`
	// Connection scheme (default: http)
	Scheme string `json:"scheme,omitempty"`
	// Skip ssl verification (default: true)
	SslVerify bool `json:"ssl_verify,omitempty"`
	// If you want to configure SSL/TLS version, you can specify ssl_version parameter. [SSLv23, TLSv1, TLSv1_1, TLSv1_2]
	SslVersion string `json:"ssl_version,omitempty"`
	// Enable Logstash log format.(default: false)
	LogstashFormat bool `json:"logstash_format,omitempty"`
	// Adds a @timestamp field to the log, following all settings logstash_format does, except without the restrictions on index_name. This allows one to log to an alias in Elasticsearch and utilize the rollover API.
	IncludeTimestamp bool `json:"include_timestamp,omitempty"`
	// Set the Logstash prefix.(default: true)
	LogstashPrefix string `json:"logstash_prefix,omitempty"`
	// Set the Logstash prefix separator.(default: -)
	LogstashPrefixSeparator string `json:"logstash_prefix_separator,omitempty"`
	// Set the Logstash date format.(default: %Y.%m.%d)
	LogstashDateformat string `json:"logstash_dateformat,omitempty"`
	// This param is to set a pipeline id of your elasticsearch to be added into the request, you can configure ingest node.
	Pipeline string `json:"pipeline,omitempty"`
	// The format of the time stamp field (@timestamp or what you specify with time_key). This parameter only has an effect when logstash_format is true as it only affects the name of the index we write to.
	TimeKeyFormat string `json:"time_key_format,omitempty"`
	// Should the record not include a time_key, define the degree of sub-second time precision to preserve from the time portion of the routed event.
	TimePrecision string `json:"time_precision,omitempty"`
	// By default, when inserting records in Logstash format, @timestamp is dynamically created with the time at log ingestion. If you'd like to use a custom time, include an @timestamp with your record.
	TimeKey string `json:"time_key,omitempty"`
	// By default, the records inserted into index logstash-YYMMDD with UTC (Coordinated Universal Time). This option allows to use local time if you describe utc_index to false.(default: true)
	UtcIndex bool `json:"utc_index,omitempty"`
	// Tell this plugin to find the index name to write to in the record under this key in preference to other mechanisms. Key can be specified as path to nested record using dot ('.') as a separator. https://github.com/uken/fluent-plugin-elasticsearch#target_index_key
	TargetIndexKey string `json:"target_index_key,omitempty"`
	// Similar to target_index_key config, find the type name to write to in the record under this key (or nested record). If key not found in record - fallback to type_name.(default: true)
	TargetTypeKey string `json:"target_type_key,omitempty"`
	// The name of the template to define. If a template by the name given is already present, it will be left unchanged, unless template_overwrite is set, in which case the template will be updated.
	TemplateName string `json:"template_name,omitempty"`

	// The path to the file containing the template to install.
	TemplateFile string `json:"template_file,omitempty"`

	// Specify index templates in form of hash. Can contain multiple templates.
	Templates string `json:"templates,omitempty"`
	// Specify the string and its value to be replaced in form of hash. Can contain multiple key value pair that would be replaced in the specified template_file. This setting only creates template and to add rollover index please check the rollover_index configuration.
	CustomizeTemplate string `json:"customize_template,omitempty"`
	// Specify this as true when an index with rollover capability needs to be created.(default: false) https://github.com/uken/fluent-plugin-elasticsearch#rollover_index
	RolloverIndex bool `json:"rollover_index,omitempty"`
	// Specify this to override the index date pattern for creating a rollover index.(default: now/d)
	IndexDatePattern string `json:"index_date_pattern,omitempty"`
	// Specify the deflector alias which would be assigned to the rollover index created. This is useful in case of using the Elasticsearch rollover API
	DeflectorAlias string `json:"deflector_alias,omitempty"`
	// Specify the index prefix for the rollover index to be created.
	IndexPrefix string `json:"index_prefix,omitempty"`
	// Specify the application name for the rollover index to be created.(default: default)
	ApplicationName string `json:"application_name,omitempty"`
	// Always update the template, even if it already exists.(default: false)
	TemplateOverwrite bool `json:"template_overwrite,omitempty"`
	// You can specify times of retry putting template.(default: 10)
	MaxRetryPuttingTemplate string `json:"max_retry_putting_template,omitempty"`
	// Indicates whether to fail when max_retry_putting_template is exceeded. If you have multiple output plugin, you could use this property to do not fail on fluentd statup.(default: true)
	FailOnPuttingTemplateRetryExceed bool `json:"fail_on_putting_template_retry_exceed,omitempty"`

	// You can specify times of retry obtaining Elasticsearch version.(default: 15)
	MaxRetryGetEsVersion string `json:"max_retry_get_es_version,omitempty"`

	// You can specify HTTP request timeout.(default: 5s)
	RequestTimeout string `json:"request_timeout,omitempty"`
	// You can tune how the elasticsearch-transport host reloading feature works.(default: true)
	ReloadConnections bool `json:"reload_connections,omitempty"`
	//Indicates that the elasticsearch-transport will try to reload the nodes addresses if there is a failure while making the request, this can be useful to quickly remove a dead node from the list of addresses.(default: false)
	ReloadOnFailure bool `json:"reload_on_failure,omitempty"`
	// You can set in the elasticsearch-transport how often dead connections from the elasticsearch-transport's pool will be resurrected.(default: 60s)
	ResurrectAfter string `json:"resurrect_after,omitempty"`

	// This will add the Fluentd tag in the JSON record.(default: false)
	IncludeTagKey bool `json:"include_tag_key,omitempty"`
	// This will add the Fluentd tag in the JSON record.(default: tag)
	TagKey string `json:"tag_key,omitempty"`

	// https://github.com/uken/fluent-plugin-elasticsearch#id_key
	IdKey string `json:"id_key,omitempty"`
	// Similar to parent_key config, will add _routing into elasticsearch command if routing_key is set and the field does exist in input event.
	RoutingKey string `json:"routing_key,omitempty"`
	// Remove keys on update will not update the configured keys in elasticsearch when a record is being updated. This setting only has any effect if the write operation is update or upsert.
	RemoveKeysOnUpdate string `json:"remove_keys_on_update,omitempty"`
	// This setting allows remove_keys_on_update to be configured with a key in each record, in much the same way as target_index_key works.
	RemoveKeysOnUpdateKey string `json:"remove_keys_on_update_key,omitempty"`
	// This setting allows custom routing of messages in response to bulk request failures. The default behavior is to emit failed records using the same tag that was provided.
	RetryTag string `json:"retry_tag,omitempty"`
	// The write_operation can be any of: (index,create,update,upsert)(default: index)
	WriteOperation string `json:"write_operation,omitempty"`
	// Indicates that the plugin should reset connection on any error (reconnect on next send). By default it will reconnect only on "host unreachable exceptions". We recommended to set this true in the presence of elasticsearch shield.(default: false)
	ReconnectOnError bool `json:"reconnect_on_error,omitempty"`
	// This is debugging purpose option to enable to obtain transporter layer log. (default: false)
	WithTransporterLog bool `json:"with_transporter_log,omitempty"`
	// With content_type application/x-ndjson, elasticsearch plugin adds application/x-ndjson as Content-Type in payload. (default: application/json)
	ContentType string `json:"content_type,omitempty"`
	//With this option set to true, Fluentd manifests the index name in the request URL (rather than in the request body). You can use this option to enforce an URL-based access control.
	IncludeIndexInUrl bool `json:"include_index_in_url,omitempty"`
	// With logstash_format true, elasticsearch plugin parses timestamp field for generating index name. If the record has invalid timestamp value, this plugin emits an error event to @ERROR label with time_parse_error_tag configured tag.
	TimeParseErrorTag string `json:"time_parse_error_tag,omitempty"`
	// With http_backend typhoeus, elasticsearch plugin uses typhoeus faraday http backend. Typhoeus can handle HTTP keepalive. (default: excon)
	HttpBackend string `json:"http_backend,omitempty"`

	// With default behavior, Elasticsearch client uses Yajl as JSON encoder/decoder. Oj is the alternative high performance JSON encoder/decoder. When this parameter sets as true, Elasticsearch client uses Oj as JSON encoder/decoder. (default: fqlse)
	OreferOjSerializer bool `json:"prefer_oj_serializer,omitempty"`

	// +docLink:"Buffer,./buffer.md"
	Buffer *Buffer `json:"buffer,omitempty"`
}

func (e *ElasticsearchOutput) ToDirective(secretLoader secret.SecretLoader) (types.Directive, error) {
	elasticsearch := &types.OutputPlugin{
		PluginMeta: types.PluginMeta{
			Type:      "elasticsearch",
			Directive: "match",
			Tag:       "**",
		},
	}
	if params, err := types.NewStructToStringMapper(secretLoader).StringsMap(e); err != nil {
		return nil, err
	} else {
		elasticsearch.Params = params
	}
	if e.Buffer != nil {
		if buffer, err := e.Buffer.ToDirective(secretLoader); err != nil {
			return nil, err
		} else {
			elasticsearch.SubDirectives = append(elasticsearch.SubDirectives, buffer)
		}
	}
	return elasticsearch, nil
}
