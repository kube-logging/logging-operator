---
title: MongoDB
weight: 200
generated_file: true
---

# Sending messages from a local network to an MongoDB database
## Overview

Based on the [MongoDB destination of the AxoSyslog syslog-ng distribution](https://axoflow.com/docs/axosyslog-core/chapter-destinations/configuring-destinations-mongodb/).

Available in Logging operator version 4.4 and later.

## Example

{{< highlight yaml >}}
apiVersion: logging.banzaicloud.io/v1beta1
kind: SyslogNGOutput
metadata:
  name: mongodb
  namespace: default
spec:
  mongodb:
    collection: syslog
    uri: "mongodb://mongodb-endpoint/syslog?wtimeoutMS=60000&socketTimeoutMS=60000&connectTimeoutMS=60000"
    value_pairs: scope("selected-macros" "nv-pairs")
{{</ highlight >}}

For more information, see the [documentation of the AxoSyslog syslog-ng distribution](https://axoflow.com/docs/axosyslog-core/chapter-destinations/configuring-destinations-mongodb/).


## Configuration
## MongoDB

###  (Batch, required) {#mongodb-}

Batching parameters 


###  (Bulk, required) {#mongodb-}

Bulk operation related options 


### collection (string, required) {#mongodb-collection}

The name of the MongoDB collection where the log messages are stored (collections are similar to SQL tables). Note that the name of the collection must not start with a dollar sign ($), and that it may contain dot (.) characters. 


### dir (string, optional) {#mongodb-dir}

Defines the folder where the disk-buffer files are stored. 


### disk_buffer (*DiskBuffer, optional) {#mongodb-disk_buffer}

This option enables putting outgoing messages into the disk buffer of the destination to avoid message loss in case of a system failure on the destination side. For details, see the [Syslog-ng DiskBuffer options](../disk_buffer/).

Default: false

### log-fifo-size (int, optional) {#mongodb-log-fifo-size}

The number of messages that the output queue can store. 


### persist_name (string, optional) {#mongodb-persist_name}

If you receive the following error message during syslog-ng startup, set the `persist-name()` option of the duplicate drivers: `Error checking the uniqueness of the persist names, please override it with persist-name option. Shutting down.` See the [documentation of the AxoSyslog syslog-ng distribution](https://axoflow.com/docs/axosyslog-core/chapter-destinations/configuring-destinations-http-nonjava/reference-destination-http-nonjava/#persist-name) for more information. 


### retries (int, optional) {#mongodb-retries}

The number of times syslog-ng OSE attempts to send a message to this destination. If syslog-ng OSE could not send a message, it will try again until the number of attempts reaches retries, then drops the message. 


### time_reopen (int, optional) {#mongodb-time_reopen}

The time to wait in seconds before a dead connection is reestablished.

Default: 60

### uri (*secret.Secret, optional) {#mongodb-uri}

Connection string used for authentication. See the [documentation of the AxoSyslog syslog-ng distribution](https://axoflow.com/docs/axosyslog-core/chapter-destinations/configuring-destinations-mongodb/reference-destination-mongodb/#mongodb-option-uri) 


### value_pairs (ValuePairs, optional) {#mongodb-value_pairs}

Creates structured name-value pairs from the data and metadata of the log message.

Default: `"scope("selected-macros" "nv-pairs")"`

### write_concern (RawString, optional) {#mongodb-write_concern}

Description: Sets the write concern mode of the MongoDB operations, for both bulk and single mode. For details, see the [documentation of the AxoSyslog syslog-ng distribution](https://axoflow.com/docs/axosyslog-core/chapter-destinations/configuring-destinations-mongodb/reference-destination-mongodb/#mongodb-option-write-concern). 



## Bulk

Bulk operation related options.
For details, see the [documentation of the AxoSyslog syslog-ng distribution](https://axoflow.com/docs/axosyslog-core/chapter-destinations/configuring-destinations-mongodb/reference-destination-mongodb/#mongodb-option-bulk).

### bulk (*bool, optional) {#bulk-bulk}

Enables bulk insert mode. If disabled, each messages is inserted individually.

Default: yes

### bulk_bypass_validation (*bool, optional) {#bulk-bulk_bypass_validation}

If set to yes, it disables MongoDB bulk operations validation mode.

Default: no

### bulk_unordered (*bool, optional) {#bulk-bulk_unordered}

Description: Enables unordered bulk operations mode.

Default: no


## ValuePairs

TODO move this to a common module once it is used in more places

### exclude (RawString, optional) {#valuepairs-exclude}


### key (RawString, optional) {#valuepairs-key}


### pair (RawString, optional) {#valuepairs-pair}


### scope (RawString, optional) {#valuepairs-scope}



