---
title: Loggly output
weight: 200
generated_file: true
---

# Loggly output plugin for syslog-ng
## Overview
 The `loggly()` destination sends log messages to the [Loggly](https://www.loggly.com/) Logging-as-a-Service provider.
 You can send log messages over TCP, or encrypted with [TLS for syslog-ng outputs](/docs/configuration/plugins/syslog-ng-outputs/tls/).

 ## Prerequisites

 You need a Loggly account and your user token to use this output.

## Configuration
## Loggly

Documentation: https://github.com/syslog-ng/syslog-ng/blob/master/scl/loggly/loggly.conf

### host (string, optional) {#loggly-host}

Address of the destination host 

Default: -

### tag (string, optional) {#loggly-tag}

Event tag [more information](https://documentation.solarwinds.com/en/success_center/loggly/content/admin/tags.htm) 

Default: -

### token (*secret.Secret, required) {#loggly-token}

Your Customer Token that you received from Loggly [more information](https://axoflow.com/docs/axosyslog-core/chapter-destinations/configuring-destinations-loggly/) 

Default: -

###  (SyslogOutput, required) {#loggly-}

syslog output configuration 

Default: -


