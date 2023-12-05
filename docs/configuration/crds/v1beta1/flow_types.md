---
title: FlowSpec
weight: 200
generated_file: true
---

## FlowSpec

FlowSpec is the Kubernetes spec for Flows

### filters ([]Filter, optional) {#flowspec-filters}


### flowLabel (string, optional) {#flowspec-flowlabel}


### globalOutputRefs ([]string, optional) {#flowspec-globaloutputrefs}


### includeLabelInRouter (*bool, optional) {#flowspec-includelabelinrouter}


### localOutputRefs ([]string, optional) {#flowspec-localoutputrefs}


### loggingRef (string, optional) {#flowspec-loggingref}


### match ([]Match, optional) {#flowspec-match}


### outputRefs ([]string, optional) {#flowspec-outputrefs}

Deprecated 


### selectors (map[string]string, optional) {#flowspec-selectors}

Deprecated 



## Match

### select (*Select, optional) {#match-select}


### exclude (*Exclude, optional) {#match-exclude}



## Select

### container_names ([]string, optional) {#select-container_names}


### hosts ([]string, optional) {#select-hosts}


### labels (map[string]string, optional) {#select-labels}



## Exclude

### container_names ([]string, optional) {#exclude-container_names}


### hosts ([]string, optional) {#exclude-hosts}


### labels (map[string]string, optional) {#exclude-labels}



## Filter

Filter definition for FlowSpec

### concat (*filter.Concat, optional) {#filter-concat}


### dedot (*filter.DedotFilterConfig, optional) {#filter-dedot}


### detectExceptions (*filter.DetectExceptions, optional) {#filter-detectexceptions}


### elasticsearch_genid (*filter.ElasticsearchGenId, optional) {#filter-elasticsearch_genid}


### enhanceK8s (*filter.EnhanceK8s, optional) {#filter-enhancek8s}


### geoip (*filter.GeoIP, optional) {#filter-geoip}


### grep (*filter.GrepConfig, optional) {#filter-grep}


### kube_events_timestamp (*filter.KubeEventsTimestampConfig, optional) {#filter-kube_events_timestamp}


### parser (*filter.ParserConfig, optional) {#filter-parser}


### prometheus (*filter.PrometheusConfig, optional) {#filter-prometheus}


### record_modifier (*filter.RecordModifier, optional) {#filter-record_modifier}


### record_transformer (*filter.RecordTransformer, optional) {#filter-record_transformer}


### stdout (*filter.StdOutFilterConfig, optional) {#filter-stdout}


### sumologic (*filter.SumoLogic, optional) {#filter-sumologic}


### tag_normaliser (*filter.TagNormaliser, optional) {#filter-tag_normaliser}


### throttle (*filter.Throttle, optional) {#filter-throttle}


### useragent (*filter.UserAgent, optional) {#filter-useragent}



## FlowStatus

FlowStatus defines the observed state of Flow

### active (*bool, optional) {#flowstatus-active}


### problems ([]string, optional) {#flowstatus-problems}


### problemsCount (int, optional) {#flowstatus-problemscount}



## Flow

Flow Kubernetes object

###  (metav1.TypeMeta, required) {#flow-}


### metadata (metav1.ObjectMeta, optional) {#flow-metadata}


### spec (FlowSpec, optional) {#flow-spec}


### status (FlowStatus, optional) {#flow-status}



## FlowList

FlowList contains a list of Flow

###  (metav1.TypeMeta, required) {#flowlist-}


### metadata (metav1.ListMeta, optional) {#flowlist-metadata}


### items ([]Flow, required) {#flowlist-items}



