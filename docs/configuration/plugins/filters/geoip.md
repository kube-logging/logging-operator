---
title: Geo IP
weight: 200
generated_file: true
---

# Fluentd GeoIP filter
## Overview
 Fluentd Filter plugin to add information about geographical location of IP addresses with Maxmind GeoIP databases.
 More information at https://github.com/y-ken/fluent-plugin-geoip

## Configuration
## GeoIP

### geoip_lookup_keys (string, optional) {#geoip-geoip_lookup_keys}

Specify one or more geoip lookup field which has ip address  

Default:  host

### geoip_database (string, optional) {#geoip-geoip_database}

Specify optional geoip database (using bundled GeoLiteCity databse by default) 

Default: -

### geoip2_database (string, optional) {#geoip-geoip2_database}

Specify optional geoip2 database (using bundled GeoLite2-City.mmdb by default) 

Default: -

### backend_library (string, optional) {#geoip-backend_library}

Specify backend library (geoip2_c, geoip, geoip2_compat) 

Default: -

### skip_adding_null_record (*bool, optional) {#geoip-skip_adding_null_record}

To avoid get stacktrace error with `[null, null]` array for elasticsearch. 

Default: true

### records ([]Record, optional) {#geoip-records}

Records are represented as maps: `key: value` 

Default: -


 ## Example `GeoIP` filter configurations
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
	localOutputRefs:
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
