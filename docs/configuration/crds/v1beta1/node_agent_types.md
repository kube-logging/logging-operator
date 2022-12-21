## NodeAgent

### name (string, optional) {#nodeagent-name}

NodeAgent unique name. 

Default: -

### profile (string, optional) {#nodeagent-profile}

Specify the Logging-Operator nodeAgents profile. It can be linux or windows .  

Default: linux

### metadata (types.MetaBase, optional) {#nodeagent-metadata}

Default: -

### nodeAgentFluentbit (*NodeAgentFluentbit, optional) {#nodeagent-nodeagentfluentbit}

Default: -


## NodeAgentFluentbit

### enabled (*bool, optional) {#nodeagentfluentbit-enabled}

Default: -

### daemonSet (*typeoverride.DaemonSet, optional) {#nodeagentfluentbit-daemonset}

Default: -

### serviceAccount (*typeoverride.ServiceAccount, optional) {#nodeagentfluentbit-serviceaccount}

Default: -

### tls (*FluentbitTLS, optional) {#nodeagentfluentbit-tls}

Default: -

### targetHost (string, optional) {#nodeagentfluentbit-targethost}

Default: -

### targetPort (int32, optional) {#nodeagentfluentbit-targetport}

Default: -

### flush (int32, optional) {#nodeagentfluentbit-flush}

Set the flush time in seconds.nanoseconds. The engine loop uses a Flush timeout to define when is required to flush the records ingested by input plugins through the defined output plugins. (default: 1) 

Default: 1

### grace (int32, optional) {#nodeagentfluentbit-grace}

Set the grace time in seconds as Integer value. The engine loop uses a Grace timeout to define wait time on exit (default: 5) 

Default: 5

### logLevel (string, optional) {#nodeagentfluentbit-loglevel}

Set the logging verbosity level. Allowed values are: error, warn, info, debug and trace. Values are accumulative, e.g: if 'debug' is set, it will include error, warning, info and debug.  Note that trace mode is only available if Fluent Bit was built with the WITH_TRACE option enabled. (default: info) 

Default: info

### coroStackSize (int32, optional) {#nodeagentfluentbit-corostacksize}

Set the coroutines stack size in bytes. The value must be greater than the page size of the running system. Don't set too small value (say 4096), or coroutine threads can overrun the stack buffer. Do not change the default value of this parameter unless you know what you are doing. (default: 24576) 

Default: 24576

### metrics (*Metrics, optional) {#nodeagentfluentbit-metrics}

Default: -

### metricsService (*typeoverride.Service, optional) {#nodeagentfluentbit-metricsservice}

Default: -

### security (*Security, optional) {#nodeagentfluentbit-security}

Default: -

### positiondb (volume.KubernetesVolume, optional) {#nodeagentfluentbit-positiondb}

[volume.KubernetesVolume](https://github.com/banzaicloud/operator-tools/tree/master/docs/types) 

Default: -

### containersPath (string, optional) {#nodeagentfluentbit-containerspath}

Default: -

### varLogsPath (string, optional) {#nodeagentfluentbit-varlogspath}

Default: -

### extraVolumeMounts ([]*VolumeMount, optional) {#nodeagentfluentbit-extravolumemounts}

Default: -

### inputTail (InputTail, optional) {#nodeagentfluentbit-inputtail}

Default: -

### filterAws (*FilterAws, optional) {#nodeagentfluentbit-filteraws}

Default: -

### filterKubernetes (FilterKubernetes, optional) {#nodeagentfluentbit-filterkubernetes}

Default: -

### disableKubernetesFilter (*bool, optional) {#nodeagentfluentbit-disablekubernetesfilter}

Default: -

### bufferStorage (BufferStorage, optional) {#nodeagentfluentbit-bufferstorage}

Default: -

### bufferStorageVolume (volume.KubernetesVolume, optional) {#nodeagentfluentbit-bufferstoragevolume}

[volume.KubernetesVolume](https://github.com/banzaicloud/operator-tools/tree/master/docs/types) 

Default: -

### customConfigSecret (string, optional) {#nodeagentfluentbit-customconfigsecret}

Default: -

### podPriorityClassName (string, optional) {#nodeagentfluentbit-podpriorityclassname}

Default: -

### livenessDefaultCheck (*bool, optional) {#nodeagentfluentbit-livenessdefaultcheck}

Default: true

### network (*FluentbitNetwork, optional) {#nodeagentfluentbit-network}

Default: -

### forwardOptions (*ForwardOptions, optional) {#nodeagentfluentbit-forwardoptions}

Default: -

### enableUpstream (*bool, optional) {#nodeagentfluentbit-enableupstream}

Default: -


