# [Parser Filter](https://docs.fluentd.org/filter/parser)
## Overview
 Parses" string field in event records and mutates its

## Configuration
### ParserConfig
| Variable Name | Type | Required | Default | Description |
|---|---|---|---|---|
| key_name | string | No | - | Specify field name in the record to parse. If you leave empty the Container Runtime default will be used.<br> |
| reserve_time | bool | No | - | Keep original event time in parsed result.<br> |
| reserve_data | bool | No | - | Keep original key-value pair in parsed result.<br> |
| remove_key_name_field | bool | No | - | Remove key_name field when parsing is succeeded<br> |
| replace_invalid_sequence | bool | No | - | If true, invalid string is replaced with safe characters and re-parse it.<br> |
| inject_key_prefix | string | No | - | Store parsed values with specified key name prefix.<br> |
| hash_value_field | string | No | - | Store parsed values as a hash value in a field.<br> |
| emit_invalid_record_to_error | bool | No | - | Emit invalid record to @ERROR label. Invalid cases are: key not exist, format is not matched, unexpected error<br> |
| parsers | []ParseSection | No | - | Deprecated, use parse<br> |
| parse | ParseSection | No | - | [Parse Section](#Parse-Section)<br> |
### Parse Section
| Variable Name | Type | Required | Default | Description |
|---|---|---|---|---|
| type | string | No | - | Parse type: apache2, apache_error, nginx, syslog, csv, tsv, ltsv, json, multiline, none, logfmt<br> |
| expression | string | No | - | Regexp expression to evaluate<br> |
| time_key | string | No | - | Specify time field for event time. If the event doesn't have this field, current time is used.<br> |
| null_value_pattern | string | No | - | Specify null value pattern.<br> |
| null_empty_string | bool | No | - | If true, empty string field is replaced with nil<br> |
| estimate_current_event | bool | No | - | If true, use Fluent::EventTime.now(current time) as a timestamp when time_key is specified.<br> |
| keep_time_key | bool | No | - | If true, keep time field in the record.<br> |
| types | string | No | - | Types casting the fields to proper types example: field1:type, field2:type<br> |
| time_format | string | No | - | Process value using specified format. This is available only when time_type is string<br> |
| time_type | string | No |  string | Parse/format value according to this type available values: float, unixtime, string <br> |
| local_time | bool | No |  true | Ff true, use local time. Otherwise, UTC is used. This is exclusive with utc. <br> |
| utc | bool | No |  false | If true, use UTC. Otherwise, local time is used. This is exclusive with localtime <br> |
| timezone | string | No |  nil | Use specified timezone. one can parse/format the time value in the specified timezone. <br> |
| patterns | []ParseSection | No | - | Only available when using type: multi_format<br>[Parse Section](#Parse-Section)<br> |
| format | string | No | - | Only available when using type: multi_format<br> |
 #### Example `Parser` filter configurations
 ```yaml
apiVersion: logging.banzaicloud.io/v1beta1
kind: Flow
metadata:
  name: demo-flow
spec:
  filters:
    - parse:
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
  outputRefs:
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
