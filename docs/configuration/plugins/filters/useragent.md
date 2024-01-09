---
title: User Agent
weight: 200
generated_file: true
---

# Fluentd UserAgent filter
## Overview
 Fluentd Filter plugin to parse user-agent
 More information at https://github.com/bungoume/fluent-plugin-ua-parser

## Configuration
## UserAgent

### delete_key (bool, optional) {#useragent-delete_key}

Delete input key

Default: false

### flatten (bool, optional) {#useragent-flatten}

Join hashed data by '_'

Default: false

### key_name (string, optional) {#useragent-key_name}

Target key name

Default: user_agent

### out_key (string, optional) {#useragent-out_key}

Output prefix key name

Default: ua




## Example `UserAgent` filter configurations

{{< highlight yaml >}}
apiVersion: logging.banzaicloud.io/v1beta1
kind: Flow
metadata:
  name: demo-flow
spec:
  filters:
    - useragent:
        key_name: my_agent
        delete_key: true
        out_key: ua_fields
        flatten: true
  selectors: {}
  localOutputRefs:
    - demo-output
{{</ highlight >}}

Fluentd config result:

{{< highlight xml >}}
<filter **>
  @type ua_parser
  @id test_useragent
  key_name my_agent
  delete_key true
  out_key ua_fields
  flatten true
</filter>
{{</ highlight >}}


---
