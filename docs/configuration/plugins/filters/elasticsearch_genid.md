## ElasticsearchGenId

### hash_id_key (string, optional) {#elasticsearchgenid-hash_id_key}

You can specify generated hash storing key. 

Default: -

### hash_type (string, optional) {#elasticsearchgenid-hash_type}

You can specify hash algorithm. Support algorithms md5, sha1, sha256, sha512. Default: sha1 

Default: -

### include_tag_in_seed (bool, optional) {#elasticsearchgenid-include_tag_in_seed}

You can specify to use tag for hash generation seed. 

Default: -

### include_time_in_seed (bool, optional) {#elasticsearchgenid-include_time_in_seed}

You can specify to use time for hash generation seed. 

Default: -

### record_keys (string, optional) {#elasticsearchgenid-record_keys}

You can specify keys which are record in events for hash generation seed. This parameter should be used with use_record_as_seed parameter in practice. 

Default: -

### separator (string, optional) {#elasticsearchgenid-separator}

You can specify separator charactor to creating seed for hash generation. 

Default: -

### use_entire_record (bool, optional) {#elasticsearchgenid-use_entire_record}

You can specify to use entire record in events for hash generation seed. 

Default: -

### use_record_as_seed (bool, optional) {#elasticsearchgenid-use_record_as_seed}

You can specify to use record in events for hash generation seed. This parameter should be used with record_keys parameter in practice. 

Default: -


