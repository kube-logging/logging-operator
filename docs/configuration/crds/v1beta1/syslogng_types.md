---
title: SyslogNGSpec
weight: 200
generated_file: true
---

## SyslogNGSpec

SyslogNGSpec defines the desired state of SyslogNG

### bufferVolumeMetrics (*BufferMetrics, optional) {#syslogngspec-buffervolumemetrics}

Default: -

### bufferVolumeMetricsService (*typeoverride.Service, optional) {#syslogngspec-buffervolumemetricsservice}

Default: -

### configCheckPod (*typeoverride.PodSpec, optional) {#syslogngspec-configcheckpod}

Default: -

### globalOptions (*GlobalOptions, optional) {#syslogngspec-globaloptions}

Default: -

### jsonKeyDelim (string, optional) {#syslogngspec-jsonkeydelim}

Default: -

### jsonKeyPrefix (string, optional) {#syslogngspec-jsonkeyprefix}

Default: -

### logIWSize (int, optional) {#syslogngspec-logiwsize}

Default: -

### maxConnections (int, optional) {#syslogngspec-maxconnections}

Default: -

### metrics (*Metrics, optional) {#syslogngspec-metrics}

Default: -

### metricsService (*typeoverride.Service, optional) {#syslogngspec-metricsservice}

Default: -

### readinessDefaultCheck (ReadinessDefaultCheck, optional) {#syslogngspec-readinessdefaultcheck}

Default: -

### serviceAccount (*typeoverride.ServiceAccount, optional) {#syslogngspec-serviceaccount}

Default: -

### service (*typeoverride.Service, optional) {#syslogngspec-service}

Default: -

### skipRBACCreate (bool, optional) {#syslogngspec-skiprbaccreate}

Default: -

### statefulSet (*typeoverride.StatefulSet, optional) {#syslogngspec-statefulset}

Default: -

### tls (SyslogNGTLS, optional) {#syslogngspec-tls}

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

### log_level (*string, optional) {#globaloptions-log_level}

Log level https://axoflow.com/docs/axosyslog-core/chapter-global-options/reference-options/#global-options-log-level 

Default: -

### stats (*Stats, optional) {#globaloptions-stats}

Stats https://axoflow.com/docs/axosyslog-core/chapter-global-options/reference-options/#global-option-stats 

Default: -

### stats_freq (*int, optional) {#globaloptions-stats_freq}

deprecated use stats/freq from 4.1+ 

Default: -

### stats_level (*int, optional) {#globaloptions-stats_level}

deprecated use stats/level from 4.1+ 

Default: -


## Stats

### freq (*int, optional) {#stats-freq}

Default: -

### level (*int, optional) {#stats-level}

Default: -


