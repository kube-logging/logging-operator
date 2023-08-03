---
title: SQS
weight: 200
generated_file: true
---

# [SQS Output](https://github.com/ixixi/fluent-plugin-sqs)
## Overview
 Fluentd output plugin for SQS.

## Configuration
## Output Config

### sqs_url (string, optional) {#output config-sqs_url}

SQS queue url e.g. https://sqs.us-west-2.amazonaws.com/123456789012/myqueue 

Default: -

### queue_name (string, optional) {#output config-queue_name}

SQS queue name - required if sqs_url is not set 

Default: -

### aws_key_id (*secret.Secret, optional) {#output config-aws_key_id}

AWS access key id 

Default: -

### aws_sec_key (*secret.Secret, optional) {#output config-aws_sec_key}

AWS secret key 

Default: -

### create_queue (*bool, optional) {#output config-create_queue}

Create SQS queue  

Default:  true

### region (string, optional) {#output config-region}

AWS region  

Default:  ap-northeast-1

### message_group_id (string, optional) {#output config-message_group_id}

Message group id for FIFO queue 

Default: -

### delay_seconds (int, optional) {#output config-delay_seconds}

Delivery delay seconds  

Default:  0

### include_tag (*bool, optional) {#output config-include_tag}

Include tag  

Default:  true

### tag_property_name (string, optional) {#output config-tag_property_name}

Tags property name in json  

Default:  '__tag'

### buffer (*Buffer, optional) {#output config-buffer}

[Buffer](../buffer/) 

Default: -

### slow_flush_log_threshold (string, optional) {#output config-slow_flush_log_threshold}

The threshold for chunk flush performance check. Parameter type is float, not time, default: 20.0 (seconds) If chunk flush takes longer time than this threshold, fluentd logs warning message and increases metric fluentd_output_status_slow_flush_count. 

Default: -


 ## Example `SQS` output configurations
 ```yaml
 apiVersion: logging.banzaicloud.io/v1beta1
 kind: Output
 metadata:

	name: sqs-output-sample

 spec:

	sqs:
	  queue_name: some-aws-sqs-queue
	  create_queue: false
	  region: us-east-1

 ```

 #### Fluentd Config Result
 ```

	<match **>
	    @type sqs
	    @id test_sqs
	    queue_name some-aws-sqs-queue
	    create_queue false
	    region us-east-1
	</match>

 ```

---
