---
title: MQTT
weight: 200
generated_file: true
---

# Sending messages from a local network to an MQTT broker
## Overview

Sends messages from a local network to an MQTT broker. For details on the available options of the output, see the [documentation of the AxoSyslog syslog-ng distribution](https://axoflow.com/docs/axosyslog-core/chapter-destinations/destination-mqtt-intro/).

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
    topic: test/demo
{{</ highlight >}}


## Configuration
## MQTT

### address (string, optional) {#mqtt-address}

Address of the destination host 


### fallback-topic (string, optional) {#mqtt-fallback-topic}

fallback-topic is used when syslog-ng cannot post a message to the originally defined topic (which can include invalid characters coming from templates). 


### qos (int, optional) {#mqtt-qos}

qos stands for quality of service and can take three values in the MQTT world. Its default value is 0, where there is no guarantee that the message is ever delivered. 


### template (string, optional) {#mqtt-template}

Template where you can configure the message template sent to the MQTT broker. By default, the template is: `$ISODATE $HOST $MSGHDR$MSG` 


### topic (string, optional) {#mqtt-topic}

Topic defines in which topic syslog-ng stores the log message. You can also use templates here, and use, for example, the $HOST macro in the topic name hierarchy. 



