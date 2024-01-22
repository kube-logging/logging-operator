---
title: SyslogNGFlowSpec
weight: 200
generated_file: true
---

## SyslogNGFlowSpec

SyslogNGFlowSpec is the Kubernetes spec for SyslogNGFlows

### filters ([]SyslogNGFilter, optional) {#syslogngflowspec-filters}


### globalOutputRefs ([]string, optional) {#syslogngflowspec-globaloutputrefs}


### localOutputRefs ([]string, optional) {#syslogngflowspec-localoutputrefs}


### loggingRef (string, optional) {#syslogngflowspec-loggingref}


### match (*SyslogNGMatch, optional) {#syslogngflowspec-match}


### outputMetrics ([]filter.MetricsProbe, optional) {#syslogngflowspec-outputmetrics}

Output metrics are applied before the log reaches the destination and contain output metadata like: `name,` `namespace` and `scope`. Scope shows whether the output is a local or global one. Available in Logging operator version 4.5 and later. 



## SyslogNGFilter

Filter definition for SyslogNGFlowSpec

### id (string, optional) {#syslogngfilter-id}


### match (*filter.MatchConfig, optional) {#syslogngfilter-match}


### parser (*filter.ParserConfig, optional) {#syslogngfilter-parser}


### rewrite ([]filter.RewriteConfig, optional) {#syslogngfilter-rewrite}



## SyslogNGFlow

Flow Kubernetes object

###  (metav1.TypeMeta, required) {#syslogngflow-}


### metadata (metav1.ObjectMeta, optional) {#syslogngflow-metadata}


### spec (SyslogNGFlowSpec, optional) {#syslogngflow-spec}


### status (SyslogNGFlowStatus, optional) {#syslogngflow-status}



## SyslogNGFlowList

FlowList contains a list of Flow

###  (metav1.TypeMeta, required) {#syslogngflowlist-}


### metadata (metav1.ListMeta, optional) {#syslogngflowlist-metadata}


### items ([]SyslogNGFlow, required) {#syslogngflowlist-items}



