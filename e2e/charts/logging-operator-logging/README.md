# logging-operator-logging

![version: 4.0.0-rc18-1](https://img.shields.io/badge/version-4.0.0--rc18--1-informational?style=flat-square) ![type: application](https://img.shields.io/badge/type-application-informational?style=flat-square) ![app version: 4.0.0-rc17](https://img.shields.io/badge/app%20version-4.0.0--rc17-informational?style=flat-square) ![kube version: >=1.16.0-0](https://img.shields.io/badge/kube%20version->=1.16.0--0-informational?style=flat-square) [![artifact hub](https://img.shields.io/badge/artifact%20hub-logging--operator--logging-informational?style=flat-square)](https://artifacthub.io/packages/helm/kube-logging/logging-operator-logging)

A Helm chart to configure logging resource for the Logging operator.

**Homepage:** <https://kube-logging.github.io>

## TL;DR;

```bash
helm repo add kube-logging https://kube-logging.github.io/helm-charts
helm install --generate-name --wait kube-logging/logging-operator-logging
```

## Values

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| loggingRef | string | `""` |  |
| flowConfigCheckDisabled | bool | `false` |  |
| flowConfigOverride | string | `""` |  |
| nameOverride | string | `""` |  |
| fullnameOverride | string | `""` |  |
| enableRecreateWorkloadOnImmutableFieldChange | bool | `false` | Permit deletion and recreation of resources on update of immutable field. |
| tls.enabled | bool | `true` | Enable secure connection between fluentd and fluent-bit |
| tls.fluentdSecretName | string | `""` | Specified secret name, which contain tls certs |
| tls.fluentbitSecretName | string | `""` | Specified secret name, which contain tls certs |
| tls.sharedKey | string | `""` |  |
| fluentbit | object | `{}` | Fluent-bit configurations https://banzaicloud.com/docs/one-eye/logging-operator/configuration/crds/v1beta1/fluentbit_types/ |
| fluentd | object | `{}` | Fluentd configurations https://banzaicloud.com/docs/one-eye/logging-operator/configuration/crds/v1beta1/fluentd_types/ |
| nodeAgents | object | `{}` | Node agents definitions |
| skipInvalidResources | bool | `false` | Whether to skip invalid Flow and ClusterFlow resources |
| watchNamespaces | list | `[]` | Limit namespaces from where to read Flow and Output specs |
| clusterDomain | string | `"cluster.local"` | Cluster domain name to be used when templating URLs to services |
| controlNamespace | string | `""` | Control namespace that contains ClusterOutput and ClusterFlow resources |
| allowClusterResourcesFromAllNamespaces | bool | `false` | Allow configuration of cluster resources from any namespace |
| defaultFlow | object | `{}` | Default flow |
| globalFilters | list | `[]` | Global filters |
| clusterFlows | list | `[]` | ClusterFlows to deploy |
| clusterOutputs | list | `[]` | ClusterOutputs to deploy |
| eventTailer | object | `{}` | EventTailer config |
| hostTailer | object | `{}` | HostTailer config |
| scc.enabled | bool | `false` | OpenShift SecurityContextConstraints enabled |
