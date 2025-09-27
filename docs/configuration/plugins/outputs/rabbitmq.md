---
title: RabbitMQ
weight: 200
generated_file: true
---

# RabbitMQ plugin for Fluentd
## Overview

Sends logs to RabbitMQ Queues. For details, see [https://github.com/nttcom/fluent-plugin-rabbitmq](https://github.com/nttcom/fluent-plugin-rabbitmq).

## Example output configurations

```yaml
spec:
  rabbitmq:
    host: rabbitmq-master.prod.svc.cluster.local
    buffer:
      tags: "[]"
      flush_interval: 10s
```


## Configuration
## Output Config

### app_id (int, optional) {#output config-app_id}


### automatically_recover (bool, optional) {#output config-automatically_recover}


### buffer (*Buffer, optional) {#output config-buffer}

[Buffer](../buffer/) 


### connection_timeout (int, optional) {#output config-connection_timeout}


### content_encoding (string, optional) {#output config-content_encoding}


### content_type (string, optional) {#output config-content_type}


### continuation_timeout (int, optional) {#output config-continuation_timeout}


### exchange (string, required) {#output config-exchange}


### exchange_durable (string, optional) {#output config-exchange_durable}


### exchange_no_declare (string, optional) {#output config-exchange_no_declare}


### exchange_type (string, required) {#output config-exchange_type}


### expiration (int, optional) {#output config-expiration}


### format (*Format, optional) {#output config-format}

[Format](../format/) 


### frame_max (int, optional) {#output config-frame_max}


### heartbeat (int, optional) {#output config-heartbeat}


### host (string, optional) {#output config-host}


### hosts ([]string, optional) {#output config-hosts}


### id_key (string, optional) {#output config-id_key}


### message_type (string, optional) {#output config-message_type}


### network_recovery_interval (int, optional) {#output config-network_recovery_interval}


### pass (*secret.Secret, optional) {#output config-pass}


### persistent (bool, optional) {#output config-persistent}


### port (int, optional) {#output config-port}


### priority (int, optional) {#output config-priority}


### recovery_attempts (int, optional) {#output config-recovery_attempts}


### routing_key (string, optional) {#output config-routing_key}


### tls (bool, optional) {#output config-tls}


### tls_ca_certificates ([]string, optional) {#output config-tls_ca_certificates}


### tls_cert (string, optional) {#output config-tls_cert}


### tls_key (string, optional) {#output config-tls_key}


### timestamp (string, optional) {#output config-timestamp}


### user (*secret.Secret, optional) {#output config-user}


### vhost (string, optional) {#output config-vhost}


### verify_peer (bool, optional) {#output config-verify_peer}



