---
title: NodeAgent
weight: 200
generated_file: true
---

## NodeAgent

NodeAgent

###  (metav1.TypeMeta, required) {#nodeagent-}


### metadata (metav1.ObjectMeta, optional) {#nodeagent-metadata}


### spec (NodeAgentSpec, optional) {#nodeagent-spec}


### status (NodeAgentStatus, optional) {#nodeagent-status}



## NodeAgentSpec

NodeAgentSpec

###  (NodeAgentConfig, required) {#nodeagentspec-}

InlineNodeAgent 


### loggingRef (string, optional) {#nodeagentspec-loggingref}



## NodeAgentConfig

### nodeAgentFluentbit (*NodeAgentFluentbit, optional) {#nodeagentconfig-nodeagentfluentbit}


### metadata (types.MetaBase, optional) {#nodeagentconfig-metadata}


### profile (string, optional) {#nodeagentconfig-profile}



## NodeAgentStatus

NodeAgentStatus


## NodeAgentList

NodeAgentList

###  (metav1.TypeMeta, required) {#nodeagentlist-}


### metadata (metav1.ListMeta, optional) {#nodeagentlist-metadata}


### items ([]NodeAgent, required) {#nodeagentlist-items}



## InlineNodeAgent

InlineNodeAgent
@deprecated, replaced by NodeAgent

###  (NodeAgentConfig, required) {#inlinenodeagent-}


### name (string, optional) {#inlinenodeagent-name}

InlineNodeAgent unique name. 



## NodeAgentFluentbit

### bufferStorage (BufferStorage, optional) {#nodeagentfluentbit-bufferstorage}


### bufferStorageVolume (volume.KubernetesVolume, optional) {#nodeagentfluentbit-bufferstoragevolume}

[volume.KubernetesVolume](https://github.com/cisco-open/operator-tools/tree/master/docs/types) 


### containersPath (string, optional) {#nodeagentfluentbit-containerspath}


### coroStackSize (int32, optional) {#nodeagentfluentbit-corostacksize}

Set the coroutines stack size in bytes. The value must be greater than the page size of the running system. Don't set too small value (say 4096), or coroutine threads can overrun the stack buffer. Do not change the default value of this parameter unless you know what you are doing. (default: 24576) 

Default: 24576

### customConfigSecret (string, optional) {#nodeagentfluentbit-customconfigsecret}


### daemonSet (*typeoverride.DaemonSet, optional) {#nodeagentfluentbit-daemonset}


### disableKubernetesFilter (*bool, optional) {#nodeagentfluentbit-disablekubernetesfilter}


### enableUpstream (*bool, optional) {#nodeagentfluentbit-enableupstream}


### enabled (*bool, optional) {#nodeagentfluentbit-enabled}


### extraVolumeMounts ([]*VolumeMount, optional) {#nodeagentfluentbit-extravolumemounts}


### filterAws (*FilterAws, optional) {#nodeagentfluentbit-filteraws}


### filterKubernetes (FilterKubernetes, optional) {#nodeagentfluentbit-filterkubernetes}


### flush (int32, optional) {#nodeagentfluentbit-flush}

Set the flush time in seconds.nanoseconds. The engine loop uses a Flush timeout to define when is required to flush the records ingested by input plugins through the defined output plugins. (default: 1) 

Default: 1

### forwardOptions (*ForwardOptions, optional) {#nodeagentfluentbit-forwardoptions}


### grace (int32, optional) {#nodeagentfluentbit-grace}

Set the grace time in seconds as Integer value. The engine loop uses a Grace timeout to define wait time on exit (default: 5) 

Default: 5

### inputTail (InputTail, optional) {#nodeagentfluentbit-inputtail}


### livenessDefaultCheck (*bool, optional) {#nodeagentfluentbit-livenessdefaultcheck}

Default: true

### logLevel (string, optional) {#nodeagentfluentbit-loglevel}

Set the logging verbosity level. Allowed values are: error, warn, info, debug and trace. Values are accumulative, e.g: if 'debug' is set, it will include error, warning, info and debug.  Note that trace mode is only available if Fluent Bit was built with the WITH_TRACE option enabled. (default: info) 

Default: info

### metrics (*Metrics, optional) {#nodeagentfluentbit-metrics}


### metricsService (*typeoverride.Service, optional) {#nodeagentfluentbit-metricsservice}


### network (*FluentbitNetwork, optional) {#nodeagentfluentbit-network}


### podPriorityClassName (string, optional) {#nodeagentfluentbit-podpriorityclassname}


### positiondb (volume.KubernetesVolume, optional) {#nodeagentfluentbit-positiondb}

[volume.KubernetesVolume](https://github.com/cisco-open/operator-tools/tree/master/docs/types) 


### security (*Security, optional) {#nodeagentfluentbit-security}


### serviceAccount (*typeoverride.ServiceAccount, optional) {#nodeagentfluentbit-serviceaccount}


### tls (*FluentbitTLS, optional) {#nodeagentfluentbit-tls}


### targetHost (string, optional) {#nodeagentfluentbit-targethost}


### targetPort (int32, optional) {#nodeagentfluentbit-targetport}


### varLogsPath (string, optional) {#nodeagentfluentbit-varlogspath}



