---
title: Authentication config for syslog-ng outputs
weight: 200
generated_file: true
---

# Authentication config for syslog-ng outputs
## Overview
 More info at TODO

## Configuration
## Auth

### alts (*ALTS, optional) {#auth-alts}

Application Layer Transport Security (ALTS) is a simple to use authentication, only available within Google’s infrastructure. 

Default: -

### adc (*ADC, optional) {#auth-adc}

Application Default Credentials (ADC). 

Default: -

### insecure (*Insecure, optional) {#auth-insecure}

This is the default method, authentication is disabled (auth(insecure())). 

Default: -

### tls (*TLS, optional) {#auth-tls}

This option sets various options related to TLS encryption, for example, key/certificate files and trusted CA locations. TLS can be used only with tcp-based transport protocols. For details, see [TLS for syslog-ng outputs](../tls/) and the [syslog-ng documentation](https://axoflow.com/docs/axosyslog-core/chapter-encrypted-transport-tls/tlsoptions). 

Default: -


## ADC


## Insecure


## ALTS

### target-service-accounts ([]string, optional) {#alts-target-service-accounts}

Default: -


