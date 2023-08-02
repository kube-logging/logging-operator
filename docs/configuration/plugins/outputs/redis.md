---
title: Redis
weight: 200
generated_file: true
---

# Redis plugin for Fluentd
## Overview
 Sends logs to Redis endpoints.
 More info at https://github.com/fluent-plugins-nursery/fluent-plugin-redis

 ## Example output configurations
 ```yaml
 spec:

	redis:
	  host: redis-master.prod.svc.cluster.local
	  buffer:
	    tags: "[]"
	    flush_interval: 10s

 ```

## Configuration
## Output Config

### host (string, optional) {#output config-host}

Host Redis endpoint  

Default:  localhost

### port (int, optional) {#output config-port}

Port of the Redis server  

Default:  6379

### db_number (int, optional) {#output config-db_number}

DbNumber database number is optional.  

Default:  0

### password (*secret.Secret, optional) {#output config-password}

Redis Server password 

Default: -

### insert_key_prefix (string, optional) {#output config-insert_key_prefix}

insert_key_prefix  

Default:  "${tag}"

### strftime_format (string, optional) {#output config-strftime_format}

strftime_format Users can set strftime format.  

Default:  "%s"

### allow_duplicate_key (bool, optional) {#output config-allow_duplicate_key}

allow_duplicate_key Allow insert key duplicate. It will work as update values.  

Default:  false

### ttl (int, optional) {#output config-ttl}

ttl If 0 or negative value is set, ttl is not set in each key. 

Default: -

### format (*Format, optional) {#output config-format}

[Format](../format/) 

Default: -

### buffer (*Buffer, optional) {#output config-buffer}

[Buffer](../buffer/) 

Default: -

### slow_flush_log_threshold (string, optional) {#output config-slow_flush_log_threshold}

The threshold for chunk flush performance check. Parameter type is float, not time, default: 20.0 (seconds) If chunk flush takes longer time than this threshold, fluentd logs warning message and increases metric fluentd_output_status_slow_flush_count. 

Default: -


