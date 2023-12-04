---
title: Rewrite
weight: 200
generated_file: true
---

# [Rewrite](https://axoflow.com/docs/axosyslog-core/chapter-manipulating-messages/modifying-messages/)
## Overview

Rewrite filters can be used to modify record contents. Logging operator currently supports the following rewrite functions:

- [group_unset](#groupunset)
- [rename](#rename)
- [set](#set)
- [substitute](#subst)
- [unset](#unset)

> Note: All rewrite functions support an optional `condition` which has the same syntax as the [match filter](../match/).

For details on how rewrite rules work in syslog-ng, see the [documentation of the AxoSyslog syslog-ng distribution](https://axoflow.com/docs/axosyslog-core/chapter-manipulating-messages/modifying-messages/).

## Group unset {#groupunset}

The `group_unset` function removes from the record a group of fields matching a pattern.

{{< highlight yaml >}}
  filters:
  - rewrite:
    - group_unset:
        pattern: "json.kubernetes.annotations.*"
{{</ highlight >}}

For details, see the [documentation of the AxoSyslog syslog-ng distribution](https://axoflow.com/docs/axosyslog-core/chapter-manipulating-messages/modifying-messages/rewrite-unset/).

## Rename

The `rename` function changes the name of an existing field name.

{{< highlight yaml >}}
  filters:
  - rewrite:
    - rename:
        oldName: "json.kubernetes.labels.app"
        newName: "json.kubernetes.labels.app.kubernetes.io/name"
{{</ highlight >}}

For details, see the [documentation of the AxoSyslog syslog-ng distribution](https://axoflow.com/docs/axosyslog-core/chapter-manipulating-messages/modifying-messages/rewrite-rename/).

## Set

The `set` function sets the value of a field.

{{< highlight yaml >}}
  filters:
  - rewrite:
    - set:
        field: "json.kubernetes.cluster"
        value: "prod-us"
{{</ highlight >}}

For details, see the [documentation of the AxoSyslog syslog-ng distribution](https://axoflow.com/docs/axosyslog-core/chapter-manipulating-messages/modifying-messages/rewrite-set/).

## Substitute (subst) {#subst}

The `subst` function replaces parts of a field with a replacement value based on a pattern.

{{< highlight yaml >}}
  filters:
  - rewrite:
    - subst:
        pattern: "\d\d\d\d-\d\d\d\d-\d\d\d\d-\d\d\d\d"
        replace: "[redacted bank card number]"
        field: "MESSAGE"
{{</ highlight >}}

The function also supports the `type` and `flags` fields for specifying pattern type and flags as described in the [match expression regexp function](../match/).

For details, see the [documentation of the AxoSyslog syslog-ng distribution](https://axoflow.com/docs/axosyslog-core/chapter-manipulating-messages/modifying-messages/rewrite-replace/).

## Unset

You can unset macros or fields of the message.

> Note: Unsetting a field completely deletes any previous value of the field.

{{< highlight yaml >}}
  filters:
  - rewrite:
    - unset:
        field: "json.kubernetes.cluster"
{{</ highlight >}}

For details, see the [documentation of the AxoSyslog syslog-ng distribution](https://axoflow.com/docs/axosyslog-core/chapter-manipulating-messages/modifying-messages/rewrite-unset/).


## Configuration
## RewriteConfig

### group_unset (*GroupUnsetConfig, optional) {#rewriteconfig-group_unset}


### rename (*RenameConfig, optional) {#rewriteconfig-rename}


### set (*SetConfig, optional) {#rewriteconfig-set}


### subst (*SubstituteConfig, optional) {#rewriteconfig-subst}


### unset (*UnsetConfig, optional) {#rewriteconfig-unset}



## RenameConfig

For details, see the [documentation of the AxoSyslog syslog-ng distribution](https://axoflow.com/docs/axosyslog-core/chapter-manipulating-messages/modifying-messages/rewrite-rename/).

### condition (*MatchExpr, optional) {#renameconfig-condition}


### newName (string, required) {#renameconfig-newname}


### oldName (string, required) {#renameconfig-oldname}



## SetConfig

For details, see the [documentation of the AxoSyslog syslog-ng distribution](https://axoflow.com/docs/axosyslog-core/chapter-manipulating-messages/modifying-messages/rewrite-set/).

### condition (*MatchExpr, optional) {#setconfig-condition}


### field (string, required) {#setconfig-field}


### value (string, required) {#setconfig-value}



## SubstituteConfig

For details, see the [documentation of the AxoSyslog syslog-ng distribution](https://axoflow.com/docs/axosyslog-core/chapter-manipulating-messages/modifying-messages/rewrite-set/).

### condition (*MatchExpr, optional) {#substituteconfig-condition}


### field (string, required) {#substituteconfig-field}


### flags ([]string, optional) {#substituteconfig-flags}


### pattern (string, required) {#substituteconfig-pattern}


### replace (string, required) {#substituteconfig-replace}


### type (string, optional) {#substituteconfig-type}



## UnsetConfig

For details, see the [documentation of the AxoSyslog syslog-ng distribution](https://axoflow.com/docs/axosyslog-core/chapter-manipulating-messages/modifying-messages/rewrite-unset/).

### condition (*MatchExpr, optional) {#unsetconfig-condition}


### field (string, required) {#unsetconfig-field}



## GroupUnsetConfig

For details, see the [documentation of the AxoSyslog syslog-ng distribution](https://axoflow.com/docs/axosyslog-core/chapter-manipulating-messages/modifying-messages/rewrite-unset/).

### condition (*MatchExpr, optional) {#groupunsetconfig-condition}


### pattern (string, required) {#groupunsetconfig-pattern}



