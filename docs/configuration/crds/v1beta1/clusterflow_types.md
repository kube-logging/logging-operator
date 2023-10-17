---
title: ClusterFlow
weight: 200
generated_file: true
---

## ClusterFlow

ClusterFlow is the Schema for the clusterflows API

###  (metav1.TypeMeta, required) {#clusterflow-}

Default: -

### metadata (metav1.ObjectMeta, optional) {#clusterflow-metadata}

Default: -

### spec (ClusterFlowSpec, optional) {#clusterflow-spec}

Name of the logging cluster to be attached 

Default: -

### status (FlowStatus, optional) {#clusterflow-status}

Default: -


## ClusterMatch

### select (*ClusterSelect, optional) {#clustermatch-select}

Default: -

### exclude (*ClusterExclude, optional) {#clustermatch-exclude}

Default: -


## ClusterSelect

### container_names ([]string, optional) {#clusterselect-container_names}

Default: -

### hosts ([]string, optional) {#clusterselect-hosts}

Default: -

### labels (map[string]string, optional) {#clusterselect-labels}

Default: -

### namespaces ([]string, optional) {#clusterselect-namespaces}

Default: -


## ClusterExclude

### container_names ([]string, optional) {#clusterexclude-container_names}

Default: -

### hosts ([]string, optional) {#clusterexclude-hosts}

Default: -

### labels (map[string]string, optional) {#clusterexclude-labels}

Default: -

### namespaces ([]string, optional) {#clusterexclude-namespaces}

Default: -


## ClusterFlowSpec

ClusterFlowSpec is the Kubernetes spec for ClusterFlows

### filters ([]Filter, optional) {#clusterflowspec-filters}

Default: -

### flowLabel (string, optional) {#clusterflowspec-flowlabel}

Default: -

### globalOutputRefs ([]string, optional) {#clusterflowspec-globaloutputrefs}

Default: -

### includeLabelInRouter (*bool, optional) {#clusterflowspec-includelabelinrouter}

Default: -

### loggingRef (string, optional) {#clusterflowspec-loggingref}

Default: -

### match ([]ClusterMatch, optional) {#clusterflowspec-match}

Default: -

### outputRefs ([]string, optional) {#clusterflowspec-outputrefs}

Deprecated 

Default: -

### selectors (map[string]string, optional) {#clusterflowspec-selectors}

Deprecated 

Default: -


## ClusterFlowList

ClusterFlowList contains a list of ClusterFlow

###  (metav1.TypeMeta, required) {#clusterflowlist-}

Default: -

### metadata (metav1.ListMeta, optional) {#clusterflowlist-metadata}

Default: -

### items ([]ClusterFlow, required) {#clusterflowlist-items}

Default: -


