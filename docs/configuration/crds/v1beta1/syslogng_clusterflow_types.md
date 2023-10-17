---
title: SyslogNGClusterFlow
weight: 200
generated_file: true
---

## SyslogNGClusterFlow

SyslogNGClusterFlow is the Schema for the syslog-ng clusterflows API

###  (metav1.TypeMeta, required) {#syslogngclusterflow-}

Default: -

### metadata (metav1.ObjectMeta, optional) {#syslogngclusterflow-metadata}

Default: -

### spec (SyslogNGClusterFlowSpec, optional) {#syslogngclusterflow-spec}

Default: -

### status (SyslogNGFlowStatus, optional) {#syslogngclusterflow-status}

Default: -


## SyslogNGClusterFlowSpec

SyslogNGClusterFlowSpec is the Kubernetes spec for Flows

### filters ([]SyslogNGFilter, optional) {#syslogngclusterflowspec-filters}

Default: -

### globalOutputRefs ([]string, optional) {#syslogngclusterflowspec-globaloutputrefs}

Default: -

### loggingRef (string, optional) {#syslogngclusterflowspec-loggingref}

Default: -

### match (*SyslogNGMatch, optional) {#syslogngclusterflowspec-match}

Default: -


## SyslogNGClusterFlowList

SyslogNGClusterFlowList contains a list of SyslogNGClusterFlow

###  (metav1.TypeMeta, required) {#syslogngclusterflowlist-}

Default: -

### metadata (metav1.ListMeta, optional) {#syslogngclusterflowlist-metadata}

Default: -

### items ([]SyslogNGClusterFlow, required) {#syslogngclusterflowlist-items}

Default: -


