---
title: Record Modifier
weight: 200
generated_file: true
---

# [Record Modifier](https://github.com/repeatedly/fluent-plugin-record-modifier)
## Overview
 Modify each event record.

## Configuration
## RecordModifier

### char_encoding (string, optional) {#recordmodifier-char_encoding}

Fluentd including some plugins treats logs as a BINARY by default to forward. To overide that, use a target encoding or a from:to encoding here. 


### prepare_value (string, optional) {#recordmodifier-prepare_value}

Prepare values for filtering in configure phase. Prepared values can be used in `<record>`. You can write any ruby code. 


### records ([]Record, optional) {#recordmodifier-records}

Add records docs at: https://github.com/repeatedly/fluent-plugin-record-modifier Records are represented as maps: `key: value` 


### remove_keys (string, optional) {#recordmodifier-remove_keys}

A comma-delimited list of keys to delete 


### replaces ([]Replace, optional) {#recordmodifier-replaces}

Replace specific value for keys 


### whitelist_keys (string, optional) {#recordmodifier-whitelist_keys}

This is exclusive with remove_keys 





## Example `Record Modifier` filter configurations

{{< highlight yaml >}}
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
  localOutputRefs:
    - demo-output
{{</ highlight >}}

Fluentd config result:

{{< highlight xml >}}
<filter **>
  @type record_modifier
  @id test_record_modifier
  <record>
    foo bar
  </record>
</filter>
{{</ highlight >}}


---
## [Replace Directive](https://github.com/repeatedly/fluent-plugin-record-modifier#replace_keys_value)

Specify replace rule. This directive contains three parameters.

### expression (string, required) {#[replace directive](https://github.com/repeatedly/fluent-plugin-record-modifier#replace_keys_value)-expression}

Regular expression 


### key (string, required) {#[replace directive](https://github.com/repeatedly/fluent-plugin-record-modifier#replace_keys_value)-key}

Key to search for 


### replace (string, required) {#[replace directive](https://github.com/repeatedly/fluent-plugin-record-modifier#replace_keys_value)-replace}

Value to replace with 



