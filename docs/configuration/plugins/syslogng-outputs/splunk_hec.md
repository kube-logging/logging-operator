---
title: SplunkHEC
weight: 200
generated_file: true
---

# Sending messages over Splunk HEC
## Overview

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

More information at https://axoflow.com/docs/axosyslog-core/chapter-destinations/syslog-ng-with-splunk/



## Configuration
## SplunkHECOutput

###  (HTTPOutput, required) {#splunkhecoutput-}

Default: -

### content_type (string, optional) {#splunkhecoutput-content_type}

Additional HTTP request content-type option. 

Default: -

### default_index (string, optional) {#splunkhecoutput-default_index}

Fallback option for index field. See [syslog-ng docs](https://axoflow.com/docs/axosyslog-core/chapter-destinations/syslog-ng-with-splunk/) 

Default: -

### default_source (string, optional) {#splunkhecoutput-default_source}

Fallback option for source field. 

Default: -

### default_sourcetype (string, optional) {#splunkhecoutput-default_sourcetype}

Fallback option for sourcetype field. 

Default: -

### event (string, optional) {#splunkhecoutput-event}

event() accepts a template, which declares the content of the log message sent to Splunk. Default value: ${MSG} 

Default: -

### extra_headers ([]string, optional) {#splunkhecoutput-extra_headers}

Additional HTTP request headers. 

Default: -

### extra_queries ([]string, optional) {#splunkhecoutput-extra_queries}

Additional HTTP request query options. 

Default: -

### fields (string, optional) {#splunkhecoutput-fields}

Additional indexing metadata for Splunk. 

Default: -

### host (string, optional) {#splunkhecoutput-host}

Sets the host field. 

Default: -

### index (string, optional) {#splunkhecoutput-index}

Splunk index where the messages will be stored. 

Default: -

### source (string, optional) {#splunkhecoutput-source}

Sets the source field. 

Default: -

### sourcetype (string, optional) {#splunkhecoutput-sourcetype}

Sets the sourcetype field. 

Default: -

### time (string, optional) {#splunkhecoutput-time}

Sets the time field. 

Default: -

### token (secret.Secret, optional) {#splunkhecoutput-token}

The token that syslog-ng OSE uses to authenticate on the event collector. 

Default: -


