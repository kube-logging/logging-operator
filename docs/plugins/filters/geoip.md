# Fluentd GeoIP filter
## Overview
 Fluentd Filter plugin to add information about geographical location of IP addresses with Maxmind GeoIP databases.
 More information at https://github.com/y-ken/fluent-plugin-geoip

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
 #### Example `GeoIP` filter configurations
 ```yaml
apiVersion: logging.banzaicloud.io/v1beta1
kind: Flow
metadata:
  name: demo-flow
spec:
  filters:
    - geoip:
        geoip_lookup_keys: remote_addr
        records:
          - city: ${city.names.en["remote_addr"]}
            location_array: '''[${location.longitude["remote"]},${location.latitude["remote"]}]'''
            country: ${country.iso_code["remote_addr"]}
            country_name: ${country.names.en["remote_addr"]}
            postal_code:  ${postal.code["remote_addr"]}
  selectors: {}
  outputRefs:
    - demo-output
 ```

 #### Fluentd Config Result
 ```yaml
<filter **>
  @type geoip
  @id test_geoip
  geoip_lookup_keys remote_addr
  skip_adding_null_record true
  <record>
    city ${city.names.en["remote_addr"]}
    country ${country.iso_code["remote_addr"]}
    country_name ${country.names.en["remote_addr"]}
    location_array '[${location.longitude["remote"]},${location.latitude["remote"]}]'
    postal_code ${postal.code["remote_addr"]}
  </record>
</filter>
 ```

---
