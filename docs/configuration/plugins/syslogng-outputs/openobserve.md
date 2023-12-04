---
title: Openobserve
weight: 200
generated_file: true
---

# Sending messages over Openobserve
## Overview

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
More information at https://axoflow.com/docs/axosyslog-core/chapter-destinations/openobserve/


## Configuration
## OpenobserveOutput

###  (HTTPOutput, required) {#openobserveoutput-}

Default: -

### organization (string, optional) {#openobserveoutput-organization}

Name of the organization in Openobserve. 

Default: -

### stream (string, optional) {#openobserveoutput-stream}

Name of the stream in Openobserve. 

Default: -

### record (string, optional) {#openobserveoutput-record}

Arguments to the `$format-json()` template function. Default: --scope rfc5424 --exclude DATE --key ISODATE @timestamp=${ISODATE}" 

Default: -


