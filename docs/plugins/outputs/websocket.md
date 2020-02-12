# Websocket output plugin for Fluentd
## Overview
 This plugin works as websocket server which can output JSON string or MessagePack binary.

 More info and examples at https://github.com/banzaicloud/fluent-plugin-websocket

## Configuration
### Output Config
| Variable Name | Type | Required | Default | Description |
|---|---|---|---|---|
| host | string | No |  0.0.0.0 (ANY) | WebSocket server IP address. <br> |
| port | int | No |  8080 | WebSocket server port. <br> |
| use_msgpack | *bool | No |  false | Send MessagePack format binary. Otherwise, you send JSON format text. <br> |
| add_time | *bool | No |  false | Add timestamp to the data. <br> |
| add_tag | *bool | No |  true | Add fluentd tag to the data. <br> |
| buffered_messages | int | No |  0 | The number of messages to be buffered. The new connection receives them. <br> |
| token | *secret.Secret | No |  nil | Authentication token. Passed as get param. If set to nil, authentication is disabled. <br> |
