---
title: SplunkHEC
weight: 200
generated_file: true
---

# Sending messages over Splunk HEC
## Overview

Based on the [Splunk destination of AxoSyslog core](https://axoflow.com/docs/axosyslog-core/chapter-destinations/syslog-ng-with-splunk/).

Available in Logging operator version 4.4 and later.

## Example

{{< highlight yaml >}}
apiVersion: logging.banzaicloud.io/v1beta1
kind: SyslogNGOutput
metadata:
  name: splunkhec
spec:
  splunk_hec_event:
    url: "https://splunk-endpoint"
    token:
      valueFrom:
          secretKeyRef:
            name: splunk-hec
            key: token
{{</ highlight >}}


## Configuration
## SplunkHECOutput

###  (HTTPOutput, required) {#splunkhecoutput-}


### content_type (string, optional) {#splunkhecoutput-content_type}

Additional HTTP request content-type option. 


### default_index (string, optional) {#splunkhecoutput-default_index}

Fallback option for index field. For details, see the [documentation of the AxoSyslog syslog-ng distribution](https://axoflow.com/docs/axosyslog-core/chapter-destinations/syslog-ng-with-splunk/). 


### default_source (string, optional) {#splunkhecoutput-default_source}

Fallback option for source field. 


### default_sourcetype (string, optional) {#splunkhecoutput-default_sourcetype}

Fallback option for sourcetype field. 


### event (string, optional) {#splunkhecoutput-event}

event() accepts a template, which declares the content of the log message sent to Splunk. Default value: `${MSG}` 


### extra_headers ([]string, optional) {#splunkhecoutput-extra_headers}

Additional HTTP request headers. 


### extra_queries ([]string, optional) {#splunkhecoutput-extra_queries}

Additional HTTP request query options. 


### fields (string, optional) {#splunkhecoutput-fields}

Additional indexing metadata for Splunk. 


### host (string, optional) {#splunkhecoutput-host}

Sets the host field. 


### index (string, optional) {#splunkhecoutput-index}

Splunk index where the messages will be stored. 


### source (string, optional) {#splunkhecoutput-source}

Sets the source field. 


### sourcetype (string, optional) {#splunkhecoutput-sourcetype}

Sets the sourcetype field. 


### time (string, optional) {#splunkhecoutput-time}

Sets the time field. 


### token (secret.Secret, optional) {#splunkhecoutput-token}

The token that syslog-ng OSE uses to authenticate on the event collector. 



