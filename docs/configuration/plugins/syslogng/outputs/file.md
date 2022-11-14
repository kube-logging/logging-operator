---
title: File
weight: 200
generated_file: true
---

# File output plugin for syslog-ng
## Overview
 Storing messages in plain-text files
 More info at https://www.syslog-ng.com/technical-documents/doc/syslog-ng-open-source-edition/3.37/administration-guide/56#TOPIC-1829124

## Configuration
## FileOutput

Documentation: https://www.syslog-ng.com/technical-documents/doc/syslog-ng-open-source-edition/3.37/administration-guide/36#TOPIC-1829044

### path (string, required) {#fileoutput-path}

Store file path 

Default: -

### create_dirs (bool, optional) {#fileoutput-create_dirs}

Enable creating non-existing directories.  

Default:  false

### dir_group (string, optional) {#fileoutput-dir_group}

The group of the directories created by syslog-ng. To preserve the original properties of an existing directory, use the option without specifying an attribute: dir-group().  

Default:  Use the global settings

### dir_owner (string, optional) {#fileoutput-dir_owner}

The owner of the directories created by syslog-ng. To preserve the original properties of an existing directory, use the option without specifying an attribute: dir-owner().  

Default:  Use the global settings

### dir_perm (int, optional) {#fileoutput-dir_perm}

The permission mask of directories created by syslog-ng. Log directories are only created if a file after macro expansion refers to a non-existing directory, and directory creation is enabled (see also the create-dirs() option). For octal numbers prefix the number with 0, for example use 0755 for rwxr-xr-x. 

Default:  Use the global settings

### disk_buffer (*DiskBuffer, optional) {#fileoutput-disk_buffer}

This option enables putting outgoing messages into the disk buffer of the destination to avoid message loss in case of a system failure on the destination side.   

Default:  false

### template (string, optional) {#fileoutput-template}

Default: -

### persist_name (string, optional) {#fileoutput-persist_name}

Default: -


