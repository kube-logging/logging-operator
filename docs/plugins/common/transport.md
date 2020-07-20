---
title: Transport
weight: 200
generated_file: true
---

### Transport
| Variable Name | Type | Required | Default | Description |
|---|---|---|---|---|
| protocol | string | No | - | Protocol Default: :tcp<br> |
| version | string | No | - | Version Default: 'TLSv1_2'<br> |
| ciphers | string | No | - | Ciphers Default: "ALL:!aNULL:!eNULL:!SSLv2"<br> |
| insecure | bool | No | - | Use secure connection when use tls) Default: false<br> |
| ca_path | string | No | - | Specify path to CA certificate file<br> |
| cert_path | string | No | - | Specify path to Certificate file<br> |
| private_key_path | string | No | - | Specify path to private Key file<br> |
| private_key_passphrase | string | No | - | public CA private key passphrase contained path<br> |
| client_cert_auth | bool | No | - | When this is set Fluentd will check all incoming HTTPS requests<br>for a client certificate signed by the trusted CA, requests that<br>don't supply a valid client certificate will fail.<br> |
| ca_cert_path | string | No | - | Specify private CA contained path<br> |
| ca_private_key_path | string | No | - | private CA private key contained path<br> |
| ca_private_key_passphrase | string | No | - | private CA private key passphrase contained path<br> |
