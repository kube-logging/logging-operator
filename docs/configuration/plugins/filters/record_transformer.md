---
title: Record Transformer
weight: 200
generated_file: true
---

# [Record Transformer](https://docs.fluentd.org/filter/record_transformer)
## Overview
 Mutates/transforms incoming event streams.

## Configuration
## RecordTransformer

### remove_keys (string, optional) {#recordtransformer-remove_keys}

A comma-delimited list of keys to delete 

Default: -

### keep_keys (string, optional) {#recordtransformer-keep_keys}

A comma-delimited list of keys to keep. 

Default: -

### renew_record (bool, optional) {#recordtransformer-renew_record}

Create new Hash to transform incoming data  

Default:  false

### renew_time_key (string, optional) {#recordtransformer-renew_time_key}

Specify field name of the record to overwrite the time of events. Its value must be unix time. 

Default: -

### enable_ruby (bool, optional) {#recordtransformer-enable_ruby}

When set to true, the full Ruby syntax is enabled in the ${...} expression.  

Default:  false

### auto_typecast (bool, optional) {#recordtransformer-auto_typecast}

Use original value type.  

Default:  true

### records ([]Record, optional) {#recordtransformer-records}

Add records docs at: https://docs.fluentd.org/filter/record_transformer Records are represented as maps: `key: value` 

Default: -


 ## Example `Record Transformer` filter configurations
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
	localOutputRefs:
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
