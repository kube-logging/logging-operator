---
title: ClusterOutput
weight: 200
generated_file: true
---

## ClusterOutput

ClusterOutput is the Schema for the clusteroutputs API

###  (metav1.TypeMeta, required) {#clusteroutput-}


### metadata (metav1.ObjectMeta, optional) {#clusteroutput-metadata}


### spec (ClusterOutputSpec, required) {#clusteroutput-spec}


### status (OutputStatus, optional) {#clusteroutput-status}



## ClusterOutputSpec

ClusterOutputSpec contains Kubernetes spec for ClusterOutput

###  (OutputSpec, required) {#clusteroutputspec-}


### enabledNamespaces ([]string, optional) {#clusteroutputspec-enablednamespaces}



## ClusterOutputList

ClusterOutputList contains a list of ClusterOutput

###  (metav1.TypeMeta, required) {#clusteroutputlist-}


### metadata (metav1.ListMeta, optional) {#clusteroutputlist-metadata}


### items ([]ClusterOutput, required) {#clusteroutputlist-items}



