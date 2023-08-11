---
title: SyslogNGSpec
weight: 200
generated_file: true
---

## SyslogNGSpec

SyslogNGSpec defines the desired state of SyslogNG

### tls (SyslogNGTLS, optional) {#syslogngspec-tls}

Default: -

### readinessDefaultCheck (ReadinessDefaultCheck, optional) {#syslogngspec-readinessdefaultcheck}

Default: -

### skipRBACCreate (bool, optional) {#syslogngspec-skiprbaccreate}

Default: -

### statefulSet (*typeoverride.StatefulSet, optional) {#syslogngspec-statefulset}

Default: -

### service (*typeoverride.Service, optional) {#syslogngspec-service}

Default: -

### serviceAccount (*typeoverride.ServiceAccount, optional) {#syslogngspec-serviceaccount}

Default: -

### configCheckPod (*typeoverride.PodSpec, optional) {#syslogngspec-configcheckpod}

Default: -

### metrics (*Metrics, optional) {#syslogngspec-metrics}

Default: -

### metricsService (*typeoverride.Service, optional) {#syslogngspec-metricsservice}

Default: -

### bufferVolumeMetrics (*BufferMetrics, optional) {#syslogngspec-buffervolumemetrics}

Default: -

### bufferVolumeMetricsService (*typeoverride.Service, optional) {#syslogngspec-buffervolumemetricsservice}

Default: -

### globalOptions (*GlobalOptions, optional) {#syslogngspec-globaloptions}

Default: -

### jsonKeyPrefix (string, optional) {#syslogngspec-jsonkeyprefix}

Default: -

### jsonKeyDelim (string, optional) {#syslogngspec-jsonkeydelim}

Default: -

### maxConnections (int, optional) {#syslogngspec-maxconnections}

Default: -

### logIWSize (int, optional) {#syslogngspec-logiwsize}

Default: -


## SyslogNGTLS

SyslogNGTLS defines the TLS configs

### enabled (bool, required) {#syslogngtls-enabled}

Default: -

### secretName (string, optional) {#syslogngtls-secretname}

Default: -

### sharedKey (string, optional) {#syslogngtls-sharedkey}

Default: -


## GlobalOptions

### stats_level (*int, optional) {#globaloptions-stats_level}

deprecated use stats/level from 4.1+ 

Default: -

### stats_freq (*int, optional) {#globaloptions-stats_freq}

deprecated use stats/freq from 4.1+ 

Default: -

### stats (*Stats, optional) {#globaloptions-stats}

TODO switch to this by default 

Default: -


## Stats

### level (*int, optional) {#stats-level}

Default: -

### freq (*int, optional) {#stats-freq}

Default: -


