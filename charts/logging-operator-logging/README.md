# logging-operator-logging

![version: 4.2.2](https://img.shields.io/badge/version-4.2.2-informational?style=flat-square) ![type: application](https://img.shields.io/badge/type-application-informational?style=flat-square)  ![kube version: >=1.22.0-0](https://img.shields.io/badge/kube%20version->=1.22.0--0-informational?style=flat-square) [![artifact hub](https://img.shields.io/badge/artifact%20hub-logging--operator--logging-informational?style=flat-square)](https://artifacthub.io/packages/helm/kube-logging/logging-operator-logging)

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
| nameOverride | string | `""` |  |
| fullnameOverride | string | `""` |  |
| tls.enabled | bool | `true` | Enable secure connection between fluentd and fluent-bit |
| tls.fluentdSecretName | string | `""` | Specified secret name, which contain tls certs |
| tls.fluentbitSecretName | string | `""` | Specified secret name, which contain tls certs |
| tls.sharedKey | string | `""` |  |
| loggingRef | string | `""` | Reference to the logging system. Each of the loggingRefs can manage a fluentbit daemonset and a fluentd statefulset. |
| flowConfigCheckDisabled | bool | `false` | Disable configuration check before applying new fluentd configuration. |
| skipInvalidResources | bool | `false` | Whether to skip invalid Flow and ClusterFlow resources |
| flowConfigOverride | string | `""` | Override generated config. This is a raw configuration string for troubleshooting purposes. |
| fluentbitDisabled | bool | `false` | Flag to disable fluentbit completely |
| fluentbit | object | `{}` | Fluent-bit configurations https://kube-logging.github.io/docs/configuration/crds/v1beta1/fluentbit_types/ |
| fluentdDisabled | bool | `false` | Flag to disable fluentd completely |
| fluentd | object | `{}` | Fluentd configurations https://kube-logging.github.io/docs/configuration/crds/v1beta1/fluentd_types/ |
| syslogNG | object | `{}` | Syslog-NG statefulset configuration |
| defaultFlow | object | `{}` | Default flow for unmatched logs. This Flow configuration collects all logs that didn’t matched any other Flow. |
| errorOutputRef | string | `""` | GlobalOutput name to flush ERROR events to |
| globalFilters | list | `[]` | Global filters to apply on logs before any match or filter mechanism. |
| watchNamespaces | list | `[]` | Limit namespaces to watch Flow and Output custom resources. |
| clusterDomain | string | `"cluster.local"` | Cluster domain name to be used when templating URLs to services |
| controlNamespace | string | `""` | Namespace for cluster wide configuration resources like ClusterFlow and ClusterOutput. This should be a protected namespace from regular users. Resources like fluentbit and fluentd will run in this namespace as well. |
| allowClusterResourcesFromAllNamespaces | bool | `false` | Allow configuration of cluster resources from any namespace. Mutually exclusive with ControlNamespace restriction of Cluster resources |
| nodeAgents | object | `{}` | NodeAgent Configuration |
| enableRecreateWorkloadOnImmutableFieldChange | bool | `false` | EnableRecreateWorkloadOnImmutableFieldChange enables the operator to recreate the fluentbit daemonset and the fluentd statefulset (and possibly other resource in the future) in case there is a change in an immutable field that otherwise couldn’t be managed with a simple update. |
| clusterFlows | list | `[]` | ClusterFlows to deploy |
| clusterOutputs | list | `[]` | ClusterOutputs to deploy |
| eventTailer | object | `{}` | EventTailer config |
| hostTailer | object | `{}` | HostTailer config |
| scc.enabled | bool | `false` | OpenShift SecurityContextConstraints enabled |
