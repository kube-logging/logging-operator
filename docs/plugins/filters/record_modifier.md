---
title: Record Modifier
weight: 200
---

# [Record Modifier](https://github.com/repeatedly/fluent-plugin-record-modifier)
## Overview
 Modify each event record.

## Configuration
### RecordModifier
| Variable Name | Type | Required | Default | Description |
|---|---|---|---|---|
| prepare_value | string | No | - | Prepare values for filtering in configure phase. Prepared values can be used in <record>. You can write any ruby code.<br> |
| char_encoding | string | No | - | Fluentd including some plugins treats logs as a BINARY by default to forward. To overide that, use a target encoding or a from:to encoding here.<br> |
| remove_keys | string | No | - | A comma-delimited list of keys to delete<br> |
| whitelist_keys | string | No | - | This is exclusive with remove_keys<br> |
| replaces | []Replace | No | - | Replace specific value for keys<br> |
| records | []Record | No | - | Add records docs at: https://github.com/repeatedly/fluent-plugin-record-modifier<br>Records are represented as maps: `key: value`<br> |
 #### Example `Record Modifier` filter configurations
 ```yaml
apiVersion: logging.banzaicloud.io/v1beta1
kind: Flow
metadata:
  name: demo-flow
spec:
  filters:
    - record_modifier:
        records:
        - foo: "bar"
  selectors: {}
  outputRefs:
    - demo-output
 ```

 #### Fluentd Config Result
 ```yaml
<filter **>
  @type record_modifier
  @id test_record_modifier
  <record>
    foo bar
  </record>
</filter>
 ```

---
### [Replace Directive](https://github.com/repeatedly/fluent-plugin-record-modifier#replace_keys_value)
#### Specify replace rule. This directive contains three parameters.

| Variable Name | Type | Required | Default | Description |
|---|---|---|---|---|
| key | string | Yes | - | Key to search for<br> |
| expression | string | Yes | - | Regular expression<br> |
| replace | string | Yes | - | Value to replace with<br> |
