---
title: LogicMonitor Logs
weight: 200
generated_file: true
---

# LogicMonitor Logs output plugin for Fluentd
## Overview

LogicMonitor Logs output plugin for Fluentd

Sends log records to LogicMonitor Logs via the LM API.

For details, see [https://github.com/logicmonitor/lm-logs-fluentd](https://github.com/logicmonitor/lm-logs-fluentd).

## Example output configurations

```yaml
spec:
  lmLogs:
    company_name: mycompany
    access_id:
      valueFrom:
        secretKeyRef:
          name: lm-credentials
          key: access_id
    access_key:
      valueFrom:
        secretKeyRef:
          name: lm-credentials
          key: access_key
    resource_mapping: '{"kubernetes.host": "system.hostname"}'
    flush_interval: 60s
    debug: false
```


## Configuration
## LogicMonitorLogs

### access_id (*secret.Secret, optional) {#logicmonitorlogs-access_id}

LM API Token access ID [Secret](../secret/) 


### access_key (*secret.Secret, optional) {#logicmonitorlogs-access_key}

LM API Token access key [Secret](../secret/) 


### bearer_token (*secret.Secret, optional) {#logicmonitorlogs-bearer_token}

LM API Bearer Token. Either specify access_id and access_key both or bearer_token. If all specified, LMv1 token(access_id and access_key) will be used for authentication with LogicMonitor [Secret](../secret/) 


### buffer (*Buffer, optional) {#logicmonitorlogs-buffer}

[Buffer](../buffer/) 


### company_domain (string, optional) {#logicmonitorlogs-company_domain}

LogicMonitor account domain. For eg. for url test.logicmonitor.com, company_domain is logicmonitor.com (default: logicmonitor.com) 

Default: logicmonitor.com

### company_name (string, required) {#logicmonitorlogs-company_name}

LogicMonitor account name 


### debug (*bool, optional) {#logicmonitorlogs-debug}

When true, logs more information to the fluentd console 


### device_less_logs (*bool, optional) {#logicmonitorlogs-device_less_logs}

When true, do not map log with any resource. record must have service when true

Default: false

### flush_interval (string, optional) {#logicmonitorlogs-flush_interval}

Defines the time in seconds to wait before sending batches of logs to LogicMonitor (default: 60s) 

Default: 60s

### force_encoding (string, optional) {#logicmonitorlogs-force_encoding}

Specify charset when logs contains invalid utf-8 characters 


### format (*Format, optional) {#logicmonitorlogs-format}

[Format](../format/) 


### http_proxy (string, optional) {#logicmonitorlogs-http_proxy}

http proxy string eg. http://user:pass@proxy.server:port 


### include_metadata (*bool, optional) {#logicmonitorlogs-include_metadata}

When true, appends additional metadata to the log

Default: false

### resource_mapping (string, optional) {#logicmonitorlogs-resource_mapping}

The mapping that defines the source of the log event to the LM resource. In this case, the <event_key> in the incoming event is mapped to the value of <lm_property> 



