---
title: Elasticsearch
weight: 200
generated_file: true
---

# Sending messages over Elasticsearch
## Overview

## Example

{{< highlight yaml >}}
apiVersion: logging.banzaicloud.io/v1beta1
kind: SyslogNGOutput
metadata:
  name: elasticsearch
spec:
  elasticsearch:
    url: "https://elastic-search-endpoint:9200/_bulk"
    index: "indexname"
    type: ""
    user: "username"
    password:
      valueFrom:
        secretKeyRef:
          name: elastic
          key: password
{{</ highlight >}}
More information at https://axoflow.com/docs/axosyslog-core/chapter-destinations/configuring-destinations-elasticsearch-http/


## Configuration
## ElasticsearchOutput

###  (HTTPOutput, required) {#elasticsearchoutput-}

Default: -

### custom_id (string, optional) {#elasticsearchoutput-custom_id}

The document ID. If no ID is specified, a document ID is automatically generated. 

Default: -

### index (string, optional) {#elasticsearchoutput-index}

Name of the data stream, index, or index alias to perform the action on. 

Default: -

### logstash_suffix (string, optional) {#elasticsearchoutput-logstash_suffix}

Set the suffix for logs in logstash format. Default: "${YEAR}.${MONTH}.${DAY}" 

Default: -

### logstash_prefix (string, optional) {#elasticsearchoutput-logstash_prefix}

Set the prefix for logs in logstash format. If set, then Index field will be ignored. 

Default: -

### logstash_prefix_separator (string, optional) {#elasticsearchoutput-logstash_prefix_separator}

Set the separator between LogstashPrefix and LogStashDateformat. Default: "-" 

Default: -

### type (*string, optional) {#elasticsearchoutput-type}

The document type associated with the operation. Elasticsearch indices now support a single document type: _doc 

Default: -


