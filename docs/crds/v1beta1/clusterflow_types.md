---
title: ClusterFlow
weight: 200
generated_file: true
---

### ClusterFlow
#### ClusterFlow is the Schema for the clusterflows API

| Variable Name | Type | Required | Default | Description |
|---|---|---|---|---|
|  | metav1.TypeMeta | Yes | - |  |
| metadata | metav1.ObjectMeta | No | - |  |
| spec | ClusterFlowSpec | No | - | Name of the logging cluster to be attached<br> |
| status | FlowStatus | No | - |  |
### ClusterMatch
| Variable Name | Type | Required | Default | Description |
|---|---|---|---|---|
| select | *ClusterSelect | No | - |  |
| exclude | *ClusterExclude | No | - |  |
### ClusterSelect
| Variable Name | Type | Required | Default | Description |
|---|---|---|---|---|
| namespaces | []string | No | - |  |
| labels | map[string]string | No | - |  |
| hosts | []string | No | - |  |
| container_names | []string | No | - |  |
### ClusterExclude
| Variable Name | Type | Required | Default | Description |
|---|---|---|---|---|
| namespaces | []string | No | - |  |
| labels | map[string]string | No | - |  |
| hosts | []string | No | - |  |
| container_names | []string | No | - |  |
### ClusterFlowSpec
#### FlowSpec is the Kubernetes spec for Flows

| Variable Name | Type | Required | Default | Description |
|---|---|---|---|---|
| selectors | map[string]string | No | - | Deprecated<br> |
| match | []ClusterMatch | No | - |  |
| filters | []Filter | No | - |  |
| loggingRef | string | No | - |  |
| outputRefs | []string | No | - | Deprecated<br> |
| globalOutputRefs | []string | No | - |  |
### ClusterFlowList
#### ClusterFlowList contains a list of ClusterFlow

| Variable Name | Type | Required | Default | Description |
|---|---|---|---|---|
|  | metav1.TypeMeta | Yes | - |  |
| metadata | metav1.ListMeta | No | - |  |
| items | []ClusterFlow | Yes | - |  |
