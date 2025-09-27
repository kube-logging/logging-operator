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

Application Id 


### automatically_recover (bool, optional) {#output config-automatically_recover}

Automatic network failure recovery 


### buffer (*Buffer, optional) {#output config-buffer}

[Buffer](../buffer/) 


### connection_timeout (int, optional) {#output config-connection_timeout}

Connection Timeout in seconds 


### content_encoding (string, optional) {#output config-content_encoding}

Message content encoding 


### content_type (string, optional) {#output config-content_type}

Message content type 


### continuation_timeout (int, optional) {#output config-continuation_timeout}

Continuation Timeout in seconds 


### exchange (string, required) {#output config-exchange}

Name of the exchange 


### exchange_durable (bool, optional) {#output config-exchange_durable}

Exchange durability 


### exchange_no_declare (string, optional) {#output config-exchange_no_declare}

Weather to declare exchange or not 


### exchange_type (string, required) {#output config-exchange_type}

Type of the exchange 


### expiration (int, optional) {#output config-expiration}

Message message time-to-live in seconds 


### format (*Format, optional) {#output config-format}

[Format](../format/) 


### frame_max (int, optional) {#output config-frame_max}

Maximum permissible size of a frame 


### heartbeat (int, optional) {#output config-heartbeat}

Heartbeat Timeout in seconds 


### host (string, optional) {#output config-host}

Host 


### hosts ([]string, optional) {#output config-hosts}

Hosts 


### id_key (string, optional) {#output config-id_key}

Id to specify message_id 


### message_type (string, optional) {#output config-message_type}

Message type 


### network_recovery_interval (int, optional) {#output config-network_recovery_interval}

Network Recovery Interval in seconds 


### pass (*secret.Secret, optional) {#output config-pass}

Pass 


### persistent (bool, optional) {#output config-persistent}

Messages are persistent to disk 


### port (int, optional) {#output config-port}

Port 


### priority (int, optional) {#output config-priority}

Message priority 


### recovery_attempts (int, optional) {#output config-recovery_attempts}

Recovery Attempts 


### routing_key (string, optional) {#output config-routing_key}

Routing key to route messages 


### tls (bool, optional) {#output config-tls}

Enable TLS or not 


### tls_ca_certificates ([]string, optional) {#output config-tls_ca_certificates}

Path to TLS CA certificates files 


### tls_cert (string, optional) {#output config-tls_cert}

Path to TLS certificate file 


### tls_key (string, optional) {#output config-tls_key}

Path to TLS key file 


### timestamp (bool, optional) {#output config-timestamp}

Time of record is used as timestamp in AMQP message 


### user (*secret.Secret, optional) {#output config-user}

User 


### vhost (string, optional) {#output config-vhost}

VHost 


### verify_peer (bool, optional) {#output config-verify_peer}

Verify Peer or not 



