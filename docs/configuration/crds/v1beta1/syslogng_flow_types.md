---
title: SyslogNGFlowSpec
weight: 200
generated_file: true
---

## SyslogNGFlowSpec

SyslogNGFlowSpec is the Kubernetes spec for SyslogNGFlows

### match (*SyslogNGMatch, optional) {#syslogngflowspec-match}

Default: -

### filters ([]SyslogNGFilter, optional) {#syslogngflowspec-filters}

Default: -

### loggingRef (string, optional) {#syslogngflowspec-loggingref}

Default: -

### globalOutputRefs ([]string, optional) {#syslogngflowspec-globaloutputrefs}

Default: -

### localOutputRefs ([]string, optional) {#syslogngflowspec-localoutputrefs}

Default: -

### outputMetrics ([]filter.MetricsProbe, optional) {#syslogngflowspec-outputmetrics}

Default: -


## SyslogNGFilter

Filter definition for SyslogNGFlowSpec

### id (string, optional) {#syslogngfilter-id}

Default: -

### match (*filter.MatchConfig, optional) {#syslogngfilter-match}

Default: -

### rewrite ([]filter.RewriteConfig, optional) {#syslogngfilter-rewrite}

Default: -

### parser (*filter.ParserConfig, optional) {#syslogngfilter-parser}

Default: -


## SyslogNGFlow

Flow Kubernetes object

###  (metav1.TypeMeta, required) {#syslogngflow-}

Default: -

### metadata (metav1.ObjectMeta, optional) {#syslogngflow-metadata}

Default: -

### spec (SyslogNGFlowSpec, optional) {#syslogngflow-spec}

Default: -

### status (SyslogNGFlowStatus, optional) {#syslogngflow-status}

Default: -


## SyslogNGFlowList

FlowList contains a list of Flow

###  (metav1.TypeMeta, required) {#syslogngflowlist-}

Default: -

### metadata (metav1.ListMeta, optional) {#syslogngflowlist-metadata}

Default: -

### items ([]SyslogNGFlow, required) {#syslogngflowlist-items}

Default: -


