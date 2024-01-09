---
title: Match
weight: 200
generated_file: true
---

# Match
## Overview

Match filters can be used to select the log records to process. These filters have the same options and syntax as [syslog-ng flow match expressions]({{< relref "/docs/configuration/plugins/syslogng-filters/match.md" >}}).

{{< highlight yaml >}}
filters:
- match:
    or:
    - regexp:
        value: json.kubernetes.labels.app.kubernetes.io/name
        pattern: apache
        type: string
    - regexp:
        value: json.kubernetes.labels.app.kubernetes.io/name
        pattern: nginx
        type: string
{{</ highlight >}}


## Configuration
## MatchExpr

### and ([]MatchExpr, optional) {#matchexpr-and}


### not (*MatchExpr, optional) {#matchexpr-not}


### or ([]MatchExpr, optional) {#matchexpr-or}


### regexp (*RegexpMatchExpr, optional) {#matchexpr-regexp}

[Regexp Directive](#Regexp-Directive) 



## Regexp Directive


Specify filtering rule. For details, see the [AxoSyslog Core documentation](https://axoflow.com/docs/axosyslog-core/chapter-manipulating-messages/customizing-message-format/reference-template-functions/#template-function-list)


### flags ([]string, optional) {#regexp directive-flags}

Pattern flags. For details, see the [AxoSyslog Core documentation](https://axoflow.com/docs/axosyslog-core/chapter-manipulating-messages/regular-expressions/reference-regexp-types/regexp-flags-options/) 


### pattern (string, required) {#regexp directive-pattern}

Pattern expression to evaluate 


### template (string, optional) {#regexp directive-template}

Specify a template of the record fields to match against. 


### type (string, optional) {#regexp directive-type}

Pattern type. For details, see the [AxoSyslog Core documentation](https://axoflow.com/docs/axosyslog-core/chapter-manipulating-messages/regular-expressions/reference-regexp-types/regexp-type-options/) 


### value (string, optional) {#regexp directive-value}

Specify a field name of the record to match against the value of. 





## Example `Regexp` filter configurations

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
syslog-ng config result:

```shell
log {
    source(main_input);
    filter {
      match("^5\d\d$" value("first"));
    };
    destination(output_default_demo-output);
};
```


---
