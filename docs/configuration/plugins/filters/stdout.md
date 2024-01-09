---
title: StdOut
weight: 200
generated_file: true
---

# [Stdout Filter](https://docs.fluentd.org/filter/stdout)
## Overview
 Fluentd Filter plugin to print events to stdout

## Configuration
## StdOutFilterConfig

### output_type (string, optional) {#stdoutfilterconfig-output_type}

This is the option of stdout format. 





## Example `StdOut` filter configurations

{{< highlight yaml >}}
apiVersion: logging.banzaicloud.io/v1beta1
kind: Flow
metadata:
  name: demo-flow
spec:
  filters:
    - stdout:
        output_type: json
  selectors: {}
  localOutputRefs:
    - demo-output
{{</ highlight >}}

Fluentd config result:

{{< highlight yaml >}}
<filter **>
  @type stdout
  @id test_stdout
  output_type json
</filter>
{{</ highlight >}}


---
