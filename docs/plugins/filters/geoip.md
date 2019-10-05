# Fluentd GeoIP filter
## Overview
Fluentd Filter plugin to add information about geographical location of IP addresses with Maxmind GeoIP databases.
More information at https://github.com/y-ken/fluent-plugin-geoip

#### Example record configurations
```
spec:
filters:
- tag_normaliser:
format: ${namespace_name}.${pod_name}.${container_name}
- parser:
key_name: message
remove_key_name_field: true
parsers:
- type: nginx
- geoip:
records:
- city: ${city.names.en["remote_addr"]}
latitude: ${location.latitude["remote_addr"]}
longitude: ${location.longitude["remote_addr"]}
country: ${country.iso_code["remote_addr"]}
country_name: ${country.names.en["remote_addr"]}
postal_code:  ${postal.code["remote_addr"]}
```

## Configuration
### GeoIP
| Variable Name | Type | Required | Default | Description |
|---|---|---|---|---|
| geoip_lookup_keys | string | No |  host | Specify one or more geoip lookup field which has ip address <br> |
| geoip_database | string | No | - | Specify optional geoip database (using bundled GeoLiteCity databse by default)<br> |
| geoip_2_database | string | No | - | Specify optional geoip2 database (using bundled GeoLite2-City.mmdb by default)<br> |
| backend_library | string | No | - | Specify backend library (geoip2_c, geoip, geoip2_compat)<br> |
| skip_adding_null_record | bool | No | true | To avoid get stacktrace error with `[null, null]` array for elasticsearch.<br> |
| records | []Record | No | - | Records are represented as maps: `key: value`<br> |
