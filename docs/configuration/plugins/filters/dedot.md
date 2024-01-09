---
title: Dedot
weight: 200
generated_file: true
---

# [Dedot Filter](https://github.com/lunardial/fluent-plugin-dedot_filter)
## Overview
 Fluentd Filter plugin to de-dot field name for elasticsearch.

## Configuration
## DedotFilterConfig

### de_dot_nested (bool, optional) {#dedotfilterconfig-de_dot_nested}

Will cause the plugin to recourse through nested structures (hashes and arrays), and remove dots in those key-names too.

Default: false

### de_dot_separator (string, optional) {#dedotfilterconfig-de_dot_separator}

Separator

Default: _




## Example `Dedot` filter configurations

{{< highlight yaml >}}
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
  localOutputRefs:
    - demo-output
{{</ highlight >}}

Fluentd config result:

{{< highlight xml >}}
<filter **>
  @type dedot
  @id test_dedot
  de_dot_nested true
  de_dot_separator -
</filter>
{{</ highlight >}}


---
