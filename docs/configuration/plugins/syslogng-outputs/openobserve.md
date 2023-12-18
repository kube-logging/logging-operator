---
title: Openobserve
weight: 200
generated_file: true
---

# Sending messages over Openobserve
## Overview

Send messages to [OpenObserve](https://openobserve.ai/docs/api/ingestion/logs/json/) using its [Logs Ingestion - JSON API](https://openobserve.ai/docs/api/ingestion/logs/json/). This API accepts multiple records in batch in JSON format.


## Example

{{< highlight yaml >}}
apiVersion: logging.banzaicloud.io/v1beta1
kind: SyslogNGOutput
metadata:
  name: openobserve
spec:
  openobserve:
    url: "https://some-openobserve-endpoint"
    port: 5040
    organization: "default"
    stream: "default"
    user: "username"
    password:
      valueFrom:
        secretKeyRef:
          name: openobserve
          key: password
{{</ highlight >}}

For details on the available options of the output, see the [documentation of the AxoSyslog syslog-ng distribution](https://axoflow.com/docs/axosyslog-core/chapter-destinations/openobserve/).


## Configuration
## OpenobserveOutput

###  (HTTPOutput, required) {#openobserveoutput-}


### organization (string, optional) {#openobserveoutput-organization}

Name of the organization in Openobserve. 


### record (string, optional) {#openobserveoutput-record}

Arguments to the `$format-json()` template function. Default: `"--scope rfc5424 --exclude DATE --key ISODATE @timestamp=${ISODATE}"` 


### stream (string, optional) {#openobserveoutput-stream}

Name of the stream in Openobserve. 



