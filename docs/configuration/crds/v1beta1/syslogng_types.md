---
title: SyslogNGSpec
weight: 200
generated_file: true
---

## SyslogNGSpec

SyslogNGSpec defines the desired state of SyslogNG

### bufferVolumeMetrics (*BufferMetrics, optional) {#syslogngspec-buffervolumemetrics}


### bufferVolumeMetricsService (*typeoverride.Service, optional) {#syslogngspec-buffervolumemetricsservice}


### configCheckPod (*typeoverride.PodSpec, optional) {#syslogngspec-configcheckpod}


### globalOptions (*GlobalOptions, optional) {#syslogngspec-globaloptions}


### jsonKeyDelim (string, optional) {#syslogngspec-jsonkeydelim}


### jsonKeyPrefix (string, optional) {#syslogngspec-jsonkeyprefix}


### logIWSize (int, optional) {#syslogngspec-logiwsize}


### maxConnections (int, optional) {#syslogngspec-maxconnections}


### metrics (*Metrics, optional) {#syslogngspec-metrics}


### metricsService (*typeoverride.Service, optional) {#syslogngspec-metricsservice}


### readinessDefaultCheck (ReadinessDefaultCheck, optional) {#syslogngspec-readinessdefaultcheck}


### serviceAccount (*typeoverride.ServiceAccount, optional) {#syslogngspec-serviceaccount}


### service (*typeoverride.Service, optional) {#syslogngspec-service}


### skipRBACCreate (bool, optional) {#syslogngspec-skiprbaccreate}


### sourceDateParser (*SourceDateParser, optional) {#syslogngspec-sourcedateparser}

Parses date automatically from the timestamp registered by the container runtime. Note: json key prefix and delimiter are respected 


### sourceMetrics ([]filter.MetricsProbe, optional) {#syslogngspec-sourcemetrics}


### statefulSet (*typeoverride.StatefulSet, optional) {#syslogngspec-statefulset}


### tls (SyslogNGTLS, optional) {#syslogngspec-tls}



## SourceDateParser

### format (*string, optional) {#sourcedateparser-format}

Default: "%FT%T.%f%z" 


### template (*string, optional) {#sourcedateparser-template}

Default(depending on JSONKeyPrefix): "${json.time}" 



## SyslogNGTLS

SyslogNGTLS defines the TLS configs

### enabled (bool, required) {#syslogngtls-enabled}


### secretName (string, optional) {#syslogngtls-secretname}


### sharedKey (string, optional) {#syslogngtls-sharedkey}



## GlobalOptions

### log_level (*string, optional) {#globaloptions-log_level}

See the [AxoSyslog Core documentation](https://axoflow.com/docs/axosyslog-core/chapter-global-options/reference-options/#global-options-log-level). 


### stats (*Stats, optional) {#globaloptions-stats}

See the [AxoSyslog Core documentation](https://axoflow.com/docs/axosyslog-core/chapter-global-options/reference-options/#global-option-stats). 


### stats_freq (*int, optional) {#globaloptions-stats_freq}

Deprecated. Use stats/freq from 4.1+ 


### stats_level (*int, optional) {#globaloptions-stats_level}

Deprecated. Use stats/level from 4.1+ 



## Stats

### freq (*int, optional) {#stats-freq}


### level (*int, optional) {#stats-level}



