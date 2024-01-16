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

Output metrics are applied before the log reaches the destination and contain output metadata like: `name,` `namespace` and `scope`. Scope shows whether the output is a local or global one. Available in Logging operator version 4.5 and later. 



## SyslogNGClusterFlowList

SyslogNGClusterFlowList contains a list of SyslogNGClusterFlow

###  (metav1.TypeMeta, required) {#syslogngclusterflowlist-}


### metadata (metav1.ListMeta, optional) {#syslogngclusterflowlist-metadata}


### items ([]SyslogNGClusterFlow, required) {#syslogngclusterflowlist-items}



