---
title: Amazon Kinesis
weight: 200
generated_file: true
---

# Kinesis Stream output plugin for Fluentd
## Overview
  More info at https://github.com/awslabs/aws-fluent-plugin-kinesis#configuration-kinesis_streams

 #### Example output configurations
 ```
 spec:
   kinesisStream:
     stream_name: example-stream-name
     region: us-east-1
     format:
       type: json
 ```

## Configuration
### KinesisStream
#### Send your logs to a Kinesis Stream

| Variable Name | Type | Required | Default | Description |
|---|---|---|---|---|
| stream_name | string | Yes | - | Name of the stream to put data.<br> |
| partition_key | string | No | - | A key to extract partition key from JSON object. Default nil, which means partition key will be generated randomly.<br> |
| aws_key_id | *secret.Secret | No | - | AWS access key id. This parameter is required when your agent is not running on EC2 instance with an IAM Role.<br> |
| aws_sec_key | *secret.Secret | No | - | AWS secret key. This parameter is required when your agent is not running on EC2 instance with an IAM Role.<br> |
| aws_ses_token | *secret.Secret | No | - | AWS session token. This parameter is optional, but can be provided if using MFA or temporary credentials when your agent is not running on EC2 instance with an IAM Role.<br> |
| aws_iam_retries | int | No | - | The number of attempts to make (with exponential backoff) when loading instance profile credentials from the EC2 metadata service using an IAM role. Defaults to 5 retries.<br> |
| assume_role_credentials | *KinesisStreamAssumeRoleCredentials | No | - | Typically, you can use AssumeRole for cross-account access or federation.<br> |
| process_credentials | *KinesisStreamProcessCredentials | No | - | This loads AWS access credentials from an external process.<br> |
| region | string | No | - | AWS region of your stream. It should be in form like us-east-1, us-west-2. Default nil, which means try to find from environment variable AWS_REGION.<br> |
| retries_on_batch_request | int | No | - | The plugin will put multiple records to Amazon Kinesis Data Streams in batches using PutRecords. A set of records in a batch may fail for reasons documented in the Kinesis Service API Reference for PutRecords. Failed records will be retried retries_on_batch_request times<br> |
| reset_backoff_if_success | bool | No | - | Boolean, default true. If enabled, when after retrying, the next retrying checks the number of succeeded records on the former batch request and reset exponential backoff if there is any success. Because batch request could be composed by requests across shards, simple exponential backoff for the batch request wouldn't work some cases.<br> |
| batch_request_max_count | int | No | - | Integer, default 500. The number of max count of making batch request from record chunk. It can't exceed the default value because it's API limit.<br> |
| batch_request_max_size | int | No | - | Integer. The number of max size of making batch request from record chunk. It can't exceed the default value because it's API limit.<br> |
| format | *Format | No | - | [Format](../format/)<br> |
| buffer | *Buffer | No | - | [Buffer](../buffer/)<br> |
### Assume Role Credentials
#### assume_role_credentials

| Variable Name | Type | Required | Default | Description |
|---|---|---|---|---|
| role_arn | string | Yes | - | The Amazon Resource Name (ARN) of the role to assume<br> |
| role_session_name | string | Yes | - | An identifier for the assumed role session<br> |
| policy | string | No | - | An IAM policy in JSON format<br> |
| duration_seconds | string | No | - | The duration, in seconds, of the role session (900-3600)<br> |
| external_id | string | No | - | A unique identifier that is used by third parties when assuming roles in their customers' accounts.<br> |
### Process Credentials
#### process_credentials

| Variable Name | Type | Required | Default | Description |
|---|---|---|---|---|
| process | string | Yes | - | Command more info: https://docs.aws.amazon.com/sdk-for-ruby/v3/api/Aws/ProcessCredentials.html<br> |
