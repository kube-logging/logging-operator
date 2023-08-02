---
title: Loggly output
weight: 200
generated_file: true
---

# Loggly output plugin for syslog-ng
## Overview
 The `loggly()` destination sends log messages to the [Loggly](https://www.loggly.com/) Logging-as-a-Service provider. You can send log messages over TCP, or encrypted with TLS. For details, see the [syslog-ng documentation](https://www.syslog-ng.com/technical-documents/doc/syslog-ng-open-source-edition/3.37/administration-guide/43#TOPIC-1829072).

 ## Prerequisites

 You need a Loggly account and your user token to use this output.

## Configuration
## Loggly

### host (string, optional) {#loggly-host}

Address of the destination host 

Default: -

### tag (string, optional) {#loggly-tag}

Event tag [more information](https://documentation.solarwinds.com/en/success_center/loggly/content/admin/tags.htm) 

Default: -

### token (*secret.Secret, required) {#loggly-token}

Your Customer Token that you received from Loggly [more information](https://www.syslog-ng.com/technical-documents/doc/syslog-ng-open-source-edition/3.37/administration-guide/43#loggly-option-token) 

Default: -

###  (SyslogOutput, required) {#loggly-}

syslog output configuration 

Default: -


