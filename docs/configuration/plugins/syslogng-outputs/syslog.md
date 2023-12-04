---
title: Syslog (RFC5424) output
weight: 200
generated_file: true
---

# Syslog output configuration
## Overview

The `syslog` output sends log records over a socket using the Syslog protocol (RFC 5424). Based on the [syslog destination of AxoSyslog core](https://axoflow.com/docs/axosyslog-core/chapter-destinations/configuring-destinations-syslog/).


{{< highlight yaml >}}
kind: SyslogNGOutput
metadata:
  name: test
  namespace: default
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
{{</ highlight >}}

The following example also configures disk-based buffering for the output. For details, see the [Syslog-ng DiskBuffer options](../disk_buffer/).

{{< highlight yaml >}}
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
{{</ highlight >}}

For details on the available options of the output, see the [documentation of the AxoSyslog syslog-ng distribution](https://axoflow.com/docs/axosyslog-core/chapter-destinations/configuring-destinations-syslog/).


## Configuration
## SyslogOutput

Documentation: https://axoflow.com/docs/axosyslog-core/chapter-destinations/configuring-destinations-syslog/

### close_on_input (*bool, optional) {#syslogoutput-close_on_input}

By default, syslog-ng OSE closes destination sockets if it receives any input from the socket (for example, a reply). If this option is set to no, syslog-ng OSE just ignores the input, but does not close the socket. For details, see the [documentation of the AxoSyslog syslog-ng distribution](https://axoflow.com/docs/axosyslog-core/chapter-destinations/configuring-destinations-syslog/reference-destination-syslog-chapter/#close-on-input). 


### disk_buffer (*DiskBuffer, optional) {#syslogoutput-disk_buffer}

Enables putting outgoing messages into the disk buffer of the destination to avoid message loss in case of a system failure on the destination side. For details, see the [Syslog-ng DiskBuffer options](../disk_buffer/). 


### flags ([]string, optional) {#syslogoutput-flags}

Flags influence the behavior of the destination driver. For details, see the [documentation of the AxoSyslog syslog-ng distribution](https://axoflow.com/docs/axosyslog-core/chapter-destinations/configuring-destinations-syslog/reference-destination-syslog-chapter/#flags). 


### flush_lines (int, optional) {#syslogoutput-flush_lines}

Specifies how many lines are flushed to a destination at a time. For details, see the [documentation of the AxoSyslog syslog-ng distribution](https://axoflow.com/docs/axosyslog-core/chapter-destinations/configuring-destinations-syslog/reference-destination-syslog-chapter/#flush-lines). 


### host (string, optional) {#syslogoutput-host}

Address of the destination host 


### persist_name (string, optional) {#syslogoutput-persist_name}

Unique name for the syslog-ng driver. If you receive the following error message during syslog-ng startup, set the `persist-name()` option of the duplicate drivers: `Error checking the uniqueness of the persist names, please override it with persist-name option. Shutting down.` See the [documentation of the AxoSyslog syslog-ng distribution](https://axoflow.com/docs/axosyslog-core/chapter-destinations/configuring-destinations-http-nonjava/reference-destination-http-nonjava/#persist-name) for more information. 


### port (int, optional) {#syslogoutput-port}

The port number to connect to. For details, see the [documentation of the AxoSyslog syslog-ng distribution](https://axoflow.com/docs/axosyslog-core/chapter-destinations/configuring-destinations-syslog/reference-destination-syslog-chapter/#port-or-destport). 


### so_keepalive (*bool, optional) {#syslogoutput-so_keepalive}

Enables keep-alive messages, keeping the socket open. For details, see the [documentation of the AxoSyslog syslog-ng distribution](https://axoflow.com/docs/axosyslog-core/chapter-destinations/configuring-destinations-syslog/reference-destination-syslog-chapter/#so-keepalive). 


### suppress (int, optional) {#syslogoutput-suppress}

Specifies the number of seconds syslog-ng waits for identical messages. For details, see the [documentation of the AxoSyslog syslog-ng distribution](https://axoflow.com/docs/axosyslog-core/chapter-destinations/configuring-destinations-syslog/reference-destination-syslog-chapter/#suppress). 


### tls (*TLS, optional) {#syslogoutput-tls}

Sets various options related to TLS encryption, for example, key/certificate files and trusted CA locations. For details, see the [documentation of the AxoSyslog syslog-ng distribution](https://axoflow.com/docs/axosyslog-core/chapter-destinations/configuring-destinations-syslog/reference-destination-syslog-chapter/#tls). 


### ts_format (string, optional) {#syslogoutput-ts_format}

Override the global timestamp format (set in the global ts-format() parameter) for the specific destination. For details, see the [documentation of the AxoSyslog syslog-ng distribution](https://axoflow.com/docs/axosyslog-core/chapter-destinations/configuring-destinations-syslog/reference-destination-syslog-chapter/#ts-format). 


### template (string, optional) {#syslogoutput-template}

Specifies a template defining the logformat to be used in the destination. For details, see the [documentation of the AxoSyslog syslog-ng distribution](https://axoflow.com/docs/axosyslog-core/chapter-destinations/configuring-destinations-syslog/reference-destination-syslog-chapter/#template).

Default: 0

### template_escape (*bool, optional) {#syslogoutput-template_escape}

Turns on escaping for the ', ", and backspace characters in templated output files. For details, see the [documentation of the AxoSyslog syslog-ng distribution](https://axoflow.com/docs/axosyslog-core/chapter-destinations/configuring-destinations-syslog/reference-destination-syslog-chapter/#template-escape). 


### transport (string, optional) {#syslogoutput-transport}

Specifies the protocol used to send messages to the destination server. For details, see the [documentation of the AxoSyslog syslog-ng distribution](https://axoflow.com/docs/axosyslog-core/chapter-destinations/configuring-destinations-syslog/reference-destination-syslog-chapter/#transport). 



