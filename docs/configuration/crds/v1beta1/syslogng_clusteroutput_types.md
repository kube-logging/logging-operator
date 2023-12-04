---
title: SyslogNGClusterOutput
weight: 200
generated_file: true
---

## SyslogNGClusterOutput

SyslogNGClusterOutput is the Schema for the syslog-ng clusteroutputs API

###  (metav1.TypeMeta, required) {#syslogngclusteroutput-}


### metadata (metav1.ObjectMeta, optional) {#syslogngclusteroutput-metadata}


### spec (SyslogNGClusterOutputSpec, required) {#syslogngclusteroutput-spec}


### status (SyslogNGOutputStatus, optional) {#syslogngclusteroutput-status}



## SyslogNGClusterOutputSpec

SyslogNGClusterOutputSpec contains Kubernetes spec for SyslogNGClusterOutput

###  (SyslogNGOutputSpec, required) {#syslogngclusteroutputspec-}


### enabledNamespaces ([]string, optional) {#syslogngclusteroutputspec-enablednamespaces}



## SyslogNGClusterOutputList

SyslogNGClusterOutputList contains a list of SyslogNGClusterOutput

###  (metav1.TypeMeta, required) {#syslogngclusteroutputlist-}


### metadata (metav1.ListMeta, optional) {#syslogngclusteroutputlist-metadata}


### items ([]SyslogNGClusterOutput, required) {#syslogngclusteroutputlist-items}



