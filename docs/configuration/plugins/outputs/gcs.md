---
title: Google Cloud Storage
weight: 200
generated_file: true
---

# Google Cloud Storage
## Overview
 Store logs in Google Cloud Storage. For details, see [https://github.com/kube-logging/fluent-plugin-gcs](https://github.com/kube-logging/fluent-plugin-gcs).

 ## Example
 ```yaml
 spec:

	gcs:
	  project: logging-example
	  bucket: banzai-log-test
	  path: logs/${tag}/%Y/%m/%d/

 ```

## Configuration
## GCSOutput

### project (string, required) {#gcsoutput-project}

Project identifier for GCS 

Default: -

### keyfile (string, optional) {#gcsoutput-keyfile}

Path of GCS service account credentials JSON file 

Default: -

### credentials_json (*secret.Secret, optional) {#gcsoutput-credentials_json}

GCS service account credentials in JSON format [Secret](../secret/) 

Default: -

### client_retries (int, optional) {#gcsoutput-client_retries}

Number of times to retry requests on server error 

Default: -

### client_timeout (int, optional) {#gcsoutput-client_timeout}

Default timeout to use in requests 

Default: -

### bucket (string, required) {#gcsoutput-bucket}

Name of a GCS bucket 

Default: -

### object_key_format (string, optional) {#gcsoutput-object_key_format}

Format of GCS object keys  

Default:  %{path}%{time_slice}_%{index}.%{file_extension}

### path (string, optional) {#gcsoutput-path}

Path prefix of the files on GCS 

Default: -

### store_as (string, optional) {#gcsoutput-store_as}

Archive format on GCS: gzip json text  

Default:  gzip

### transcoding (bool, optional) {#gcsoutput-transcoding}

Enable the decompressive form of transcoding 

Default: -

### auto_create_bucket (bool, optional) {#gcsoutput-auto_create_bucket}

Create GCS bucket if it does not exists  

Default:  true

### hex_random_length (int, optional) {#gcsoutput-hex_random_length}

Max length of `%{hex_random}` placeholder(4-16)  

Default:  4

### overwrite (bool, optional) {#gcsoutput-overwrite}

Overwrite already existing path  

Default:  false

### acl (string, optional) {#gcsoutput-acl}

Permission for the object in GCS: auth_read owner_full owner_read private project_private public_read 

Default: -

### storage_class (string, optional) {#gcsoutput-storage_class}

Storage class of the file: dra nearline coldline multi_regional regional standard 

Default: -

### encryption_key (string, optional) {#gcsoutput-encryption_key}

Customer-supplied, AES-256 encryption key 

Default: -

### object_metadata ([]ObjectMetadata, optional) {#gcsoutput-object_metadata}

User provided web-safe keys and arbitrary string values that will returned with requests for the file as "x-goog-meta-" response headers. [Object Metadata](#objectmetadata) 

Default: -

### format (*Format, optional) {#gcsoutput-format}

[Format](../format/) 

Default: -

### buffer (*Buffer, optional) {#gcsoutput-buffer}

[Buffer](../buffer/) 

Default: -

### slow_flush_log_threshold (string, optional) {#gcsoutput-slow_flush_log_threshold}

The threshold for chunk flush performance check. Parameter type is float, not time, default: 20.0 (seconds) If chunk flush takes longer time than this threshold, fluentd logs warning message and increases metric fluentd_output_status_slow_flush_count. 

Default: -


## ObjectMetadata

### key (string, required) {#objectmetadata-key}

Key 

Default: -

### value (string, required) {#objectmetadata-value}

Value 

Default: -


