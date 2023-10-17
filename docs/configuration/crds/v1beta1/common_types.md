---
title: Common
weight: 200
generated_file: true
---

## ImageSpec

ImageSpec struct hold information about image specification

### imagePullSecrets ([]corev1.LocalObjectReference, optional) {#imagespec-imagepullsecrets}

Default: -

### pullPolicy (string, optional) {#imagespec-pullpolicy}

Default: -

### repository (string, optional) {#imagespec-repository}

Default: -

### tag (string, optional) {#imagespec-tag}

Default: -


## Metrics

Metrics defines the service monitor endpoints

### interval (string, optional) {#metrics-interval}

Default: -

### path (string, optional) {#metrics-path}

Default: -

### port (int32, optional) {#metrics-port}

Default: -

### prometheusAnnotations (bool, optional) {#metrics-prometheusannotations}

Default: -

### prometheusRules (bool, optional) {#metrics-prometheusrules}

Default: -

### serviceMonitor (bool, optional) {#metrics-servicemonitor}

Default: -

### serviceMonitorConfig (ServiceMonitorConfig, optional) {#metrics-servicemonitorconfig}

Default: -

### timeout (string, optional) {#metrics-timeout}

Default: -


## BufferMetrics

BufferMetrics defines the service monitor endpoints

###  (Metrics, required) {#buffermetrics-}

Default: -

### mount_name (string, optional) {#buffermetrics-mount_name}

Default: -


## ServiceMonitorConfig

ServiceMonitorConfig defines the ServiceMonitor properties

### additionalLabels (map[string]string, optional) {#servicemonitorconfig-additionallabels}

Default: -

### honorLabels (bool, optional) {#servicemonitorconfig-honorlabels}

Default: -

### metricRelabelings ([]*v1.RelabelConfig, optional) {#servicemonitorconfig-metricrelabelings}

Default: -

### relabelings ([]*v1.RelabelConfig, optional) {#servicemonitorconfig-relabelings}

Default: -

### scheme (string, optional) {#servicemonitorconfig-scheme}

Default: -

### tlsConfig (*v1.TLSConfig, optional) {#servicemonitorconfig-tlsconfig}

Default: -


## Security

Security defines Fluentd, FluentbitAgent deployment security properties

### podSecurityContext (*corev1.PodSecurityContext, optional) {#security-podsecuritycontext}

Default: -

### podSecurityPolicyCreate (bool, optional) {#security-podsecuritypolicycreate}

Default: -

### roleBasedAccessControlCreate (*bool, optional) {#security-rolebasedaccesscontrolcreate}

Default: -

### securityContext (*corev1.SecurityContext, optional) {#security-securitycontext}

Default: -

### serviceAccount (string, optional) {#security-serviceaccount}

Default: -


## ReadinessDefaultCheck

ReadinessDefaultCheck Enable default readiness checks

### bufferFileNumber (bool, optional) {#readinessdefaultcheck-bufferfilenumber}

Default: -

### bufferFileNumberMax (int32, optional) {#readinessdefaultcheck-bufferfilenumbermax}

Default: -

### bufferFreeSpace (bool, optional) {#readinessdefaultcheck-bufferfreespace}

Enable default Readiness check it'll fail if the buffer volume free space exceeds the `readinessDefaultThreshold` percentage (90%). 

Default: -

### bufferFreeSpaceThreshold (int32, optional) {#readinessdefaultcheck-bufferfreespacethreshold}

Default: -

### failureThreshold (int32, optional) {#readinessdefaultcheck-failurethreshold}

Default: -

### initialDelaySeconds (int32, optional) {#readinessdefaultcheck-initialdelayseconds}

Default: -

### periodSeconds (int32, optional) {#readinessdefaultcheck-periodseconds}

Default: -

### successThreshold (int32, optional) {#readinessdefaultcheck-successthreshold}

Default: -

### timeoutSeconds (int32, optional) {#readinessdefaultcheck-timeoutseconds}

Default: -


