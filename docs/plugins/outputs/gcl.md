# Google Cloud Logging for Fluentd
## Overview
  More info at https://cloud.google.com/logging/docs/agent/configuration#cloud-fluentd-config
>Example: [Google Cloud Logging Output Deployment](../../../docs/example-gcl.md)

 #### Example output configurations
 ```
 spec:
  googleClouds:
    num_threads: 8
    use_grpc: true
    partial_success: true
    autoformat_stackdriver_trace: true
    buffer:
      timekey: 10m
      timekey_wait: 30s
      timekey_use_utc: true*/
 ```

## Configuration
### Output Config
| Variable Name | Type | Required | Default | Description |
|---|---|---|---|---|
| num_threads | int | No | - | The number of simultaneous log flushes that can be processed by the output plugin.<br> |
| use_grpc | bool | No |  true | Whether to use gRPC instead of REST/JSON to communicate to the Logging API. With gRPC enabled, CPU usage will typically be lower. <br> |
| partial_success | bool | No |  true | Whether to support partial success for logs ingestion. If true, invalid log entries in a full set are dropped, and valid log entries are successfully ingested into the Logging API. If false, the full set would be dropped if it contained any invalid log entries.  <br> |
| autoformat_stackdriver_trace | bool | No |  true | When set to true, the trace will be reformatted if the value of structured payload field logging.googleapis.com/trace matches ResourceTrace traceId format. Details of the autoformatting can be found under Special fields in structured payloads. <br> |
| buffer | *Buffer | No | - | [Buffer](./buffer.md)<br> |
| format | *Format | No | - | [Format](./format.md)<br> |
| auth_method | string | No | - | <br> |
| private_key_email | string | No | - | <br> |
| private_key_path | string | No | - | <br> |
| private_key_passphrase | string | No | - | <br> |
