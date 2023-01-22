---
title: Loggly output
weight: 200
generated_file: true
---

# Loggly output plugin for syslog-ng
## Overview
 More info at https://www.syslog-ng.com/technical-documents/doc/syslog-ng-open-source-edition/3.37/administration-guide/43#TOPIC-1829072

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


