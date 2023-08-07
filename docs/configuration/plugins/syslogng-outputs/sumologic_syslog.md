---
title: Sumo Logic Syslog
weight: 200
generated_file: true
---

# Storing messages in Sumo Logic over syslog
## Overview
 More info at https://www.syslog-ng.com/technical-documents/doc/syslog-ng-open-source-edition/3.37/administration-guide/56#TOPIC-1829122

## Configuration
## SumologicSyslogOutput

### port (int, optional) {#sumologicsyslogoutput-port}

This option sets the port number of the Sumo Logic server to connect to.  

Default:  6514

### deployment (string, optional) {#sumologicsyslogoutput-deployment}

This option specifies your Sumo Logic deployment.https://help.sumologic.com/APIs/General-API-Information/Sumo-Logic-Endpoints-by-Deployment-and-Firewall-Security   

Default:  empty

### tag (string, optional) {#sumologicsyslogoutput-tag}

This option specifies the list of tags to add as the tags fields of Sumo Logic messages. If not specified, syslog-ng OSE automatically adds the tags already assigned to the message. If you set the tag() option, only the tags you specify will be added to the messages.  

Default:  tag

### token (int, optional) {#sumologicsyslogoutput-token}

The Cloud Syslog Cloud Token that you received from the Sumo Logic service while configuring your cloud syslog source. https://help.sumologic.com/03Send-Data/Sources/02Sources-for-Hosted-Collectors/Cloud-Syslog-Source#configure-a-cloud%C2%A0syslog%C2%A0source 

Default: -

### tls (*TLS, optional) {#sumologicsyslogoutput-tls}

This option sets various options related to TLS encryption, for example, key/certificate files and trusted CA locations. TLS can be used only with tcp-based transport protocols. For details, see [TLS for syslog-ng outputs](../tls/) and the [syslog-ng documentation](https://www.syslog-ng.com/technical-documents/doc/syslog-ng-open-source-edition/3.37/administration-guide/73#TOPIC-1829193). 

Default: -

### disk_buffer (*DiskBuffer, optional) {#sumologicsyslogoutput-disk_buffer}

This option enables putting outgoing messages into the disk buffer of the destination to avoid message loss in case of a system failure on the destination side. For details, see the [Syslog-ng DiskBuffer options](../disk_buffer/).  

Default:  false

### persist_name (string, optional) {#sumologicsyslogoutput-persist_name}

Default: -


