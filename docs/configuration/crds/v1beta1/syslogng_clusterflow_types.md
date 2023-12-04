---
title: SyslogNGClusterFlow
weight: 200
generated_file: true
---

## SyslogNGClusterFlow

SyslogNGClusterFlow is the Schema for the syslog-ng clusterflows API

###  (metav1.TypeMeta, required) {#syslogngclusterflow-}


### metadata (metav1.ObjectMeta, optional) {#syslogngclusterflow-metadata}


### spec (SyslogNGClusterFlowSpec, optional) {#syslogngclusterflow-spec}


### status (SyslogNGFlowStatus, optional) {#syslogngclusterflow-status}



## SyslogNGClusterFlowSpec

SyslogNGClusterFlowSpec is the Kubernetes spec for Flows

### filters ([]SyslogNGFilter, optional) {#syslogngclusterflowspec-filters}


### globalOutputRefs ([]string, optional) {#syslogngclusterflowspec-globaloutputrefs}


### loggingRef (string, optional) {#syslogngclusterflowspec-loggingref}


### match (*SyslogNGMatch, optional) {#syslogngclusterflowspec-match}


### outputMetrics ([]filter.MetricsProbe, optional) {#syslogngclusterflowspec-outputmetrics}



## SyslogNGClusterFlowList

SyslogNGClusterFlowList contains a list of SyslogNGClusterFlow

###  (metav1.TypeMeta, required) {#syslogngclusterflowlist-}


### metadata (metav1.ListMeta, optional) {#syslogngclusterflowlist-metadata}


### items ([]SyslogNGClusterFlow, required) {#syslogngclusterflowlist-items}



