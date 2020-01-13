# CloudWatch output plugin for Fluentd
## Overview
This plugin has been designed to output logs or metrics to Amazon CloudWatch.
More info at https://github.com/fluent-plugins-nursery/fluent-plugin-cloudwatch-logs

 #### Example output configurations
 ```
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
          key: awsSecretAccesKey
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
### Output Config
| Variable Name | Type | Required | Default | Description |
|---|---|---|---|---|
| auto_create_stream | bool | No |  false | Create log group and stream automatically. <br> |
| aws_key_id | *secret.Secret | No | - | AWS access key id<br>[Secret](./secret.md)<br> |
| aws_sec_key | *secret.Secret | No | - | AWS secret key.<br>[Secret](./secret.md)<br> |
| aws_instance_profile_credentials_retries | int | No |  nil | Instance Profile Credentials call retries <br> |
| aws_use_sts | bool | No | - | Enable AssumeRoleCredentials to authenticate, rather than the default credential hierarchy. See 'Cross-Account Operation' below for more detail.<br> |
| aws_sts_role_arn | string | No | - | The role ARN to assume when using cross-account sts authentication<br> |
| aws_sts_session_name | string | No |  'fluentd' | The session name to use with sts authentication  <br> |
| concurrency | int | No |  1 | Use to set the number of threads pushing data to CloudWatch. <br> |
| endpoint | string | No | - | Use this parameter to connect to the local API endpoint (for testing)<br> |
| http_proxy | string | No | - | Use to set an optional HTTP proxy<br> |
| include_time_key | bool | No |  UTC | Include time key as part of the log entry <br> |
| json_handler | string | No | - | Name of the library to be used to handle JSON data. For now, supported libraries are json (default) and yajl<br> |
| localtime | bool | No | - | Use localtime timezone for include_time_key output (overrides UTC default)<br> |
| log_group_aws_tags | string | No | - | Set a hash with keys and values to tag the log group resource<br> |
| log_group_aws_tags_key | string | No | - | Specified field of records as AWS tags for the log group<br> |
| log_group_name | string | No | - | Name of log group to store logs<br> |
| log_group_name_key | string | No | - | Specified field of records as log group name<br> |
| log_rejected_request | string | No |  false | Output rejected_log_events_info request log. <br> |
| log_stream_name | string | No | - | Name of log stream to store logs<br> |
| log_stream_name_key | string | No | - | Specified field of records as log stream name<br> |
| max_events_per_batch | int | No |  10000 | Maximum number of events to send at once <br> |
| max_message_length | int | No | - | Maximum length of the message<br> |
| message_keys | string | No | - | Keys to send messages as events<br> |
| put_log_events_disable_retry_limit | bool | No | - | If true, put_log_events_retry_limit will be ignored<br> |
| put_log_events_retry_limit | int | No | - | Maximum count of retry (if exceeding this, the events will be discarded)<br> |
| put_log_events_retry_wait | string | No | - | Time before retrying PutLogEvents (retry interval increases exponentially like put_log_events_retry_wait * (2 ^ retry_count))<br> |
| region | string | Yes | - | AWS Region<br> |
| remove_log_group_aws_tags_key | string | No | - | Remove field specified by log_group_aws_tags_key<br> |
| remove_log_group_name_key | string | No | - | Remove field specified by log_group_name_key<br> |
| remove_log_stream_name_key | string | No | - | Remove field specified by log_stream_name_key<br> |
| remove_retention_in_days | string | No | - | Remove field specified by retention_in_days<br> |
| retention_in_days | string | No | - | Use to set the expiry time for log group when created with auto_create_stream. (default to no expiry)<br> |
| retention_in_days_key | string | No | - | Use specified field of records as retention period<br> |
| use_tag_as_group | bool | No | - | Use tag as a group name<br> |
| use_tag_as_stream | bool | No | - | Use tag as a stream name<br> |
| buffer | *Buffer | No | - | [Buffer](./buffer.md)<br> |
| format | *Format | No | - | [Format](./format.md)<br> |
