---
title: MongoDB
weight: 200
generated_file: true
---

# Sending messages from a local network to an MongoDB database
## Overview

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

More information at https://axoflow.com/docs/axosyslog-core/chapter-destinations/configuring-destinations-mongodb/


## Configuration
## MongoDB

### collection (string, required) {#mongodb-collection}

The name of the MongoDB collection where the log messages are stored (collections are similar to SQL tables). Note that the name of the collection must not start with a dollar sign ($), and that it may contain dot (.) characters. 

Default: -

### dir (string, optional) {#mongodb-dir}

Defines the folder where the disk-buffer files are stored. 

Default: -

### disk_buffer (*DiskBuffer, optional) {#mongodb-disk_buffer}

This option enables putting outgoing messages into the disk buffer of the destination to avoid message loss in case of a system failure on the destination side. For details, see the [Syslog-ng DiskBuffer options](../disk_buffer/).

Default: false

### uri (*secret.Secret, optional) {#mongodb-uri}

Connection string used for authentication. See [syslog-ng docs](https://axoflow.com/docs/axosyslog-core/chapter-destinations/configuring-destinations-mongodb/reference-destination-mongodb/#mongodb-option-uri) 

Default: -

### value_pairs (ValuePairs, optional) {#mongodb-value_pairs}

Creates structured name-value pairs from the data and metadata of the log message.

Default: "scope("selected-macros" "nv-pairs")"

###  (Batch, required) {#mongodb-}

Batching parameters 

Default: -

###  (Bulk, required) {#mongodb-}

Bulk operation related options 

Default: -

### log-fifo-size (int, optional) {#mongodb-log-fifo-size}

The number of messages that the output queue can store. 

Default: -

### persist_name (string, optional) {#mongodb-persist_name}

If you receive the following error message during AxoSyslog startup, set the persist-name() option of the duplicate drivers: `Error checking the uniqueness of the persist names, please override it with persist-name option. Shutting down.` See [syslog-ng docs](https://axoflow.com/docs/axosyslog-core/chapter-destinations/configuring-destinations-http-nonjava/reference-destination-http-nonjava/#persist-name) for more information. 

Default: -

### retries (int, optional) {#mongodb-retries}

The number of times syslog-ng OSE attempts to send a message to this destination. If syslog-ng OSE could not send a message, it will try again until the number of attempts reaches retries, then drops the message. 

Default: -

### time_reopen (int, optional) {#mongodb-time_reopen}

The time to wait in seconds before a dead connection is reestablished.

Default: 60

### write_concern (RawString, optional) {#mongodb-write_concern}

Description: Sets the write concern mode of the MongoDB operations, for both bulk and single mode. See [syslog-ng docs] https://axoflow.com/docs/axosyslog-core/chapter-destinations/configuring-destinations-mongodb/reference-destination-mongodb/#mongodb-option-write-concern 

Default: -


## Bulk

Bulk operation related options
See [syslog-ng docs] https://axoflow.com/docs/axosyslog-core/chapter-destinations/configuring-destinations-mongodb/reference-destination-mongodb/#mongodb-option-bulk

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

### scope (RawString, optional) {#valuepairs-scope}

Default: -

### exclude (RawString, optional) {#valuepairs-exclude}

Default: -

### key (RawString, optional) {#valuepairs-key}

Default: -

### pair (RawString, optional) {#valuepairs-pair}

Default: -


## RawString

### raw_string (string, optional) {#rawstring-raw_string}

Default: -


