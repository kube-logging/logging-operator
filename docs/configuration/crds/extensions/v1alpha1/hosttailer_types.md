---
title: HostTailer
weight: 200
generated_file: true
---

## HostTailerSpec

HostTailerSpec defines the desired state of HostTailer

### enableRecreateWorkloadOnImmutableFieldChange (bool, optional) {#hosttailerspec-enablerecreateworkloadonimmutablefieldchange}

EnableRecreateWorkloadOnImmutableFieldChange enables the operator to recreate the daemonset (and possibly other resource in the future) in case there is a change in an immutable field that otherwise couldn't be managed with a simple update. 

Default: -

### fileTailers ([]FileTailer, optional) {#hosttailerspec-filetailers}

List of [file tailers](#filetailer). 

Default: -

### image (tailer.ImageSpec, optional) {#hosttailerspec-image}

Default: -

### systemdTailers ([]SystemdTailer, optional) {#hosttailerspec-systemdtailers}

List of [systemd tailers](#systemdtailer). 

Default: -

### workloadOverrides (*types.PodSpecBase, optional) {#hosttailerspec-workloadoverrides}

Override podSpec fields for the given daemonset 

Default: -

### workloadMetaOverrides (*types.MetaBase, optional) {#hosttailerspec-workloadmetaoverrides}

Override metadata of the created resources 

Default: -


## HostTailerStatus

HostTailerStatus defines the observed state of [HostTailer](#hosttailer).


## HostTailer

HostTailer is the Schema for the hosttailers API

###  (metav1.TypeMeta, required) {#hosttailer-}

Default: -

### metadata (metav1.ObjectMeta, optional) {#hosttailer-metadata}

Default: -

### spec (HostTailerSpec, optional) {#hosttailer-spec}

Default: -

### status (HostTailerStatus, optional) {#hosttailer-status}

Default: -


## HostTailerList

HostTailerList contains a list of [HostTailers](#hosttailer).

###  (metav1.TypeMeta, required) {#hosttailerlist-}

Default: -

### metadata (metav1.ListMeta, optional) {#hosttailerlist-metadata}

Default: -

### items ([]HostTailer, required) {#hosttailerlist-items}

Default: -


## FileTailer

FileTailer configuration options

### buffer_chunk_size (string, optional) {#filetailer-buffer_chunk_size}

Set the buffer chunk size per active filetailer 

Default: -

### buffer_max_size (string, optional) {#filetailer-buffer_max_size}

Set the limit of the buffer size per active filetailer 

Default: -

### containerOverrides (*types.ContainerBase, optional) {#filetailer-containeroverrides}

Override container fields for the given tailer 

Default: -

### disabled (bool, optional) {#filetailer-disabled}

Disable tailing the file 

Default: -

### image (*tailer.ImageSpec, optional) {#filetailer-image}

Override image field for the given trailer 

Default: -

### name (string, required) {#filetailer-name}

Name for the tailer 

Default: -

### path (string, optional) {#filetailer-path}

Path to the loggable file 

Default: -

### read_from_head (bool, optional) {#filetailer-read_from_head}

Start reading from the head of new log files 

Default: -

### skip_long_lines (string, optional) {#filetailer-skip_long_lines}

Skip long line when exceeding Buffer_Max_Size 

Default: -


## SystemdTailer

SystemdTailer configuration options

### containerOverrides (*types.ContainerBase, optional) {#systemdtailer-containeroverrides}

Override container fields for the given tailer 

Default: -

### disabled (bool, optional) {#systemdtailer-disabled}

Disable component 

Default: -

### image (*tailer.ImageSpec, optional) {#systemdtailer-image}

Override image field for the given trailer 

Default: -

### maxEntries (int, optional) {#systemdtailer-maxentries}

Maximum entries to read when starting to tail logs to avoid high pressure 

Default: -

### name (string, required) {#systemdtailer-name}

Name for the tailer 

Default: -

### path (string, optional) {#systemdtailer-path}

Override systemd log path 

Default: -

### systemdFilter (string, optional) {#systemdtailer-systemdfilter}

Filter to select systemd unit example: kubelet.service 

Default: -


