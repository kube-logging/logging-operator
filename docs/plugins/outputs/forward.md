---
title: Forward
weight: 200
generated_file: true
---

### ForwardOutput
| Variable Name | Type | Required | Default | Description |
|---|---|---|---|---|
| servers | []FluentdServer | Yes | - | Server definitions at least one is required<br>[Server](#fluentd-server)<br> |
| require_ack_response | bool | No | - | Change the protocol to at-least-once. The plugin waits the ack from destination's in_forward plugin.<br> |
| ack_response_timeout | int | No |  190 | This option is used when require_ack_response is true. This default value is based on popular tcp_syn_retries. <br> |
| send_timeout | int | No |  60 | The timeout time when sending event logs. <br> |
| connect_timeout | int | No | - | The timeout time for socket connect. When the connection timed out during establishment, Errno::ETIMEDOUT is raised.<br> |
| recover_wait | int | No |  10 | The wait time before accepting a server fault recovery. <br> |
| heartbeat_type | string | No | - | The transport protocol to use for heartbeats. Set "none" to disable heartbeat. [transport, tcp, udp, none]<br> |
| heartbeat_interval | int | No |  1 | The interval of the heartbeat packer. <br> |
| phi_failure_detector | bool | No |  true | Use the "Phi accrual failure detector" to detect server failure. <br> |
| phi_threshold | int | No |  16 | The threshold parameter used to detect server faults. <br>`phi_threshold` is deeply related to `heartbeat_interval`. If you are using longer `heartbeat_interval`, please use the larger `phi_threshold`. Otherwise you will see frequent detachments of destination servers. The default value 16 is tuned for `heartbeat_interval` 1s.<br> |
| hard_timeout | int | No |  60 | The hard timeout used to detect server failure. The default value is equal to the send_timeout parameter. <br> |
| expire_dns_cache | int | No | - | Set TTL to expire DNS cache in seconds. Set 0 not to use DNS Cache. (defult: 0)<br> |
| dns_round_robin | bool | No | - | Enable client-side DNS round robin. Uniform randomly pick an IP address to send data when a hostname has several IP addresses.<br>`heartbeat_type udp` is not available with `dns_round_robin true`. Use `heartbeat_type tcp` or `heartbeat_type none`.<br> |
| ignore_network_errors_at_startup | bool | No | - | Ignore DNS resolution and errors at startup time.<br> |
| tls_version | string | No |  TLSv1_2 | The default version of TLS transport. [TLSv1_1, TLSv1_2] <br> |
| tls_ciphers | string | No |  ALL:!aNULL:!eNULL:!SSLv2 | The cipher configuration of TLS transport. <br> |
| tls_insecure_mode | bool | No |  false | Skip all verification of certificates or not. <br> |
| tls_allow_self_signed_cert | bool | No |  false | Allow self signed certificates or not. <br> |
| tls_verify_hostname | bool | No |  true | Verify hostname of servers and certificates or not in TLS transport. <br> |
| tls_cert_path | *secret.Secret | No | - | The additional CA certificate path for TLS.<br> |
| tls_client_cert_path | *secret.Secret | No | - | The client certificate path for TLS<br> |
| tls_client_private_key_path | *secret.Secret | No | - | The client private key path for TLS.<br> |
| tls_client_private_key_passphrase | *secret.Secret | No | - | The client private key passphrase for TLS.<br> |
| tls_cert_thumbprint | string | No | - | The certificate thumbprint for searching from Windows system certstore This parameter is for Windows only.<br> |
| tls_cert_logical_store_name | string | No | - | The certificate logical store name on Windows system certstore. This parameter is for Windows only.<br> |
| tls_cert_use_enterprise_store | bool | No | - | Enable to use certificate enterprise store on Windows system certstore. This parameter is for Windows only.<br> |
| keepalive | bool | No |  false | Enable keepalive connection. <br> |
| keepalive_timeout | int | No |  0 | Expired time of keepalive. Default value is nil, which means to keep connection as long as possible. <br> |
| security | *common.Security | No | - | [Security](../../common/security/)<br> |
| verify_connection_at_startup | bool | No |  false | Verify that a connection can be made with one of out_forward nodes at the time of startup. <br> |
| buffer | *Buffer | No | - | [Buffer](../buffer/)<br> |
### Fluentd Server
#### server

| Variable Name | Type | Required | Default | Description |
|---|---|---|---|---|
| host | string | Yes | - | The IP address or host name of the server.<br> |
| name | string | No | - | The name of the server. Used for logging and certificate verification in TLS transport (when host is address).<br> |
| port | int | No |  24224 | The port number of the host. Note that both TCP packets (event stream) and UDP packets (heartbeat message) are sent to this port. <br> |
| shared_key | *secret.Secret | No | - | The shared key per server.<br> |
| username | *secret.Secret | No | - | The username for authentication.<br> |
| password | *secret.Secret | No | - | The password for authentication.<br> |
| standby | bool | No | - | Marks a node as the standby node for an Active-Standby model between Fluentd nodes. When an active node goes down, the standby node is promoted to an active node. The standby node is not used by the out_forward plugin until then.<br> |
| weight | int | No |  60 | The load balancing weight. If the weight of one server is 20 and the weight of the other server is 30, events are sent in a 2:3 ratio. .<br> |
