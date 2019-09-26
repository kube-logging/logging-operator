### Fluentd Filter plugin to de-dot field name for elasticsearch.
#### More info at https://github.com/lunardial/fluent-plugin-dedot_filter

| Variable Name | Type | Required | Default | Description |
|---|---|---|---|---|
| de_dot_nested | bool | No | true | Will cause the plugin to recurse through nested structures (hashes and arrays), and remove dots in those key-names too.<br> |
| de_dot_separator | string | No | _ | Separator <br> |
