---
title: MQTT
weight: 200
generated_file: true
---

# Sending messages from a local network to an MQTT broker
## Overview

 ## Prerequisites

 ## Example

 {{< highlight yaml >}}
 apiVersion: logging.banzaicloud.io/v1beta1
 kind: SyslogNGOutput
 metadata:

	name: mqtt
	namespace: default

 spec:

	mqtt:
	  address: tcp://mosquitto:1883
	  template: |
	    $(format-json --subkeys json~ --key-delimiter ~)
	  topic: test/demo

 {{</ highlight >}}

## Configuration
## MQTT

### address (string, optional) {#mqtt-address}

Address of the destination host 

Default: -

### topic (string, optional) {#mqtt-topic}

Topic defines in which topic syslog-ng stores the log message. You can also use templates here, and use, for example, the $HOST macro in the topic name hierarchy. 

Default: -

### fallback-topic (string, optional) {#mqtt-fallback-topic}

fallback-topic is used when syslog-ng cannot post a message to the originally defined topic (which can include invalid characters coming from templates). 

Default: -

### template (string, optional) {#mqtt-template}

Template where you can configure the message template sent to the MQTT broker. By default, the template is: “$ISODATE $HOST $MSGHDR$MSG” 

Default: -

### qos (int, optional) {#mqtt-qos}

qos stands for quality of service and can take three values in the MQTT world. Its default value is 0, where there is no guarantee that the message is ever delivered. 

Default: -


