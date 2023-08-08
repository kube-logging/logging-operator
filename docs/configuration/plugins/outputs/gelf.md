---
title: GELF
weight: 200
generated_file: true
---

# [GELF Output](https://github.com/hotschedules/fluent-plugin-gelf-hs)
## Overview
 Fluentd output plugin for GELF.

## Configuration
## Output Config

### host (string, required) {#output config-host}

Destination host 

Default: -

### port (int, required) {#output config-port}

Destination host port 

Default: -

### protocol (string, optional) {#output config-protocol}

Transport Protocol  

Default:  "udp"

### tls (*bool, optional) {#output config-tls}

Enable TlS  

Default:  false

### tls_options (map[string]string, optional) {#output config-tls_options}

TLS Options  - for options see https://github.com/graylog-labs/gelf-rb/blob/72916932b789f7a6768c3cdd6ab69a3c942dbcef/lib/gelf/transport/tcp_tls.rb#L7-L12 

Default:  {}


 ## Example `GELF` output configurations
 ```yaml
 apiVersion: logging.banzaicloud.io/v1beta1
 kind: Output
 metadata:

	name: gelf-output-sample

 spec:

	gelf:
	  host: gelf-host
	  port: 12201

 ```

 ## Fluentd Config Result
 ```

	 <match **>
		@type gelf
		@id test_gelf
		host gelf-host
		port 12201
	 </match>

 ```

---
