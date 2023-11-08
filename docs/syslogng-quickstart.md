# Logging Operator Syslog-ng quickstart

This document describes how to get started with the preview version of the Logging Operator using syslog-ng as the forward agent.

## What is new?

There are 4 new CRDs to manage Logging Operator with syslog-ng:
- SyslogNGFlow
- SyslogNGClusterFlow
- SyslogNGOutput
- SyslogNGClusterOutput

These resources works identical to their counterparts used to manage Logging Operator with Fluentd, but are tailored to features available via syslog-ng.

Also, the Logging CRD has been extended with a section for configuring syslog-ng under `logging.spec.syslogNG`.

## Flows

SyslogNGFlow and SyslogNGClusterFlow resources have almost the same structure as Flow and ClusterFlow resources with the main differences explained in the following sections.

## Custom JSON delimiter

As `syslog-ng` by default uses `.` (dots) to separate key-value pairs during JSON parsing there can be unintended splits if there is dots in the key name.

Example input:
```json
{
  "ts2": 1662035132,
  "stream": "stdout",
  "logtag": "F",
  "kubernetes": {
    "labels": {
      "pod-template-hash": "57988cbd65",
      "app.kubernetes.io/name": "log-generator"
    },
    "host": "ip-192-168-4-79.eu-west-1.compute.internal",
    "container_name": "log-generator"
  },
  "cluster": "xxxxx"
}
```
during JSON rendering syslog-ng will split keys at dots

Example output:
```json
{
  "ts2": 1662035132,
  "stream": "stdout",
  "logtag": "F",
  "kubernetes": {
    "labels": {
      "pod-template-hash": "57988cbd65",
      "app": {
        "kubernetes": {
          "io/name": "log-generator"
        }
      },
      "host": "ip-192-168-4-79.eu-west-1.compute.internal",
      "container_name": "log-generator"
    },
    "cluster": "xxxxx"
  }
}
```

If you have many dots in your keys you can change the delimiter in the `Logging` resource.
```yaml
apiVersion: logging.banzaicloud.io/v1beta1
kind: Logging
metadata:
  name: logging-demo
spec:
  controlNamespace: default
  fluentbit:
    ...
  syslogNG:
    jsonKeyDelim: ';'
```

> However you need to be extra carefull to follow through this change on all resources!

### Flow
#### Match
You MUST use the new delimiter on each match statement
```yaml
  match:
    and:
    - regexp:
        pattern: one-eye-log-generator
        type: string
        value: json;kubernetes;labels;app.kubernetes.io/instance
    - regexp:
        pattern: log-generator
        type: string
        value: json;kubernetes;labels;app.kubernetes.io/name
```
#### Filters
Every field definition should use the new delimiter!
```yaml
filters:
  - parser:
      regexp:
        patterns:
        - '^(?<remote>[^ ]*) (?<host>[^ ]*) (?<user>[^ ]*) \[(?<time>[^\]]*)\] "(?<method>\S+)(?:
          +(?<path>[^\"]*?)(?: +\S*)?)?" (?<code>[^ ]*) (?<size>[^ ]*)(?: "(?<referer>[^\"]*)"
          "(?<agent>[^\"]*)"(?:\s+(?<http_x_forwarded_for>[^ ]+))?)?$'
        prefix: json;
        template: ${json;message}
  - rewrite:
    - set:
        field: json;cluster
        value: xxxxx
    - unset:
        field: json;message
```

### Output

In the `Output` template you should also be aware to set the special **delimiter**.

```yaml
spec:
  file:
    template: |
      $(format-json --subkeys json; --key-delimiter ;)
```

### Routing

The `match` field used to define routing rules has become more powerful by supporting [syslog-ng's *filter expressions*](https://www.syslog-ng.com/technical-documents/doc/syslog-ng-open-source-edition/3.37/administration-guide/65#TOPIC-1829159).
With this solution there are no restrictions on routing, you can filter both on metadata and log content as well.
The syntax is slightly different though to accomodate for all available options.

The following example filters for specific Pod labels
```yaml
  match:
    and:
    - regexp:
        value: json.kubernetes.labels.app.kubernetes.io/instance
        pattern: one-eye-log-generator
        type: string
    - regexp:
        value: json.kubernetes.labels.app.kubernetes.io/name
        pattern: log-generator
        type: string
```

Using `not` in a match statement

```yaml
match:
  not:
    and:
    - regexp:
        value: json.kubernetes.labels.app.kubernetes.io/instance
        pattern: one-eye-log-generator
        type: string
    - regexp:
        value: json.kubernetes.labels.app.kubernetes.io/name
        pattern: log-generator
        type: string
```

> Note: You need to use the `json.` prefix in field names.

Fields can be referenced using *dot notation*, e.g. in `{"kubernetes": {"namespace_name": "default"}}` the `namespace_name` field can be referenced using `json.kubernetes.namespace_name`.

Match expressions are basically a combination of filtering functions using the `and`, `or`, and `not` boolean operators.
Currently, only a pattern matching function is supported (called [`match`](https://www.syslog-ng.com/technical-documents/doc/syslog-ng-open-source-edition/3.37/administration-guide/68#TOPIC-1829171) in syslog-ng parlance, but renamed to `regexp` in the CRD to avoid confusion).

The `match` field can have one of the following options:
```yaml
  match:
    and: <list of nested match expressions>  // Logical AND between expressions
    or: <list of nested match expressions>   // Logical OR between expressions
    not: <nested match expression>           // Logical NOT of an expression
    regexp: ... // Pattern matching on a field's value or a templated value
```

The `regexp` field can have the following fields
```yaml
  regexp:
    pattern: <a pattern string>                            // Pattern match against, e.g. "my-app-\d+". The pattern's type is determined by the type field.
    value: <a field reference>                             // Reference to a field whose value to match. If this field is set, the template field cannot be used.
    template: <a templeted string combining field values>  // Template expression whose value to match. If this field is set, the value field cannot be used. For more info, see https://www.syslog-ng.com/technical-documents/doc/syslog-ng-open-source-edition/3.37/administration-guide/74#TOPIC-1829197
    type: <pattern type>                                   // Pattern type. Default is PCRE. For more info, see https://www.syslog-ng.com/technical-documents/doc/syslog-ng-open-source-edition/3.37/administration-guide/81#TOPIC-1829223
    flags: <list of flags>                                 // Pattern flags. For more info, see https://www.syslog-ng.com/technical-documents/doc/syslog-ng-open-source-edition/3.37/administration-guide/81#TOPIC-1829224
```

### Filters

Logging Operator currently supports the following filters for syslog-ng:
- match
- parser
- rewrite

#### Match filter
Match filters can be used to further narrow down the set of log records to process.
These filters have the same options and syntax as syslog-ng flow match expressions described above.

```yaml
  filters:
  - match:
      or:
      - regexp:
          value: json.kubernetes.labels.app.kubernetes.io/name
          pattern: apache
          type: string
      - regexp:
          value: json.kubernetes.labels.app.kubernetes.io/name
          pattern: nginx
          type: string
```

#### Parser filter
Parser filters can be used to extract key-value pairs from message data.
Logging Operator currently supports the following parsers:
- regexp
- syslog-parser

##### Regexp parser
The regexp parser can parse fields from a message with the help of regular expressions.

```yaml
  filters:
  - parser:
      regexp:
        patterns:
        - ".*test_field -> (?<test_field>.*)$"
        prefix: .regexp.
```

For more info, see https://www.syslog-ng.com/technical-documents/doc/syslog-ng-open-source-edition/3.36/administration-guide/90

#### Syslog-parser

The syslog parser can parse syslog messages. Docs: https://www.syslog-ng.com/technical-documents/doc/syslog-ng-open-source-edition/3.37/administration-guide/83#TOPIC-1829231

```yaml
filters:
-  parser:
      syslog-parser: {}
```

#### Rewrite filter
Rewrite filters can be used to modify record contents.
Logging Operator currently supports the following rewrite functions:
- rename
- set
- substitute
- unset
- group_unset

> Note: All rewrite functions support an optional `condition` which has the same syntax as match expressions described above.

##### Rename
The `rename` function changes an existing field's name.
```yaml
  filters:
  - rewrite:
    - rename:
        oldName: "json.kubernetes.labels.app"
        newName: "json.kubernetes.labels.app.kubernetes.io/name"
```

##### Set
The `set` function sets the value of a field.

```yaml
  filters:
  - rewrite:
    - set:
        field: "json.kubernetes.cluster"
        value: "prod-us"
```

##### Substitute
The `subst` function replaces parts of a field with a replacement value based on a pattern.

```yaml
  filters:
  - rewrite:
    - subst:
        pattern: "\d\d\d\d-\d\d\d\d-\d\d\d\d-\d\d\d\d"
        replace: "[redacted bank card number]"
        field: "MESSAGE"
```

The function also supports the `type` and `flags` fields for specifying pattern type and flags as described above for match expression regexp function.

#### Group Unset

The `group_unset` function removes keys based on a patterns. If an object have multiple sub-object you have to use this instead of `unset`

```yaml
  filters:
  - rewrite:
    - group_unset:
        pattern: json.kubernetes.annotations.*
```

## Outputs

SyslogNGOutput and SyslogNGClusterOutput resources have almost the same structure as Output and ClusterOutput resources with the main difference being the number and kind of supported destinations.

Logging Operator currently supports the following kinds of outputs for syslog-ng:
- [file](#file-output)
- [http](#http-output)
- [loggly](#loggly-output)
- [sumologic-http](#sumologic-http-output)
- [syslog](#syslog-output)

### File output

The `file` output stores log records to a plain text file.

```yaml
apiVersion: logging.banzaicloud.io/v1beta1
kind: SyslogNGOutput
  spec:
    file:
      path: /mnt/archive/logs/${YEAR}/${MONTH}/${DAY}/app.log
      create_dirs: true
    template: |
      $(format-json --subkeys json. --exclude json.kubernetes.labels.* json.kubernetes.labels=literal($(format-flat-json --subkeys json.kubernetes.labels.)))
```

### HTTP output

The `http` output sends log records over HTTP to the specified URL.

#### Parameters
```yaml
  url:            # Specifies the hostname or IP address and optionally the port number of the web service that can receive log data via HTTP. Use a colon (:) after the address to specify the port number of the server. For example: http://127.0.0.1:8000
  headers:        # Custom HTTP headers to include in the request, for example, headers("HEADER1: header1", "HEADER2: header2").  (default: empty)
  time_reopen:    # The time to wait in seconds before a dead connection is reestablished. (default: 60)
  tls:            # This option sets various options related to TLS encryption, for example, key/certificate files and trusted CA locations. TLS can be used only with tcp-based transport protocols. For details, see https://www.syslog-ng.com/technical-documents/doc/syslog-ng-open-source-edition/3.37/administration-guide/73#TOPIC-1829193
  disk_buffer:    # This option enables putting outgoing messages into the disk buffer of the destination to avoid message loss in case of a system failure on the destination side.
  batch-lines:    # Specifies how many lines are flushed to a destination in one batch. The syslog-ng OSE application waits for this number of lines to accumulate and sends them off in a single batch. Increasing this number increases throughput as more messages are sent in a single batch, but also increases message latency. For example, if you set batch-lines() to 100, syslog-ng OSE waits for 100 messages.
  batch-bytes:    # Sets the maximum size of payload in a batch. If the size of the messages reaches this value, syslog-ng OSE sends the batch to the destination even if the number of messages is less than the value of the batch-lines() option. Note that if the batch-timeout() option is enabled and the queue becomes empty, syslog-ng OSE flushes the messages only if batch-timeout() expires, or the batch reaches the limit set in batch-bytes().
  batch-timeout:  # Specifies the time syslog-ng OSE waits for lines to accumulate in the output buffer. The syslog-ng OSE application sends batches to the destinations evenly. The timer starts when the first message arrives to the buffer, so if only few messages arrive, syslog-ng OSE sends messages to the destination at most once every batch-timeout() milliseconds.
  body:           # The body of the HTTP request, for example, body("${ISODATE} ${MESSAGE}"). You can use strings, macros, and template functions in the body. If not set, it will contain the message received from the source by default.
  body-prefix:    # The string syslog-ng OSE puts at the beginning of the body of the HTTP request, before the log message.
  body-suffix:    # The string syslog-ng OSE puts to the end of the body of the HTTP request, after the log message.
  delimiter:      # By default, syslog-ng OSE separates the log messages of the batch with a newline character.
  method:         # Specifies the HTTP method to use when sending the message to the server. POST | PUT
  retries:        # The number of times syslog-ng OSE attempts to send a message to this destination. If syslog-ng OSE could not send a message, it will try again until the number of attempts reaches retries, then drops the message.
  user:           # The username that syslog-ng OSE uses to authenticate on the server where it sends the messages.
  password:       # The password that syslog-ng OSE uses to authenticate on the server where it sends the messages.
  user-agent:     # The value of the USER-AGENT header in the messages sent to the server.
  workers:        # Specifies the number of worker threads (at least 1) that syslog-ng OSE uses to send messages to the server. Increasing the number of worker threads can drastically improve the performance of the destination.
```

```yaml
apiVersion: logging.banzaicloud.io/v1beta1
kind: SyslogNGOutput
metadata:
  name: test-http
  namespace: default
spec:
  http:
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
    headers:
    - 'X-App-Header: example'
    tls:
      use-system-cert-store: true
```

### Loggly Output

The `loggly` output is a specialised `syslog` output. Extended with `token` and `tag`.

```yaml
apiVersion: logging.banzaicloud.io/v1beta1
kind: SyslogNGOutput
  spec:
    loggly:
      tag: "test-tag"
      token:
        valueFrom:
          secretKeyRef:
            name: loggly-secret
            key: token
      ...SyslogParams...
```

Message templating example:
```yaml
apiVersion: logging.banzaicloud.io/v1beta1
kind: SyslogNGOutput
  spec:
    loggly:
      tag: "test-tag"
      token:
        valueFrom:
          secretKeyRef:
            name: loggly-secret
            key: token
      template: "$(format-json
            --subkeys json.
            --exclude json.kubernetes.labels.*
            json.kubernetes.labels=literal($(format-flat-json --subkeys json.kubernetes.labels.)))\n"
      ...SyslogParams...
```

Loggly config in syslog: https://github.com/syslog-ng/syslog-ng/blob/master/scl/loggly/loggly.conf

### Sumologic-http output

The `sumologic-http` output sends log records over HTTP to Sumologic.

#### Parameters
```yaml
  body: # Body content template to send
  deployment: # Deployment code for sumologic. More info: https://help.sumologic.com/APIs/General-API-Information/Sumo-Logic-Endpoints-by-Deployment-and-Firewall-Security
  collector: # Sumo Logic service token (secret)
  headers: # Extra headers for Sumologic like X-Sumo-Name
  tls: # Required TLS configuration for Sumologic. Minimal config is use-system-cert-store: true
  disk_buffer: # Disk buffer parameters
  batch-lines: # Collect messages into batches number of lines (recommended)
  batch-bytes: # Collect messages into batches size of batch
  batch-timeout: # Time out for sending batch if no input available
```

```yaml
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
```

### Syslog output

The `syslog` output sends log records over a socket using the Syslog protocol (RFC 5424).

```yaml
apiVersion: logging.banzaicloud.io/v1beta1
kind: SyslogNGOutput
  spec:
    syslog:
      host: 10.12.34.56
      transport: tls
      tls:
        ca_file:
          mountFrom:
            secretKeyRef:
              name: tls-secret
              key: ca.crt
        cert_file:
          mountFrom:
            secretKeyRef:
              name: tls-secret
              key: tls.crt
        key_file:
          mountFrom:
            secretKeyRef:
              name: tls-secret
              key: tls.key
```

### Disk buffer

All syslog-ng outputs support buffering to disk.

#### Parameters
```yaml
  reliable:        # if set to true, syslog-ng cannot lose logs in case of a reload, restart, unreachable destination, or crash
  compaction:      # if set to true, syslog-ng prunes the unused space in the LogMessage representation
  dir:             # the folder where disk-buffer files are stored
  disk_buf_size:   # required; the maximum size of the disk-buffer in bytes
  mem_buf_length:  # the number of messages stored in overflow queue
  mem_buf_size:    # the size of messages in bytes that is used in the memory part of the disk buffer
  q_out_size:      # the number of messages stored in the output buffer of the destination
```

For more details, see [the official docs](https://www.syslog-ng.com/technical-documents/doc/syslog-ng-open-source-edition/3.37/administration-guide/disk-buffer).

---

Example kubernetes message:
```json
{"user":"-","ts":"2022-08-02T07:41:34.090763Z","time":"02/Aug/2022:07:41:34 +0000","stream":"stdout","source":"/var/log/log-generator","size":"1628","remote":"85.151.230.190","referer":"-","path":"/index.html","method":"POST","logtag":"F","kubernetes":{"pod_name":"one-eye-log-generator-57988cbd65-gpgc4","pod_id":"010d4598-2e34-4165-accd-2b77e4fc4bb6","namespace_name":"default","labels":{"pod-template-hash":"57988cbd65","app.kubernetes.io/name":"log-generator","app.kubernetes.io/instance":"one-eye-log-generator"},"host":"ip-192-168-6-51.eu-west-1.compute.internal","docker_id":"89e7cf414b5e6ff1fbee977d62cbc96794d2debd6a52803857d5dbad57d4f772","container_name":"log-generator","container_image":"ghcr.io/kube-logging/log-generator:0.3.20","container_hash":"ghcr.io/kube-logging/log-generator@sha256:b031138718194a17fdac2964bacf9543f96b037a65cd50138a5754ddb7897bb5"},"http_x_forwarded_for":"\"-\"","host":"-","code":"403","cluster":"xxxxx","agent":"Mozilla/5.0 (iPhone; CPU iPhone OS 10_3_3 like Mac OS X) AppleWebKit/603.3.8 (KHTML, like Gecko) Version/9.0 Mobile/13B143 Safari/601.1"}
```

## Complex example

```yaml
apiVersion: logging.banzaicloud.io/v1beta1
kind: SyslogNGFlow
metadata:
  name: TestFlow
  namespace: default
spec:
  match:
    and:
    - regexp:
        value: json.kubernetes.labels.app.kubernetes.io/instance
        pattern: one-eye-log-generator
        type: string
    - regexp:
        value:  json.kubernetes.labels.app.kubernetes.io/name
        pattern: log-generator
        type: string
  filters:
  -  parser:
       regexp:
         patterns:
         - '^(?<remote>[^ ]*) (?<host>[^ ]*) (?<user>[^ ]*) \[(?<time>[^\]]*)\] "(?<method>\S+)(?: +(?<path>[^\"]*?)(?: +\S*)?)?" (?<code>[^ ]*) (?<size>[^ ]*)(?: "(?<referer>[^\"]*)" "(?<agent>[^\"]*)"(?:\s+(?<http_x_forwarded_for>[^ ]+))?)?$'
         template: ${json.message}
         prefix: json.
  - rewrite:
    -  set:
         field: json.cluster
         value: xxxxx
    -  unset:
         field: json.message
    -  set:
         field: json.source
         value: /var/log/log-generator
         condition:
           regexp:
             value:  json.kubernetes.container_name
             pattern: log-generator
             type: string
  localOutputRefs:
    - syslog-output
```

```yaml
apiVersion: logging.banzaicloud.io/v1beta1
kind: SyslogNGOutput
metadata:
  name: test
  namespace: default
spec:
  syslog:
    host: 10.20.9.89
    port: 601
    disk_buffer:
      disk_buf_size: 512000000
      dir: /buffer
      reliable: true
    template: "$(format-json
                --subkeys json.
                --exclude json.kubernetes.labels.*
                json.kubernetes.labels=literal($(format-flat-json --subkeys json.kubernetes.labels.)))\n"
    tls:
      ca_file:
        mountFrom:
          secretKeyRef:
            key: ca.crt
            name: syslog-tls-cert
      cert_file:
        mountFrom:
          secretKeyRef:
            key: tls.crt
            name: syslog-tls-cert
      key_file:
        mountFrom:
          secretKeyRef:
            key: tls.key
            name: syslog-tls-cert
    transport: tls
```


```yaml
apiVersion: logging.banzaicloud.io/v1beta1
kind: Logging
  name: test
spec:
  syslogNG:
    statefulSet:
      spec:
        template:
          spec:
            containers:
            - name: syslog-ng
              volumeMounts:
              - mountPath: /buffers
                name: buffer
        volumeClaimTemplates:
        - metadata:
            name: buffer
          spec:
            accessModes:
            - ReadWriteOnce
            resources:
              requests:
                storage: 10Gi
```