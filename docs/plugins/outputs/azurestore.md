---
title: Azure Storage
weight: 200
generated_file: true
---

# Azure Storage output plugin for Fluentd
## Overview
Azure Storage output plugin buffers logs in local file and upload them to Azure Storage periodically.
More info at https://github.com/htgc/fluent-plugin-azurestorage

## Configuration
### Output Config
| Variable Name | Type | Required | Default | Description |
|---|---|---|---|---|
| path | string | No | - | Path prefix of the files on Azure<br> |
| azure_storage_account | *secret.Secret | Yes | - | Your azure storage account<br>[Secret](../secret/)<br> |
| azure_storage_access_key | *secret.Secret | Yes | - | Your azure storage access key<br>[Secret](../secret/)<br> |
| azure_container | string | Yes | - | Your azure storage container<br> |
| azure_storage_type | string | No |  blob | Azure storage type currently only "blob" supported <br> |
| azure_object_key_format | string | No |  %{path}%{time_slice}_%{index}.%{file_extension} | Object key format <br> |
| store_as | string | No |  gzip | Store as: gzip, json, text, lzo, lzma2 <br> |
| auto_create_container | bool | No |  true | Automatically create container if not exists<br> |
| format | string | No | json | Compat format type: out_file, json, ltsv (default: out_file)<br> |
| buffer | *Buffer | No | - | [Buffer](../buffer/)<br> |
