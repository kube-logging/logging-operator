---
title: syslog-ng syslog (rfc5424) output
weight: 200
generated_file: true
---

# Syslog output configuration
## Overview
 More info at https://www.syslog-ng.com/technical-documents/doc/syslog-ng-open-source-edition/3.37/administration-guide/32#kanchor2338

## Configuration
## SyslogOutput

Documentation: https://www.syslog-ng.com/technical-documents/doc/syslog-ng-open-source-edition/3.37/administration-guide/56#TOPIC-1829124

### host (string, optional) {#syslogoutput-host}

Address of the destination host 

Default: -

### port (int, optional) {#syslogoutput-port}

The port number to connect to. [more information](https://www.syslog-ng.com/technical-documents/doc/syslog-ng-open-source-edition/3.37/administration-guide/56#kanchor895) 

Default: -

### transport (string, optional) {#syslogoutput-transport}

Specifies the protocol used to send messages to the destination server. [more information]() [more information](https://www.syslog-ng.com/technical-documents/doc/syslog-ng-open-source-edition/3.37/administration-guide/56#kanchor911) 

Default: -

### close_on_input (*bool, optional) {#syslogoutput-close_on_input}

By default, syslog-ng OSE closes destination sockets if it receives any input from the socket (for example, a reply). If this option is set to no, syslog-ng OSE just ignores the input, but does not close the socket. [more information](https://www.syslog-ng.com/technical-documents/doc/syslog-ng-open-source-edition/3.37/administration-guide/56#kanchor859) 

Default: -

### flags ([]string, optional) {#syslogoutput-flags}

Flags influence the behavior of the destination driver. [more information](https://www.syslog-ng.com/technical-documents/doc/syslog-ng-open-source-edition/3.37/administration-guide/56#kanchor877) 

Default: -

### flush_lines (int, optional) {#syslogoutput-flush_lines}

Specifies how many lines are flushed to a destination at a time. [more information](https://www.syslog-ng.com/technical-documents/doc/syslog-ng-open-source-edition/3.37/administration-guide/56#kanchor880) 

Default: -

### so_keepalive (*bool, optional) {#syslogoutput-so_keepalive}

Enables keep-alive messages, keeping the socket open. [more information](https://www.syslog-ng.com/technical-documents/doc/syslog-ng-open-source-edition/3.37/administration-guide/56#kanchor897) 

Default: -

### suppress (int, optional) {#syslogoutput-suppress}

Specifies the number of seconds syslog-ng waits for identical messages. [more information](https://www.syslog-ng.com/technical-documents/doc/syslog-ng-open-source-edition/3.37/administration-guide/56#kanchor901) 

Default: -

### template (string, optional) {#syslogoutput-template}

Specifies a template defining the logformat to be used in the destination. [more information](https://www.syslog-ng.com/technical-documents/doc/syslog-ng-open-source-edition/3.37/administration-guide/56#kanchor905)  

Default:  0

### template_escape (*bool, optional) {#syslogoutput-template_escape}

Turns on escaping for the ', ", and backspace characters in templated output files. [more information](https://www.syslog-ng.com/technical-documents/doc/syslog-ng-open-source-edition/3.37/administration-guide/56#kanchor906) 

Default: -

### tls (*TLS, optional) {#syslogoutput-tls}

Sets various options related to TLS encryption, for example, key/certificate files and trusted CA locations. [more information](https://www.syslog-ng.com/technical-documents/doc/syslog-ng-open-source-edition/3.37/administration-guide/56#kanchor910) 

Default: -

### ts_format (string, optional) {#syslogoutput-ts_format}

Override the global timestamp format (set in the global ts-format() parameter) for the specific destination. [more information](https://www.syslog-ng.com/technical-documents/doc/syslog-ng-open-source-edition/3.37/administration-guide/56#kanchor912) 

Default: -

### disk_buffer (*DiskBuffer, optional) {#syslogoutput-disk_buffer}

Enables putting outgoing messages into the disk buffer of the destination to avoid message loss in case of a system failure on the destination side. [more information](https://www.syslog-ng.com/technical-documents/doc/syslog-ng-open-source-edition/3.37/administration-guide/56#kanchor860) 

Default: -

### persist_name (string, optional) {#syslogoutput-persist_name}

Unique name for the syslog-ng driver [more information](https://www.syslog-ng.com/technical-documents/doc/syslog-ng-open-source-edition/3.16/administration-guide/persist-name) 

Default: -


