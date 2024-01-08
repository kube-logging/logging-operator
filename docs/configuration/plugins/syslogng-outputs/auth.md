---
title: Authentication config for syslog-ng outputs
weight: 200
generated_file: true
---

# Authentication config for syslog-ng outputs
## Overview
 GRPC-based outputs use this configuration instead of the simple `tls` field found at most HTTP based destinations. For details, see the documentation of a related syslog-ng destination, for example, [Grafana Loki](https://axoflow.com/docs/axosyslog-core/chapter-destinations/destination-loki/#auth).

## Configuration
## Auth

Authentication settings. Only one authentication method can be set. Default: Insecure

### adc (*ADC, optional) {#auth-adc}

Application Default Credentials (ADC). 


### alts (*ALTS, optional) {#auth-alts}

Application Layer Transport Security (ALTS) is a simple to use authentication, only available within Google’s infrastructure. 


### insecure (*Insecure, optional) {#auth-insecure}

This is the default method, authentication is disabled (`auth(insecure())`). 


### tls (*GrpcTLS, optional) {#auth-tls}

This option sets various options related to TLS encryption, for example, key/certificate files and trusted CA locations. TLS can be used only with tcp-based transport protocols. For details, see [TLS for syslog-ng outputs](../tls/) and the [documentation of the AxoSyslog syslog-ng distribution](https://axoflow.com/docs/axosyslog-core/chapter-encrypted-transport-tls/tlsoptions). 



## ADC


## Insecure


## ALTS

### target-service-accounts ([]string, optional) {#alts-target-service-accounts}



