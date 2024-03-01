---
title: VMware Log Intelligence
weight: 200
generated_file: true
---

# VMware Log Intelligence output plugin for Fluentd
## Overview

For details, see [https://github.com/vmware/fluent-plugin-vmware-log-intelligence](https://github.com/vmware/fluent-plugin-vmware-log-intelligence).
## Example output configurations
```yaml
spec:
  vmwarelogintelligence:
    endpoint_url: https://data.upgrade.symphony-dev.com/le-mans/v1/streams/ingestion-pipeline-stream
	verify_ssl: true
	http_compress: false
	headers:
	  Content-Type: "application/json"
	  Authorization: "Bearer 12345"
	  structure: simple
	buffer:
	  chunk_limit_records: 300
	  flush_interval: 3s
	  retry_max_times: 3
```


## Configuration
## VMwareLogIntelligence

### buffer (*Buffer, optional) {#vmwarelogintelligence-buffer}

[Buffer](../buffer/) 


### endpoint_url (string, required) {#vmwarelogintelligence-endpoint_url}

Log Intelligence endpoint to send logs to https://github.com/vmware/fluent-plugin-vmware-log-intelligence?tab=readme-ov-file#label-endpoint_url 


### format (*Format, optional) {#vmwarelogintelligence-format}

[Format](../format/) 


### http_compress (*bool, optional) {#vmwarelogintelligence-http_compress}

Compress http request https://github.com/vmware/fluent-plugin-vmware-log-intelligence?tab=readme-ov-file#label-http_compress 


### headers (LogIntelligenceHeaders, required) {#vmwarelogintelligence-headers}

Required headers for sending logs to VMware Log Intelligence https://github.com/vmware/fluent-plugin-vmware-log-intelligence?tab=readme-ov-file#label-3Cheaders-3E 


### verify_ssl (bool, required) {#vmwarelogintelligence-verify_ssl}

Verify SSL (default: true) https://github.com/vmware/fluent-plugin-vmware-log-intelligence?tab=readme-ov-file#label-verify_ssl 

Default: true


## VMwareLogIntelligenceHeaders

headers
https://github.com/vmware/fluent-plugin-vmware-log-intelligence?tab=readme-ov-file#label-3Cheaders-3E

### authorization (*secret.Secret, required) {#vmwarelogintelligenceheaders-authorization}

Authorization Bearer token for http request to VMware Log Intelligence [Secret](../secret/) 


### content_type (string, required) {#vmwarelogintelligenceheaders-content_type}

Content Type for http request to VMware Log Intelligence 

Default: application/json

### structure (string, required) {#vmwarelogintelligenceheaders-structure}

Structure for http request to VMware Log Intelligence 

Default: simple


## LogIntelligenceHeadersOut

LogIntelligenceHeadersOut is used to convert the input LogIntelligenceHeaders to a fluentd
output that uses the correct key names for the VMware Log Intelligence plugin. This allows the
Ouput to accept the config is snake_case (as other output plugins do) but output the fluentd
<headers> config with the proper key names (ie. content_type -> Content-Type)

### Authorization (*secret.Secret, required) {#logintelligenceheadersout-authorization}

Authorization Bearer token for http request to VMware Log Intelligence 


### Content-Type (string, required) {#logintelligenceheadersout-content-type}

Content Type for http request to VMware Log Intelligence 

Default: application/json

### structure (string, required) {#logintelligenceheadersout-structure}

Structure for http request to VMware Log Intelligence 

Default: simple


