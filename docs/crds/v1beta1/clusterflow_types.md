### ClusterFlow
#### ClusterFlow is the Schema for the clusterflows API

| Variable Name | Type | Required | Default | Description |
|---|---|---|---|---|
|  | metav1.TypeMeta | Yes | - |  |
| metadata | metav1.ObjectMeta | No | - |  |
| spec | FlowSpec | No | - | Name of the logging cluster to be attached<br> |
| status | FlowStatus | No | - |  |
### ClusterFlowList
#### ClusterFlowList contains a list of ClusterFlow

| Variable Name | Type | Required | Default | Description |
|---|---|---|---|---|
|  | metav1.TypeMeta | Yes | - |  |
| metadata | metav1.ListMeta | No | - |  |
| items | []ClusterFlow | Yes | - |  |
