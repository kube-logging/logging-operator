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

### regexp ([]RegexpSection, optional) {#grepconfig-regexp}

[Regexp Directive](#Regexp-Directive) 

Default: -

### exclude ([]ExcludeSection, optional) {#grepconfig-exclude}

[Exclude Directive](#Exclude-Directive) 

Default: -

### or ([]OrSection, optional) {#grepconfig-or}

[Or Directive](#Or-Directive) 

Default: -

### and ([]AndSection, optional) {#grepconfig-and}

[And Directive](#And-Directive) 

Default: -


## Regexp Directive

Specify filtering rule (as described in the [Fluentd documentation](https://docs.fluentd.org/filter/grep#less-than-regexp-greater-than-directive)). This directive contains two parameters.

### key (string, required) {#regexp directive-key}

Specify field name in the record to parse. 

Default: -

### pattern (string, required) {#regexp directive-pattern}

Pattern expression to evaluate 

Default: -


 ## Example `Regexp` filter configurations
 ```yaml
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

 ```

 #### Fluentd Config Result
 ```yaml

	<filter **>
	  @type grep
	  @id demo-flow_1_grep
	  <regexp>
	    key first
	    pattern /^5\d\d$/
	  </regexp>
	</filter>

 ```

---
## Exclude Directive

Specify filtering rule to reject events (as described in the [Fluentd documentation](https://docs.fluentd.org/filter/grep#less-than-exclude-greater-than-directive)). This directive contains two parameters.

### key (string, required) {#exclude directive-key}

Specify field name in the record to parse. 

Default: -

### pattern (string, required) {#exclude directive-pattern}

Pattern expression to evaluate 

Default: -


 ## Example `Exclude` filter configurations
 ```yaml
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

 ```

 #### Fluentd Config Result
 ```yaml

	<filter **>
	  @type grep
	  @id demo-flow_0_grep
	  <exclude>
	    key first
	    pattern /^5\d\d$/
	  </exclude>
	</filter>

 ```

---
## Or Directive

Specify filtering rule (as described in the [Fluentd documentation](https://docs.fluentd.org/filter/grep#less-than-or-greater-than-directive). This directive contains either `regexp` or `exclude` directive.

### regexp ([]RegexpSection, optional) {#or directive-regexp}

[Regexp Directive](#Regexp-Directive) 

Default: -

### exclude ([]ExcludeSection, optional) {#or directive-exclude}

[Exclude Directive](#Exclude-Directive) 

Default: -


 ## Example `Or` filter configurations
 ```yaml
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

 ```

 #### Fluentd Config Result
 ```yaml

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

 ```

---
## And Directive

Specify filtering rule (as described in the [Fluentd documentation](https://docs.fluentd.org/filter/grep#less-than-and-greater-than-directive). This directive contains either `regexp` or `exclude` directive.

### regexp ([]RegexpSection, optional) {#and directive-regexp}

[Regexp Directive](#Regexp-Directive) 

Default: -

### exclude ([]ExcludeSection, optional) {#and directive-exclude}

[Exclude Directive](#Exclude-Directive) 

Default: -


 ## Example `And` filter configurations
 ```yaml
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

 ```

 #### Fluentd Config Result
 ```yaml

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

 ```

---
