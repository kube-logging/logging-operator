---
title: Sumo Logic HTTP
weight: 200
generated_file: true
---

# Storing messages in Sumo Logic over http
## Overview
 More info at https://www.syslog-ng.com/technical-documents/doc/syslog-ng-open-source-edition/3.37/administration-guide/55

## Configuration
## SumologicHTTPOutput

### collector (*secret.Secret, optional) {#sumologichttpoutput-collector}

The Cloud Syslog Cloud Token that you received from the Sumo Logic service while configuring your cloud syslog source.  

Default:  empty

### deployment (string, optional) {#sumologichttpoutput-deployment}

This option specifies your Sumo Logic deployment.https://help.sumologic.com/APIs/General-API-Information/Sumo-Logic-Endpoints-by-Deployment-and-Firewall-Security   

Default:  empty

### headers ([]string, optional) {#sumologichttpoutput-headers}

Custom HTTP headers to include in the request, for example, headers("HEADER1: header1", "HEADER2: header2").   

Default:  empty

### time_reopen (int, optional) {#sumologichttpoutput-time_reopen}

The time to wait in seconds before a dead connection is reestablished.  

Default:  60

### tls (*TLS, optional) {#sumologichttpoutput-tls}

This option sets various options related to TLS encryption, for example, key/certificate files and trusted CA locations. TLS can be used only with tcp-based transport protocols. For details, see https://www.syslog-ng.com/technical-documents/doc/syslog-ng-open-source-edition/3.37/administration-guide/73#TOPIC-1829193 

Default: -

### disk_buffer (*DiskBuffer, optional) {#sumologichttpoutput-disk_buffer}

This option enables putting outgoing messages into the disk buffer of the destination to avoid message loss in case of a system failure on the destination side.   

Default:  false

### body (string, optional) {#sumologichttpoutput-body}

Default: -

### batch-lines (int, optional) {#sumologichttpoutput-batch-lines}

Default: -

### batch-bytes (int, optional) {#sumologichttpoutput-batch-bytes}

Default: -

### batch-timeout (int, optional) {#sumologichttpoutput-batch-timeout}

Default: -

### persist_name (string, optional) {#sumologichttpoutput-persist_name}

Default: -


