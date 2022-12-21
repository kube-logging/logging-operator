---
title: ClusterOutput
weight: 200
generated_file: true
---

## ClusterOutput

ClusterOutput is the Schema for the clusteroutputs API

###  (metav1.TypeMeta, required) {#clusteroutput-}

Default: -

### metadata (metav1.ObjectMeta, optional) {#clusteroutput-metadata}

Default: -

### spec (ClusterOutputSpec, required) {#clusteroutput-spec}

Default: -

### status (OutputStatus, optional) {#clusteroutput-status}

Default: -


## ClusterOutputSpec

ClusterOutputSpec contains Kubernetes spec for ClusterOutput

###  (OutputSpec, required) {#clusteroutputspec-}

Default: -

### enabledNamespaces ([]string, optional) {#clusteroutputspec-enablednamespaces}

Default: -


## ClusterOutputList

ClusterOutputList contains a list of ClusterOutput

###  (metav1.TypeMeta, required) {#clusteroutputlist-}

Default: -

### metadata (metav1.ListMeta, optional) {#clusteroutputlist-metadata}

Default: -

### items ([]ClusterOutput, required) {#clusteroutputlist-items}

Default: -


