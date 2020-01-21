# [Stdout Filter](https://docs.fluentd.org/filter/stdout)
## Overview
 Fluentd Filter plugin to print events to stdout

## Configuration
### StdOutFilterConfig
| Variable Name | Type | Required | Default | Description |
|---|---|---|---|---|
| output_type | string | No | - | This is the option of stdout format.<br> |
 #### Example `StdOut` filter configurations
 ```yaml
apiVersion: logging.banzaicloud.io/v1beta1
kind: Flow
metadata:
  name: demo-flow
spec:
  filters:
    - stdout:
        output_type: json
  selectors: {}
  outputRefs:
    - demo-output
 ```

 #### Fluentd Config Result
 ```yaml
<filter **>
  @type stdout
  @id test_stdout
  output_type json
</filter>
 ```

---
