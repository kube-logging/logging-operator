---
title: Mattermost
weight: 200
generated_file: true
---

# Mattermost plugin for Fluentd
## Overview
 Sends logs to Mattermost via webhooks.
 More info at https://github.com/levigo-systems/fluent-plugin-mattermost

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

### webhook_url (*secret.Secret, required) {#output config-webhook_url}

webhook_url Incoming Webhook URI (Required for Incoming Webhook mode). 

Default: -

### channel_id (string, optional) {#output config-channel_id}

channel_id the id of the channel where you want to receive the information. 

Default: -

### message_color (string, optional) {#output config-message_color}

message_color color of the message you are sending, the format is hex code.  

Default:  #A9A9A9

### message_title (string, optional) {#output config-message_title}

message_title title you want to add to the message.  

Default:  fluent_title_default

### message (string, optional) {#output config-message}

message The message you want to send, can be a static message, which you add at this point, or you can receive the fluent infos with the %s 

Default: -

### enable_tls (*bool, optional) {#output config-enable_tls}

enable_tls you can set the communication channel if it uses tls.  

Default:  true

### ca_path (*secret.Secret, optional) {#output config-ca_path}

ca_path you can set the path of the certificates. 

Default: -


