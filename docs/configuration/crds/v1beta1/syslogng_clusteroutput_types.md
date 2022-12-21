---
title: SyslogNGClusterOutput
weight: 200
generated_file: true
---

## SyslogNGClusterOutput

SyslogNGClusterOutput is the Schema for the syslog-ng clusteroutputs API

###  (metav1.TypeMeta, required) {#syslogngclusteroutput-}

Default: -

### metadata (metav1.ObjectMeta, optional) {#syslogngclusteroutput-metadata}

Default: -

### spec (SyslogNGClusterOutputSpec, required) {#syslogngclusteroutput-spec}

Default: -

### status (SyslogNGOutputStatus, optional) {#syslogngclusteroutput-status}

Default: -


## SyslogNGClusterOutputSpec

SyslogNGClusterOutputSpec contains Kubernetes spec for SyslogNGClusterOutput

###  (SyslogNGOutputSpec, required) {#syslogngclusteroutputspec-}

Default: -

### enabledNamespaces ([]string, optional) {#syslogngclusteroutputspec-enablednamespaces}

Default: -


## SyslogNGClusterOutputList

SyslogNGClusterOutputList contains a list of SyslogNGClusterOutput

###  (metav1.TypeMeta, required) {#syslogngclusteroutputlist-}

Default: -

### metadata (metav1.ListMeta, optional) {#syslogngclusteroutputlist-metadata}

Default: -

### items ([]SyslogNGClusterOutput, required) {#syslogngclusteroutputlist-items}

Default: -


