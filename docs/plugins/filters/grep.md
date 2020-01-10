# [Grep Filter](https://docs.fluentd.org/filter/grep)
## Overview
 The grep filter plugin "greps" events by the values of specified fields.

## Configuration
### GrepConfig
| Variable Name | Type | Required | Default | Description |
|---|---|---|---|---|
| regexp | []RegexpSection | No | - | [Regexp Directive](#Regexp-Directive)<br> |
| exclude | []ExcludeSection | No | - | [Exclude Directive](#Exclude-Directive)<br> |
| or | []OrSection | No | - | [Or Directive](#Or-Directive)<br> |
| and | []AndSection | No | - | [And Directive](#And-Directive)<br> |
### [Regexp Directive](https://docs.fluentd.org/filter/grep#less-than-regexp-greater-than-directive)
#### Specify filtering rule. This directive contains two parameters.

| Variable Name | Type | Required | Default | Description |
|---|---|---|---|---|
| key | string | Yes | - | Specify field name in the record to parse.<br> |
| pattern | string | Yes | - | Pattern expression to evaluate<br> |
 #### Example `Regexp` filter configurations
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
  outputRefs:
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
### [Exclude Directive](https://docs.fluentd.org/filter/grep#less-than-exclude-greater-than-directive)
#### Specify filtering rule to reject events. This directive contains two parameters.

| Variable Name | Type | Required | Default | Description |
|---|---|---|---|---|
| key | string | Yes | - | Specify field name in the record to parse.<br> |
| pattern | string | Yes | - | Pattern expression to evaluate<br> |
 #### Example `Exclude` filter configurations
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
  outputRefs:
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
### [Or Directive](https://docs.fluentd.org/filter/grep#less-than-or-greater-than-directive)
#### Specify filtering rule. This directive contains either `regexp` or `exclude` directive.

| Variable Name | Type | Required | Default | Description |
|---|---|---|---|---|
| regexp | []RegexpSection | No | - | [Regexp Directive](#Regexp-Directive)<br> |
| exclude | []ExcludeSection | No | - | [Exclude Directive](#Exclude-Directive)<br> |
 #### Example `Or` filter configurations
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
  outputRefs:
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
### [And Directive](https://docs.fluentd.org/filter/grep#less-than-and-greater-than-directive)
#### Specify filtering rule. This directive contains either `regexp` or `exclude` directive.

| Variable Name | Type | Required | Default | Description |
|---|---|---|---|---|
| regexp | []RegexpSection | No | - | [Regexp Directive](#Regexp-Directive)<br> |
| exclude | []ExcludeSection | No | - | [Exclude Directive](#Exclude-Directive)<br> |
 #### Example `And` filter configurations
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
  outputRefs:
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
