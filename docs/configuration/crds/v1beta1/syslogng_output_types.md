---
title: SyslogNGOutputSpec
weight: 200
generated_file: true
---

## SyslogNGOutputSpec

SyslogNGOutputSpec defines the desired state of SyslogNGOutput

### loggingRef (string, optional) {#syslogngoutputspec-loggingref}

Default: -

### loggly (*output.Loggly, optional) {#syslogngoutputspec-loggly}

Default: -

### syslog (*output.SyslogOutput, optional) {#syslogngoutputspec-syslog}

Default: -

### file (*output.FileOutput, optional) {#syslogngoutputspec-file}

Default: -

### mqtt (*output.MQTT, optional) {#syslogngoutputspec-mqtt}

Default: -

### sumologic-http (*output.SumologicHTTPOutput, optional) {#syslogngoutputspec-sumologic-http}

Default: -

### sumologic-syslog (*output.SumologicSyslogOutput, optional) {#syslogngoutputspec-sumologic-syslog}

Default: -

### http (*output.HTTPOutput, optional) {#syslogngoutputspec-http}

Default: -

### logscale (*output.LogScaleOutput, optional) {#syslogngoutputspec-logscale}

Default: -


## SyslogNGOutput

SyslogNGOutput is the Schema for the syslog-ng outputs API

###  (metav1.TypeMeta, required) {#syslogngoutput-}

Default: -

### metadata (metav1.ObjectMeta, optional) {#syslogngoutput-metadata}

Default: -

### spec (SyslogNGOutputSpec, optional) {#syslogngoutput-spec}

Default: -

### status (SyslogNGOutputStatus, optional) {#syslogngoutput-status}

Default: -


## SyslogNGOutputList

SyslogNGOutputList contains a list of SyslogNGOutput

###  (metav1.TypeMeta, required) {#syslogngoutputlist-}

Default: -

### metadata (metav1.ListMeta, optional) {#syslogngoutputlist-metadata}

Default: -

### items ([]SyslogNGOutput, required) {#syslogngoutputlist-items}

Default: -


