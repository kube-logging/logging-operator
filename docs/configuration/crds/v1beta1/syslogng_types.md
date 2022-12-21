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

Default: -

### stats_freq (*int, optional) {#globaloptions-stats_freq}

Default: -


