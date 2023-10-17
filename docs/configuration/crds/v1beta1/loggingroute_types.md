---
title: LoggingRouteSpec
weight: 200
generated_file: true
---

## LoggingRouteSpec

LoggingRouteSpec defines the desired state of LoggingRoute

### source (string, required) {#loggingroutespec-source}

Source identifies the logging that this policy applies to 

Default: -

### targets (metav1.LabelSelector, required) {#loggingroutespec-targets}

Targets refers to the list of logging resources specified by a label selector to forward logs to Filtering of namespaces will happen based on the watchNamespaces and watchNamespaceSelector fields of the target logging resource 

Default: -


## LoggingRouteStatus

LoggingRouteStatus defines the actual state of the LoggingRoute

### notices ([]string, optional) {#loggingroutestatus-notices}

Enumerate non-blocker issues the user should pay attention to 

Default: -

### noticesCount (int, optional) {#loggingroutestatus-noticescount}

Summarize the number of notices for the CLI output 

Default: -

### problems ([]string, optional) {#loggingroutestatus-problems}

Enumerate problems that prohibits this route to take effect and populate the tenants field 

Default: -

### problemsCount (int, optional) {#loggingroutestatus-problemscount}

Summarize the number of problems for the CLI output 

Default: -

### tenants ([]Tenant, optional) {#loggingroutestatus-tenants}

Enumerate all loggings with all the destination namespaces expanded 

Default: -


## Tenant

### name (string, required) {#tenant-name}

Default: -

### namespaces ([]string, optional) {#tenant-namespaces}

Default: -


## LoggingRoute

LoggingRoute (experimental)
Connects a log collector with log aggregators from other logging domains and routes relevant logs based on watch namespaces

###  (metav1.TypeMeta, required) {#loggingroute-}

Default: -

### metadata (metav1.ObjectMeta, optional) {#loggingroute-metadata}

Default: -

### spec (LoggingRouteSpec, optional) {#loggingroute-spec}

Default: -

### status (LoggingRouteStatus, optional) {#loggingroute-status}

Default: -


## LoggingRouteList

###  (metav1.TypeMeta, required) {#loggingroutelist-}

Default: -

### metadata (metav1.ListMeta, optional) {#loggingroutelist-metadata}

Default: -

### items ([]LoggingRoute, required) {#loggingroutelist-items}

Default: -


