---
title: Amazon CloudWatch
weight: 200
generated_file: true
---

# CloudWatch output plugin for Fluentd
## Overview
 This plugin has been designed to output logs or metrics to Amazon CloudWatch.
 More info at [https://github.com/fluent-plugins-nursery/fluent-plugin-cloudwatch-logs](https://github.com/fluent-plugins-nursery/fluent-plugin-cloudwatch-logs).

 ## Example output configurations
 ```yaml
 spec:

	cloudwatch:
	  aws_key_id:
	    valueFrom:
	      secretKeyRef:
	        name: logging-s3
	        key: awsAccessKeyId
	  aws_sec_key:
	    valueFrom:
	      secretKeyRef:
	        name: logging-s3
	        key: awsSecretAccessKey
	  log_group_name: operator-log-group
	  log_stream_name: operator-log-stream
	  region: us-east-1
	  auto_create_stream true
	  buffer:
	    timekey: 30s
	    timekey_wait: 30s
	    timekey_use_utc: true

 ```

## Configuration
## Output Config

### auto_create_stream (bool, optional) {#output config-auto_create_stream}

Create log group and stream automatically.  

Default:  false

### aws_key_id (*secret.Secret, optional) {#output config-aws_key_id}

AWS access key id [Secret](../secret/) 

Default: -

### aws_sec_key (*secret.Secret, optional) {#output config-aws_sec_key}

AWS secret key. [Secret](../secret/) 

Default: -

### aws_instance_profile_credentials_retries (int, optional) {#output config-aws_instance_profile_credentials_retries}

Instance Profile Credentials call retries  

Default:  nil

### aws_use_sts (bool, optional) {#output config-aws_use_sts}

Enable AssumeRoleCredentials to authenticate, rather than the default credential hierarchy. See 'Cross-Account Operation' below for more detail. 

Default: -

### aws_sts_role_arn (string, optional) {#output config-aws_sts_role_arn}

The role ARN to assume when using cross-account sts authentication 

Default: -

### aws_sts_session_name (string, optional) {#output config-aws_sts_session_name}

The session name to use with sts authentication   

Default:  'fluentd'

### concurrency (int, optional) {#output config-concurrency}

Use to set the number of threads pushing data to CloudWatch.  

Default:  1

### endpoint (string, optional) {#output config-endpoint}

Use this parameter to connect to the local API endpoint (for testing) 

Default: -

### http_proxy (string, optional) {#output config-http_proxy}

Use to set an optional HTTP proxy 

Default: -

### include_time_key (bool, optional) {#output config-include_time_key}

Include time key as part of the log entry  

Default:  UTC

### json_handler (string, optional) {#output config-json_handler}

Name of the library to be used to handle JSON data. For now, supported libraries are json (default) and yajl 

Default: -

### localtime (bool, optional) {#output config-localtime}

Use localtime timezone for include_time_key output (overrides UTC default) 

Default: -

### log_group_aws_tags (string, optional) {#output config-log_group_aws_tags}

Set a hash with keys and values to tag the log group resource 

Default: -

### log_group_aws_tags_key (string, optional) {#output config-log_group_aws_tags_key}

Specified field of records as AWS tags for the log group 

Default: -

### log_group_name (string, optional) {#output config-log_group_name}

Name of log group to store logs 

Default: -

### log_group_name_key (string, optional) {#output config-log_group_name_key}

Specified field of records as log group name 

Default: -

### log_rejected_request (string, optional) {#output config-log_rejected_request}

Output rejected_log_events_info request log.  

Default:  false

### log_stream_name (string, optional) {#output config-log_stream_name}

Name of log stream to store logs 

Default: -

### log_stream_name_key (string, optional) {#output config-log_stream_name_key}

Specified field of records as log stream name 

Default: -

### max_events_per_batch (int, optional) {#output config-max_events_per_batch}

Maximum number of events to send at once  

Default:  10000

### max_message_length (int, optional) {#output config-max_message_length}

Maximum length of the message 

Default: -

### message_keys (string, optional) {#output config-message_keys}

Keys to send messages as events 

Default: -

### put_log_events_disable_retry_limit (bool, optional) {#output config-put_log_events_disable_retry_limit}

If true, put_log_events_retry_limit will be ignored 

Default: -

### put_log_events_retry_limit (int, optional) {#output config-put_log_events_retry_limit}

Maximum count of retry (if exceeding this, the events will be discarded) 

Default: -

### put_log_events_retry_wait (string, optional) {#output config-put_log_events_retry_wait}

Time before retrying PutLogEvents (retry interval increases exponentially like put_log_events_retry_wait * (2 ^ retry_count)) 

Default: -

### region (string, required) {#output config-region}

AWS Region 

Default: -

### remove_log_group_aws_tags_key (string, optional) {#output config-remove_log_group_aws_tags_key}

Remove field specified by log_group_aws_tags_key 

Default: -

### remove_log_group_name_key (string, optional) {#output config-remove_log_group_name_key}

Remove field specified by log_group_name_key 

Default: -

### remove_log_stream_name_key (string, optional) {#output config-remove_log_stream_name_key}

Remove field specified by log_stream_name_key 

Default: -

### remove_retention_in_days (string, optional) {#output config-remove_retention_in_days}

Remove field specified by retention_in_days 

Default: -

### retention_in_days (string, optional) {#output config-retention_in_days}

Use to set the expiry time for log group when created with auto_create_stream. (default to no expiry) 

Default: -

### retention_in_days_key (string, optional) {#output config-retention_in_days_key}

Use specified field of records as retention period 

Default: -

### use_tag_as_group (bool, optional) {#output config-use_tag_as_group}

Use tag as a group name 

Default: -

### use_tag_as_stream (bool, optional) {#output config-use_tag_as_stream}

Use tag as a stream name 

Default: -

### buffer (*Buffer, optional) {#output config-buffer}

[Buffer](../buffer/) 

Default: -

### slow_flush_log_threshold (string, optional) {#output config-slow_flush_log_threshold}

The threshold for chunk flush performance check. Parameter type is float, not time, default: 20.0 (seconds) If chunk flush takes longer time than this threshold, fluentd logs warning message and increases metric fluentd_output_status_slow_flush_count. 

Default: -

### format (*Format, optional) {#output config-format}

[Format](../format/) 

Default: -


