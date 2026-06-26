---
title: Raw
weight: 200
generated_file: true
---

# Raw
## Overview
 Configure custom or unexposed Fluentd filters via raw configuration. This allows you to specify any configuration that is not supported by the operator. The configuration should be in the format of a Fluentd filter configuration.

## Example `Raw` filter configurations

### Configure a custom filter via raw configuration

{{< highlight yaml >}}
apiVersion: logging.banzaicloud.io/v1beta1
kind: Flow
metadata:
  name: demo-flow
spec:
  filters:
    - raw:
        config: |
          @type my_filter
          <my_section>
            foo bar
        	tags ["web", "api", "db"]
          </my_section>
  selectors: {}
  localOutputRefs:
    - demo-output
{{</ highlight >}}

Fluentd Config Result

{{< highlight xml >}}
<filter **>
  @type my_filter
  @id test
  <my_section>
    foo bar
    tags ["web", "api", "db"]
  </my_section>
</filter>
{{</ highlight >}}

### Configure an unexposed filter via raw configuration

{{< highlight yaml >}}
apiVersion: logging.banzaicloud.io/v1beta1
kind: Flow
metadata:
  name: demo-flow
spec:
  filters:
    - raw:
        config: |
          @type ua_parser
          flatten
          key_name ua_string
  selectors: {}
  localOutputRefs:
    - demo-output
{{</ highlight >}}

Fluentd Config Result

{{< highlight xml >}}
<filter **>
  @type ua_parser
  @id test
  flatten
  key_name ua_string
</filter>
{{</ highlight >}}



## Configuration
## Raw

### config (string, optional) {#raw-config}

Raw configuration for the filter. 



