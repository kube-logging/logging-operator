---
title: TLS config for syslog-ng outputs
weight: 200
generated_file: true
---

# TLS config for syslog-ng outputs
## Overview
 More info at https://www.syslog-ng.com/technical-documents/doc/syslog-ng-open-source-edition/3.37/administration-guide/32#kanchor2338

## Configuration
## TLS

### ca_dir (*secret.Secret, optional) {#tls-ca_dir}

The name of a directory that contains a set of trusted CA certificates in PEM format. [more information](https://www.syslog-ng.com/technical-documents/doc/syslog-ng-open-source-edition/3.37/administration-guide/73#kanchor3142) 

Default: -

### ca_file (*secret.Secret, optional) {#tls-ca_file}

The name of a file that contains a set of trusted CA certificates in PEM format. (Optional) [more information](https://www.syslog-ng.com/technical-documents/doc/syslog-ng-open-source-edition/3.37/administration-guide/73#kanchor3144) 

Default: -

### key_file (*secret.Secret, optional) {#tls-key_file}

The name of a file that contains an unencrypted private key in PEM format, suitable as a TLS key. [more information](https://www.syslog-ng.com/technical-documents/doc/syslog-ng-open-source-edition/3.37/administration-guide/73#kanchor3163) 

Default: -

### cert_file (*secret.Secret, optional) {#tls-cert_file}

Name of a file, that contains an X.509 certificate (or a certificate chain) in PEM format, suitable as a TLS certificate, matching the private key set in the key-file() option. [more information](https://www.syslog-ng.com/technical-documents/doc/syslog-ng-open-source-edition/3.37/administration-guide/73#kanchor3146) 

Default: -

### peer_verify (string, optional) {#tls-peer_verify}

Verification method of the peer. [more information](https://www.syslog-ng.com/technical-documents/doc/syslog-ng-open-source-edition/3.37/administration-guide/73#tls-options-peer-verify) 

Default: -

### use-system-cert-store (*bool, optional) {#tls-use-system-cert-store}

Use the certificate store of the system for verifying HTTPS certificates. [more information](https://curl.se/docs/sslcerts.html) 

Default: -

### cipher-suite (string, optional) {#tls-cipher-suite}

Description: Specifies the cipher, hash, and key-exchange algorithms used for the encryption, for example, ECDHE-ECDSA-AES256-SHA384. The list of available algorithms depends on the version of OpenSSL used to compile syslog-ng OSE 

Default: -


