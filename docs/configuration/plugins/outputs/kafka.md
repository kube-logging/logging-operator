---
title: Kafka
weight: 200
generated_file: true
---

# Kafka output plugin for Fluentd
## Overview


For details, see [https://github.com/fluent/fluent-plugin-kafka](https://github.com/fluent/fluent-plugin-kafka).

For an example deployment, see [Transport Nginx Access Logs into Kafka with Logging Operator](../../../../quickstarts/kafka-nginx/).

## Example output configurations

```yaml
spec:
  kafka:
    brokers: kafka-headless.kafka.svc.cluster.local:29092
    default_topic: topic
    sasl_over_ssl: false
    format:
      type: json
    buffer:
      tags: topic
      timekey: 1m
      timekey_wait: 30s
      timekey_use_utc: true
```


## Configuration
## Kafka

Send your logs to Kafka.
Setting use_rdkafka to true opts for rdkafka2, which offers higher performance compared to ruby-kafka.
(Note: requires fluentd image version v1.16-4.9-full or higher)
-[more info](https://github.com/fluent/fluent-plugin-kafka#output-plugin)

### ack_timeout (int, optional) {#kafka-ack_timeout}

How long the producer waits for acks. The unit is seconds

Default: nil => Uses default of ruby-kafka library

### brokers (string, required) {#kafka-brokers}

The list of all seed brokers, with their host and port information. 


### buffer (*Buffer, optional) {#kafka-buffer}

[Buffer](../buffer/) 


### client_id (string, optional) {#kafka-client_id}

Client ID

Default: "kafka"

### compression_codec (string, optional) {#kafka-compression_codec}

The codec the producer uses to compress messages . The available options are gzip and snappy.

Default: nil

### default_message_key (string, optional) {#kafka-default_message_key}

The name of default message key .

Default: nil

### default_partition_key (string, optional) {#kafka-default_partition_key}

The name of default partition key .

Default: nil

### default_topic (string, optional) {#kafka-default_topic}

The name of default topic .

Default: nil

### discard_kafka_delivery_failed (bool, optional) {#kafka-discard_kafka_delivery_failed}

Discard the record where Kafka DeliveryFailed occurred

Default: false

### exclude_partion_key (bool, optional) {#kafka-exclude_partion_key}

Exclude Partition key

Default: false

### exclude_topic_key (bool, optional) {#kafka-exclude_topic_key}

Exclude Topic key

Default: false

### format (*Format, required) {#kafka-format}

[Format](../format/) 


### get_kafka_client_log (bool, optional) {#kafka-get_kafka_client_log}

Get Kafka Client log

Default: false

### headers (map[string]string, optional) {#kafka-headers}

Headers

Default: {}

### headers_from_record (map[string]string, optional) {#kafka-headers_from_record}

Headers from Record

Default: {}

### idempotent (bool, optional) {#kafka-idempotent}

Idempotent

Default: false

### kafka_agg_max_bytes (int, optional) {#kafka-kafka_agg_max_bytes}

Maximum value of total message size to be included in one batch transmission. .

Default: 4096

### kafka_agg_max_messages (int, optional) {#kafka-kafka_agg_max_messages}

Maximum number of messages to include in one batch transmission. .

Default: nil

### keytab (*secret.Secret, optional) {#kafka-keytab}


### max_send_limit_bytes (int, optional) {#kafka-max_send_limit_bytes}

Max byte size to send message to avoid MessageSizeTooLarge. Messages over the limit will be dropped

Default: no limit

### max_send_retries (int, optional) {#kafka-max_send_retries}

Number of times to retry sending of messages to a leader

Default: 1

### message_key_key (string, optional) {#kafka-message_key_key}

Message Key

Default: "message_key"

### partition_key (string, optional) {#kafka-partition_key}

Partition

Default: "partition"

### partition_key_key (string, optional) {#kafka-partition_key_key}

Partition Key

Default: "partition_key"

### password (*secret.Secret, optional) {#kafka-password}

Password when using PLAIN/SCRAM SASL authentication 


### principal (string, optional) {#kafka-principal}


### rdkafka_options (RdkafkaOptions, optional) {#kafka-rdkafka_options}


### required_acks (int, optional) {#kafka-required_acks}

The number of acks required per request .

Default: -1

### ssl_ca_cert (*secret.Secret, optional) {#kafka-ssl_ca_cert}

CA certificate 


### ssl_ca_certs_from_system (*bool, optional) {#kafka-ssl_ca_certs_from_system}

System's CA cert store

Default: false

### ssl_client_cert (*secret.Secret, optional) {#kafka-ssl_client_cert}

Client certificate 


### ssl_client_cert_chain (*secret.Secret, optional) {#kafka-ssl_client_cert_chain}

Client certificate chain 


### ssl_client_cert_key (*secret.Secret, optional) {#kafka-ssl_client_cert_key}

Client certificate key 


### ssl_verify_hostname (*bool, optional) {#kafka-ssl_verify_hostname}

Verify certificate hostname 


### sasl_over_ssl (*bool, optional) {#kafka-sasl_over_ssl}

SASL over SSL

Default: true

### scram_mechanism (string, optional) {#kafka-scram_mechanism}

If set, use SCRAM authentication with specified mechanism. When unset, default to PLAIN authentication 


### slow_flush_log_threshold (string, optional) {#kafka-slow_flush_log_threshold}

The threshold for chunk flush performance check. Parameter type is float, not time, default: 20.0 (seconds) If chunk flush takes longer time than this threshold, Fluentd logs a warning message and increases the  `fluentd_output_status_slow_flush_count` metric. 


### topic_key (string, optional) {#kafka-topic_key}

Topic Key

Default: "topic"

### use_default_for_unknown_topic (bool, optional) {#kafka-use_default_for_unknown_topic}

Use default for unknown topics

Default: false

### use_rdkafka (bool, optional) {#kafka-use_rdkafka}

Use rdkafka2 instead of the legacy kafka2 output plugin. This plugin requires fluentd image version v1.16-4.9-full or higher. 


### username (*secret.Secret, optional) {#kafka-username}

Username when using PLAIN/SCRAM SASL authentication 



## RdkafkaOptions

RdkafkaOptions represents the global configuration properties for librdkafka.

### allow.auto.create.topics (bool, optional) {#rdkafkaoptions-allow.auto.create.topics}

Allow automatic topic creation on the broker when subscribing to or assigning non-existent topics. The broker must also be configured with `auto.create.topics.enable=true` for this configuration to take effect. Note: the default value (true) for the producer is different from the default value (false) for the consumer. Further, the consumer default value is different from the Java consumer (true), and this property is not supported by the Java producer. Requires broker version >= 0.11.0.0, for older broker versions only the broker configuration applies. 


### api.version.fallback.ms (int, optional) {#rdkafkaoptions-api.version.fallback.ms}

Dictates how long the `broker.version.fallback` fallback is used in the case the ApiVersionRequest fails. 


### api.version.request (bool, optional) {#rdkafkaoptions-api.version.request}

Request broker's supported API versions to adjust functionality to available protocol features. If set to false, or the ApiVersionRequest fails, the fallback version `broker.version.fallback` will be used. **NOTE**: Depends on broker version >=0.10.0. If the request is not supported by (an older) broker the `broker.version.fallback` fallback is used. 


### api.version.request.timeout.ms (int, optional) {#rdkafkaoptions-api.version.request.timeout.ms}

Timeout for broker API version requests. 


### background_event_cb (string, optional) {#rdkafkaoptions-background_event_cb}

Background queue event callback (set with rd_kafka_conf_set_background_event_cb()) 


### bootstrap.servers (string, optional) {#rdkafkaoptions-bootstrap.servers}

Alias for `metadata.broker.list`: Initial list of brokers as a CSV list of broker host or host:port. The application may also use `rd_kafka_brokers_add()` to add brokers during runtime. 


### broker.address.family (string, optional) {#rdkafkaoptions-broker.address.family}

Allowed broker IP address families: any, v4, v6 


### broker.address.ttl (int, optional) {#rdkafkaoptions-broker.address.ttl}

How long to cache the broker address resolving results (milliseconds). 


### broker.version.fallback (string, optional) {#rdkafkaoptions-broker.version.fallback}

Older broker versions (before 0.10.0) provide no way for a client to query for supported protocol features (ApiVersionRequest, see `api.version.request`) making it impossible for the client to know what features it may use. As a workaround a user may set this property to the expected broker version and the client will automatically adjust its feature set accordingly if the ApiVersionRequest fails (or is disabled). The fallback broker version will be used for `api.version.fallback.ms`. Valid values are: 0.9.0, 0.8.2, 0.8.1, 0.8.0. Any other value >= 0.10, such as 0.10.2.1, enables ApiVersionRequests. 


### builtin.features (string, optional) {#rdkafkaoptions-builtin.features}

Indicates the builtin features for this build of librdkafka. An application can either query this value or attempt to set it with its list of required features to check for library support. 


### client.id (string, optional) {#rdkafkaoptions-client.id}

Client identifier. 


### closesocket_cb (string, optional) {#rdkafkaoptions-closesocket_cb}

Socket close callback 


### connect_cb (string, optional) {#rdkafkaoptions-connect_cb}

Socket connect callback 


### connections.max.idle.ms (int, optional) {#rdkafkaoptions-connections.max.idle.ms}

Close broker connections after the specified time of inactivity. Disable with 0. If this property is left at its default value some heuristics are performed to determine a suitable default value, this is currently limited to identifying brokers on Azure (see librdkafka issue #3109 for more info). 


### debug (string, optional) {#rdkafkaoptions-debug}

A comma-separated list of debug contexts to enable. Detailed Producer debugging: broker,topic,msg. Consumer: consumer,cgrp,topic,fetch 


### default_topic_conf (string, optional) {#rdkafkaoptions-default_topic_conf}

Default topic configuration for automatically subscribed topics 


### enable.random.seed (bool, optional) {#rdkafkaoptions-enable.random.seed}

If enabled librdkafka will initialize the PRNG with srand(current_time.milliseconds) on the first invocation of rd_kafka_new() (required only if rand_r() is not available on your platform). If disabled the application must call srand() prior to calling rd_kafka_new(). 


### enable.ssl.certificate.verification (bool, optional) {#rdkafkaoptions-enable.ssl.certificate.verification}

Enable OpenSSL's builtin broker (server) certificate verification. This verification can be extended by the application by implementing a certificate_verify_cb. 


### enable.sasl.oauthbearer.unsecure.jwt (bool, optional) {#rdkafkaoptions-enable.sasl.oauthbearer.unsecure.jwt}

Enable the builtin unsecure JWT OAUTHBEARER token handler if no oauthbearer_refresh_cb has been set. This builtin handler should only be used for development or testing, and not in production. 


### enabled_events (int, optional) {#rdkafkaoptions-enabled_events}

See `rd_kafka_conf_set_events()` 


### error_cb (string, optional) {#rdkafkaoptions-error_cb}

Error callback (set with rd_kafka_conf_set_error_cb()) 


### interceptors (string, optional) {#rdkafkaoptions-interceptors}

Interceptors added through rd_kafka_conf_interceptor_add_..() and any configuration handled by interceptors. 


### internal.termination.signal (int, optional) {#rdkafkaoptions-internal.termination.signal}

Signal that librdkafka will use to quickly terminate on rd_kafka_destroy(). If this signal is not set then there will be a delay before rd_kafka_wait_destroyed() returns true as internal threads are timing out their system calls. If this signal is set however the delay will be minimal. The application should mask this signal as an internal signal handler is installed. 


### log_cb (string, optional) {#rdkafkaoptions-log_cb}

Log callback (set with rd_kafka_conf_set_log_cb()) 


### log.connection.close (bool, optional) {#rdkafkaoptions-log.connection.close}

Log broker disconnects. It might be useful to turn this off when interacting with 0.9 brokers with an aggressive `connections.max.idle.ms` value. 


### log_level (int, optional) {#rdkafkaoptions-log_level}

Logging level (syslog(3) levels) 


### log.queue (bool, optional) {#rdkafkaoptions-log.queue}

Disable spontaneous log_cb from internal librdkafka threads, instead enqueue log messages on queue set with `rd_kafka_set_log_queue()` and serve log callbacks or events through the standard poll APIs. **NOTE**: Log messages will linger in a temporary queue until the log queue has been set. 


### log.thread.name (bool, optional) {#rdkafkaoptions-log.thread.name}

Print internal thread name in log messages (useful for debugging librdkafka internals) 


### max.in.flight (int, optional) {#rdkafkaoptions-max.in.flight}

Alias for `max.in.flight.requests.per.connection`: Maximum number of in-flight requests per broker connection. This is a generic property applied to all broker communication, however it is primarily relevant to produce requests. In particular, note that other mechanisms limit the number of outstanding consumer fetch request per broker to one. 


### max.in.flight.requests.per.connection (int, optional) {#rdkafkaoptions-max.in.flight.requests.per.connection}

Maximum number of in-flight requests per broker connection. This is a generic property applied to all broker communication, however it is primarily relevant to produce requests. In particular, note that other mechanisms limit the number of outstanding consumer fetch request per broker to one. 


### message.copy.max.bytes (int, optional) {#rdkafkaoptions-message.copy.max.bytes}

Maximum size for message to be copied to buffer. Messages larger than this will be passed by reference (zero-copy) at the expense of larger iovecs. 


### message.max.bytes (int, optional) {#rdkafkaoptions-message.max.bytes}

Maximum Kafka protocol request message size. Due to differing framing overhead between protocol versions the producer is unable to reliably enforce a strict max message limit at produce time and may exceed the maximum size by one message in protocol ProduceRequests, the broker will enforce the the topic's `max.message.bytes` limit (see Apache Kafka documentation). 


### metadata.broker.list (string, optional) {#rdkafkaoptions-metadata.broker.list}

Initial list of brokers as a CSV list of broker host or host:port. The application may also use `rd_kafka_brokers_add()` to add brokers during runtime. 


### metadata.max.age.ms (int, optional) {#rdkafkaoptions-metadata.max.age.ms}

Metadata cache max age. Defaults to topic.metadata.refresh.interval.ms * 3 


### oauthbearer_token_refresh_cb (string, optional) {#rdkafkaoptions-oauthbearer_token_refresh_cb}

SASL/OAUTHBEARER token refresh callback (set with rd_kafka_conf_set_oauthbearer_token_refresh_cb(), triggered by rd_kafka_poll(), et.al. This callback will be triggered when it is time to refresh the client's OAUTHBEARER token. Also see rd_kafka_conf_enable_sasl_queue(). 


### opaque (string, optional) {#rdkafkaoptions-opaque}

Application opaque (set with rd_kafka_conf_set_opaque()) 


### open_cb (string, optional) {#rdkafkaoptions-open_cb}

File open callback to provide race-free CLOEXEC 


### plugin.library.paths (string, optional) {#rdkafkaoptions-plugin.library.paths}

List of plugin libraries to load (; separated). The library search path is platform dependent (see dlopen(3) for Unix and LoadLibrary() for Windows). If no filename extension is specified the platform-specific extension (such as .dll or .so) will be appended automatically. 


### receive.message.max.bytes (int, optional) {#rdkafkaoptions-receive.message.max.bytes}

Maximum Kafka protocol response message size. This serves as a safety precaution to avoid memory exhaustion in case of protocol hickups. This value must be at least `fetch.max.bytes`  + 512 to allow for protocol overhead; the value is adjusted automatically unless the configuration property is explicitly set. 


### reconnect.backoff.max.ms (int, optional) {#rdkafkaoptions-reconnect.backoff.max.ms}

The maximum time to wait before reconnecting to a broker after the connection has been closed. 


### reconnect.backoff.ms (int, optional) {#rdkafkaoptions-reconnect.backoff.ms}

The initial time to wait before reconnecting to a broker after the connection has been closed. The time is increased exponentially until `reconnect.backoff.max.ms` is reached. -25% to +50% jitter is applied to each reconnect backoff. A value of 0 disables the backoff and reconnects immediately. 


### resolve_cb (string, optional) {#rdkafkaoptions-resolve_cb}

Address resolution callback (set with rd_kafka_conf_set_resolve_cb()) 


### ssl.ca.location (string, optional) {#rdkafkaoptions-ssl.ca.location}

File or directory path to CA certificate(s) for verifying the broker's key. Defaults: On Windows the system's CA certificates are automatically looked up in the Windows Root certificate store. On Mac OSX this configuration defaults to `probe`. It is recommended to install openssl using Homebrew, to provide CA certificates. On Linux install the distribution's ca-certificates package. If OpenSSL is statically linked or `ssl.ca.location` is set to `probe` a list of standard paths will be probed and the first one found will be used as the default CA certificate location path. If OpenSSL is dynamically linked the OpenSSL library's default path will be used (see `OPENSSLDIR` in `openssl version -a`). 


### ssl.ca.pem (string, optional) {#rdkafkaoptions-ssl.ca.pem}

CA certificate string (PEM format) for verifying the broker's key. 


### ssl.certificate.location (string, optional) {#rdkafkaoptions-ssl.certificate.location}

Path to client's public key (PEM) used for authentication. 


### ssl.certificate.pem (string, optional) {#rdkafkaoptions-ssl.certificate.pem}

Client's public key string (PEM format) used for authentication. 


### ssl.cipher.suites (string, optional) {#rdkafkaoptions-ssl.cipher.suites}

A cipher suite is a named combination of authentication, encryption, MAC and key exchange algorithm used to negotiate the security settings for a network connection using TLS or SSL network protocol. See manual page for `ciphers(1)` and `SSL_CTX_set_cipher_list(3). 


### ssl.crl.location (string, optional) {#rdkafkaoptions-ssl.crl.location}

Path to CRL for verifying broker's certificate validity. 


### ssl.curves.list (string, optional) {#rdkafkaoptions-ssl.curves.list}

The supported-curves extension in the TLS ClientHello message specifies the curves (standard/named, or 'explicit' GF(2^k) or GF(p)) the client is willing to have the server use. See manual page for `SSL_CTX_set1_curves_list(3)`. OpenSSL >= 1.0.2 required. 


### ssl.endpoint.identification.algorithm (string, optional) {#rdkafkaoptions-ssl.endpoint.identification.algorithm}

Endpoint identification algorithm to validate broker hostname using broker certificate. https - Server (broker) hostname verification as specified in RFC2818. none - No endpoint verification. OpenSSL >= 1.0.2 required. 


### ssl.engine.id (string, optional) {#rdkafkaoptions-ssl.engine.id}

OpenSSL engine id is the name used for loading engine. 


### ssl.engine.location (string, optional) {#rdkafkaoptions-ssl.engine.location}

**DEPRECATED** Path to OpenSSL engine library. OpenSSL >= 1.1.x required. DEPRECATED: OpenSSL engine support is deprecated and should be replaced by OpenSSL 3 providers. 


### ssl.key.location (string, optional) {#rdkafkaoptions-ssl.key.location}

Path to client's private key (PEM) used for authentication. 


### ssl.key.password (string, optional) {#rdkafkaoptions-ssl.key.password}

Private key passphrase (for use with `ssl.key.location` and `set_ssl_cert()`). 


### ssl.key.pem (string, optional) {#rdkafkaoptions-ssl.key.pem}

Client's private key string (PEM format) used for authentication. 


### ssl.keystore.location (string, optional) {#rdkafkaoptions-ssl.keystore.location}

Path to client's keystore (PKCS#12) used for authentication. 


### ssl.keystore.password (string, optional) {#rdkafkaoptions-ssl.keystore.password}

Client's keystore (PKCS#12) password. 


### ssl.providers (string, optional) {#rdkafkaoptions-ssl.providers}

Comma-separated list of OpenSSL 3.0.x implementation providers. E.g., "default,legacy". 


### ssl.sigalgs.list (string, optional) {#rdkafkaoptions-ssl.sigalgs.list}

The client uses the TLS ClientHello signature_algorithms extension to indicate to the server which signature/hash algorithm pairs may be used in digital signatures. See manual page for `SSL_CTX_set1_sigalgs_list(3)`. OpenSSL >= 1.0.2 required. 


### sasl.kerberos.keytab (string, optional) {#rdkafkaoptions-sasl.kerberos.keytab}

Path to Kerberos keytab file. This configuration property is only used as a variable in sasl.kerberos.kinit.cmd as  ... -t "%{sasl.kerberos.keytab}". 


### sasl.kerberos.kinit.cmd (string, optional) {#rdkafkaoptions-sasl.kerberos.kinit.cmd}

Shell command to refresh or acquire the client's Kerberos ticket. This command is executed on client creation and every sasl.kerberos.min.time.before.relogin (0=disable). 


### sasl.kerberos.min.time.before.relogin (int, optional) {#rdkafkaoptions-sasl.kerberos.min.time.before.relogin}

Minimum time in milliseconds between key refresh attempts. Disable automatic key refresh by setting this property to 0. 


### sasl.kerberos.principal (string, optional) {#rdkafkaoptions-sasl.kerberos.principal}

This client's Kerberos principal name. (Not supported on Windows, will use the logon user's principal). 


### sasl.kerberos.service.name (string, optional) {#rdkafkaoptions-sasl.kerberos.service.name}

Kerberos principal name that Kafka runs as, not including /hostname@REALM. 


### sasl.mechanisms (string, optional) {#rdkafkaoptions-sasl.mechanisms}

SASL mechanism to use for authentication. Supported: GSSAPI, PLAIN, SCRAM-SHA-256, SCRAM-SHA-512, OAUTHBEARER. NOTE: Despite the name only one mechanism must be configured. 


### sasl.oauthbearer.client.id (string, optional) {#rdkafkaoptions-sasl.oauthbearer.client.id}

Public identifier for the application. Must be unique across all clients that the authorization server handles. Only used when sasl.oauthbearer.method is set to "oidc". 


### sasl.oauthbearer.client.secret (string, optional) {#rdkafkaoptions-sasl.oauthbearer.client.secret}

Client secret only known to the application and the authorization server. This should be a sufficiently random string that is not guessable. Only used when sasl.oauthbearer.method is set to "oidc". 


### sasl.oauthbearer.config (string, optional) {#rdkafkaoptions-sasl.oauthbearer.config}

SASL/OAUTHBEARER configuration. The format is implementation-dependent and must be parsed accordingly. The default unsecured token implementation (see https://tools.ietf.org/html/rfc7515#appendix-A.5) recognizes space-separated name=value pairs with valid names including principalClaimName, principal, scopeClaimName, scope, and lifeSeconds. The default value for principalClaimName is "sub", the default value for scopeClaimName is "scope", and the default value for lifeSeconds is 3600. The scope value is CSV format with the default value being no/empty scope. For example: principalClaimName=azp principal=admin scopeClaimName=roles scope=role1,role2 lifeSeconds=600. In addition, SASL extensions can be communicated to the broker via extension_NAME=value. For example: principal=admin extension_traceId=123. 


### sasl.oauthbearer.extensions (string, optional) {#rdkafkaoptions-sasl.oauthbearer.extensions}

Allow additional information to be provided to the broker. Comma-separated list of key=value pairs. E.g., "supportFeatureX=true,organizationId=sales-emea".Only used when sasl.oauthbearer.method is set to "oidc". 


### sasl.oauthbearer.method (string, optional) {#rdkafkaoptions-sasl.oauthbearer.method}

Set to "default" or "oidc" to control which login method to be used. If set to "oidc", the following properties must also be specified: sasl.oauthbearer.client.id, sasl.oauthbearer.client.secret, and sasl.oauthbearer.token.endpoint.url. 


### sasl.oauthbearer.scope (string, optional) {#rdkafkaoptions-sasl.oauthbearer.scope}

Client use this to specify the scope of the access request to the broker. Only used when sasl.oauthbearer.method is set to "oidc". 


### sasl.oauthbearer.token.endpoint.url (string, optional) {#rdkafkaoptions-sasl.oauthbearer.token.endpoint.url}

OAuth/OIDC issuer token endpoint HTTP(S) URI used to retrieve token. Only used when sasl.oauthbearer.method is set to "oidc". 


### sasl.password (string, optional) {#rdkafkaoptions-sasl.password}

SASL password for use with the PLAIN and SASL-SCRAM-.. mechanism. 


### sasl.username (string, optional) {#rdkafkaoptions-sasl.username}

SASL username for use with the PLAIN and SASL-SCRAM-.. mechanisms. 


### security.protocol (string, optional) {#rdkafkaoptions-security.protocol}

Protocol used to communicate with brokers. 


### socket.blocking.max.ms (int, optional) {#rdkafkaoptions-socket.blocking.max.ms}

DEPRECATED No longer used. 


### socket_cb (string, optional) {#rdkafkaoptions-socket_cb}

Socket creation callback to provide race-free CLOEXEC 


### socket.connection.setup.timeout.ms (int, optional) {#rdkafkaoptions-socket.connection.setup.timeout.ms}

Maximum time allowed for broker connection setup (TCP connection setup as well SSL and SASL handshake). If the connection to the broker is not fully functional after this the connection will be closed and retried. 


### socket.keepalive.enable (bool, optional) {#rdkafkaoptions-socket.keepalive.enable}

Enable TCP keep-alives (SO_KEEPALIVE) on broker sockets 


### socket.max.fails (int, optional) {#rdkafkaoptions-socket.max.fails}

Disconnect from broker when this number of send failures (e.g., timed out requests) is reached. Disable with 0. WARNING: It is highly recommended to leave this setting at its default value of 1 to avoid the client and broker to become desynchronized in case of request timeouts. NOTE: The connection is automatically re-established. 


### socket.nagle.disable (bool, optional) {#rdkafkaoptions-socket.nagle.disable}

Disable the Nagle algorithm (TCP_NODELAY) on broker sockets. 


### socket.receive.buffer.bytes (int, optional) {#rdkafkaoptions-socket.receive.buffer.bytes}

Broker socket receive buffer size. System default is used if 0. 


### socket.send.buffer.bytes (int, optional) {#rdkafkaoptions-socket.send.buffer.bytes}

Broker socket send buffer size. System default is used if 0. 


### socket.timeout.ms (int, optional) {#rdkafkaoptions-socket.timeout.ms}

Default timeout for network requests. Producer: ProduceRequests will use the lesser value of `socket.timeout.ms` and remaining `message.timeout.ms` for the first message in the batch. Consumer: FetchRequests will use `fetch.wait.max.ms` + `socket.timeout.ms`. Admin: Admin requests will use `socket.timeout.ms` or explicitly set `rd_kafka_AdminOptions_set_operation_timeout()` value. 


### statistics.interval.ms (int, optional) {#rdkafkaoptions-statistics.interval.ms}

librdkafka statistics emit interval. The application also needs to register a stats callback using `rd_kafka_conf_set_stats_cb()`. The granularity is 1000ms. A value of 0 disables statistics. 


### stats_cb (string, optional) {#rdkafkaoptions-stats_cb}

Statistics callback (set with rd_kafka_conf_set_stats_cb()) 


### throttle_cb (string, optional) {#rdkafkaoptions-throttle_cb}

Throttle callback (set with rd_kafka_conf_set_throttle_cb()) 


### topic.blacklist (string, optional) {#rdkafkaoptions-topic.blacklist}

Topic blacklist, a comma-separated list of regular expressions for matching topic names that should be ignored in broker metadata information as if the topics did not exist. 


### topic.metadata.propagation.max.ms (int, optional) {#rdkafkaoptions-topic.metadata.propagation.max.ms}

Apache Kafka topic creation is asynchronous and it takes some time for a new topic to propagate throughout the cluster to all brokers. If a client requests topic metadata after manual topic creation but before the topic has been fully propagated to the broker the client is requesting metadata from, the topic will seem to be non-existent and the client will mark the topic as such, failing queued produced messages with `ERR__UNKNOWN_TOPIC`. This setting delays marking a topic as non-existent until the configured propagation max time has passed. The maximum propagation time is calculated from the time the topic is first referenced in the client, e.g., on produce(). 


### topic.metadata.refresh.fast.interval.ms (int, optional) {#rdkafkaoptions-topic.metadata.refresh.fast.interval.ms}

When a topic loses its leader a new metadata request will be enqueued immediately and then with this initial interval, exponentially increasing upto `retry.backoff.max.ms`, until the topic metadata has been refreshed. If not set explicitly, it will be defaulted to `retry.backoff.ms`. This is used to recover quickly from transitioning leader brokers. 


### topic.metadata.refresh.interval.ms (int, optional) {#rdkafkaoptions-topic.metadata.refresh.interval.ms}

Period of time in milliseconds at which topic and broker metadata is refreshed in order to proactively discover any new brokers, topics, partitions or partition leader changes. Use -1 to disable the intervalled refresh (not recommended). If there are no locally referenced topics (no topic objects created, no messages produced, no subscription or no assignment) then only the broker list will be refreshed every interval but no more often than every 10s. 


### topic.metadata.refresh.sparse (bool, optional) {#rdkafkaoptions-topic.metadata.refresh.sparse}

Sparse metadata requests (consumes less network bandwidth) 



