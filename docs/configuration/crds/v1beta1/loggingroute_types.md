## LoggingRouteSpec

### source (string, required) {#loggingroutespec-source}

Source identifies the logging that this policy applies to 

Default: -

### targets (metav1.LabelSelector, required) {#loggingroutespec-targets}

Targets refers to the list of logging resources specified by a label selector to forward logs to Filtering of namespaces will happen based on the watchNamespaces and watchNamespaceSelector fields of the target logging resource 

Default: -


## LoggingRouteStatus

### tenants ([]Tenant, optional) {#loggingroutestatus-tenants}

Enumerate all loggings with all the destination namespaces expanded 

Default: -

### problems ([]string, optional) {#loggingroutestatus-problems}

Enumerate problems that prohibits this policy to take effect 

Default: -


## Tenant

### name (string, required) {#tenant-name}

Default: -

### namespaces ([]string, optional) {#tenant-namespaces}

Default: -


## LoggingRoute

LoggingRoute defines

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


