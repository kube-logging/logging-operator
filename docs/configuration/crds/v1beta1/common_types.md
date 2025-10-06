---
title: Common
weight: 200
generated_file: true
---

## BasicImageSpec

BasicImageSpec struct hold basic information about image specification

### repository (string, optional) {#basicimagespec-repository}


### tag (string, optional) {#basicimagespec-tag}



## ImageSpec

ImageSpec struct hold information about image specification

### imagePullSecrets ([]corev1.LocalObjectReference, optional) {#imagespec-imagepullsecrets}


### pullPolicy (string, optional) {#imagespec-pullpolicy}


### repository (string, optional) {#imagespec-repository}


### tag (string, optional) {#imagespec-tag}



## Metrics

Metrics defines the service monitor endpoints

### enabled (*bool, optional) {#metrics-enabled}

Enabled controls whether the metrics endpoint should be exposed. Defaults to false. When false, the metrics HTTP server will not be started and no metrics port will be exposed. 


### interval (string, optional) {#metrics-interval}


### path (string, optional) {#metrics-path}


### port (int32, optional) {#metrics-port}


### prometheusAnnotations (bool, optional) {#metrics-prometheusannotations}


### prometheusRules (bool, optional) {#metrics-prometheusrules}


### prometheusRulesOverride ([]PrometheusRulesOverride, optional) {#metrics-prometheusrulesoverride}


### serviceMonitor (bool, optional) {#metrics-servicemonitor}


### serviceMonitorConfig (ServiceMonitorConfig, optional) {#metrics-servicemonitorconfig}


### timeout (string, optional) {#metrics-timeout}



## PrometheusRulesOverride

### alert (string, optional) {#prometheusrulesoverride-alert}

Name of the alert. Must be a valid label value. Only one of `record` and `alert` must be set. 


### annotations (map[string]string, optional) {#prometheusrulesoverride-annotations}

Annotations to add to each alert. Only valid for alerting rules. 


### expr (*intstr.IntOrString, optional) {#prometheusrulesoverride-expr}

PromQL expression to evaluate. 


### for (*v1.Duration, optional) {#prometheusrulesoverride-for}

Alerts are considered firing once they have been returned for this long. +optional 


### keep_firing_for (*v1.NonEmptyDuration, optional) {#prometheusrulesoverride-keep_firing_for}

KeepFiringFor defines how long an alert will continue firing after the condition that triggered it has cleared. +optional 


### labels (map[string]string, optional) {#prometheusrulesoverride-labels}

Labels to add or overwrite. 


### record (string, optional) {#prometheusrulesoverride-record}

Name of the time series to output to. Must be a valid metric name. Only one of `record` and `alert` must be set. 



## BufferMetrics

BufferMetrics defines the service monitor endpoints

###  (Metrics, required) {#buffermetrics-}


### mount_name (string, optional) {#buffermetrics-mount_name}



## ServiceMonitorConfig

ServiceMonitorConfig defines the ServiceMonitor properties

### additionalLabels (map[string]string, optional) {#servicemonitorconfig-additionallabels}


### honorLabels (bool, optional) {#servicemonitorconfig-honorlabels}


### metricRelabelings ([]v1.RelabelConfig, optional) {#servicemonitorconfig-metricrelabelings}


### relabelings ([]v1.RelabelConfig, optional) {#servicemonitorconfig-relabelings}


### scheme (string, optional) {#servicemonitorconfig-scheme}


### tlsConfig (*v1.TLSConfig, optional) {#servicemonitorconfig-tlsconfig}



## Security

Security defines Fluentd, FluentbitAgent deployment security properties

### createOpenShiftSCC (*bool, optional) {#security-createopenshiftscc}


### podSecurityContext (*corev1.PodSecurityContext, optional) {#security-podsecuritycontext}


### podSecurityPolicyCreate (bool, optional) {#security-podsecuritypolicycreate}

Warning: this is not supported anymore and does nothing 


### roleBasedAccessControlCreate (*bool, optional) {#security-rolebasedaccesscontrolcreate}


### securityContext (*corev1.SecurityContext, optional) {#security-securitycontext}


### serviceAccount (string, optional) {#security-serviceaccount}



## ReadinessDefaultCheck

ReadinessDefaultCheck Enable default readiness checks

### bufferFileNumber (bool, optional) {#readinessdefaultcheck-bufferfilenumber}


### bufferFileNumberMax (int32, optional) {#readinessdefaultcheck-bufferfilenumbermax}


### bufferFreeSpace (bool, optional) {#readinessdefaultcheck-bufferfreespace}

Enable default Readiness check it'll fail if the buffer volume free space exceeds the `readinessDefaultThreshold` percentage (90%). 


### bufferFreeSpaceThreshold (int32, optional) {#readinessdefaultcheck-bufferfreespacethreshold}


### failureThreshold (int32, optional) {#readinessdefaultcheck-failurethreshold}


### initialDelaySeconds (int32, optional) {#readinessdefaultcheck-initialdelayseconds}


### periodSeconds (int32, optional) {#readinessdefaultcheck-periodseconds}


### successThreshold (int32, optional) {#readinessdefaultcheck-successthreshold}


### timeoutSeconds (int32, optional) {#readinessdefaultcheck-timeoutseconds}



