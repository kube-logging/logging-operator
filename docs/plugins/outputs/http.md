---
title: Http
weight: 200
generated_file: true
---

# Http plugin for Fluentd
## Overview
 Sends logs to HTTP/HTTPS endpoints.
 More info at https://docs.fluentd.org/output/http.

 #### Example output configurations
 ```
 spec:
   http:
     endpoint: http://logserver.com:9000/api
     buffer:
       tags: "[]"
       flush_interval: 10s
 ```

## Configuration
### Output Config
| Variable Name | Type | Required | Default | Description |
|---|---|---|---|---|
| endpoint | string | Yes | - | Endpoint for HTTP request.<br> |
| http_method | string | No |  post | Method for HTTP request. [post, put] <br> |
| proxy | string | No | - | Proxy for HTTP request.<br> |
| content_type | string | No | - | Content-Type for HTTP request.<br> |
| json_array | bool | No |  false | Using array format of JSON. This parameter is used and valid only for json format. When json_array as true, Content-Type should be application/json and be able to use JSON data for the HTTP request body.  <br> |
| format | *Format | No | - | [Format](../format/)<br> |
| headers | map[string]string | No | - | Additional headers for HTTP request.<br> |
| open_timeout | int | No | - | Connection open timeout in seconds.<br> |
| read_timeout | int | No | - | Read timeout in seconds.<br> |
| ssl_timeout | int | No | - | TLS timeout in seconds.<br> |
| tls_version | string | No |  TLSv1_2 | The default version of TLS transport. [TLSv1_1, TLSv1_2] <br> |
| tls_ciphers | string | No |  ALL:!aNULL:!eNULL:!SSLv2 | The cipher configuration of TLS transport. <br> |
| tls_ca_cert_path | *secret.Secret | No | - | The CA certificate path for TLS.<br> |
| tls_client_cert_path | *secret.Secret | No | - | The client certificate path for TLS.<br> |
| tls_private_key_path | *secret.Secret | No | - | The client private key path for TLS.<br> |
| tls_private_key_passphrase | *secret.Secret | No | - | The client private key passphrase for TLS.<br> |
| tls_verify_mode | string | No |  peer | The verify mode of TLS. [peer, none] <br> |
| error_response_as_unrecoverable | *bool | No |  true | Raise UnrecoverableError when the response code is non success, 1xx/3xx/4xx/5xx. If false, the plugin logs error message instead of raising UnrecoverableError. <br> |
| retryable_response_codes | []int | No |  [503] | List of retryable response codes. If the response code is included in this list, the plugin retries the buffer flush. <br> |
| auth | *HTTPAuth | No | - | [HTTP auth](#http_auth)<br> |
| buffer | *Buffer | No | - | [Buffer](../buffer/)<br> |
### HTTP auth config
#### http_auth

| Variable Name | Type | Required | Default | Description |
|---|---|---|---|---|
| username | *secret.Secret | Yes | - | Username for basic authentication.<br>[Secret](../secret/)<br> |
| password | *secret.Secret | Yes | - | Password for basic authentication.<br>[Secret](../secret/)<br> |
