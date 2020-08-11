---
title: Http
weight: 200
generated_file: true
---

# Redis plugin for Fluentd
## Overview
 Sends logs to Redis endpoints.
 More info at https://github.com/fluent-plugins-nursery/fluent-plugin-redis

 #### Example output configurations
 ```
 spec:
   redis:
     host: redis-master.prod.svc.cluster.local
     buffer:
       tags: "[]"
       flush_interval: 10s
 ```

## Configuration
### Output Config
| Variable Name | Type | Required | Default | Description |
|---|---|---|---|---|
| host | string | No |  localhost | Host Redis endpoint <br> |
| port | int | No |  6379 | Port of the Redis server <br> |
| db_number | int | No |  0 | DbNumber database number is optional. <br> |
| password | *secret.Secret | No | - | Redis Server password<br> |
| insert_key_prefix | string | No |  "${tag}" | insert_key_prefix <br> |
| strftime_format | string | No |  "%s" | strftime_format Users can set strftime format. <br> |
| allow_duplicate_key | bool | No |  false | allow_duplicate_key Allow insert key duplicate. It will work as update values. <br> |
| ttl | int | No | - | ttl If 0 or negative value is set, ttl is not set in each key.<br> |
| format | *Format | No | - | [Format](../format/)<br> |
| buffer | *Buffer | No | - | [Buffer](../buffer/)<br> |
