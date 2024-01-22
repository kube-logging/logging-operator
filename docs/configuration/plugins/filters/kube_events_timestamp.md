---
title: Kubernetes Events Timestamp
weight: 200
generated_file: true
---

# [Kubernetes Events Timestamp Filter](https://github.com/kube-logging/fluentd-filter-kube-events-timestamp)
## Overview
 Fluentd Filter plugin to select particular timestamp into an additional field

## Configuration
## KubeEventsTimestampConfig

### mapped_time_key (string, optional) {#kubeeventstimestampconfig-mapped_time_key}

Added time field name

Default: triggerts

### timestamp_fields ([]string, optional) {#kubeeventstimestampconfig-timestamp_fields}

Time field names in order of relevance

Default: event.eventTime, event.lastTimestamp, event.firstTimestamp




## Example `Kubernetes Events Timestamp` filter configurations

{{< highlight yaml >}}
apiVersion: logging.banzaicloud.io/v1beta1
kind: Flow
metadata:
  name: es-flow
spec:
  filters:
    - kube_events_timestamp:
        timestamp_fields:
          - "event.eventTime"
          - "event.lastTimestamp"
          - "event.firstTimestamp"
        mapped_time_key: mytimefield
  selectors: {}
  localOutputRefs:
    - es-output
{{</ highlight >}}

Fluentd config result:

{{< highlight xml >}}
 <filter **>
 @type kube_events_timestamp
 @id test-kube-events-timestamp
 timestamp_fields ["event.eventTime","event.lastTimestamp","event.firstTimestamp"]
 mapped_time_key mytimefield
 </filter>
{{</ highlight >}}


---
