---
title: ClusterFlow
weight: 200
generated_file: true
---

## ClusterFlow

ClusterFlow is the Schema for the clusterflows API

###  (metav1.TypeMeta, required) {#clusterflow-}


### metadata (metav1.ObjectMeta, optional) {#clusterflow-metadata}


### spec (ClusterFlowSpec, optional) {#clusterflow-spec}

Name of the logging cluster to be attached 


### status (FlowStatus, optional) {#clusterflow-status}



## ClusterMatch

### select (*ClusterSelect, optional) {#clustermatch-select}


### exclude (*ClusterExclude, optional) {#clustermatch-exclude}



## ClusterSelect

### container_names ([]string, optional) {#clusterselect-container_names}


### hosts ([]string, optional) {#clusterselect-hosts}


### labels (map[string]string, optional) {#clusterselect-labels}


### namespaces ([]string, optional) {#clusterselect-namespaces}



## ClusterExclude

### container_names ([]string, optional) {#clusterexclude-container_names}


### hosts ([]string, optional) {#clusterexclude-hosts}


### labels (map[string]string, optional) {#clusterexclude-labels}


### namespaces ([]string, optional) {#clusterexclude-namespaces}



## ClusterFlowSpec

ClusterFlowSpec is the Kubernetes spec for ClusterFlows

### filters ([]Filter, optional) {#clusterflowspec-filters}


### flowLabel (string, optional) {#clusterflowspec-flowlabel}


### globalOutputRefs ([]string, optional) {#clusterflowspec-globaloutputrefs}


### includeLabelInRouter (*bool, optional) {#clusterflowspec-includelabelinrouter}


### loggingRef (string, optional) {#clusterflowspec-loggingref}


### match ([]ClusterMatch, optional) {#clusterflowspec-match}


### outputRefs ([]string, optional) {#clusterflowspec-outputrefs}

Deprecated 


### selectors (map[string]string, optional) {#clusterflowspec-selectors}

Deprecated 



## ClusterFlowList

ClusterFlowList contains a list of ClusterFlow

###  (metav1.TypeMeta, required) {#clusterflowlist-}


### metadata (metav1.ListMeta, optional) {#clusterflowlist-metadata}


### items ([]ClusterFlow, required) {#clusterflowlist-items}



