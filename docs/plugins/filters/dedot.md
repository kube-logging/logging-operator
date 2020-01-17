# [Dedot Filter](https://github.com/lunardial/fluent-plugin-dedot_filter)
## Overview
 Fluentd Filter plugin to de-dot field name for elasticsearch.

## Configuration
### DedotFilterConfig
| Variable Name | Type | Required | Default | Description |
|---|---|---|---|---|
| de_dot_nested | bool | No |  false | Will cause the plugin to recurse through nested structures (hashes and arrays), and remove dots in those key-names too.<br> |
| de_dot_separator | string | No | _ | Separator <br> |
 #### Example `Dedot` filter configurations
 ```yaml
apiVersion: logging.banzaicloud.io/v1beta1
kind: Flow
metadata:
  name: demo-flow
spec:
  filters:
    - dedot:
        de_dot_separator: "-"
        de_dot_nested: true
  selectors: {}
  outputRefs:
    - demo-output
 ```

 #### Fluentd Config Result
 ```yaml
<filter **>
  @type dedot
  @id test_dedot
  de_dot_nested true
  de_dot_separator -
</filter>
 ```

---
