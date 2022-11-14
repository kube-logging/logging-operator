---
title: Syslog-ng DiskBuffer
weight: 200
generated_file: true
---

# disk_buffer
## Overview
 More info at https://www.syslog-ng.com/technical-documents/doc/syslog-ng-open-source-edition/3.37/administration-guide/32#kanchor2338

## Configuration
## DiskBuffer

Documentation: https://www.syslog-ng.com/technical-documents/doc/syslog-ng-open-source-edition/3.37/administration-guide/56#TOPIC-1829124

### disk_buf_size (int64, required) {#diskbuffer-disk_buf_size}

This is a required option. The maximum size of the disk-buffer in bytes. The minimum value is 1048576 bytes. 

Default: -

### reliable (bool, required) {#diskbuffer-reliable}

If set to yes, syslog-ng OSE cannot lose logs in case of reload/restart, unreachable destination or syslog-ng OSE crash. This solution provides a slower, but reliable disk-buffer option. 

Default: -

### compaction (*bool, optional) {#diskbuffer-compaction}

Prunes the unused space in the LogMessage representation 

Default: -

### dir (string, optional) {#diskbuffer-dir}

Description: Defines the folder where the disk-buffer files are stored. 

Default: -

### mem_buf_length (*int64, optional) {#diskbuffer-mem_buf_length}

Use this option if the option reliable() is set to no. This option contains the number of messages stored in overflow queue. 

Default: -

### mem_buf_size (*int64, optional) {#diskbuffer-mem_buf_size}

Use this option if the option reliable() is set to yes. This option contains the size of the messages in bytes that is used in the memory part of the disk buffer. 

Default: -

### q_out_size (*int64, optional) {#diskbuffer-q_out_size}

The number of messages stored in the output buffer of the destination. 

Default: -


