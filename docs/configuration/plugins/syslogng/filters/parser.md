---
title: Syslog-NG parser
weight: 200
generated_file: true
---

# [Syslog-NG Parser Filter](https://www.syslog-ng.com/technical-documents/doc/syslog-ng-open-source-edition/3.36/administration-guide/90)
## Overview

## Configuration
## [Parser](https://www.syslog-ng.com/technical-documents/doc/syslog-ng-open-source-edition/3.36/administration-guide/82#TOPIC-1768819)

### regexp (*RegexpParser, optional) {#[parser](https://www.syslog-ng.com/technical-documents/doc/syslog-ng-open-source-edition/3.36/administration-guide/82#topic-1768819)-regexp}

Default: -

### syslog-parser (*SyslogParser, optional) {#[parser](https://www.syslog-ng.com/technical-documents/doc/syslog-ng-open-source-edition/3.36/administration-guide/82#topic-1768819)-syslog-parser}

Default: -


## [Regexp parser](https://www.syslog-ng.com/technical-documents/doc/syslog-ng-open-source-edition/3.36/administration-guide/90)

### patterns ([]string, required) {#[regexp parser](https://www.syslog-ng.com/technical-documents/doc/syslog-ng-open-source-edition/3.36/administration-guide/90)-patterns}

The regular expression patterns that you want to find a match. regexp-parser() supports multiple patterns, and stops the processing at the first successful match. 

Default: -

### prefix (string, optional) {#[regexp parser](https://www.syslog-ng.com/technical-documents/doc/syslog-ng-open-source-edition/3.36/administration-guide/90)-prefix}

Insert a prefix before the name part of the parsed name-value pairs to help further processing. 

Default: -

### template (string, optional) {#[regexp parser](https://www.syslog-ng.com/technical-documents/doc/syslog-ng-open-source-edition/3.36/administration-guide/90)-template}

Specify a template of the record fields to match against. 

Default: -

### flags ([]string, optional) {#[regexp parser](https://www.syslog-ng.com/technical-documents/doc/syslog-ng-open-source-edition/3.36/administration-guide/90)-flags}

Pattern flags 

Default: -


## SyslogParser

### flags ([]string, optional) {#syslogparser-flags}

Pattern flags 

Default: -


