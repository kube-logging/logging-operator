---
title: Loki
weight: 200
generated_file: true
---

# Sending messages to Loki over HTTP
## Overview
 More info at https://axoflow.com/docs/axosyslog-core/chapter-destinations/syslog-ng-with-loki/

## Configuration
## LokiOutput

### labels (filter.ArrowMap, optional) {#lokioutput-labels}

Label mapping from kubernetes labels to Loki labels. 

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

###  (Batch, required) {#lokioutput-}

Batching parameters 

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

The timestamp that will be applied to the outgoing messages . Loki does not accept events, in which the timestamp is not monotonically increasing. 

Default:  current

### template (string, optional) {#lokioutput-template}

Template for customizing the log message format. 

Default: -


