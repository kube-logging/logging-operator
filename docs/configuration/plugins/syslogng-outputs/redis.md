---
title: Redis
weight: 200
generated_file: true
---

# Sending messages from a local network to the Redis server
## Overview

Based on the [Redis destination of AxoSyslog core](https://axoflow.com/docs/axosyslog-core/chapter-destinations/configuring-destinations-redis/).

Available in Logging operator version 4.4 and later.

## Example

{{< highlight yaml >}}
apiVersion: logging.banzaicloud.io/v1beta1
kind: SyslogNGOutput
metadata:
  name: redis
  namespace: default
spec:
  redis:
    host: 127.0.0.1
	port: 6379
	retries: 3
	throttle: 0
	time-reopen: 60
	workers: 1
{{</ highlight >}}

For details on the available options of the output, see the [documentation of the AxoSyslog syslog-ng distribution](https://axoflow.com/docs/axosyslog-core/chapter-destinations/configuring-destinations-redis/).


## Configuration
## RedisOutput

###  (Batch, required) {#redisoutput-}

Batching parameters 


### auth (*secret.Secret, optional) {#redisoutput-auth}

The password used for authentication on a password-protected Redis server. 


### command_and_arguments ([]string, optional) {#redisoutput-command_and_arguments}

The Redis command to execute, for example, LPUSH, INCR, or HINCRBY. Using the HINCRBY command with an increment value of 1 allows you to create various statistics. For example, the `command("HINCRBY" "${HOST}/programs" "${PROGRAM}" "1")` command counts the number of log messages on each host for each program.

Default: ""

### disk_buffer (*DiskBuffer, optional) {#redisoutput-disk_buffer}

This option enables putting outgoing messages into the disk buffer of the destination to avoid message loss in case of a system failure on the destination side. For details, see the [Syslog-ng DiskBuffer options](../disk_buffer/).

Default: false

### host (string, optional) {#redisoutput-host}

The hostname or IP address of the Redis server.

Default: 127.0.0.1

### log-fifo-size (int, optional) {#redisoutput-log-fifo-size}

The number of messages that the output queue can store. 


### persist_name (string, optional) {#redisoutput-persist_name}

Persistname 


### port (int, optional) {#redisoutput-port}

The port number of the Redis server.

Default: 6379

### command (StringList, optional) {#redisoutput-command}

Internal rendered form of the CommandAndArguments field 


### retries (int, optional) {#redisoutput-retries}

If syslog-ng OSE cannot send a message, it will try again until the number of attempts reaches `retries()`.

Default: 3

### throttle (int, optional) {#redisoutput-throttle}

Sets the maximum number of messages sent to the destination per second. Use this output-rate-limiting functionality only when using disk-buffer as well to avoid the risk of losing messages. Specifying 0 or a lower value sets the output limit to unlimited.

Default: 0

### time-reopen (int, optional) {#redisoutput-time-reopen}

The time to wait in seconds before a dead connection is reestablished.

Default: 60

### workers (int, optional) {#redisoutput-workers}

Specifies the number of worker threads (at least 1) that syslog-ng OSE uses to send messages to the server. Increasing the number of worker threads can drastically improve the performance of the destination.

Default: 1


## StringList

### string-list ([]string, optional) {#stringlist-string-list}



