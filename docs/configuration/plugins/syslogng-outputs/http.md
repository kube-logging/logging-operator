---
title: HTTP
weight: 200
generated_file: true
---

# Sending messages over HTTP
## Overview

A simple example sending logs over HTTP to a fluentbit HTTP endpoint:
{{< highlight yaml >}}
kind: SyslogNGOutput
apiVersion: logging.banzaicloud.io/v1beta1
metadata:
  name: http
spec:
  http:
    #URL of the ingest endpoint
    url: http://fluentbit-endpoint:8080/tag
    method: POST
    headers:
      - "Content-type: application/json"
{{</ highlight >}}

A more complex example to demonstrate sending logs to OpenObserve
{{< highlight yaml >}}
kind: SyslogNGOutput
apiVersion: logging.banzaicloud.io/v1beta1
metadata:
  name: openobserve
spec:
  http:
    #URL of the ingest endpoint
    url: https://openobserve-endpoint/api/default/log-generator/_json
    user: "username"
    password:
      valueFrom:
        secretKeyRef:
          name: openobserve
          key: password
    method: POST
    # Parameters for sending logs in batches
    batch-lines: 5000
    batch-bytes: 4096
    batch-timeout: 300
    headers:
      - "Connection: keep-alive"
    # Disable TLS peer verification for demo
    tls:
      peer_verify: "no"
    body-prefix: "["
    body-suffix: "]"
    delimiter: ","
    body: "${MESSAGE}"
{{</ highlight >}}

More information at: https://axoflow.com/docs/axosyslog-core/chapter-destinations/configuring-destinations-http-nonjava/


## Configuration
## HTTPOutput

### url (string, optional) {#httpoutput-url}

Specifies the hostname or IP address and optionally the port number of the web service that can receive log data via HTTP. Use a colon (:) after the address to specify the port number of the server. For example: http://127.0.0.1:8000 

Default: -

### headers ([]string, optional) {#httpoutput-headers}

Custom HTTP headers to include in the request, for example, headers("HEADER1: header1", "HEADER2: header2").

Default: empty

### time_reopen (int, optional) {#httpoutput-time_reopen}

The time to wait in seconds before a dead connection is reestablished.

Default: 60

### tls (*TLS, optional) {#httpoutput-tls}

This option sets various options related to TLS encryption, for example, key/certificate files and trusted CA locations. TLS can be used only with tcp-based transport protocols. For details, see [TLS for syslog-ng outputs](/docs/configuration/plugins/syslog-ng-outputs/tls/). 

Default: -

### disk_buffer (*DiskBuffer, optional) {#httpoutput-disk_buffer}

This option enables putting outgoing messages into the disk buffer of the destination to avoid message loss in case of a system failure on the destination side. For details, see the [Syslog-ng DiskBuffer options](../disk_buffer/).

Default: false

###  (Batch, required) {#httpoutput-}

Batching parameters 

Default: -

### body (string, optional) {#httpoutput-body}

The body of the HTTP request, for example, body("${ISODATE} ${MESSAGE}"). You can use strings, macros, and template functions in the body. If not set, it will contain the message received from the source by default. 

Default: -

### body-prefix (string, optional) {#httpoutput-body-prefix}

The string syslog-ng OSE puts at the beginning of the body of the HTTP request, before the log message. 

Default: -

### body-suffix (string, optional) {#httpoutput-body-suffix}

The string syslog-ng OSE puts to the end of the body of the HTTP request, after the log message. 

Default: -

### delimiter (string, optional) {#httpoutput-delimiter}

By default, syslog-ng OSE separates the log messages of the batch with a newline character. 

Default: -

### method (string, optional) {#httpoutput-method}

Specifies the HTTP method to use when sending the message to the server. POST | PUT 

Default: -

### retries (int, optional) {#httpoutput-retries}

The number of times syslog-ng OSE attempts to send a message to this destination. If syslog-ng OSE could not send a message, it will try again until the number of attempts reaches retries, then drops the message. 

Default: -

### user (string, optional) {#httpoutput-user}

The username that syslog-ng OSE uses to authenticate on the server where it sends the messages. 

Default: -

### password (secret.Secret, optional) {#httpoutput-password}

The password that syslog-ng OSE uses to authenticate on the server where it sends the messages. 

Default: -

### user-agent (string, optional) {#httpoutput-user-agent}

The value of the USER-AGENT header in the messages sent to the server. 

Default: -

### workers (int, optional) {#httpoutput-workers}

Specifies the number of worker threads (at least 1) that syslog-ng OSE uses to send messages to the server. Increasing the number of worker threads can drastically improve the performance of the destination. 

Default: -

### persist_name (string, optional) {#httpoutput-persist_name}

If you receive the following error message during AxoSyslog startup, set the persist-name() option of the duplicate drivers: `Error checking the uniqueness of the persist names, please override it with persist-name option. Shutting down.` See [syslog-ng docs](https://axoflow.com/docs/axosyslog-core/chapter-destinations/configuring-destinations-http-nonjava/reference-destination-http-nonjava/#persist-name) for more information. 

Default: -

### log-fifo-size (int, optional) {#httpoutput-log-fifo-size}

The number of messages that the output queue can store. 

Default: -

### timeout (int, optional) {#httpoutput-timeout}

Sets the maximum number of messages sent to the destination per second. Use this output-rate-limiting functionality only when using disk-buffer as well to avoid the risk of losing messages. Specifying 0 or a lower value sets the output limit to unlimited. 

Default: -

### response-action (filter.RawArrowMap, optional) {#httpoutput-response-action}

Specifies what AxoSyslog does with the log message, based on the response code received from the HTTP server. See [syslog-ng docs](https://axoflow.com/docs/axosyslog-core/chapter-destinations/configuring-destinations-http-nonjava/reference-destination-http-nonjava/#response-action) for more information. 

Default: -


## Batch

### batch-lines (int, optional) {#batch-batch-lines}

Description: Specifies how many lines are flushed to a destination in one batch. The syslog-ng OSE application waits for this number of lines to accumulate and sends them off in a single batch. Increasing this number increases throughput as more messages are sent in a single batch, but also increases message latency. For example, if you set batch-lines() to 100, syslog-ng OSE waits for 100 messages. 

Default: -

### batch-bytes (int, optional) {#batch-batch-bytes}

Description: Sets the maximum size of payload in a batch. If the size of the messages reaches this value, syslog-ng OSE sends the batch to the destination even if the number of messages is less than the value of the batch-lines() option. Note that if the batch-timeout() option is enabled and the queue becomes empty, syslog-ng OSE flushes the messages only if batch-timeout() expires, or the batch reaches the limit set in batch-bytes(). 

Default: -

### batch-timeout (int, optional) {#batch-batch-timeout}

Description: Specifies the time syslog-ng OSE waits for lines to accumulate in the output buffer. The syslog-ng OSE application sends batches to the destinations evenly. The timer starts when the first message arrives to the buffer, so if only few messages arrive, syslog-ng OSE sends messages to the destination at most once every batch-timeout() milliseconds. 

Default: -


