---
title: Syslog
weight: 200
generated_file: true
---

# [Syslog Output](https://github.com/cloudfoundry/fluent-plugin-syslog_rfc5424)
## Overview
 Fluentd output plugin for remote syslog with RFC5424 headers logs.

## Configuration
## SyslogOutputConfig

### host (string, required) {#syslogoutputconfig-host}

Destination host address 

Default: -

### port (int, optional) {#syslogoutputconfig-port}

Destination host port  

Default:  "514"

### transport (string, optional) {#syslogoutputconfig-transport}

Transport Protocol  

Default:  "tls"

### insecure (*bool, optional) {#syslogoutputconfig-insecure}

skip ssl validation  

Default:  false

### verify_fqdn (*bool, optional) {#syslogoutputconfig-verify_fqdn}

verify_fqdn  

Default:  nil

### enable_system_cert_store (*bool, optional) {#syslogoutputconfig-enable_system_cert_store}

cert_store to set ca_certificate for ssl context 

Default: -

### trusted_ca_path (*secret.Secret, optional) {#syslogoutputconfig-trusted_ca_path}

file path to ca to trust 

Default: -

### client_cert_path (*secret.Secret, optional) {#syslogoutputconfig-client_cert_path}

file path for private_key_path 

Default: -

### private_key_path (*secret.Secret, optional) {#syslogoutputconfig-private_key_path}

file path for private_key_path 

Default: -

### private_key_passphrase (*secret.Secret, optional) {#syslogoutputconfig-private_key_passphrase}

PrivateKeyPassphrase for private key   

Default:  "nil"

### allow_self_signed_cert (*bool, optional) {#syslogoutputconfig-allow_self_signed_cert}

allow_self_signed_cert for mutual tls  

Default:  false

### fqdn (string, optional) {#syslogoutputconfig-fqdn}

Fqdn  

Default:  "nil"

### version (string, optional) {#syslogoutputconfig-version}

TLS Version   

Default:  "TLSv1_2"

### format (*FormatRfc5424, optional) {#syslogoutputconfig-format}

[Format](../format_rfc5424/) 

Default: -

### buffer (*Buffer, optional) {#syslogoutputconfig-buffer}

[Buffer](../buffer/) 

Default: -

### slow_flush_log_threshold (string, optional) {#syslogoutputconfig-slow_flush_log_threshold}

The threshold for chunk flush performance check. Parameter type is float, not time, default: 20.0 (seconds) If chunk flush takes longer time than this threshold, fluentd logs warning message and increases metric fluentd_output_status_slow_flush_count. 

Default: -


 ## Example `File` output configurations
 ```yaml
 apiVersion: logging.banzaicloud.io/v1beta1
 kind: Output
 metadata:

	name: demo-output

 spec:

	syslog:
	  host: SYSLOG-HOST
	  port: 123
	  format:
	    app_name_field: example.custom_field_1
	    proc_id_field: example.custom_field_2
	  buffer:
	    timekey: 1m
	    timekey_wait: 10s
	    timekey_use_utc: true

 ```

 #### Fluentd Config Result
 ```

	 <match **>
		@type syslog_rfc5424
		@id test_syslog
		host SYSLOG-HOST
		port 123
	 <format>
	   @type syslog_rfc5424
	   app_name_field example.custom_field_1
	   proc_id_field example.custom_field_2
	 </format>
		<buffer tag,time>
		  @type file
		  path /buffers/test_file.*.buffer
		  retry_forever true
		  timekey 1m
		  timekey_use_utc true
		  timekey_wait 30s
		</buffer>
	 </match>

 ```

---
