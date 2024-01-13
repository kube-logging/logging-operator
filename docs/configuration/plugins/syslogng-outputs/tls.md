---
title: TLS config for syslog-ng outputs
weight: 200
generated_file: true
---

# TLS config for syslog-ng outputs
## Overview
 For details on how TLS configuration works in syslog-ng, see the [AxoSyslog Core documentation](https://axoflow.com/docs/axosyslog-core/chapter-encrypted-transport-tls/).

## Configuration
## TLS

### ca_dir (*secret.Secret, optional) {#tls-ca_dir}

The name of a directory that contains a set of trusted CA certificates in PEM format. For details, see the [AxoSyslog Core documentation](https://axoflow.com/docs/axosyslog-core/chapter-encrypted-transport-tls/tlsoptions/#ca-dir). 


### ca_file (*secret.Secret, optional) {#tls-ca_file}

The name of a file that contains a set of trusted CA certificates in PEM format. (Optional) For details, see the [AxoSyslog Core documentation](https://axoflow.com/docs/axosyslog-core/chapter-encrypted-transport-tls/tlsoptions/#ca-file). 


### cert_file (*secret.Secret, optional) {#tls-cert_file}

Name of a file, that contains an X.509 certificate (or a certificate chain) in PEM format, suitable as a TLS certificate, matching the private key set in the key-file() option. For details, see the [AxoSyslog Core documentation](https://axoflow.com/docs/axosyslog-core/chapter-encrypted-transport-tls/tlsoptions/#cert-file). 


### cipher-suite (string, optional) {#tls-cipher-suite}

Description: Specifies the cipher, hash, and key-exchange algorithms used for the encryption, for example, ECDHE-ECDSA-AES256-SHA384. The list of available algorithms depends on the version of OpenSSL used to compile syslog-ng. 


### key_file (*secret.Secret, optional) {#tls-key_file}

The name of a file that contains an unencrypted private key in PEM format, suitable as a TLS key. For details, see the [AxoSyslog Core documentation](https://axoflow.com/docs/axosyslog-core/chapter-encrypted-transport-tls/tlsoptions/#key-file). 


### peer_verify (*bool, optional) {#tls-peer_verify}

Verification method of the peer. For details, see the [AxoSyslog Core documentation](https://axoflow.com/docs/axosyslog-core/chapter-encrypted-transport-tls/tlsoptions/#tls-options-peer-verify). 


### ssl_version (string, optional) {#tls-ssl_version}

Configure required TLS version. Accepted values: [sslv3, tlsv1, tlsv1_0, tlsv1_1, tlsv1_2, tlsv1_3] 


### use-system-cert-store (*bool, optional) {#tls-use-system-cert-store}

Use the certificate store of the system for verifying HTTPS certificates. For details, see the [AxoSyslog Core documentation](https://curl.se/docs/sslcerts.html). 



## GrpcTLS

### ca_file (*secret.Secret, optional) {#grpctls-ca_file}

The name of a file that contains a set of trusted CA certificates in PEM format. (Optional) For details, see the [AxoSyslog Core documentation](https://axoflow.com/docs/axosyslog-core/chapter-encrypted-transport-tls/tlsoptions/#ca-file). 


### cert_file (*secret.Secret, optional) {#grpctls-cert_file}

Name of a file, that contains an X.509 certificate (or a certificate chain) in PEM format, suitable as a TLS certificate, matching the private key set in the key-file() option. For details, see the [AxoSyslog Core documentation](https://axoflow.com/docs/axosyslog-core/chapter-encrypted-transport-tls/tlsoptions/#cert-file). 


### key_file (*secret.Secret, optional) {#grpctls-key_file}

The name of a file that contains an unencrypted private key in PEM format, suitable as a TLS key. For details, see the [AxoSyslog Core documentation](https://axoflow.com/docs/axosyslog-core/chapter-encrypted-transport-tls/tlsoptions/#key-file). 



