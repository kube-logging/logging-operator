---
title: LoggingRouteSpec
weight: 200
generated_file: true
---

## LoggingRouteSpec

LoggingRouteSpec defines the desired state of LoggingRoute

### source (string, required) {#loggingroutespec-source}

Source identifies the logging that this policy applies to 


### targets (metav1.LabelSelector, required) {#loggingroutespec-targets}

Targets refers to the list of logging resources specified by a label selector to forward logs to. Filtering of namespaces will happen based on the watchNamespaces and watchNamespaceSelector fields of the target logging resource. 



## LoggingRouteStatus

LoggingRouteStatus defines the actual state of the LoggingRoute

### notices ([]string, optional) {#loggingroutestatus-notices}

Enumerate non-blocker issues the user should pay attention to 


### noticesCount (int, optional) {#loggingroutestatus-noticescount}

Summarize the number of notices for the CLI output 


### problems ([]string, optional) {#loggingroutestatus-problems}

Enumerate problems that prohibits this route to take effect and populate the tenants field 


### problemsCount (int, optional) {#loggingroutestatus-problemscount}

Summarize the number of problems for the CLI output 


### tenants ([]Tenant, optional) {#loggingroutestatus-tenants}

Enumerate all loggings with all the destination namespaces expanded 



## Tenant

### name (string, required) {#tenant-name}


### namespaces ([]string, optional) {#tenant-namespaces}



## LoggingRoute

LoggingRoute (experimental)
Connects a log collector with log aggregators from other logging domains and routes relevant logs based on watch namespaces

###  (metav1.TypeMeta, required) {#loggingroute-}


### metadata (metav1.ObjectMeta, optional) {#loggingroute-metadata}


### spec (LoggingRouteSpec, optional) {#loggingroute-spec}


### status (LoggingRouteStatus, optional) {#loggingroute-status}



## LoggingRouteList

###  (metav1.TypeMeta, required) {#loggingroutelist-}


### metadata (metav1.ListMeta, optional) {#loggingroutelist-metadata}


### items ([]LoggingRoute, required) {#loggingroutelist-items}



