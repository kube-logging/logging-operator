---
title: Amazon Kinesis
weight: 200
generated_file: true
---

# Kinesis Firehose output plugin for Fluentd
## Overview


For details, see [https://github.com/awslabs/aws-fluent-plugin-kinesis#configuration-kinesis_firehose](https://github.com/awslabs/aws-fluent-plugin-kinesis#configuration-kinesis_firehose).

## Example output configurations

```yaml
spec:
  kinesisFirehose:
    delivery_stream_name: example-stream-name
    region: us-east-1
    format:
      type: json
```


## Configuration
## KinesisFirehose

Send your logs to a Kinesis Firehose

### aws_iam_retries (int, optional) {#kinesisfirehose-aws_iam_retries}

The number of attempts to make (with exponential backoff) when loading instance profile credentials from the EC2 metadata service using an IAM role. Defaults to 5 retries. 


### aws_key_id (*secret.Secret, optional) {#kinesisfirehose-aws_key_id}

AWS access key id. This parameter is required when your agent is not running on EC2 instance with an IAM Role. 


### aws_sec_key (*secret.Secret, optional) {#kinesisfirehose-aws_sec_key}

AWS secret key. This parameter is required when your agent is not running on EC2 instance with an IAM Role. 


### aws_ses_token (*secret.Secret, optional) {#kinesisfirehose-aws_ses_token}

AWS session token. This parameter is optional, but can be provided if using MFA or temporary credentials when your agent is not running on EC2 instance with an IAM Role. 


### append_new_line (*bool, optional) {#kinesisfirehose-append_new_line}

If it is enabled, the plugin adds new line character (\n) to each serialized record. Before appending \n, plugin calls chomp and removes separator from the end of each record as chomp_record is true. Therefore, you don't need to enable chomp_record option when you use kinesis_firehose output with default configuration (append_new_line is true). If you want to set append_new_line false, you can choose chomp_record false (default) or true (compatible format with plugin v2). (Default:true) 


### assume_role_credentials (*KinesisFirehoseAssumeRoleCredentials, optional) {#kinesisfirehose-assume_role_credentials}

Typically, you can use AssumeRole for cross-account access or federation. 


### batch_request_max_count (int, optional) {#kinesisfirehose-batch_request_max_count}

Integer, default 500. The number of max count of making batch request from record chunk. It can't exceed the default value because it's API limit. 


### batch_request_max_size (int, optional) {#kinesisfirehose-batch_request_max_size}

Integer. The number of max size of making batch request from record chunk. It can't exceed the default value because it's API limit. 


### buffer (*Buffer, optional) {#kinesisfirehose-buffer}

[Buffer](../buffer/) 


### delivery_stream_name (string, required) {#kinesisfirehose-delivery_stream_name}

Name of the delivery stream to put data. 


### format (*Format, optional) {#kinesisfirehose-format}

[Format](../format/) 


### process_credentials (*KinesisFirehoseProcessCredentials, optional) {#kinesisfirehose-process_credentials}

This loads AWS access credentials from an external process. 


### region (string, optional) {#kinesisfirehose-region}

AWS region of your stream. It should be in form like us-east-1, us-west-2. Default nil, which means try to find from environment variable AWS_REGION. 


### reset_backoff_if_success (bool, optional) {#kinesisfirehose-reset_backoff_if_success}

Boolean, default true. If enabled, when after retrying, the next retrying checks the number of succeeded records on the former batch request and reset exponential backoff if there is any success. Because batch request could be composed by requests across shards, simple exponential backoff for the batch request wouldn't work some cases. 


### retries_on_batch_request (int, optional) {#kinesisfirehose-retries_on_batch_request}

The plugin will put multiple records to Amazon Kinesis Data Streams in batches using PutRecords. A set of records in a batch may fail for reasons documented in the Kinesis Service API Reference for PutRecords. Failed records will be retried retries_on_batch_request times 


### slow_flush_log_threshold (string, optional) {#kinesisfirehose-slow_flush_log_threshold}

The threshold for chunk flush performance check. Parameter type is float, not time, default: 20.0 (seconds) If chunk flush takes longer time than this threshold, fluentd logs warning message and increases metric fluentd_output_status_slow_flush_count. 



## Assume Role Credentials

assume_role_credentials

### duration_seconds (string, optional) {#assume role credentials-duration_seconds}

The duration, in seconds, of the role session (900-3600) 


### external_id (string, optional) {#assume role credentials-external_id}

A unique identifier that is used by third parties when assuming roles in their customers' accounts. 


### policy (string, optional) {#assume role credentials-policy}

An IAM policy in JSON format 


### role_arn (string, required) {#assume role credentials-role_arn}

The Amazon Resource Name (ARN) of the role to assume 


### role_session_name (string, required) {#assume role credentials-role_session_name}

An identifier for the assumed role session 



## Process Credentials

process_credentials

### process (string, required) {#process credentials-process}

Command more info: https://docs.aws.amazon.com/sdk-for-ruby/v3/api/Aws/ProcessCredentials.html 



