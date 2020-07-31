---
title: Amazon Elasticsearch
weight: 200
generated_file: true
---

# Amazon Elasticsearch output plugin for Fluentd
## Overview
  More info at https://github.com/atomita/fluent-plugin-aws-elasticsearch-service

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
### Amazon Elasticsearch
#### Send your logs to a Amazon Elasticsearch Service

| Variable Name | Type | Required | Default | Description |
|---|---|---|---|---|
| logstash_format | bool | No | - | logstash_format<br> |
| include_tag_key | bool | No | - | include_tag_key<br> |
| tag_key | string | No | - | tag_key<br> |
| flush_interval | string | No | - | flush_interval<br> |
| endpoint | *EndpointCredentials | No | - | AWS Endpoint Credentials<br> |
| format | *Format | No | - | [Format](../format/)<br> |
| buffer | *Buffer | No | - | [Buffer](../buffer/)<br> |
### Endpoint Credentials
#### endpoint

| Variable Name | Type | Required | Default | Description |
|---|---|---|---|---|
| region | string | No | - | AWS region. It should be in form like us-east-1, us-west-2. Default nil, which means try to find from environment variable AWS_REGION.<br> |
| url | string | No | - | AWS connection url.<br> |
| access_key_id | *secret.Secret | No | - | AWS access key id. This parameter is required when your agent is not running on EC2 instance with an IAM Role.<br> |
| secret_access_key | *secret.Secret | No | - | AWS secret key. This parameter is required when your agent is not running on EC2 instance with an IAM Role.<br> |
| assume_role_arn | *secret.Secret | No | - | Typically, you can use AssumeRole for cross-account access or federation.<br> |
| ecs_container_credentials_relative_uri | *secret.Secret | No | - | Set with AWS_CONTAINER_CREDENTIALS_RELATIVE_URI environment variable value<br> |
| assume_role_session_name | *secret.Secret | No | - | AssumeRoleWithWebIdentity https://docs.aws.amazon.com/STS/latest/APIReference/API_AssumeRoleWithWebIdentity.html<br> |
| assume_role_web_identity_token_file | *secret.Secret | No | - | AssumeRoleWithWebIdentity https://docs.aws.amazon.com/STS/latest/APIReference/API_AssumeRoleWithWebIdentity.html<br> |
| sts_credentials_region | *secret.Secret | No | - | By default, the AWS Security Token Service (AWS STS) is available as a global service, and all AWS STS requests go to a single endpoint at https://sts.amazonaws.com. AWS recommends using Regional AWS STS endpoints instead of the global endpoint to reduce latency, build in redundancy, and increase session token validity. https://docs.aws.amazon.com/IAM/latest/UserGuide/id_credentials_temp_enable-regions.html<br> |
