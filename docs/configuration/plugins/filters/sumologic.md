---
title: SumoLogic
weight: 200
generated_file: true
---

# Sumo Logic collection solution for Kubernetes
## Overview
 More info at https://github.com/SumoLogic/sumologic-kubernetes-collection

## Configuration
## SumoLogic

### collector_key_name (string, optional) {#sumologic-collector_key_name}

CollectorKey Name

Default: `_collector`

### collector_value (string, optional) {#sumologic-collector_value}

Collector Value

Default: "undefined"

### exclude_container_regex (string, optional) {#sumologic-exclude_container_regex}

Exclude Container Regex

Default: ""

### exclude_facility_regex (string, optional) {#sumologic-exclude_facility_regex}

Exclude Facility Regex

Default: ""

### exclude_host_regex (string, optional) {#sumologic-exclude_host_regex}

Exclude Host Regex

Default: ""

### exclude_namespace_regex (string, optional) {#sumologic-exclude_namespace_regex}

Exclude Namespace Regex

Default: ""

### exclude_pod_regex (string, optional) {#sumologic-exclude_pod_regex}

Exclude Pod Regex

Default: ""

### exclude_priority_regex (string, optional) {#sumologic-exclude_priority_regex}

Exclude Priority Regex

Default: ""

### exclude_unit_regex (string, optional) {#sumologic-exclude_unit_regex}

Exclude Unit Regex

Default: ""

### log_format (string, optional) {#sumologic-log_format}

Log Format

Default: json

### source_category (string, optional) {#sumologic-source_category}

Source Category

Default: `%{namespace}/%{pod_name}`

### source_category_key_name (string, optional) {#sumologic-source_category_key_name}

Source CategoryKey Name

Default: `_sourceCategory`

### source_category_prefix (string, optional) {#sumologic-source_category_prefix}

Source Category Prefix

Default: kubernetes/

### source_category_replace_dash (string, optional) {#sumologic-source_category_replace_dash}

Source Category Replace Dash

Default: "/"

### source_host (string, optional) {#sumologic-source_host}

Source Host

Default: ""

### source_host_key_name (string, optional) {#sumologic-source_host_key_name}

Source HostKey Name

Default: `_sourceHost`

### source_name (string, optional) {#sumologic-source_name}

Source Name

Default: `%{namespace}.%{pod}.%{container}`

### source_name_key_name (string, optional) {#sumologic-source_name_key_name}

Source NameKey Name

Default: `_sourceName`

### tracing_annotation_prefix (string, optional) {#sumologic-tracing_annotation_prefix}

Tracing Annotation Prefix

Default: `pod_annotation_`

### tracing_container_name (string, optional) {#sumologic-tracing_container_name}

Tracing Container Name

Default: "container_name"

### tracing_format (*bool, optional) {#sumologic-tracing_format}

Tracing Format

Default: false

### tracing_host (string, optional) {#sumologic-tracing_host}

Tracing Host

Default: "hostname"

### tracing_label_prefix (string, optional) {#sumologic-tracing_label_prefix}

Tracing Label Prefix

Default: `pod_label_`

### tracing_namespace (string, optional) {#sumologic-tracing_namespace}

Tracing Namespace

Default: "namespace"

### tracing_pod (string, optional) {#sumologic-tracing_pod}

Tracing Pod

Default: "pod"

### tracing_pod_id (string, optional) {#sumologic-tracing_pod_id}

Tracing Pod ID

Default: "pod_id"




## Example `Parser` filter configurations

{{< highlight yaml >}}
apiVersion: logging.banzaicloud.io/v1beta1
kind: Flow
metadata:
  name: demo-flow
spec:
  filters:
    - sumologic:
        source_name: "elso"
  selectors: {}
  localOutputRefs:
    - demo-output
{{</ highlight >}}

Fluentd config result:

{{< highlight xml >}}
<filter **>
  @type kubernetes_sumologic
  @id test_sumologic
  source_name elso
</filter>
{{</ highlight >}}


---
