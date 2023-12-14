---
title: ElasticSearch GenId
weight: 200
generated_file: true
---

# ElasticSearch GenId
## Overview

## Example `Elasticsearch Genid` filter configurations

{{< highlight yaml >}}
apiVersion: logging.banzaicloud.io/v1beta1
kind: Flow
metadata:
 name: demo-flow
spec:
 filters:
   - elasticsearch_genid:
       hash_id_key: gen_id
 selectors: {}
 localOutputRefs:
   - demo-output
{{</ highlight >}}

Fluentd Config Result

{{< highlight xml >}}
<filter **>
 @type elasticsearch_genid
 @id test_elasticsearch_genid
 hash_id_key gen_id
</filter>
{{</ highlight >}}


## Configuration
## ElasticsearchGenId

### hash_id_key (string, optional) {#elasticsearchgenid-hash_id_key}

You can specify generated hash storing key. 


### hash_type (string, optional) {#elasticsearchgenid-hash_type}

You can specify hash algorithm. Support algorithms md5, sha1, sha256, sha512. Default: sha1 


### include_tag_in_seed (bool, optional) {#elasticsearchgenid-include_tag_in_seed}

You can specify to use tag for hash generation seed. 


### include_time_in_seed (bool, optional) {#elasticsearchgenid-include_time_in_seed}

You can specify to use time for hash generation seed. 


### record_keys (string, optional) {#elasticsearchgenid-record_keys}

You can specify keys which are record in events for hash generation seed. This parameter should be used with use_record_as_seed parameter in practice. 


### separator (string, optional) {#elasticsearchgenid-separator}

You can specify separator charactor to creating seed for hash generation. 


### use_entire_record (bool, optional) {#elasticsearchgenid-use_entire_record}

You can specify to use entire record in events for hash generation seed. 


### use_record_as_seed (bool, optional) {#elasticsearchgenid-use_record_as_seed}

You can specify to use record in events for hash generation seed. This parameter should be used with record_keys parameter in practice. 



