---
title: Alibaba Cloud
weight: 200
generated_file: true
---

# Aliyun OSS plugin for Fluentd
## Overview
 **Fluent OSS output plugin** buffers event logs in local files and uploads them to OSS periodically in background threads.

 This plugin splits events by using the timestamp of event logs. For example,  a log '2019-04-09 message Hello' is reached, and then another log '2019-04-10 message World' is reached in this order, the former is stored in "20190409.gz" file, and latter in "20190410.gz" file.

 **Fluent OSS input plugin** reads data from OSS periodically.

 This plugin uses MNS on the same region of the OSS bucket. We must setup MNS and OSS event notification before using this plugin.

 [This document](https://help.aliyun.com/document_detail/52656.html) shows how to setup MNS and OSS event notification.

 This plugin will poll events from MNS queue and extract object keys from these events, and then will read those objects from OSS.
 More info at https://github.com/aliyun/fluent-plugin-oss

## Configuration
## Output Config

### endpoint (string, required) {#output config-endpoint}

OSS endpoint to connect to' 

Default: -

### bucket (string, required) {#output config-bucket}

Your bucket name 

Default: -

### access_key_id (*secret.Secret, required) {#output config-access_key_id}

Your access key id [Secret](../secret/) 

Default: -

### aaccess_key_secret (*secret.Secret, required) {#output config-aaccess_key_secret}

Your access secret key [Secret](../secret/) 

Default: -

### path (string, optional) {#output config-path}

Path prefix of the files on OSS  

Default:  fluent/logs

### upload_crc_enable (bool, optional) {#output config-upload_crc_enable}

Upload crc enabled  

Default:  true

### download_crc_enable (bool, optional) {#output config-download_crc_enable}

Download crc enabled  

Default:  true

### open_timeout (int, optional) {#output config-open_timeout}

Timeout for open connections  

Default:  10

### read_timeout (int, optional) {#output config-read_timeout}

Timeout for read response  

Default:  120

### oss_sdk_log_dir (string, optional) {#output config-oss_sdk_log_dir}

OSS SDK log directory  

Default:  /var/log/td-agent

### key_format (string, optional) {#output config-key_format}

The format of OSS object keys  

Default:  %{path}/%{time_slice}_%{index}_%{thread_id}.%{file_extension}

### store_as (string, optional) {#output config-store_as}

Archive format on OSS: gzip, json, text, lzo, lzma2  

Default:  gzip

### auto_create_bucket (bool, optional) {#output config-auto_create_bucket}

desc 'Create OSS bucket if it does not exists  

Default:  false

### overwrite (bool, optional) {#output config-overwrite}

Overwrite already existing path  

Default:  false

### check_bucket (bool, optional) {#output config-check_bucket}

Check bucket if exists or not  

Default:  true

### check_object (bool, optional) {#output config-check_object}

Check object before creation  

Default:  true

### hex_random_length (int, optional) {#output config-hex_random_length}

The length of `%{hex_random}` placeholder(4-16)  

Default:  4

### index_format (string, optional) {#output config-index_format}

`sprintf` format for `%{index}`  

Default:  %d

### warn_for_delay (string, optional) {#output config-warn_for_delay}

Given a threshold to treat events as delay, output warning logs if delayed events were put into OSS 

Default: -

### format (*Format, optional) {#output config-format}

[Format](../format/) 

Default: -

### buffer (*Buffer, optional) {#output config-buffer}

[Buffer](../buffer/) 

Default: -

### slow_flush_log_threshold (string, optional) {#output config-slow_flush_log_threshold}

The threshold for chunk flush performance check. Parameter type is float, not time, default: 20.0 (seconds) If chunk flush takes longer time than this threshold, fluentd logs warning message and increases metric fluentd_output_status_slow_flush_count. 

Default: -


