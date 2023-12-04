---
title: LogScale
weight: 200
generated_file: true
---

# Storing messages in Falcon LogScale
## Overview
 The `LogScale` output sends log records over HTTP to Falcon's LogScale.


Based on the [LogScale destination of AxoSyslog core](https://axoflow.com/docs/axosyslog-core/chapter-destinations/crowdstrike-falcon/). Sends log records over HTTP to Falcon's LogScale.

{{< highlight yaml >}}
apiVersion: logging.banzaicloud.io/v1beta1
kind: SyslogNGOutput
metadata:
  name: test-logscale
  namespace: logging
spec:
  logscale:
    token:
      valueFrom:
        secretKeyRef:
          key: token
          name: logscale-token
    timezone: "UTC"
    batch_lines: 1000
    disk_buffer:
      disk_buf_size: 512000000
      dir: /buffers
      reliable: true
{{</ highlight >}}


## Configuration
## LogScaleOutput

### attributes (string, optional) {#logscaleoutput-attributes}

A JSON object representing key-value pairs for the Event. These key-value pairs adds structure to Events, making it easier to search. Attributes can be nested JSON objects, however, we recommend limiting the amount of nesting.

Default: `"--scope rfc5424 --exclude MESSAGE --exclude DATE --leave-initial-dot"`

### batch_bytes (int, optional) {#logscaleoutput-batch_bytes}


### batch_lines (int, optional) {#logscaleoutput-batch_lines}


### batch_timeout (int, optional) {#logscaleoutput-batch_timeout}


### body (string, optional) {#logscaleoutput-body}


### content_type (string, optional) {#logscaleoutput-content_type}

This field specifies the content type of the log records being sent to Falcon's LogScale.

Default: `"application/json"`

### disk_buffer (*DiskBuffer, optional) {#logscaleoutput-disk_buffer}

This option enables putting outgoing messages into the disk buffer of the destination to avoid message loss in case of a system failure on the destination side. For details, see the [Syslog-ng DiskBuffer options](../disk_buffer/).

Default: false

### extra_headers (string, optional) {#logscaleoutput-extra_headers}

This field represents additional headers that can be included in the HTTP request when sending log records to Falcon's LogScale.

Default: empty

### persist_name (string, optional) {#logscaleoutput-persist_name}


### rawstring (string, optional) {#logscaleoutput-rawstring}

The raw string representing the Event. The default display for an Event in LogScale is the rawstring. If you do not provide the rawstring field, then the response defaults to a JSON representation of the attributes field.

Default: empty

### timezone (string, optional) {#logscaleoutput-timezone}

The timezone is only required if you specify the timestamp in milliseconds. The timezone specifies the local timezone for the event. Note that you must still specify the timestamp in UTC time. 


### token (*secret.Secret, optional) {#logscaleoutput-token}

An [Ingest Token](https://library.humio.com/data-analysis/ingesting-data-tokens.html) is a unique string that identifies a repository and allows you to send data to that repository.

Default: empty

### url (*secret.Secret, optional) {#logscaleoutput-url}

Ingester URL is the URL of the Humio cluster you want to send data to.

Default: https://cloud.humio.com


