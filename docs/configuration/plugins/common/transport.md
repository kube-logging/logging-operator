---
title: Transport
weight: 200
generated_file: true
---

## Transport

### ca_cert_path (string, optional) {#transport-ca_cert_path}

Specify private CA contained path 

Default: -

### ca_path (string, optional) {#transport-ca_path}

Specify path to CA certificate file 

Default: -

### ca_private_key_passphrase (string, optional) {#transport-ca_private_key_passphrase}

private CA private key passphrase contained path 

Default: -

### ca_private_key_path (string, optional) {#transport-ca_private_key_path}

private CA private key contained path 

Default: -

### cert_path (string, optional) {#transport-cert_path}

Specify path to Certificate file 

Default: -

### ciphers (string, optional) {#transport-ciphers}

Ciphers Default: "ALL:!aNULL:!eNULL:!SSLv2" 

Default: -

### client_cert_auth (bool, optional) {#transport-client_cert_auth}

When this is set Fluentd will check all incoming HTTPS requests for a client certificate signed by the trusted CA, requests that don't supply a valid client certificate will fail. 

Default: -

### insecure (bool, optional) {#transport-insecure}

Use secure connection when use tls) Default: false 

Default: -

### private_key_passphrase (string, optional) {#transport-private_key_passphrase}

public CA private key passphrase contained path 

Default: -

### private_key_path (string, optional) {#transport-private_key_path}

Specify path to private Key file 

Default: -

### protocol (string, optional) {#transport-protocol}

Protocol Default: :tcp 

Default: -

### version (string, optional) {#transport-version}

Version Default: 'TLSv1_2' 

Default: -


