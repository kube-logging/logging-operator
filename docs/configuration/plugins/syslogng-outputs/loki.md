---
title: Loki
weight: 200
generated_file: true
---

# Sending messages to Loki over gRPC
## Overview
 More info at https://axoflow.com/docs/axosyslog-core/chapter-destinations/syslog-ng-with-loki/

## Configuration
## LokiOutput

### labels (filter.ArrowMap, optional) {#lokioutput-labels}

Using the Labels map, Kubernetes label to Loki label mapping can be configured. Example: {"app" : "$PROGRAM"} 

Default: -

### url (string, optional) {#lokioutput-url}

Specifies the hostname or IP address and optionally the port number of the web service that can receive log data via HTTP. Use a colon (:) after the address to specify the port number of the server. For example: http://127.0.0.1:8000 

Default: -

### time_reopen (int, optional) {#lokioutput-time_reopen}

The time to wait in seconds before a dead connection is reestablished.  

Default:  60

### disk_buffer (*DiskBuffer, optional) {#lokioutput-disk_buffer}

This option enables putting outgoing messages into the disk buffer of the destination to avoid message loss in case of a system failure on the destination side. For details, see the [Syslog-ng DiskBuffer options](../disk_buffer/).  

Default:  false

### batch-lines (int, optional) {#lokioutput-batch-lines}

Description: Specifies how many lines are flushed to a destination in one batch. The syslog-ng OSE application waits for this number of lines to accumulate and sends them off in a single batch. Increasing this number increases throughput as more messages are sent in a single batch, but also increases message latency. For example, if you set batch-lines() to 100, syslog-ng OSE waits for 100 messages. 

Default: -

### batch-timeout (int, optional) {#lokioutput-batch-timeout}

Description: Specifies the time syslog-ng OSE waits for lines to accumulate in the output buffer. The syslog-ng OSE application sends batches to the destinations evenly. The timer starts when the first message arrives to the buffer, so if only few messages arrive, syslog-ng OSE sends messages to the destination at most once every batch-timeout() milliseconds. 

Default: -

### retries (int, optional) {#lokioutput-retries}

The number of times syslog-ng OSE attempts to send a message to this destination. If syslog-ng OSE could not send a message, it will try again until the number of attempts reaches retries, then drops the message. 

Default: -

### workers (int, optional) {#lokioutput-workers}

Specifies the number of worker threads (at least 1) that syslog-ng OSE uses to send messages to the server. Increasing the number of worker threads can drastically improve the performance of the destination. 

Default: -

### persist_name (string, optional) {#lokioutput-persist_name}

If you receive the following error message during AxoSyslog startup, set the persist-name() option of the duplicate drivers: `Error checking the uniqueness of the persist names, please override it with persist-name option. Shutting down.` See [syslog-ng docs](https://axoflow.com/docs/axosyslog-core/chapter-destinations/configuring-destinations-http-nonjava/reference-destination-http-nonjava/#persist-name) for more information. 

Default: -

### log-fifo-size (int, optional) {#lokioutput-log-fifo-size}

The number of messages that the output queue can store. 

Default: -

### timestamp (string, optional) {#lokioutput-timestamp}

The timestamp that will be applied to the outgoing messages (possible values: current|received|msg default: current). Loki does not accept events, in which the timestamp is not monotonically increasing. 

Default: -

### template (string, optional) {#lokioutput-template}

Template for customizing the log message format. 

Default: -


