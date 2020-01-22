### LoggingSpec
#### LoggingSpec defines the desired state of Logging

| Variable Name | Type | Required | Default | Description |
|---|---|---|---|---|
| loggingRef | string | No | - |  |
| flowConfigCheckDisabled | bool | No | - |  |
| flowConfigOverride | string | No | - |  |
| fluentbit | *FluentbitSpec | No | - |  |
| fluentd | *FluentdSpec | No | - |  |
| watchNamespaces | []string | No | - |  |
| controlNamespace | string | Yes | - |  |
| enableRecreateWorkloadOnImmutableFieldChange | bool | No | - | EnableRecreateWorkloadOnImmutableFieldChange enables the operator to recreate the<br>fluentbit daemonset and the fluentd statefulset (and possibly other resource in the future)<br>in case there is a change in an immutable field<br>that otherwise couldn't be managed with a simple update.<br> |
### LoggingStatus
#### LoggingStatus defines the observed state of Logging

| Variable Name | Type | Required | Default | Description |
|---|---|---|---|---|
| configCheckResults | map[string]bool | No | - |  |
### Logging
#### Logging is the Schema for the loggings API

| Variable Name | Type | Required | Default | Description |
|---|---|---|---|---|
|  | metav1.TypeMeta | Yes | - |  |
| metadata | metav1.ObjectMeta | No | - |  |
| spec | LoggingSpec | No | - |  |
| status | LoggingStatus | No | - |  |
### LoggingList
#### LoggingList contains a list of Logging

| Variable Name | Type | Required | Default | Description |
|---|---|---|---|---|
|  | metav1.TypeMeta | Yes | - |  |
| metadata | metav1.ListMeta | No | - |  |
| items | []Logging | Yes | - |  |
