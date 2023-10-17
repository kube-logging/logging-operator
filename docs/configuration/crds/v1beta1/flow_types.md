---
title: FlowSpec
weight: 200
generated_file: true
---

## FlowSpec

FlowSpec is the Kubernetes spec for Flows

### filters ([]Filter, optional) {#flowspec-filters}

Default: -

### flowLabel (string, optional) {#flowspec-flowlabel}

Default: -

### globalOutputRefs ([]string, optional) {#flowspec-globaloutputrefs}

Default: -

### includeLabelInRouter (*bool, optional) {#flowspec-includelabelinrouter}

Default: -

### localOutputRefs ([]string, optional) {#flowspec-localoutputrefs}

Default: -

### loggingRef (string, optional) {#flowspec-loggingref}

Default: -

### match ([]Match, optional) {#flowspec-match}

Default: -

### outputRefs ([]string, optional) {#flowspec-outputrefs}

Deprecated 

Default: -

### selectors (map[string]string, optional) {#flowspec-selectors}

Deprecated 

Default: -


## Match

### select (*Select, optional) {#match-select}

Default: -

### exclude (*Exclude, optional) {#match-exclude}

Default: -


## Select

### container_names ([]string, optional) {#select-container_names}

Default: -

### hosts ([]string, optional) {#select-hosts}

Default: -

### labels (map[string]string, optional) {#select-labels}

Default: -


## Exclude

### container_names ([]string, optional) {#exclude-container_names}

Default: -

### hosts ([]string, optional) {#exclude-hosts}

Default: -

### labels (map[string]string, optional) {#exclude-labels}

Default: -


## Filter

Filter definition for FlowSpec

### concat (*filter.Concat, optional) {#filter-concat}

Default: -

### dedot (*filter.DedotFilterConfig, optional) {#filter-dedot}

Default: -

### detectExceptions (*filter.DetectExceptions, optional) {#filter-detectexceptions}

Default: -

### elasticsearch_genid (*filter.ElasticsearchGenId, optional) {#filter-elasticsearch_genid}

Default: -

### enhanceK8s (*filter.EnhanceK8s, optional) {#filter-enhancek8s}

Default: -

### geoip (*filter.GeoIP, optional) {#filter-geoip}

Default: -

### grep (*filter.GrepConfig, optional) {#filter-grep}

Default: -

### kube_events_timestamp (*filter.KubeEventsTimestampConfig, optional) {#filter-kube_events_timestamp}

Default: -

### parser (*filter.ParserConfig, optional) {#filter-parser}

Default: -

### prometheus (*filter.PrometheusConfig, optional) {#filter-prometheus}

Default: -

### record_modifier (*filter.RecordModifier, optional) {#filter-record_modifier}

Default: -

### record_transformer (*filter.RecordTransformer, optional) {#filter-record_transformer}

Default: -

### stdout (*filter.StdOutFilterConfig, optional) {#filter-stdout}

Default: -

### sumologic (*filter.SumoLogic, optional) {#filter-sumologic}

Default: -

### tag_normaliser (*filter.TagNormaliser, optional) {#filter-tag_normaliser}

Default: -

### throttle (*filter.Throttle, optional) {#filter-throttle}

Default: -


## FlowStatus

FlowStatus defines the observed state of Flow

### active (*bool, optional) {#flowstatus-active}

Default: -

### problems ([]string, optional) {#flowstatus-problems}

Default: -

### problemsCount (int, optional) {#flowstatus-problemscount}

Default: -


## Flow

Flow Kubernetes object

###  (metav1.TypeMeta, required) {#flow-}

Default: -

### metadata (metav1.ObjectMeta, optional) {#flow-metadata}

Default: -

### spec (FlowSpec, optional) {#flow-spec}

Default: -

### status (FlowStatus, optional) {#flow-status}

Default: -


## FlowList

FlowList contains a list of Flow

###  (metav1.TypeMeta, required) {#flowlist-}

Default: -

### metadata (metav1.ListMeta, optional) {#flowlist-metadata}

Default: -

### items ([]Flow, required) {#flowlist-items}

Default: -


