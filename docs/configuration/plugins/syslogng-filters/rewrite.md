---
title: Rewrite
weight: 200
generated_file: true
---

# [Rewrite](https://www.syslog-ng.com/technical-documents/doc/syslog-ng-open-source-edition/3.37/administration-guide/76#TOPIC-1829205)
## Overview
 Rewrite filters can be used to modify record contents. Logging operator currently supports the following rewrite functions:

 - [group_unset](#groupunset)
 - [rename](#rename)
 - [set](#set)
 - [substitute](#subst)
 - [unset](#unset)

 > Note: All rewrite functions support an optional `condition` which has the same syntax as the [match filter](../match/).

 ## Group unset {#groupunset}

 The `group_unset` function removes from the record a group of fields matching a pattern.

 {{< highlight yaml >}}

	filters:
	- rewrite:
	  - group_unset:
	      pattern: "json.kubernetes.annotations.*"

 {{</ highlight >}}

 ## Rename

 The `rename` function changes the name of an existing field name.

 {{< highlight yaml >}}

	filters:
	- rewrite:
	  - rename:
	      oldName: "json.kubernetes.labels.app"
	      newName: "json.kubernetes.labels.app.kubernetes.io/name"

 {{</ highlight >}}

 ## Set

 The `set` function sets the value of a field.

 {{< highlight yaml >}}

	filters:
	- rewrite:
	  - set:
	      field: "json.kubernetes.cluster"
	      value: "prod-us"

 {{</ highlight >}}

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

 ## Unset

 You can unset macros or fields of the message.

 > Note: Unsetting a field completely deletes any previous value of the field.

 {{< highlight yaml >}}

	filters:
	- rewrite:
	  - unset:
	      field: "json.kubernetes.cluster"

 {{</ highlight >}}

## Configuration
## RewriteConfig

### group_unset (*GroupUnsetConfig, optional) {#rewriteconfig-group_unset}

Default: -

### rename (*RenameConfig, optional) {#rewriteconfig-rename}

Default: -

### set (*SetConfig, optional) {#rewriteconfig-set}

Default: -

### subst (*SubstituteConfig, optional) {#rewriteconfig-subst}

Default: -

### unset (*UnsetConfig, optional) {#rewriteconfig-unset}

Default: -


## RenameConfig

https://www.syslog-ng.com/technical-documents/doc/syslog-ng-open-source-edition/3.37/administration-guide/78#TOPIC-1829213

### condition (*MatchExpr, optional) {#renameconfig-condition}

Default: -

### newName (string, required) {#renameconfig-newname}

Default: -

### oldName (string, required) {#renameconfig-oldname}

Default: -


## SetConfig

https://www.syslog-ng.com/technical-documents/doc/syslog-ng-open-source-edition/3.37/administration-guide/77#TOPIC-1829207

### condition (*MatchExpr, optional) {#setconfig-condition}

Default: -

### field (string, required) {#setconfig-field}

Default: -

### value (string, required) {#setconfig-value}

Default: -


## SubstituteConfig

https://www.syslog-ng.com/technical-documents/doc/syslog-ng-open-source-edition/3.37/administration-guide/77#TOPIC-1829206

### condition (*MatchExpr, optional) {#substituteconfig-condition}

Default: -

### field (string, required) {#substituteconfig-field}

Default: -

### flags ([]string, optional) {#substituteconfig-flags}

Default: -

### pattern (string, required) {#substituteconfig-pattern}

Default: -

### replace (string, required) {#substituteconfig-replace}

Default: -

### type (string, optional) {#substituteconfig-type}

Default: -


## UnsetConfig

https://www.syslog-ng.com/technical-documents/doc/syslog-ng-open-source-edition/3.37/administration-guide/78#TOPIC-1829212

### condition (*MatchExpr, optional) {#unsetconfig-condition}

Default: -

### field (string, required) {#unsetconfig-field}

Default: -


## GroupUnsetConfig

https://www.syslog-ng.com/technical-documents/doc/syslog-ng-open-source-edition/3.37/administration-guide/78#TOPIC-1829212

### condition (*MatchExpr, optional) {#groupunsetconfig-condition}

Default: -

### pattern (string, required) {#groupunsetconfig-pattern}

Default: -


