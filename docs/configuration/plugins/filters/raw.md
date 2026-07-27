---
title: Raw
weight: 200
generated_file: true
---

# Raw
## Overview
 Configure custom or unexposed Fluentd filters via raw configuration. The configuration is parsed and rendered by the operator (parameter ordering and duplicate keys are not preserved). Disabled by default, set `logging.spec.enableRawFluentdFilter=true` on the Logging resource to enable it.

 > **Security warning:** enabling the flag grants **every** user who can create a `Flow` or `ClusterFlow` in the logging domain the ability to inject arbitrary configuration (and thereby run arbitrary code, for example via `exec`-style plugins) in the shared Fluentd aggregator, which processes all tenants' logs and mounts output credentials and TLS keys. Treat it as a break-glass setting and enable it only if you trust all Flow authors in the cluster.

## Example `Raw` filter configurations

### Configure a custom filter via raw configuration

{{< highlight yaml >}}
apiVersion: logging.banzaicloud.io/v1beta1
kind: Flow
metadata:
  name: demo-flow
spec:
  filters:
    - raw:
        config: |
          @type my_filter
          <my_section>
            foo bar
            tags ["web", "api", "db"]
          </my_section>
  selectors: {}
  localOutputRefs:
    - demo-output
{{</ highlight >}}

Fluentd Config Result

{{< highlight xml >}}
<filter **>
  @type my_filter
  @id test
  <my_section>
    foo bar
    tags ["web", "api", "db"]
  </my_section>
</filter>
{{</ highlight >}}

### Configure an unexposed filter via raw configuration

{{< highlight yaml >}}
apiVersion: logging.banzaicloud.io/v1beta1
kind: Flow
metadata:
  name: demo-flow
spec:
  filters:
    - raw:
        config: |
          @type ua_parser
          flatten true
          key_name ua_string
  selectors: {}
  localOutputRefs:
    - demo-output
{{</ highlight >}}

Fluentd Config Result

{{< highlight xml >}}
<filter **>
  @type ua_parser
  @id test
  flatten true
  key_name ua_string
</filter>
{{</ highlight >}}



## Configuration
## Raw

### config (string, optional) {#raw-config}

Raw configuration for the filter. 



