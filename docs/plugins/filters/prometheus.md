---
title: Prometheus
weight: 200
---

# [Prometheus Filter](https://github.com/fluent/fluent-plugin-prometheus#prometheus-outputfilter-plugin)
## Overview
 Prometheus Filter Plugin to count Incoming Records

## Configuration
### PrometheusConfig
| Variable Name | Type | Required | Default | Description |
|---|---|---|---|---|
| metrics | []MetricSection | No | - | [Metrics Section](#Metrics-Section)<br> |
| labels | Label | No | - |  |
### Metrics Section
| Variable Name | Type | Required | Default | Description |
|---|---|---|---|---|
| name | string | Yes | - | Metrics name<br> |
| type | string | Yes | - | Metrics type [counter](https://github.com/fluent/fluent-plugin-prometheus#counter-type), [gauge](https://github.com/fluent/fluent-plugin-prometheus#gauge-type), [summary](https://github.com/fluent/fluent-plugin-prometheus#summary-type), [histogram](https://github.com/fluent/fluent-plugin-prometheus#histogram-type)<br> |
| desc | string | Yes | - | Description of metric<br> |
| key | string | No | - | Key name of record for instrumentation.<br> |
| buckets | string | No | - | Buckets of record for instrumentation<br> |
| labels | Label | No | - | Additional labels for this metric<br> |
 #### Example `Prometheus` filter configurations
 ```
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
  outputRefs:
    - demo-output
 ```

 #### Fluentd Config Result
 ```
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
 ```

---
