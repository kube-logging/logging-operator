# [Record Transformer](https://docs.fluentd.org/filter/record_transformer)
## Overview
 Mutates/transforms incoming event streams.

## Configuration
### RecordTransformer
| Variable Name | Type | Required | Default | Description |
|---|---|---|---|---|
| remove_keys | string | No | - | A comma-delimited list of keys to delete<br> |
| keep_keys | string | No | - | A comma-delimited list of keys to keep.<br> |
| renew_record | bool | No |  false | Create new Hash to transform incoming data <br> |
| renew_time_key | string | No | - | Specify field name of the record to overwrite the time of events. Its value must be unix time.<br> |
| enable_ruby | bool | No |  false | When set to true, the full Ruby syntax is enabled in the ${...} expression. <br> |
| auto_typecast | bool | No |  true | Use original value type. <br> |
| records | []Record | No | - | Add records docs at: https://docs.fluentd.org/filter/record_transformer<br>Records are represented as maps: `key: value`<br> |
 #### Example `Record Transformer` filter configurations
 ```yaml
apiVersion: logging.banzaicloud.io/v1beta1
kind: Flow
metadata:
  name: demo-flow
spec:
  filters:
    - record_transformer:
        records:
        - foo: "bar"
  selectors: {}
  outputRefs:
    - demo-output
 ```

 #### Fluentd Config Result
 ```yaml
<filter **>
  @type record_transformer
  @id test_record_transformer
  <record>
    foo bar
  </record>
</filter>
 ```

---
