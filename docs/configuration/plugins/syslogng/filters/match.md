---
title: Syslog-NG match
weight: 200
generated_file: true
---

# [Syslog-NG Match Filter](https://docs.fluentd.org/filter/grep)
## Overview
 The match filter can be used to selectively keep records

## Configuration
## MatchExpr

### and ([]MatchExpr, optional) {#matchexpr-and}

[And Directive](#And-Directive) 

Default: -

### not (*MatchExpr, optional) {#matchexpr-not}

[Not Directive](#Exclude-Directive) 

Default: -

### regexp (*RegexpMatchExpr, optional) {#matchexpr-regexp}

[Regexp Directive](#Regexp-Directive) 

Default: -

### or ([]MatchExpr, optional) {#matchexpr-or}

[Or Directive](#Or-Directive) 

Default: -


## [Regexp Directive](https://www.syslog-ng.com/technical-documents/doc/syslog-ng-open-source-edition/3.37/administration-guide/68#TOPIC-1829171) {#Regexp-Directive}

Specify filtering rule.

### pattern (string, required) {#[regexp directive](https://www.syslog-ng.com/technical-documents/doc/syslog-ng-open-source-edition/3.37/administration-guide/68#topic-1829171) {#regexp-directive}-pattern}

Pattern expression to evaluate 

Default: -

### template (string, optional) {#[regexp directive](https://www.syslog-ng.com/technical-documents/doc/syslog-ng-open-source-edition/3.37/administration-guide/68#topic-1829171) {#regexp-directive}-template}

Specify a template of the record fields to match against. 

Default: -

### value (string, optional) {#[regexp directive](https://www.syslog-ng.com/technical-documents/doc/syslog-ng-open-source-edition/3.37/administration-guide/68#topic-1829171) {#regexp-directive}-value}

Specify a field name of the record to match against the value of. 

Default: -

### flags ([]string, optional) {#[regexp directive](https://www.syslog-ng.com/technical-documents/doc/syslog-ng-open-source-edition/3.37/administration-guide/68#topic-1829171) {#regexp-directive}-flags}

Pattern flags 

Default: -

### type (string, optional) {#[regexp directive](https://www.syslog-ng.com/technical-documents/doc/syslog-ng-open-source-edition/3.37/administration-guide/68#topic-1829171) {#regexp-directive}-type}

Pattern type 

Default: -


 #### Example `Regexp` filter configurations
 ```yaml
apiVersion: logging.banzaicloud.io/v1beta1
kind: Flow
metadata:
  name: demo-flow
spec:
  filters:
    - match:
        regexp:
        - value: first
          pattern: ^5\d\d$
  match: {}
  localOutputRefs:
    - demo-output
 ```

 #### Syslog-NG Config Result
 ```
 log {
    source(main_input);
    filter {
        match("^5\d\d$" value("first"));
    };
    destination(output_default_demo-output);
 };
 ```

---
