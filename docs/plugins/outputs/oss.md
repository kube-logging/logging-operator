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
### Output Config
| Variable Name | Type | Required | Default | Description |
|---|---|---|---|---|
| endpoint | string | Yes | - | OSS endpoint to connect to'<br> |
| bucket | string | Yes | - | Your bucket name<br> |
| access_key_id | *secret.Secret | Yes | - | Your access key id<br>[Secret](../secret/)<br> |
| aaccess_key_secret | *secret.Secret | Yes | - | Your access secret key<br>[Secret](../secret/)<br> |
| path | string | No |  fluent/logs | Path prefix of the files on OSS <br> |
| upload_crc_enable | bool | No |  true | Upload crc enabled <br> |
| download_crc_enable | bool | No |  true | Download crc enabled <br> |
| open_timeout | int | No |  10 | Timeout for open connections <br> |
| read_timeout | int | No |  120 | Timeout for read response <br> |
| oss_sdk_log_dir | string | No |  /var/log/td-agent | OSS SDK log directory <br> |
| key_format | string | No |  %{path}/%{time_slice}_%{index}_%{thread_id}.%{file_extension} | The format of OSS object keys <br> |
| store_as | string | No |  gzip | Archive format on OSS: gzip, json, text, lzo, lzma2 <br> |
| auto_create_bucket | bool | No |  false | desc 'Create OSS bucket if it does not exists <br> |
| overwrite | bool | No |  false | Overwrite already existing path <br> |
| check_bucket | bool | No |  true | Check bucket if exists or not <br> |
| check_object | bool | No |  true | Check object before creation <br> |
| hex_random_length | int | No |  4 | The length of `%{hex_random}` placeholder(4-16) <br> |
| index_format | string | No |  %d | `sprintf` format for `%{index}` <br> |
| warn_for_delay | string | No | - | Given a threshold to treat events as delay, output warning logs if delayed events were put into OSS<br> |
| format | *Format | No | - | [Format](../format/)<br> |
| buffer | *Buffer | No | - | [Buffer](../buffer/)<br> |
