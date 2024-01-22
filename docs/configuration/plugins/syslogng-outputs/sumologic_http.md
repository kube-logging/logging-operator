---
title: Sumo Logic HTTP
weight: 200
generated_file: true
---

# Storing messages in Sumo Logic over http
## Overview
 The `sumologic-http` output sends log records over HTTP to Sumo Logic.

The `sumologic-http` output sends log records over HTTP to Sumo Logic. For details on the available options of the output, see the [documentation of the AxoSyslog syslog-ng distribution](https://axoflow.com/docs/axosyslog-core/chapter-destinations/destination-sumologic-intro/destination-sumologic-options/).

## Prerequisites

You need a Sumo Logic account to use this output. For details, see the [documentation of the AxoSyslog syslog-ng distribution](https://axoflow.com/docs/axosyslog-core/chapter-destinations/destination-sumologic-intro/).

## Example

{{< highlight yaml >}}
apiVersion: logging.banzaicloud.io/v1beta1
kind: SyslogNGOutput
metadata:
  name: test-sumo
  namespace: default
spec:
  sumologic-http:
    batch-lines: 1000
    disk_buffer:
      disk_buf_size: 512000000
      dir: /buffers
      reliable: true
    body: "$(format-json
                --subkeys json.
                --exclude json.kubernetes.annotations.*
                json.kubernetes.annotations=literal($(format-flat-json --subkeys json.kubernetes.annotations.))
                --exclude json.kubernetes.labels.*
                json.kubernetes.labels=literal($(format-flat-json --subkeys json.kubernetes.labels.)))"
    collector:
      valueFrom:
        secretKeyRef:
          key: token
          name: sumo-collector
    deployment: us2
    headers:
    - 'X-Sumo-Name: source-name'
    - 'X-Sumo-Category: source-category'
    tls:
      use-system-cert-store: true
{{</ highlight >}}

For details on the available options of the output, see the [documentation of the AxoSyslog syslog-ng distribution](https://axoflow.com/docs/axosyslog-core/chapter-destinations/destination-sumologic-intro/destination-sumologic-http/).


## Configuration
## SumologicHTTPOutput

### batch-bytes (int, optional) {#sumologichttpoutput-batch-bytes}


### batch-lines (int, optional) {#sumologichttpoutput-batch-lines}


### batch-timeout (int, optional) {#sumologichttpoutput-batch-timeout}


### body (string, optional) {#sumologichttpoutput-body}


### collector (*secret.Secret, optional) {#sumologichttpoutput-collector}

The Cloud Syslog Cloud Token that you received from the Sumo Logic service while configuring your cloud syslog source.

Default: empty

### deployment (string, optional) {#sumologichttpoutput-deployment}

This option specifies your Sumo Logic deployment.https://help.sumologic.com/APIs/General-API-Information/Sumo-Logic-Endpoints-by-Deployment-and-Firewall-Security

Default: empty

### disk_buffer (*DiskBuffer, optional) {#sumologichttpoutput-disk_buffer}

This option enables putting outgoing messages into the disk buffer of the destination to avoid message loss in case of a system failure on the destination side. For details, see the [Syslog-ng DiskBuffer options](../disk_buffer/).

Default: false

### headers ([]string, optional) {#sumologichttpoutput-headers}

Custom HTTP headers to include in the request, for example, `headers("HEADER1: header1", "HEADER2: header2")`.

Default: empty

### persist_name (string, optional) {#sumologichttpoutput-persist_name}


### tls (*TLS, optional) {#sumologichttpoutput-tls}

This option sets various options related to TLS encryption, for example, key/certificate files and trusted CA locations. TLS can be used only with tcp-based transport protocols. For details, see [TLS for syslog-ng outputs](../tls/) and the [documentation of the AxoSyslog syslog-ng distribution](https://axoflow.com/docs/axosyslog-core/chapter-encrypted-transport-tls/tlsoptions/). 


### time_reopen (int, optional) {#sumologichttpoutput-time_reopen}

The time to wait in seconds before a dead connection is reestablished.

Default: 60

### url (*secret.Secret, optional) {#sumologichttpoutput-url}



