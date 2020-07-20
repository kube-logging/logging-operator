---
title: ClusterOutput
weight: 200
generated_file: true
---

### ClusterOutput
#### ClusterOutput is the Schema for the clusteroutputs API

| Variable Name | Type | Required | Default | Description |
|---|---|---|---|---|
|  | metav1.TypeMeta | Yes | - |  |
| metadata | metav1.ObjectMeta | No | - |  |
| spec | ClusterOutputSpec | Yes | - |  |
| status | OutputStatus | No | - |  |
### ClusterOutputSpec
#### ClusterOutputSpec contains Kubernetes spec for CLusterOutput

| Variable Name | Type | Required | Default | Description |
|---|---|---|---|---|
|  | OutputSpec | Yes | - |  |
| enabledNamespaces | []string | No | - |  |
### ClusterOutputList
#### ClusterOutputList contains a list of ClusterOutput

| Variable Name | Type | Required | Default | Description |
|---|---|---|---|---|
|  | metav1.TypeMeta | Yes | - |  |
| metadata | metav1.ListMeta | No | - |  |
| items | []ClusterOutput | Yes | - |  |
