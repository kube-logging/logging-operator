---
title: Amazon Kinesis
weight: 200
generated_file: true
---

# Kinesis Stream output plugin for Fluentd
## Overview


For details, see [https://github.com/awslabs/aws-fluent-plugin-kinesis#configuration-kinesis_streams](https://github.com/awslabs/aws-fluent-plugin-kinesis#configuration-kinesis_streams).

## Example output configurations

```yaml
spec:
  kinesisStream:
    stream_name: example-stream-name
    region: us-east-1
    format:
      type: json
```


## Configuration
## KinesisStream

Send your logs to a Kinesis Stream

### aws_iam_retries (int, optional) {#kinesisstream-aws_iam_retries}

The number of attempts to make (with exponential backoff) when loading instance profile credentials from the EC2 metadata service using an IAM role. Defaults to 5 retries. 


### aws_key_id (*secret.Secret, optional) {#kinesisstream-aws_key_id}

AWS access key id. This parameter is required when your agent is not running on EC2 instance with an IAM Role. 


### aws_sec_key (*secret.Secret, optional) {#kinesisstream-aws_sec_key}

AWS secret key. This parameter is required when your agent is not running on EC2 instance with an IAM Role. 


### aws_ses_token (*secret.Secret, optional) {#kinesisstream-aws_ses_token}

AWS session token. This parameter is optional, but can be provided if using MFA or temporary credentials when your agent is not running on EC2 instance with an IAM Role. 


### assume_role_credentials (*KinesisStreamAssumeRoleCredentials, optional) {#kinesisstream-assume_role_credentials}

Typically, you can use AssumeRole for cross-account access or federation. 


### batch_request_max_count (int, optional) {#kinesisstream-batch_request_max_count}

Integer, default 500. The number of max count of making batch request from record chunk. It can't exceed the default value because it's API limit. 


### batch_request_max_size (int, optional) {#kinesisstream-batch_request_max_size}

Integer. The number of max size of making batch request from record chunk. It can't exceed the default value because it's API limit. 


### buffer (*Buffer, optional) {#kinesisstream-buffer}

[Buffer](../buffer/) 


### format (*Format, optional) {#kinesisstream-format}

[Format](../format/) 


### partition_key (string, optional) {#kinesisstream-partition_key}

A key to extract partition key from JSON object. Default nil, which means partition key will be generated randomly. 


### process_credentials (*KinesisStreamProcessCredentials, optional) {#kinesisstream-process_credentials}

This loads AWS access credentials from an external process. 


### region (string, optional) {#kinesisstream-region}

AWS region of your stream. It should be in form like us-east-1, us-west-2. Default nil, which means try to find from environment variable AWS_REGION. 


### reset_backoff_if_success (bool, optional) {#kinesisstream-reset_backoff_if_success}

Boolean, default true. If enabled, when after retrying, the next retrying checks the number of succeeded records on the former batch request and reset exponential backoff if there is any success. Because batch request could be composed by requests across shards, simple exponential backoff for the batch request wouldn't work some cases. 


### retries_on_batch_request (int, optional) {#kinesisstream-retries_on_batch_request}

The plugin will put multiple records to Amazon Kinesis Data Streams in batches using PutRecords. A set of records in a batch may fail for reasons documented in the Kinesis Service API Reference for PutRecords. Failed records will be retried retries_on_batch_request times 


### slow_flush_log_threshold (string, optional) {#kinesisstream-slow_flush_log_threshold}

The threshold for chunk flush performance check. Parameter type is float, not time, default: 20.0 (seconds) If chunk flush takes longer time than this threshold, fluentd logs warning message and increases metric fluentd_output_status_slow_flush_count. 


### stream_name (string, required) {#kinesisstream-stream_name}

Name of the stream to put data. 



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



