---
title: Loggly output
weight: 200
generated_file: true
---

# Loggly output plugin for syslog-ng
## Overview

The `loggly()` destination sends log messages to the [Loggly](https://www.loggly.com/) Logging-as-a-Service provider.
You can send log messages over TCP, or encrypted with [TLS for syslog-ng outputs](/docs/configuration/plugins/syslog-ng-outputs/tls/).

For details on the available options of the output, see the [documentation of the AxoSyslog syslog-ng distribution](https://axoflow.com/docs/axosyslog-core/chapter-destinations/configuring-destinations-loggly/).

## Prerequisites

You need a Loggly account and your user token to use this output.


## Configuration
## Loggly

Documentation: https://github.com/syslog-ng/syslog-ng/blob/master/scl/loggly/loggly.conf

###  (SyslogOutput, required) {#loggly-}

syslog output configuration 


### host (string, optional) {#loggly-host}

Address of the destination host. 


### tag (string, optional) {#loggly-tag}

Event tag. For details, see the [Loggy documentation](https://documentation.solarwinds.com/en/success_center/loggly/content/admin/tags.htm) 


### token (*secret.Secret, required) {#loggly-token}

Your Customer Token that you received from Loggly. For details, see the [documentation of the AxoSyslog syslog-ng distribution](https://axoflow.com/docs/axosyslog-core/chapter-destinations/configuring-destinations-loggly/) 



