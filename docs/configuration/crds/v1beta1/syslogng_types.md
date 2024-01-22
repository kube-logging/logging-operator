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

Available in Logging operator version 4.5 and later. Set the maximum number of connections for the source. For details, see [documentation of the AxoSyslog syslog-ng distribution](https://axoflow.com/docs/axosyslog-core/chapter-routing-filters/concepts-flow-control/configuring-flow-control/). 


### metrics (*Metrics, optional) {#syslogngspec-metrics}


### metricsService (*typeoverride.Service, optional) {#syslogngspec-metricsservice}


### readinessDefaultCheck (ReadinessDefaultCheck, optional) {#syslogngspec-readinessdefaultcheck}


### serviceAccount (*typeoverride.ServiceAccount, optional) {#syslogngspec-serviceaccount}


### service (*typeoverride.Service, optional) {#syslogngspec-service}


### skipRBACCreate (bool, optional) {#syslogngspec-skiprbaccreate}


### sourceDateParser (*SourceDateParser, optional) {#syslogngspec-sourcedateparser}

Available in Logging operator version 4.5 and later. Parses date automatically from the timestamp registered by the container runtime. Note: `jsonKeyPrefix` and `jsonKeyDelim` are respected. 


### sourceMetrics ([]filter.MetricsProbe, optional) {#syslogngspec-sourcemetrics}

Available in Logging operator version 4.5 and later. Create [custom log metrics for sources and outputs]({{< relref "/docs/examples/custom-syslog-ng-metrics.md" >}}). 


### statefulSet (*typeoverride.StatefulSet, optional) {#syslogngspec-statefulset}


### tls (SyslogNGTLS, optional) {#syslogngspec-tls}



## SourceDateParser



Available in Logging operator version 4.5 and later.

Parses date automatically from the timestamp registered by the container runtime.
Note: `jsonKeyPrefix` and `jsonKeyDelim` are respected.
It is disabled by default, but if enabled, then the default settings parse the timestamp written by the container runtime and parsed by Fluent Bit using the `cri` or the `docker` parser.


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



