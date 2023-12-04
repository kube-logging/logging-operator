---
title: SyslogNGOutputSpec
weight: 200
generated_file: true
---

## SyslogNGOutputSpec

SyslogNGOutputSpec defines the desired state of SyslogNGOutput

### elasticsearch (*output.ElasticsearchOutput, optional) {#syslogngoutputspec-elasticsearch}


### file (*output.FileOutput, optional) {#syslogngoutputspec-file}


### http (*output.HTTPOutput, optional) {#syslogngoutputspec-http}


### logscale (*output.LogScaleOutput, optional) {#syslogngoutputspec-logscale}


### loggingRef (string, optional) {#syslogngoutputspec-loggingref}


### loggly (*output.Loggly, optional) {#syslogngoutputspec-loggly}


### loki (*output.LokiOutput, optional) {#syslogngoutputspec-loki}

Available in Logging operator version 4.4 and later. 


### mqtt (*output.MQTT, optional) {#syslogngoutputspec-mqtt}


### mongodb (*output.MongoDB, optional) {#syslogngoutputspec-mongodb}


### openobserve (*output.OpenobserveOutput, optional) {#syslogngoutputspec-openobserve}

Available in Logging operator version 4.5 and later. 


### redis (*output.RedisOutput, optional) {#syslogngoutputspec-redis}


### s3 (*output.S3Output, optional) {#syslogngoutputspec-s3}

Available in Logging operator version 4.4 and later. 


### splunk_hec_event (*output.SplunkHECOutput, optional) {#syslogngoutputspec-splunk_hec_event}


### sumologic-http (*output.SumologicHTTPOutput, optional) {#syslogngoutputspec-sumologic-http}


### sumologic-syslog (*output.SumologicSyslogOutput, optional) {#syslogngoutputspec-sumologic-syslog}


### syslog (*output.SyslogOutput, optional) {#syslogngoutputspec-syslog}



## SyslogNGOutput

SyslogNGOutput is the Schema for the syslog-ng outputs API

###  (metav1.TypeMeta, required) {#syslogngoutput-}


### metadata (metav1.ObjectMeta, optional) {#syslogngoutput-metadata}


### spec (SyslogNGOutputSpec, optional) {#syslogngoutput-spec}


### status (SyslogNGOutputStatus, optional) {#syslogngoutput-status}



## SyslogNGOutputList

SyslogNGOutputList contains a list of SyslogNGOutput

###  (metav1.TypeMeta, required) {#syslogngoutputlist-}


### metadata (metav1.ListMeta, optional) {#syslogngoutputlist-metadata}


### items ([]SyslogNGOutput, required) {#syslogngoutputlist-items}



