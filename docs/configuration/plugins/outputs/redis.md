---
title: Redis
weight: 200
generated_file: true
---

# Redis plugin for Fluentd
## Overview

Sends logs to Redis endpoints. For details, see [https://github.com/fluent-plugins-nursery/fluent-plugin-redis](https://github.com/fluent-plugins-nursery/fluent-plugin-redis).

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

### allow_duplicate_key (bool, optional) {#output config-allow_duplicate_key}

Allow inserting key duplicate. It will work as update values.

Default: false

### buffer (*Buffer, optional) {#output config-buffer}

[Buffer](../buffer/) 


### db_number (int, optional) {#output config-db_number}

DbNumber database number is optional.

Default: 0

### format (*Format, optional) {#output config-format}

[Format](../format/) 


### host (string, optional) {#output config-host}

Host Redis endpoint

Default: localhost

### insert_key_prefix (string, optional) {#output config-insert_key_prefix}

insert_key_prefix

Default: "${tag}"

### password (*secret.Secret, optional) {#output config-password}

Redis Server password 


### port (int, optional) {#output config-port}

Port of the Redis server

Default: 6379

### slow_flush_log_threshold (string, optional) {#output config-slow_flush_log_threshold}

The threshold for chunk flush performance check. Parameter type is float, not time, default: 20.0 (seconds) If chunk flush takes longer time than this threshold, Fluentd logs a warning message and increases the `fluentd_output_status_slow_flush_count` metric. 


### strftime_format (string, optional) {#output config-strftime_format}

Users can set strftime format.

Default: "%s"

### ttl (int, optional) {#output config-ttl}

If 0 or negative value is set, ttl is not set in each key. 



