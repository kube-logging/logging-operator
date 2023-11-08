---
title: Grafana Loki
weight: 200
generated_file: true
---

# Loki output plugin 
## Overview
 Fluentd output plugin to ship logs to a Loki server.
 More info at https://grafana.com/docs/loki/latest/clients/fluentd/
 >Example: [Store Nginx Access Logs in Grafana Loki with Logging Operator](../../../../quickstarts/loki-nginx/)

 ## Example output configurations
 ```yaml
 spec:

	loki:
	  url: http://loki:3100
	  buffer:
	    timekey: 1m
	    timekey_wait: 30s
	    timekey_use_utc: true

 ```

## Configuration
## Output Config

### url (string, optional) {#output config-url}

The url of the Loki server to send logs to.

Default: https://logs-us-west1.grafana.net

### username (*secret.Secret, optional) {#output config-username}

Specify a username if the Loki server requires authentication. [Secret](../secret/) 

Default: -

### password (*secret.Secret, optional) {#output config-password}

Specify password if the Loki server requires authentication. [Secret](../secret/) 

Default: -

### cert (*secret.Secret, optional) {#output config-cert}

TLS: parameters for presenting a client certificate [Secret](../secret/) 

Default: -

### key (*secret.Secret, optional) {#output config-key}

TLS: parameters for presenting a client certificate [Secret](../secret/) 

Default: -

### ca_cert (*secret.Secret, optional) {#output config-ca_cert}

TLS: CA certificate file for server certificate verification [Secret](../secret/) 

Default: -

### insecure_tls (*bool, optional) {#output config-insecure_tls}

TLS: disable server certificate verification

Default: false

### tenant (string, optional) {#output config-tenant}

Loki is a multi-tenant log storage platform and all requests sent must include a tenant. 

Default: -

### labels (Label, optional) {#output config-labels}

Set of labels to include with every Loki stream. 

Default: -

### extra_labels (map[string]string, optional) {#output config-extra_labels}

Set of extra labels to include with every Loki stream. 

Default: -

### line_format (string, optional) {#output config-line_format}

Format to use when flattening the record to a log line: json, key_value (default: key_value) 

Default: json

### extract_kubernetes_labels (*bool, optional) {#output config-extract_kubernetes_labels}

Extract kubernetes labels as loki labels

Default: false

### remove_keys ([]string, optional) {#output config-remove_keys}

Comma separated list of needless record keys to remove

Default: []

### drop_single_key (*bool, optional) {#output config-drop_single_key}

If a record only has 1 key, then just set the log line to the value and discard the key.

Default: false

### configure_kubernetes_labels (*bool, optional) {#output config-configure_kubernetes_labels}

Configure Kubernetes metadata in a Prometheus like format

Default: false

### include_thread_label (*bool, optional) {#output config-include_thread_label}

whether to include the fluentd_thread label when multiple threads are used for flushing.

Default: true

### buffer (*Buffer, optional) {#output config-buffer}

[Buffer](../buffer/) 

Default: -

### slow_flush_log_threshold (string, optional) {#output config-slow_flush_log_threshold}

The threshold for chunk flush performance check. Parameter type is float, not time, default: 20.0 (seconds) If chunk flush takes longer time than this threshold, fluentd logs warning message and increases metric fluentd_output_status_slow_flush_count. 

Default: -


