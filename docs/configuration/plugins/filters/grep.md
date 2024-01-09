---
title: Grep
weight: 200
generated_file: true
---

# [Grep Filter](https://docs.fluentd.org/filter/grep)
## Overview
 The grep filter plugin "greps" events by the values of specified fields.

## Configuration
## GrepConfig

### and ([]AndSection, optional) {#grepconfig-and}

[And Directive](#And-Directive) 


### exclude ([]ExcludeSection, optional) {#grepconfig-exclude}

[Exclude Directive](#Exclude-Directive) 


### or ([]OrSection, optional) {#grepconfig-or}

[Or Directive](#Or-Directive) 


### regexp ([]RegexpSection, optional) {#grepconfig-regexp}

[Regexp Directive](#Regexp-Directive) 



## Regexp Directive

Specify filtering rule (as described in the [Fluentd documentation](https://docs.fluentd.org/filter/grep#less-than-regexp-greater-than-directive)). This directive contains two parameters.

### key (string, required) {#regexp directive-key}

Specify field name in the record to parse. 


### pattern (string, required) {#regexp directive-pattern}

Pattern expression to evaluate 





## Example `Regexp` filter configurations

{{< highlight yaml >}}
apiVersion: logging.banzaicloud.io/v1beta1
kind: Flow
metadata:
  name: demo-flow
spec:
  filters:
    - grep:
        regexp:
        - key: first
          pattern: /^5\d\d$/
  selectors: {}
  localOutputRefs:
    - demo-output
{{</ highlight >}}

Fluentd config result:

{{< highlight xml >}}
  <filter **>
    @type grep
    @id demo-flow_1_grep
    <regexp>
      key first
      pattern /^5\d\d$/
    </regexp>
  </filter>
{{</ highlight >}}


---
## Exclude Directive

Specify filtering rule to reject events (as described in the [Fluentd documentation](https://docs.fluentd.org/filter/grep#less-than-exclude-greater-than-directive)). This directive contains two parameters.

### key (string, required) {#exclude directive-key}

Specify field name in the record to parse. 


### pattern (string, required) {#exclude directive-pattern}

Pattern expression to evaluate 





## Example `Exclude` filter configurations

{{< highlight yaml >}}
apiVersion: logging.banzaicloud.io/v1beta1
kind: Flow
metadata:
  name: demo-flow
spec:
  filters:
    - grep:
        exclude:
        - key: first
          pattern: /^5\d\d$/
  selectors: {}
  localOutputRefs:
    - demo-output
{{</ highlight >}}

Fluentd config result:

{{< highlight xml >}}
  <filter **>
    @type grep
    @id demo-flow_0_grep
    <exclude>
      key first
      pattern /^5\d\d$/
    </exclude>
  </filter>
{{</ highlight >}}


---
## Or Directive

Specify filtering rule (as described in the [Fluentd documentation](https://docs.fluentd.org/filter/grep#less-than-or-greater-than-directive). This directive contains either `regexp` or `exclude` directive.

### exclude ([]ExcludeSection, optional) {#or directive-exclude}

[Exclude Directive](#Exclude-Directive) 


### regexp ([]RegexpSection, optional) {#or directive-regexp}

[Regexp Directive](#Regexp-Directive) 





## Example `Or` filter configurations

{{< highlight yaml >}}
apiVersion: logging.banzaicloud.io/v1beta1
kind: Flow
metadata:
  name: demo-flow
spec:
  filters:
    - grep:
        or:
          - exclude:
            - key: first
              pattern: /^5\d\d$/
            - key: second
              pattern: /\.css$/

  selectors: {}
  localOutputRefs:
    - demo-output
{{</ highlight >}}

Fluentd config result:

{{< highlight xml >}}
<or>
	<exclude>
	key first
	pattern /^5\d\d$/
	</exclude>
	<exclude>
	key second
	pattern /\.css$/
	</exclude>
</or>
{{</ highlight >}}


---
## And Directive

Specify filtering rule (as described in the [Fluentd documentation](https://docs.fluentd.org/filter/grep#less-than-and-greater-than-directive). This directive contains either `regexp` or `exclude` directive.

### exclude ([]ExcludeSection, optional) {#and directive-exclude}

[Exclude Directive](#Exclude-Directive) 


### regexp ([]RegexpSection, optional) {#and directive-regexp}

[Regexp Directive](#Regexp-Directive) 





## Example `And` filter configurations

{{< highlight yaml >}}
apiVersion: logging.banzaicloud.io/v1beta1
kind: Flow
metadata:
  name: demo-flow
spec:
  filters:
    - grep:
        and:
          - regexp:
            - key: first
              pattern: /^5\d\d$/
            - key: second
              pattern: /\.css$/

  selectors: {}
  localOutputRefs:
    - demo-output
{{</ highlight >}}

Fluentd config result:

{{< highlight xml >}}
	<and>
	  <regexp>
	    key first
	    pattern /^5\d\d$/
	  </regexp>
	  <regexp>
	    key second
	    pattern /\.css$/
	  </regexp>
	</and>
{{</ highlight >}}


---
