---
title: Forward
weight: 200
generated_file: true
---

## ForwardOutput

### ack_response_timeout (int, optional) {#forwardoutput-ack_response_timeout}

This option is used when require_ack_response is true. This default value is based on popular tcp_syn_retries.

Default: 190

### buffer (*Buffer, optional) {#forwardoutput-buffer}

[Buffer](../buffer/) 


### connect_timeout (int, optional) {#forwardoutput-connect_timeout}

The timeout time for socket connect. When the connection timed out during establishment, Errno::ETIMEDOUT is raised. 


### dns_round_robin (bool, optional) {#forwardoutput-dns_round_robin}

Enable client-side DNS round robin. Uniform randomly pick an IP address to send data when a hostname has several IP addresses. `heartbeat_type udp` is not available with `dns_round_robin true`. Use `heartbeat_type tcp` or `heartbeat_type none`. 


### expire_dns_cache (int, optional) {#forwardoutput-expire_dns_cache}

Set TTL to expire DNS cache in seconds. Set 0 not to use DNS Cache.

Default: 0

### servers ([]FluentdServer, required) {#forwardoutput-servers}

Server definitions at least one is required [Server](#fluentd-server) 


### hard_timeout (int, optional) {#forwardoutput-hard_timeout}

The hard timeout used to detect server failure. The default value is equal to the send_timeout parameter.

Default: 60

### heartbeat_interval (int, optional) {#forwardoutput-heartbeat_interval}

The interval of the heartbeat packer.

Default: 1

### heartbeat_type (string, optional) {#forwardoutput-heartbeat_type}

The transport protocol to use for heartbeats. Set "none" to disable heartbeat. [transport, tcp, udp, none] 


### ignore_network_errors_at_startup (bool, optional) {#forwardoutput-ignore_network_errors_at_startup}

Ignore DNS resolution and errors at startup time. 


### keepalive (bool, optional) {#forwardoutput-keepalive}

Enable keepalive connection.

Default: false

### keepalive_timeout (int, optional) {#forwardoutput-keepalive_timeout}

Expired time of keepalive. Default value is nil, which means to keep connection as long as possible.

Default: 0

### phi_failure_detector (bool, optional) {#forwardoutput-phi_failure_detector}

Use the "Phi accrual failure detector" to detect server failure.

Default: true

### phi_threshold (int, optional) {#forwardoutput-phi_threshold}

The threshold parameter used to detect server faults.  `phi_threshold` is deeply related to `heartbeat_interval`. If you are using longer `heartbeat_interval`, please use the larger `phi_threshold`. Otherwise you will see frequent detachments of destination servers. The default value 16 is tuned for `heartbeat_interval` 1s.

Default: 16

### recover_wait (int, optional) {#forwardoutput-recover_wait}

The wait time before accepting a server fault recovery.

Default: 10

### require_ack_response (bool, optional) {#forwardoutput-require_ack_response}

Change the protocol to at-least-once. The plugin waits the ack from destination's in_forward plugin. 


### security (*common.Security, optional) {#forwardoutput-security}

[Security](../../common/security/) 


### send_timeout (int, optional) {#forwardoutput-send_timeout}

The timeout time when sending event logs.

Default: 60

### slow_flush_log_threshold (string, optional) {#forwardoutput-slow_flush_log_threshold}

The threshold for chunk flush performance check. Parameter type is float, not time, default: 20.0 (seconds) If chunk flush takes longer time than this threshold, fluentd logs warning message and increases metric fluentd_output_status_slow_flush_count. 


### tls_allow_self_signed_cert (bool, optional) {#forwardoutput-tls_allow_self_signed_cert}

Allow self signed certificates or not.

Default: false

### tls_cert_logical_store_name (string, optional) {#forwardoutput-tls_cert_logical_store_name}

The certificate logical store name on Windows system certstore. This parameter is for Windows only. 


### tls_cert_path (*secret.Secret, optional) {#forwardoutput-tls_cert_path}

The additional CA certificate path for TLS. 


### tls_cert_thumbprint (string, optional) {#forwardoutput-tls_cert_thumbprint}

The certificate thumbprint for searching from Windows system certstore This parameter is for Windows only. 


### tls_cert_use_enterprise_store (bool, optional) {#forwardoutput-tls_cert_use_enterprise_store}

Enable to use certificate enterprise store on Windows system certstore. This parameter is for Windows only. 


### tls_ciphers (string, optional) {#forwardoutput-tls_ciphers}

The cipher configuration of TLS transport.

Default: ALL:!aNULL:!eNULL:!SSLv2

### tls_client_cert_path (*secret.Secret, optional) {#forwardoutput-tls_client_cert_path}

The client certificate path for TLS 


### tls_client_private_key_passphrase (*secret.Secret, optional) {#forwardoutput-tls_client_private_key_passphrase}

The client private key passphrase for TLS. 


### tls_client_private_key_path (*secret.Secret, optional) {#forwardoutput-tls_client_private_key_path}

The client private key path for TLS. 


### tls_insecure_mode (bool, optional) {#forwardoutput-tls_insecure_mode}

Skip all verification of certificates or not.

Default: false

### tls_verify_hostname (bool, optional) {#forwardoutput-tls_verify_hostname}

Verify hostname of servers and certificates or not in TLS transport.

Default: true

### tls_version (string, optional) {#forwardoutput-tls_version}

The default version of TLS transport. [TLSv1_1, TLSv1_2]

Default: TLSv1_2

### transport (string, optional) {#forwardoutput-transport}

The transport protocol to use [ tcp, tls ] 


### verify_connection_at_startup (bool, optional) {#forwardoutput-verify_connection_at_startup}

Verify that a connection can be made with one of out_forward nodes at the time of startup.

Default: false


## Fluentd Server

server

### host (string, required) {#fluentd server-host}

The IP address or host name of the server. 


### name (string, optional) {#fluentd server-name}

The name of the server. Used for logging and certificate verification in TLS transport (when host is address). 


### password (*secret.Secret, optional) {#fluentd server-password}

The password for authentication. 


### port (int, optional) {#fluentd server-port}

The port number of the host. Note that both TCP packets (event stream) and UDP packets (heartbeat message) are sent to this port.

Default: 24224

### shared_key (*secret.Secret, optional) {#fluentd server-shared_key}

The shared key per server. 


### standby (bool, optional) {#fluentd server-standby}

Marks a node as the standby node for an Active-Standby model between Fluentd nodes. When an active node goes down, the standby node is promoted to an active node. The standby node is not used by the out_forward plugin until then. 


### username (*secret.Secret, optional) {#fluentd server-username}

The username for authentication. 


### weight (int, optional) {#fluentd server-weight}

The load balancing weight. If the weight of one server is 20 and the weight of the other server is 30, events are sent in a 2:3 ratio. .

Default: 60


