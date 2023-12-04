---
title: EventTailer
weight: 200
generated_file: true
---

## EventTailerSpec

EventTailerSpec defines the desired state of EventTailer

### containerOverrides (*types.ContainerBase, optional) {#eventtailerspec-containeroverrides}

Override container fields for the given statefulset 


### controlNamespace (string, required) {#eventtailerspec-controlnamespace}

The resources of EventTailer will be placed into this namespace 


### image (*tailer.ImageSpec, optional) {#eventtailerspec-image}

Override image related fields for the given statefulset, highest precedence 


### positionVolume (volume.KubernetesVolume, optional) {#eventtailerspec-positionvolume}

Volume definition for tracking fluentbit file positions (optional) 


### workloadOverrides (*types.PodSpecBase, optional) {#eventtailerspec-workloadoverrides}

Override podSpec fields for the given statefulset 


### workloadMetaOverrides (*types.MetaBase, optional) {#eventtailerspec-workloadmetaoverrides}

Override metadata of the created resources 



## EventTailerStatus

EventTailerStatus defines the observed state of EventTailer


## EventTailer

EventTailer is the Schema for the eventtailers API

###  (metav1.TypeMeta, required) {#eventtailer-}


### metadata (metav1.ObjectMeta, optional) {#eventtailer-metadata}


### spec (EventTailerSpec, optional) {#eventtailer-spec}


### status (EventTailerStatus, optional) {#eventtailer-status}



## EventTailerList

EventTailerList contains a list of EventTailer

###  (metav1.TypeMeta, required) {#eventtailerlist-}


### metadata (metav1.ListMeta, optional) {#eventtailerlist-metadata}


### items ([]EventTailer, required) {#eventtailerlist-items}



