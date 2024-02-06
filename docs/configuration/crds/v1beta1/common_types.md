---
title: Common
weight: 200
generated_file: true
---

## ImageSpec

ImageSpec struct hold information about image specification

### imagePullSecrets ([]corev1.LocalObjectReference, optional) {#imagespec-imagepullsecrets}


### pullPolicy (string, optional) {#imagespec-pullpolicy}


### repository (string, optional) {#imagespec-repository}


### tag (string, optional) {#imagespec-tag}



## Metrics

Metrics defines the service monitor endpoints

### interval (string, optional) {#metrics-interval}


### path (string, optional) {#metrics-path}


### port (int32, optional) {#metrics-port}


### prometheusAnnotations (bool, optional) {#metrics-prometheusannotations}


### prometheusRules (bool, optional) {#metrics-prometheusrules}


### serviceMonitor (bool, optional) {#metrics-servicemonitor}


### serviceMonitorConfig (ServiceMonitorConfig, optional) {#metrics-servicemonitorconfig}


### timeout (string, optional) {#metrics-timeout}



## BufferMetrics

BufferMetrics defines the service monitor endpoints

###  (Metrics, required) {#buffermetrics-}


### mount_name (string, optional) {#buffermetrics-mount_name}



## ServiceMonitorConfig

ServiceMonitorConfig defines the ServiceMonitor properties

### additionalLabels (map[string]string, optional) {#servicemonitorconfig-additionallabels}


### honorLabels (bool, optional) {#servicemonitorconfig-honorlabels}


### metricRelabelings ([]*v1.RelabelConfig, optional) {#servicemonitorconfig-metricrelabelings}


### relabelings ([]*v1.RelabelConfig, optional) {#servicemonitorconfig-relabelings}


### scheme (string, optional) {#servicemonitorconfig-scheme}


### tlsConfig (*v1.TLSConfig, optional) {#servicemonitorconfig-tlsconfig}



## Security

Security defines Fluentd, FluentbitAgent deployment security properties

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



