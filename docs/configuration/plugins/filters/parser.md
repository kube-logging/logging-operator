---
title: Parser
weight: 200
generated_file: true
---

# [Parser Filter](https://docs.fluentd.org/filter/parser)
## Overview
 Parses a string field in event records and mutates its event record with the parsed result.

## Configuration
## ParserConfig

### key_name (string, optional) {#parserconfig-key_name}

Specify field name in the record to parse. If you leave empty the Container Runtime default will be used. 

Default: -

### reserve_time (bool, optional) {#parserconfig-reserve_time}

Keep original event time in parsed result. 

Default: -

### reserve_data (bool, optional) {#parserconfig-reserve_data}

Keep original key-value pair in parsed result. 

Default: -

### remove_key_name_field (bool, optional) {#parserconfig-remove_key_name_field}

Remove key_name field when parsing is succeeded 

Default: -

### replace_invalid_sequence (bool, optional) {#parserconfig-replace_invalid_sequence}

If true, invalid string is replaced with safe characters and re-parse it. 

Default: -

### inject_key_prefix (string, optional) {#parserconfig-inject_key_prefix}

Store parsed values with specified key name prefix. 

Default: -

### hash_value_field (string, optional) {#parserconfig-hash_value_field}

Store parsed values as a hash value in a field. 

Default: -

### emit_invalid_record_to_error (*bool, optional) {#parserconfig-emit_invalid_record_to_error}

Emit invalid record to @ERROR label. Invalid cases are: key not exist, format is not matched, unexpected error 

Default: -

### parse (ParseSection, optional) {#parserconfig-parse}

[Parse Section](#parse-section) 

Default: -

### parsers ([]ParseSection, optional) {#parserconfig-parsers}

Deprecated, use `parse` instead 

Default: -


## Parse Section

### type (string, optional) {#parse section-type}

Parse type: apache2, apache_error, nginx, syslog, csv, tsv, ltsv, json, multiline, none, logfmt, grok, multiline_grok 

Default: -

### expression (string, optional) {#parse section-expression}

Regexp expression to evaluate 

Default: -

### time_key (string, optional) {#parse section-time_key}

Specify time field for event time. If the event doesn't have this field, current time is used. 

Default: -

### keys (string, optional) {#parse section-keys}

Names for fields on each line. (seperated by coma) 

Default: -

### null_value_pattern (string, optional) {#parse section-null_value_pattern}

Specify null value pattern. 

Default: -

### null_empty_string (bool, optional) {#parse section-null_empty_string}

If true, empty string field is replaced with nil 

Default: -

### estimate_current_event (bool, optional) {#parse section-estimate_current_event}

If true, use Fluent::EventTime.now(current time) as a timestamp when time_key is specified. 

Default: -

### keep_time_key (bool, optional) {#parse section-keep_time_key}

If true, keep time field in the record. 

Default: -

### types (string, optional) {#parse section-types}

Types casting the fields to proper types example: field1:type, field2:type 

Default: -

### time_format (string, optional) {#parse section-time_format}

Process value using specified format. This is available only when time_type is string 

Default: -

### time_type (string, optional) {#parse section-time_type}

Parse/format value according to this type available values: float, unixtime, string  

Default:  string

### local_time (bool, optional) {#parse section-local_time}

Ff true, use local time. Otherwise, UTC is used. This is exclusive with utc.  

Default:  true

### utc (bool, optional) {#parse section-utc}

If true, use UTC. Otherwise, local time is used. This is exclusive with localtime  

Default:  false

### timezone (string, optional) {#parse section-timezone}

Use specified timezone. one can parse/format the time value in the specified timezone.  

Default:  nil

### format (string, optional) {#parse section-format}

Only available when using type: multi_format 

Default: -

### format_firstline (string, optional) {#parse section-format_firstline}

Only available when using type: multi_format 

Default: -

### delimiter (string, optional) {#parse section-delimiter}

Only available when using type: ltsv  

Default:  "\t"

### delimiter_pattern (string, optional) {#parse section-delimiter_pattern}

Only available when using type: ltsv 

Default: -

### label_delimiter (string, optional) {#parse section-label_delimiter}

Only available when using type: ltsv  

Default:  ":"

### multiline ([]string, optional) {#parse section-multiline}

The multiline parser plugin parses multiline logs. 

Default: -

### patterns ([]SingleParseSection, optional) {#parse section-patterns}

Only available when using type: multi_format [Parse Section](#parse-section) 

Default: -

### grok_pattern (string, optional) {#parse section-grok_pattern}

Only available when using type: grok, multiline_grok. The pattern of grok. You cannot specify multiple grok pattern with this. 

Default: -

### custom_pattern_path (*secret.Secret, optional) {#parse section-custom_pattern_path}

Only available when using type: grok, multiline_grok. File that includes custom grok patterns. 

Default: -

### grok_failure_key (string, optional) {#parse section-grok_failure_key}

Only available when using type: grok, multiline_grok. The key has grok failure reason. 

Default: -

### grok_name_key (string, optional) {#parse section-grok_name_key}

Only available when using type: grok, multiline_grok. The key name to store grok section's name. 

Default: -

### multiline_start_regexp (string, optional) {#parse section-multiline_start_regexp}

Only available when using type: multiline_grok The regexp to match beginning of multiline. 

Default: -

### grok_patterns ([]GrokSection, optional) {#parse section-grok_patterns}

Only available when using type: grok, multiline_grok. [Grok Section](#grok-section) Specify grok pattern series set. 

Default: -


## Parse Section (single)

### type (string, optional) {#parse section (single)-type}

Parse type: apache2, apache_error, nginx, syslog, csv, tsv, ltsv, json, multiline, none, logfmt, grok, multiline_grok 

Default: -

### expression (string, optional) {#parse section (single)-expression}

Regexp expression to evaluate 

Default: -

### time_key (string, optional) {#parse section (single)-time_key}

Specify time field for event time. If the event doesn't have this field, current time is used. 

Default: -

### null_value_pattern (string, optional) {#parse section (single)-null_value_pattern}

Specify null value pattern. 

Default: -

### null_empty_string (bool, optional) {#parse section (single)-null_empty_string}

If true, empty string field is replaced with nil 

Default: -

### estimate_current_event (bool, optional) {#parse section (single)-estimate_current_event}

If true, use Fluent::EventTime.now(current time) as a timestamp when time_key is specified. 

Default: -

### keep_time_key (bool, optional) {#parse section (single)-keep_time_key}

If true, keep time field in the record. 

Default: -

### types (string, optional) {#parse section (single)-types}

Types casting the fields to proper types example: field1:type, field2:type 

Default: -

### time_format (string, optional) {#parse section (single)-time_format}

Process value using specified format. This is available only when time_type is string 

Default: -

### time_type (string, optional) {#parse section (single)-time_type}

Parse/format value according to this type available values: float, unixtime, string  

Default:  string

### local_time (bool, optional) {#parse section (single)-local_time}

Ff true, use local time. Otherwise, UTC is used. This is exclusive with utc.  

Default:  true

### utc (bool, optional) {#parse section (single)-utc}

If true, use UTC. Otherwise, local time is used. This is exclusive with localtime  

Default:  false

### timezone (string, optional) {#parse section (single)-timezone}

Use specified timezone. one can parse/format the time value in the specified timezone.  

Default:  nil

### format (string, optional) {#parse section (single)-format}

Only available when using type: multi_format 

Default: -

### grok_pattern (string, optional) {#parse section (single)-grok_pattern}

Only available when using format: grok, multiline_grok. The pattern of grok. You cannot specify multiple grok pattern with this. 

Default: -

### custom_pattern_path (*secret.Secret, optional) {#parse section (single)-custom_pattern_path}

Only available when using format: grok, multiline_grok. File that includes custom grok patterns. 

Default: -

### grok_failure_key (string, optional) {#parse section (single)-grok_failure_key}

Only available when using format: grok, multiline_grok. The key has grok failure reason. 

Default: -

### grok_name_key (string, optional) {#parse section (single)-grok_name_key}

Only available when using format: grok, multiline_grok. The key name to store grok section's name. 

Default: -

### multiline_start_regexp (string, optional) {#parse section (single)-multiline_start_regexp}

Only available when using format: multiline_grok The regexp to match beginning of multiline. 

Default: -

### grok_patterns ([]GrokSection, optional) {#parse section (single)-grok_patterns}

Only available when using format: grok, multiline_grok. [Grok Section](#grok-section) Specify grok pattern series set. 

Default: -


## Grok Section

### name (string, optional) {#grok section-name}

The name of grok section. 

Default: -

### pattern (string, required) {#grok section-pattern}

The pattern of grok. 

Default: -

### keep_time_key (bool, optional) {#grok section-keep_time_key}

If true, keep time field in the record. 

Default: -

### time_key (string, optional) {#grok section-time_key}

Specify time field for event time. If the event doesn't have this field, current time is used. 

Default: time

### time_format (string, optional) {#grok section-time_format}

Process value using specified format. This is available only when time_type is string. 

Default: -

### timezone (string, optional) {#grok section-timezone}

Use specified timezone. one can parse/format the time value in the specified timezone. 

Default: -


 ## Example `Parser` filter configurations
 ```yaml
 apiVersion: logging.banzaicloud.io/v1beta1
 kind: Flow
 metadata:

	name: demo-flow

 spec:

	filters:
	  - parser:
	      remove_key_name_field: true
	      reserve_data: true
	      parse:
	        type: multi_format
	        patterns:
	        - format: nginx
	        - format: regexp
	          expression: /foo/
	        - format: none
	selectors: {}
	localOutputRefs:
	  - demo-output

 ```

 #### Fluentd Config Result
 ```yaml
 <filter **>

	@type parser
	@id test_parser
	key_name message
	remove_key_name_field true
	reserve_data true
	<parse>
	  @type multi_format
	  <pattern>
	    format nginx
	  </pattern>
	  <pattern>
	    expression /foo/
	    format regexp
	  </pattern>
	  <pattern>
	    format none
	  </pattern>
	</parse>

 </filter>
 ```

---
