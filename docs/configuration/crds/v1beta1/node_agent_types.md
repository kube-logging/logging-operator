---
title: NodeAgent
weight: 200
generated_file: true
---

## NodeAgent

NodeAgent

###  (metav1.TypeMeta, required) {#nodeagent-}

Default: -

### metadata (metav1.ObjectMeta, optional) {#nodeagent-metadata}

Default: -

### spec (NodeAgentSpec, optional) {#nodeagent-spec}

Default: -

### status (NodeAgentStatus, optional) {#nodeagent-status}

Default: -


## NodeAgentSpec

NodeAgentSpec

###  (NodeAgentConfig, required) {#nodeagentspec-}

InlineNodeAgent 

Default: -

### loggingRef (string, optional) {#nodeagentspec-loggingref}

Default: -


## NodeAgentConfig

### nodeAgentFluentbit (*NodeAgentFluentbit, optional) {#nodeagentconfig-nodeagentfluentbit}

Default: -

### metadata (types.MetaBase, optional) {#nodeagentconfig-metadata}

Default: -

### profile (string, optional) {#nodeagentconfig-profile}

Default: -


## NodeAgentStatus

NodeAgentStatus


## NodeAgentList

NodeAgentList

###  (metav1.TypeMeta, required) {#nodeagentlist-}

Default: -

### metadata (metav1.ListMeta, optional) {#nodeagentlist-metadata}

Default: -

### items ([]NodeAgent, required) {#nodeagentlist-items}

Default: -


## InlineNodeAgent

InlineNodeAgent
@deprecated, replaced by NodeAgent

###  (NodeAgentConfig, required) {#inlinenodeagent-}

Default: -

### name (string, optional) {#inlinenodeagent-name}

InlineNodeAgent unique name. 

Default: -


## NodeAgentFluentbit

### bufferStorage (BufferStorage, optional) {#nodeagentfluentbit-bufferstorage}

Default: -

### bufferStorageVolume (volume.KubernetesVolume, optional) {#nodeagentfluentbit-bufferstoragevolume}

[volume.KubernetesVolume](https://github.com/cisco-open/operator-tools/tree/master/docs/types) 

Default: -

### containersPath (string, optional) {#nodeagentfluentbit-containerspath}

Default: -

### coroStackSize (int32, optional) {#nodeagentfluentbit-corostacksize}

Set the coroutines stack size in bytes. The value must be greater than the page size of the running system. Don't set too small value (say 4096), or coroutine threads can overrun the stack buffer. Do not change the default value of this parameter unless you know what you are doing. (default: 24576) 

Default: 24576

### customConfigSecret (string, optional) {#nodeagentfluentbit-customconfigsecret}

Default: -

### daemonSet (*typeoverride.DaemonSet, optional) {#nodeagentfluentbit-daemonset}

Default: -

### disableKubernetesFilter (*bool, optional) {#nodeagentfluentbit-disablekubernetesfilter}

Default: -

### enableUpstream (*bool, optional) {#nodeagentfluentbit-enableupstream}

Default: -

### enabled (*bool, optional) {#nodeagentfluentbit-enabled}

Default: -

### extraVolumeMounts ([]*VolumeMount, optional) {#nodeagentfluentbit-extravolumemounts}

Default: -

### filterAws (*FilterAws, optional) {#nodeagentfluentbit-filteraws}

Default: -

### filterKubernetes (FilterKubernetes, optional) {#nodeagentfluentbit-filterkubernetes}

Default: -

### flush (int32, optional) {#nodeagentfluentbit-flush}

Set the flush time in seconds.nanoseconds. The engine loop uses a Flush timeout to define when is required to flush the records ingested by input plugins through the defined output plugins. (default: 1) 

Default: 1

### forwardOptions (*ForwardOptions, optional) {#nodeagentfluentbit-forwardoptions}

Default: -

### grace (int32, optional) {#nodeagentfluentbit-grace}

Set the grace time in seconds as Integer value. The engine loop uses a Grace timeout to define wait time on exit (default: 5) 

Default: 5

### inputTail (InputTail, optional) {#nodeagentfluentbit-inputtail}

Default: -

### livenessDefaultCheck (*bool, optional) {#nodeagentfluentbit-livenessdefaultcheck}

Default: true

### logLevel (string, optional) {#nodeagentfluentbit-loglevel}

Set the logging verbosity level. Allowed values are: error, warn, info, debug and trace. Values are accumulative, e.g: if 'debug' is set, it will include error, warning, info and debug.  Note that trace mode is only available if Fluent Bit was built with the WITH_TRACE option enabled. (default: info) 

Default: info

### metrics (*Metrics, optional) {#nodeagentfluentbit-metrics}

Default: -

### metricsService (*typeoverride.Service, optional) {#nodeagentfluentbit-metricsservice}

Default: -

### network (*FluentbitNetwork, optional) {#nodeagentfluentbit-network}

Default: -

### podPriorityClassName (string, optional) {#nodeagentfluentbit-podpriorityclassname}

Default: -

### positiondb (volume.KubernetesVolume, optional) {#nodeagentfluentbit-positiondb}

[volume.KubernetesVolume](https://github.com/cisco-open/operator-tools/tree/master/docs/types) 

Default: -

### security (*Security, optional) {#nodeagentfluentbit-security}

Default: -

### serviceAccount (*typeoverride.ServiceAccount, optional) {#nodeagentfluentbit-serviceaccount}

Default: -

### tls (*FluentbitTLS, optional) {#nodeagentfluentbit-tls}

Default: -

### targetHost (string, optional) {#nodeagentfluentbit-targethost}

Default: -

### targetPort (int32, optional) {#nodeagentfluentbit-targetport}

Default: -

### varLogsPath (string, optional) {#nodeagentfluentbit-varlogspath}

Default: -


