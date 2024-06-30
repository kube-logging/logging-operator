---
title: OTLP output
weight: 200
generated_file: true
---

# Sending fluent structured messages over OTLP GRPC
## Overview

Sends messages over OTLP GRPC. For details on the available options of the output, see the [documentation of AxoSyslog](https://axoflow.com/docs/axosyslog-core/chapter-destinations/opentelemetry/).

## Example

A simple example sending logs over OTLP GRPC to a remote OTLP endpoint:

{{< highlight yaml >}}
kind: SyslogNGOutput
apiVersion: logging.banzaicloud.io/v1beta1
metadata:
  name: otlp
spec:
  otlp:
    url: otel-server
    port: 4379
{{</ highlight >}}



## Configuration
## OTLPOutput

###  (Batch, required) {#otlpoutput-}

Batching parameters 


### auth (*Auth, optional) {#otlpoutput-auth}

Authentication configuration, see the [documentation of the AxoSyslog syslog-ng distribution](https://axoflow.com/docs/axosyslog-core/chapter-destinations/destination-syslog-ng-otlp/#auth). 


### channel_args (filter.ArrowMap, optional) {#otlpoutput-channel_args}

Add GRPC Channel arguments https://axoflow.com/docs/axosyslog-core/chapter-destinations/opentelemetry/#channel-args 


### compression (*bool, optional) {#otlpoutput-compression}

Enable or disable compression.

Default: false

### disk_buffer (*DiskBuffer, optional) {#otlpoutput-disk_buffer}

This option enables putting outgoing messages into the disk buffer of the destination to avoid message loss in case of a system failure on the destination side. For details, see the [Syslog-ng DiskBuffer options](../disk_buffer/).

Default: false

### url (string, required) {#otlpoutput-url}

Specifies the hostname or IP address and optionally the port number of the web service that can receive log data via HTTP. Use a colon (:) after the address to specify the port number of the server. For example: `http://127.0.0.1:8000` 



