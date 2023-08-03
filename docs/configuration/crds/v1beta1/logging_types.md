---
title: LoggingSpec
weight: 200
generated_file: true
---

## LoggingSpec

LoggingSpec defines the desired state of Logging

### loggingRef (string, optional) {#loggingspec-loggingref}

Reference to the logging system. Each of the `loggingRef`s can manage a fluentbit daemonset and a fluentd statefulset. 

Default: -

### flowConfigCheckDisabled (bool, optional) {#loggingspec-flowconfigcheckdisabled}

Disable configuration check before applying new fluentd configuration. 

Default: -

### skipInvalidResources (bool, optional) {#loggingspec-skipinvalidresources}

Whether to skip invalid Flow and ClusterFlow resources 

Default: -

### flowConfigOverride (string, optional) {#loggingspec-flowconfigoverride}

Override generated config. This is a *raw* configuration string for troubleshooting purposes. 

Default: -

### fluentbit (*FluentbitSpec, optional) {#loggingspec-fluentbit}

FluentbitAgent daemonset configuration. Deprecated, will be removed with next major version Migrate to the standalone NodeAgent resource 

Default: -

### fluentd (*FluentdSpec, optional) {#loggingspec-fluentd}

Fluentd statefulset configuration 

Default: -

### syslogNG (*SyslogNGSpec, optional) {#loggingspec-syslogng}

Syslog-NG statefulset configuration 

Default: -

### defaultFlow (*DefaultFlowSpec, optional) {#loggingspec-defaultflow}

Default flow for unmatched logs. This Flow configuration collects all logs that didn't matched any other Flow. 

Default: -

### errorOutputRef (string, optional) {#loggingspec-erroroutputref}

GlobalOutput name to flush ERROR events to 

Default: -

### globalFilters ([]Filter, optional) {#loggingspec-globalfilters}

Global filters to apply on logs before any match or filter mechanism. 

Default: -

### watchNamespaces ([]string, optional) {#loggingspec-watchnamespaces}

Limit namespaces to watch Flow and Output custom resources. 

Default: -

### watchNamespaceSelector (*metav1.LabelSelector, optional) {#loggingspec-watchnamespaceselector}

WatchNamespaceSelector is a LabelSelector to find matching namespaces to watch as in WatchNamespaces 

Default: -

### clusterDomain (*string, optional) {#loggingspec-clusterdomain}

Cluster domain name to be used when templating URLs to services . 

Default:  "cluster.local"

### controlNamespace (string, required) {#loggingspec-controlnamespace}

Namespace for cluster wide configuration resources like CLusterFlow and ClusterOutput. This should be a protected namespace from regular users. Resources like fluentbit and fluentd will run in this namespace as well. 

Default: -

### allowClusterResourcesFromAllNamespaces (bool, optional) {#loggingspec-allowclusterresourcesfromallnamespaces}

Allow configuration of cluster resources from any namespace. Mutually exclusive with ControlNamespace restriction of Cluster resources 

Default: -

### nodeAgents ([]*InlineNodeAgent, optional) {#loggingspec-nodeagents}

InlineNodeAgent Configuration Deprecated, will be removed with next major version 

Default: -

### enableRecreateWorkloadOnImmutableFieldChange (bool, optional) {#loggingspec-enablerecreateworkloadonimmutablefieldchange}

EnableRecreateWorkloadOnImmutableFieldChange enables the operator to recreate the fluentbit daemonset and the fluentd statefulset (and possibly other resource in the future) in case there is a change in an immutable field that otherwise couldn't be managed with a simple update. 

Default: -


## LoggingStatus

LoggingStatus defines the observed state of Logging

### configCheckResults (map[string]bool, optional) {#loggingstatus-configcheckresults}

Default: -

### problems ([]string, optional) {#loggingstatus-problems}

Default: -


## Logging

Logging is the Schema for the loggings API

###  (metav1.TypeMeta, required) {#logging-}

Default: -

### metadata (metav1.ObjectMeta, optional) {#logging-metadata}

Default: -

### spec (LoggingSpec, optional) {#logging-spec}

Default: -

### status (LoggingStatus, optional) {#logging-status}

Default: -


## LoggingList

LoggingList contains a list of Logging

###  (metav1.TypeMeta, required) {#logginglist-}

Default: -

### metadata (metav1.ListMeta, optional) {#logginglist-metadata}

Default: -

### items ([]Logging, required) {#logginglist-items}

Default: -


## DefaultFlowSpec

DefaultFlowSpec is a Flow for logs that did not match any other Flow

### filters ([]Filter, optional) {#defaultflowspec-filters}

Default: -

### outputRefs ([]string, optional) {#defaultflowspec-outputrefs}

Deprecated 

Default: -

### globalOutputRefs ([]string, optional) {#defaultflowspec-globaloutputrefs}

Default: -

### flowLabel (string, optional) {#defaultflowspec-flowlabel}

Default: -

### includeLabelInRouter (*bool, optional) {#defaultflowspec-includelabelinrouter}

Default: -


