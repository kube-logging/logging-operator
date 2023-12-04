---
title: Mattermost
weight: 200
generated_file: true
---

# Mattermost plugin for Fluentd
## Overview

Sends logs to Mattermost via webhooks.
For details, see [https://github.com/levigo-systems/fluent-plugin-mattermost](https://github.com/levigo-systems/fluent-plugin-mattermost).

## Example output configurations

```yaml
spec:
  mattermost:
    webhook_url: https://xxx.xx/hooks/xxxxxxxxxxxxxxx
    channel_id: xxxxxxxxxxxxxxx
    message_color: "#FFA500"
    enable_tls: false
```


## Configuration
## Output Config

### ca_path (*secret.Secret, optional) {#output config-ca_path}

The path of the CA certificates. 


### channel_id (string, optional) {#output config-channel_id}

The ID of the channel where you want to receive the information. 


### enable_tls (*bool, optional) {#output config-enable_tls}

You can set the communication channel if it uses TLS.

Default: true

### message (string, optional) {#output config-message}

The message you want to send. It can be a static message, which you add at this point, or you can receive the Fluentd infos with the %s 


### message_color (string, optional) {#output config-message_color}

Color of the message you are sending, in hexadecimal format.

Default: #A9A9A9

### message_title (string, optional) {#output config-message_title}

The title you want to add to the message.

Default: fluent_title_default

### webhook_url (*secret.Secret, required) {#output config-webhook_url}

Incoming Webhook URI (Required for Incoming Webhook mode). 



