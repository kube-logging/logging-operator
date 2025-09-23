---
title: LoggingSpec
weight: 200
generated_file: true
---

## LoggingSpec

LoggingSpec defines the desired state of Logging

### allowClusterResourcesFromAllNamespaces (bool, optional) {#loggingspec-allowclusterresourcesfromallnamespaces}

Allow configuration of cluster resources from any namespace. Mutually exclusive with ControlNamespace restriction of Cluster resources WARNING: Becareful when turning this on and off as it can result in some resources being orphaned. 


### clusterDomain (*string, optional) {#loggingspec-clusterdomain}

Cluster domain name to be used when templating URLs to services .

Default: "cluster.local."

### configCheck (ConfigCheck, optional) {#loggingspec-configcheck}

ConfigCheck settings that apply to both fluentd or syslog-ng. Can be overridden on the fluentd / syslog-ng level. 


### controlNamespace (string, required) {#loggingspec-controlnamespace}

Namespace for cluster wide configuration resources like ClusterFlow and ClusterOutput. This should be a protected namespace from regular users. Resources like fluentbit and fluentd will run in this namespace as well. 


### defaultFlow (*DefaultFlowSpec, optional) {#loggingspec-defaultflow}

Default flow for unmatched logs. This Flow configuration collects all logs that didn't matched any other Flow. 


### enableDockerParserCompatibilityForCRI (bool, optional) {#loggingspec-enabledockerparsercompatibilityforcri}

EnableDockerParserCompatibilityForCRI enables a log parser that is compatible with the docker parser. This has the following benefits: - automatic json log parsing using the Merge_Log feature - downstream parsers can use the `log` field instead of `message` as they did with the docker runtime - the `concat` and `parser` filters are automatically set back to use the `log` field 


### enableRecreateWorkloadOnImmutableFieldChange (bool, optional) {#loggingspec-enablerecreateworkloadonimmutablefieldchange}

EnableRecreateWorkloadOnImmutableFieldChange enables the operator to recreate the fluentbit daemonset and the fluentd statefulset (and possibly other resource in the future) in case there is a change in an immutable field that otherwise couldn't be managed with a simple update. 


### errorOutputRef (string, optional) {#loggingspec-erroroutputref}

GlobalOutput name to flush ERROR events to 


### flowConfigCheckDisabled (bool, optional) {#loggingspec-flowconfigcheckdisabled}

Disable configuration check before applying new fluentd or syslogng configuration. 


### flowConfigOverride (string, optional) {#loggingspec-flowconfigoverride}

Override generated config. This is a *raw* configuration string for troubleshooting purposes. 


### fluentbit (*FluentbitSpec, optional) {#loggingspec-fluentbit}

FluentbitAgent daemonset configuration. DEPRECATED: Migrate to the standalone FluentBitAgent resource 


### fluentd (*FluentdSpec, optional) {#loggingspec-fluentd}

Fluentd statefulset configuration 


### globalFilters ([]Filter, optional) {#loggingspec-globalfilters}

Global filters to apply on logs before any match or filter mechanism. 


### loggingRef (string, optional) {#loggingspec-loggingref}

Reference to the logging system. Each of the `loggingRef`s can manage a fluentbit daemonset and a fluentd statefulset. 


### nodeAgentNamespace (string, required) {#loggingspec-nodeagentnamespace}

Namespace to deploy Fluent Bit resources (DaemonSet, Service, ServiceAccount, config Secret, ServiceMonitors). If unset, it defaults to `controlNamespace` to preserve backward compatibility. 


### routeConfig (*RouteConfig, optional) {#loggingspec-routeconfig}

RouteConfig determines whether to use loggingRoutes or to create resources based on the logging resource that can be managed by the Telemetry Controller. 


### skipInvalidResources (bool, optional) {#loggingspec-skipinvalidresources}

Whether to skip invalid Flow and ClusterFlow resources 


### syslogNG (*SyslogNGSpec, optional) {#loggingspec-syslogng}

Syslog-NG statefulset configuration 


### watchNamespaceSelector (*metav1.LabelSelector, optional) {#loggingspec-watchnamespaceselector}

WatchNamespaceSelector is a LabelSelector to find matching namespaces to watch as in WatchNamespaces 


### watchNamespaces ([]string, optional) {#loggingspec-watchnamespaces}

Limit namespaces to watch Flow and Output custom resources. 



## ConfigCheck

### labels (map[string]string, optional) {#configcheck-labels}

Labels to use for the configcheck pods on top of labels added by the operator by default. Default values can be overwritten. 


### strategy (ConfigCheckStrategy, optional) {#configcheck-strategy}

Select the config check strategy to use. `DryRun`: Parse and validate configuration. `StartWithTimeout`: Start with given configuration and exit after specified timeout. Default: `DryRun` 


### timeoutSeconds (int, optional) {#configcheck-timeoutseconds}

Configure timeout in seconds if strategy is StartWithTimeout 



## RouteConfig

### disableLoggingRoute (bool, optional) {#routeconfig-disableloggingroute}

If DisableLoggingRoute is set to true, the logging route controller should remove the given tenant from the status of the logging resource. 


### enableTelemetryControllerRoute (bool, optional) {#routeconfig-enabletelemetrycontrollerroute}

If EnableTelemtryControllerRoute set to true, the operator will create the corresponding Tenant, Subscription, Output based on the logging resource. 


### tenantLabels (map[string]string, optional) {#routeconfig-tenantlabels}

TenantLabels is a map of labels that will be added to the tenant object so it can be matched with TelemetryController's TenantSelector ref: https://github.com/kube-logging/telemetry-controller/blob/main/api/telemetry/v1alpha1/collector_types.go 



## LoggingStatus

LoggingStatus defines the observed state of Logging

### configCheckResults (map[string]bool, optional) {#loggingstatus-configcheckresults}

Result of the config check. Under normal conditions there is a single item in the map with a bool value. 


### fluentdConfigName (string, optional) {#loggingstatus-fluentdconfigname}

Available in Logging operator version 4.5 and later. Name of the matched detached fluentd configuration object. 


### problems ([]string, optional) {#loggingstatus-problems}

Problems with the logging resource 


### problemsCount (int, optional) {#loggingstatus-problemscount}

Count of problems for printcolumn 


### syslogNGConfigName (string, optional) {#loggingstatus-syslogngconfigname}

Available in Logging operator version 4.5 and later. Name of the matched detached SyslogNG configuration object. 


### watchNamespaces ([]string, optional) {#loggingstatus-watchnamespaces}

List of namespaces that watchNamespaces + watchNamespaceSelector is resolving to. Not set means all namespaces. 



## Logging

Logging is the Schema for the loggings API

###  (metav1.TypeMeta, required) {#logging-}


### metadata (metav1.ObjectMeta, optional) {#logging-metadata}


### spec (LoggingSpec, optional) {#logging-spec}


### status (LoggingStatus, optional) {#logging-status}



## LoggingList

LoggingList contains a list of Logging

###  (metav1.TypeMeta, required) {#logginglist-}


### metadata (metav1.ListMeta, optional) {#logginglist-metadata}


### items ([]Logging, required) {#logginglist-items}



## DefaultFlowSpec

DefaultFlowSpec is a Flow for logs that did not match any other Flow

### filters ([]Filter, optional) {#defaultflowspec-filters}


### flowLabel (string, optional) {#defaultflowspec-flowlabel}


### globalOutputRefs ([]string, optional) {#defaultflowspec-globaloutputrefs}


### includeLabelInRouter (*bool, optional) {#defaultflowspec-includelabelinrouter}


### outputRefs ([]string, optional) {#defaultflowspec-outputrefs}

Deprecated 



