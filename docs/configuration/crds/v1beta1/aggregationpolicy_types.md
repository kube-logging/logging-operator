## AggregationPolicySpec

### loggingRef (string, required) {#aggregationpolicyspec-loggingref}

LoggingRef identifies the logging that this policy applies to 

Default: -

### agent (string, optional) {#aggregationpolicyspec-agent}

Agent is the name of the specific agent that this policy should be applied to Leave it empty if it should apply to all agents. 

Default: -

### watchNamespaceTargets (metav1.LabelSelector, required) {#aggregationpolicyspec-watchnamespacetargets}

WatchNamespaceTargets refers to the list of logging resources specified by a label selector to forward logs to Filtering of namespaces will happen based on the watchNamespaces and watchNamespaceSelector fields of the target logging resource 

Default: -


## AggregationPolicyStatus

### tenants ([]Tenant, optional) {#aggregationpolicystatus-tenants}

Enumerate all loggings with all the destination namespaces expanded 

Default: -

### problems ([]string, optional) {#aggregationpolicystatus-problems}

Enumerate problems that prohibits this policy to take effect 

Default: -


## Tenant

### name (string, required) {#tenant-name}

Default: -

### namespaces ([]string, optional) {#tenant-namespaces}

Default: -


## AggregationPolicy

AggregationPolicy defines

###  (metav1.TypeMeta, required) {#aggregationpolicy-}

Default: -

### metadata (metav1.ObjectMeta, optional) {#aggregationpolicy-metadata}

Default: -

### spec (AggregationPolicySpec, optional) {#aggregationpolicy-spec}

Default: -

### status (AggregationPolicyStatus, optional) {#aggregationpolicy-status}

Default: -


## AggregationPolicyList

###  (metav1.TypeMeta, required) {#aggregationpolicylist-}

Default: -

### metadata (metav1.ListMeta, optional) {#aggregationpolicylist-metadata}

Default: -

### items ([]AggregationPolicy, required) {#aggregationpolicylist-items}

Default: -


