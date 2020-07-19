---
title: Format
weight: 200
generated_file: true
---

### Format
| Variable Name | Type | Required | Default | Description |
|---|---|---|---|---|
| type | string | No |  json | Output line formatting: out_file,json,ltsv,csv,msgpack,hash,single_value <br> |
| add_newline | *bool | No |  true | When type is single_value add '\n' to the end of the message <br> |
| message_key | string | No | - | When type is single_value specify the key holding information<br> |
