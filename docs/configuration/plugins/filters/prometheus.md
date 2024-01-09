---
title: Prometheus
weight: 200
generated_file: true
---

# [Prometheus Filter](https://github.com/fluent/fluent-plugin-prometheus#prometheus-outputfilter-plugin)
## Overview
 Prometheus Filter Plugin to count Incoming Records

## Configuration
## PrometheusConfig

### labels (Label, optional) {#prometheusconfig-labels}


### metrics ([]MetricSection, optional) {#prometheusconfig-metrics}

[Metrics Section](#metrics-section) 



## Metrics Section

### buckets (string, optional) {#metrics section-buckets}

Buckets of record for instrumentation 


### desc (string, required) {#metrics section-desc}

Description of metric 


### key (string, optional) {#metrics section-key}

Key name of record for instrumentation. 


### labels (Label, optional) {#metrics section-labels}

Additional labels for this metric 


### name (string, required) {#metrics section-name}

Metrics name 


### type (string, required) {#metrics section-type}

Metrics type [counter](https://github.com/fluent/fluent-plugin-prometheus#counter-type), [gauge](https://github.com/fluent/fluent-plugin-prometheus#gauge-type), [summary](https://github.com/fluent/fluent-plugin-prometheus#summary-type), [histogram](https://github.com/fluent/fluent-plugin-prometheus#histogram-type) 





## Example `Prometheus` filter configurations

{{< highlight yaml >}}
apiVersion: logging.banzaicloud.io/v1beta1
kind: Flow
metadata:
  name: demo-flow
spec:
  filters:
    - tag_normaliser: {}
    - parser:
        remove_key_name_field: true
        reserve_data: true
        parse:
          type: nginx
    - prometheus:
        metrics:
        - name: total_counter
          desc: The total number of foo in message.
          type: counter
          labels:
            foo: bar
        labels:
          host: ${hostname}
          tag: ${tag}
          namespace: $.kubernetes.namespace
  selectors: {}
  localOutputRefs:
    - demo-output
{{</ highlight >}}

Fluentd config result:

{{< highlight xml>}}
  <filter **>
    @type prometheus
    @id logging-demo-flow_2_prometheus
    <metric>
      desc The total number of foo in message.
      name total_counter
      type counter
      <labels>
        foo bar
      </labels>
    </metric>
    <labels>
      host ${hostname}
      namespace $.kubernetes.namespace
      tag ${tag}
    </labels>
  </filter>
{{</ highlight >}}


---
