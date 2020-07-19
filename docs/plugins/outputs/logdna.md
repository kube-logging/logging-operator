# [LogDNA Output](https://github.com/logdna/fluent-plugin-logdna)
## Overview
 This plugin has been designed to output logs to LogDNA. Example Deployment: [Transport Nginx Access Logs into LogDNA with Logging Operator](../../examples/logging_output_logdna.yaml)

## Configuration
### LogDNA
#### Send your logs to LogDNA

| Variable Name | Type | Required | Default | Description |
|---|---|---|---|---|
| api_key | string | Yes | - | LogDNA Api key<br> |
| hostname | string | Yes | - | Hostname<br> |
| app | string | No | - | Application name<br> |
| buffer_chunk_limit | string | No | - | Do not increase past 8m (8MB) or your logs will be rejected by LogDNA server.<br> |
 #### Example `Regexp` filter configurations
 ```
 apiVersion: logging.banzaicloud.io/v1beta1
 kind: Output
 metadata:
   name: logdna-output-sample
 spec:
   logdna:
     api_key: xxxxxxxxxxxxxxxxxxxxxxxxxxx
     hostname: logging-operator
     app: my-app
 ```

 #### Fluentd Config Result
 ```
<match **>
	@type logdna
	@id test_logdna
	api_key xxxxxxxxxxxxxxxxxxxxxxxxxxy
	app my-app
	hostname logging-operator
</match>
 ```

---
