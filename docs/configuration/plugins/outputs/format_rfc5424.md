---
title: Format rfc5424
weight: 200
generated_file: true
---

## FormatRfc5424

### app_name_field (string, optional) {#formatrfc5424-app_name_field}

Sets app name in syslog from field in fluentd, delimited by '.'

Default: app_name

### hostname_field (string, optional) {#formatrfc5424-hostname_field}

Sets host name in syslog from field in fluentd, delimited by '.'

Default: hostname

### log_field (string, optional) {#formatrfc5424-log_field}

Sets log in syslog from field in fluentd, delimited by '.'

Default: log

### message_id_field (string, optional) {#formatrfc5424-message_id_field}

Sets msg id in syslog from field in fluentd, delimited by '.'

Default: message_id

### proc_id_field (string, optional) {#formatrfc5424-proc_id_field}

Sets proc id in syslog from field in fluentd, delimited by '.'

Default: proc_id

### rfc6587_message_size (*bool, optional) {#formatrfc5424-rfc6587_message_size}

Prepends message length for syslog transmission

Default: true

### structured_data_field (string, optional) {#formatrfc5424-structured_data_field}

Sets structured data in syslog from field in fluentd, delimited by '.' (default structured_data) 


### type (string, optional) {#formatrfc5424-type}

Output line formatting: out_file,json,ltsv,csv,msgpack,hash,single_value

Default: json


