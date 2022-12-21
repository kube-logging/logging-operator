---
title: Syslog-NG rewrite
weight: 200
generated_file: true
---

# [Syslog-NG Rewrite Filter](https://www.syslog-ng.com/technical-documents/doc/syslog-ng-open-source-edition/3.37/administration-guide/77)
## Overview
 The syslog-ng rewrite filter can be used to replace message parts.

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

### oldName (string, required) {#renameconfig-oldname}

Default: -

### newName (string, required) {#renameconfig-newname}

Default: -

### condition (*MatchExpr, optional) {#renameconfig-condition}

Default: -


## SetConfig

https://www.syslog-ng.com/technical-documents/doc/syslog-ng-open-source-edition/3.37/administration-guide/77#TOPIC-1829207

### field (string, required) {#setconfig-field}

Default: -

### value (string, required) {#setconfig-value}

Default: -

### condition (*MatchExpr, optional) {#setconfig-condition}

Default: -


## SubstituteConfig

https://www.syslog-ng.com/technical-documents/doc/syslog-ng-open-source-edition/3.37/administration-guide/77#TOPIC-1829206

### pattern (string, required) {#substituteconfig-pattern}

Default: -

### replace (string, required) {#substituteconfig-replace}

Default: -

### field (string, required) {#substituteconfig-field}

Default: -

### flags ([]string, optional) {#substituteconfig-flags}

Default: -

### type (string, optional) {#substituteconfig-type}

Default: -

### condition (*MatchExpr, optional) {#substituteconfig-condition}

Default: -


## UnsetConfig

https://www.syslog-ng.com/technical-documents/doc/syslog-ng-open-source-edition/3.37/administration-guide/78#TOPIC-1829212

### field (string, required) {#unsetconfig-field}

Default: -

### condition (*MatchExpr, optional) {#unsetconfig-condition}

Default: -


## GroupUnsetConfig

https://www.syslog-ng.com/technical-documents/doc/syslog-ng-open-source-edition/3.37/administration-guide/78#TOPIC-1829212

### pattern (string, required) {#groupunsetconfig-pattern}

Default: -

### condition (*MatchExpr, optional) {#groupunsetconfig-condition}

Default: -


