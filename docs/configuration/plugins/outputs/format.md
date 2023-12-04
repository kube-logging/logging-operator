---
title: Format
weight: 200
generated_file: true
---

# Format output records
## Overview

Specify how to format output records. For details, see [https://docs.fluentd.org/configuration/format-section](https://docs.fluentd.org/configuration/format-section).

## Example

```yaml
spec:
  format:
    path: /tmp/logs/${tag}/%Y/%m/%d.%H.%M
    format:
      type: single_value
      add_newline: true
      message_key: msg
```


## Configuration
## Format

### add_newline (*bool, optional) {#format-add_newline}

When type is single_value add '\n' to the end of the message

Default: true

### message_key (string, optional) {#format-message_key}

When type is single_value specify the key holding information 


### type (string, optional) {#format-type}

Output line formatting: out_file,json,ltsv,csv,msgpack,hash,single_value

Default: json


