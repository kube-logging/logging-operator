---
title: Amazon Elasticsearch
weight: 200
generated_file: true
---

# Amazon Elasticsearch output plugin for Fluentd
## Overview

For details, see https://github.com/atomita/fluent-plugin-aws-elasticsearch-service

## Example output configurations
{{< highlight yaml >}}
spec:
  awsElasticsearch:
    logstash_format: true
    include_tag_key: true
    tag_key: "@log_name"
    flush_interval: 1s
    endpoint:
      url: https://CLUSTER_ENDPOINT_URL
      region: eu-west-1
      access_key_id:
        value: aws-key
      secret_access_key:
        value: aws_secret
{{</ highlight >}}


## Configuration
## Amazon Elasticsearch

Send your logs to a Amazon Elasticsearch Service

###  (*ElasticsearchOutput, optional) {#amazon elasticsearch-}

ElasticSearch 


### buffer (*Buffer, optional) {#amazon elasticsearch-buffer}

[Buffer](../buffer/) 


### endpoint (*EndpointCredentials, optional) {#amazon elasticsearch-endpoint}

AWS Endpoint Credentials 


### flush_interval (string, optional) {#amazon elasticsearch-flush_interval}

flush_interval 


### format (*Format, optional) {#amazon elasticsearch-format}

[Format](../format/) 



## Endpoint Credentials

endpoint

### access_key_id (*secret.Secret, optional) {#endpoint credentials-access_key_id}

AWS access key id. This parameter is required when your agent is not running on EC2 instance with an IAM Role. 


### assume_role_arn (*secret.Secret, optional) {#endpoint credentials-assume_role_arn}

Typically, you can use AssumeRole for cross-account access or federation. 


### assume_role_session_name (*secret.Secret, optional) {#endpoint credentials-assume_role_session_name}

AssumeRoleWithWebIdentity https://docs.aws.amazon.com/STS/latest/APIReference/API_AssumeRoleWithWebIdentity.html 


### assume_role_web_identity_token_file (*secret.Secret, optional) {#endpoint credentials-assume_role_web_identity_token_file}

AssumeRoleWithWebIdentity https://docs.aws.amazon.com/STS/latest/APIReference/API_AssumeRoleWithWebIdentity.html 


### ecs_container_credentials_relative_uri (*secret.Secret, optional) {#endpoint credentials-ecs_container_credentials_relative_uri}

Set with AWS_CONTAINER_CREDENTIALS_RELATIVE_URI environment variable value 


### region (string, optional) {#endpoint credentials-region}

AWS region. It should be in form like us-east-1, us-west-2. Default nil, which means try to find from environment variable AWS_REGION. 


### secret_access_key (*secret.Secret, optional) {#endpoint credentials-secret_access_key}

AWS secret key. This parameter is required when your agent is not running on EC2 instance with an IAM Role. 


### sts_credentials_region (*secret.Secret, optional) {#endpoint credentials-sts_credentials_region}

By default, the AWS Security Token Service (AWS STS) is available as a global service, and all AWS STS requests go to a single endpoint at https://sts.amazonaws.com. AWS recommends using Regional AWS STS endpoints instead of the global endpoint to reduce latency, build in redundancy, and increase session token validity. https://docs.aws.amazon.com/IAM/latest/UserGuide/id_credentials_temp_enable-regions.html 


### url (string, optional) {#endpoint credentials-url}

AWS connection url. 



