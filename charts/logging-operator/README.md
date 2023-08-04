# logging-operator

![type: application](https://img.shields.io/badge/type-application-informational?style=flat-square) ![kube version: >=1.22.0-0](https://img.shields.io/badge/kube%20version->=1.22.0--0-informational?style=flat-square) [![artifact hub](https://img.shields.io/badge/artifact%20hub-logging--operator-informational?style=flat-square)](https://artifacthub.io/packages/helm/kube-logging/logging-operator)

Logging operator for Kubernetes based on Fluentd and Fluentbit.

**Homepage:** <https://kube-logging.github.io>

## TL;DR;

```bash
helm install --generate-name --wait oci://ghcr.io/kube-logging/helm-charts/logging-operator
```

or to install with a specific version:

```bash
helm install --generate-name --wait oci://ghcr.io/kube-logging/helm-charts/logging-operator --version $VERSION
```

## Introduction

This chart bootstraps a [Logging Operator](https://github.com/kube-logging/logging-operator) deployment on a [Kubernetes](http://kubernetes.io) cluster using the [Helm](https://helm.sh) package manager.

## Prerequisites

- Kubernetes 1.19+

## Installing CRDs

Use `createCustomResource=false` with Helm v3 to avoid trying to create CRDs from the `crds` folder and from templates at the same time.

## Values

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| replicaCount | int | `1` |  |
| image.repository | string | `"ghcr.io/kube-logging/logging-operator"` | Name of the image repository to pull the container image from. |
| image.tag | string | `""` | Image tag override for the default value (chart appVersion). |
| image.pullPolicy | string | `"IfNotPresent"` | [Image pull policy](https://kubernetes.io/docs/concepts/containers/images/#updating-images) for updating already existing images on a node. |
| env | list | `[]` |  |
| volumes | list | `[]` |  |
| volumeMounts | list | `[]` |  |
| extraArgs[0] | string | `"-enable-leader-election=true"` |  |
| imagePullSecrets | list | `[]` |  |
| nameOverride | string | `""` | A name in place of the chart name for `app:` labels. |
| fullnameOverride | string | `""` | A name to substitute for the full names of resources. |
| namespaceOverride | string | `""` | A namespace override for the app. |
| annotations | object | `{}` | Define annotations for logging-operator pods. |
| createCustomResource | bool | `false` | Deploy CRDs used by Logging Operator. |
| http.port | int | `8080` | HTTP listen port number. |
| http.service | object | `{"annotations":{},"clusterIP":"None","labels":{},"type":"ClusterIP"}` | Service definition for query http service. |
| rbac.enabled | bool | `true` | Create rbac service account and roles. |
| rbac.psp.enabled | bool | `true` | Must be used with `rbac.enabled` true. If true, creates & uses RBAC resources required in the cluster with [Pod Security Policies](https://kubernetes.io/docs/concepts/policy/pod-security-policy/) enabled. |
| rbac.psp.annotations | object | `{"seccomp.security.alpha.kubernetes.io/allowedProfileNames":"docker/default,runtime/default","seccomp.security.alpha.kubernetes.io/defaultProfileName":"runtime/default"}` | PSP annotations |
| monitoring.serviceMonitor.enabled | bool | `false` | Create a Prometheus Operator ServiceMonitor object. |
| monitoring.serviceMonitor.additionalLabels | object | `{}` |  |
| monitoring.serviceMonitor.metricRelabelings | list | `[]` |  |
| monitoring.serviceMonitor.relabelings | list | `[]` |  |
| podSecurityContext | object | `{}` | Pod SecurityContext for Logging operator. [More info](https://kubernetes.io/docs/concepts/policy/security-context/) # SecurityContext holds pod-level security attributes and common container settings. # This defaults to non root user with uid 1000 and gid 2000.	*v1.PodSecurityContext	false # ref: https://kubernetes.io/docs/tasks/configure-pod-container/security-context/ |
| securityContext | object | `{}` | Container SecurityContext for Logging operator. [More info](https://kubernetes.io/docs/concepts/policy/security-context/) |
| priorityClassName | object | `{}` | Operator priorityClassName. |
| serviceAccount.annotations | object | `{}` | Define annotations for logging-operator ServiceAccount. |
| resources | object | `{}` | CPU/Memory resource requests/limits |
| nodeSelector | object | `{}` |  |
| tolerations | list | `[]` | Node Tolerations |
| affinity | object | `{}` | Node Affinity |
| podLabels | object | `{}` | Define which Nodes the Pods are scheduled on. |
| logging | object | `{"allowClusterResourcesFromAllNamespaces":false,"clusterDomain":"cluster.local.","clusterFlows":[],"clusterOutputs":[],"controlNamespace":"","defaultFlow":{},"enableRecreateWorkloadOnImmutableFieldChange":false,"enabled":false,"errorOutputRef":"","eventTailer":{},"flowConfigCheckDisabled":false,"flowConfigOverride":"","fluentbit":{},"fluentbitDisabled":false,"fluentd":{},"fluentdDisabled":false,"globalFilters":[],"hostTailer":{},"loggingRef":"","nodeAgents":{},"skipInvalidResources":false,"syslogNG":{},"watchNamespaces":[]}` | Logging resources configuration. |
| logging.enabled | bool | `false` | Logging resources are disabled by default |
| logging.loggingRef | string | `""` | Reference to the logging system. Each of the loggingRefs can manage a fluentbit daemonset and a fluentd statefulset. |
| logging.flowConfigCheckDisabled | bool | `false` | Disable configuration check before applying new fluentd configuration. |
| logging.skipInvalidResources | bool | `false` | Whether to skip invalid Flow and ClusterFlow resources |
| logging.flowConfigOverride | string | `""` | Override generated config. This is a raw configuration string for troubleshooting purposes. |
| logging.fluentbitDisabled | bool | `false` | Flag to disable fluentbit completely |
| logging.fluentbit | object | `{}` | Fluent-bit configurations https://kube-logging.github.io/docs/configuration/crds/v1beta1/fluentbit_types/ |
| logging.fluentdDisabled | bool | `false` | Flag to disable fluentd completely |
| logging.fluentd | object | `{}` | Fluentd configurations https://kube-logging.github.io/docs/configuration/crds/v1beta1/fluentd_types/ |
| logging.syslogNG | object | `{}` | Syslog-NG statefulset configuration |
| logging.defaultFlow | object | `{}` | Default flow for unmatched logs. This Flow configuration collects all logs that didn’t match any other Flow. |
| logging.errorOutputRef | string | `""` | GlobalOutput name to flush ERROR events to |
| logging.globalFilters | list | `[]` | Global filters to apply on logs before any match or filter mechanism. |
| logging.watchNamespaces | list | `[]` | Limit namespaces to watch Flow and Output custom resources. |
| logging.clusterDomain | string | `"cluster.local."` | Cluster domain name to be used when templating URLs to services |
| logging.controlNamespace | string | `""` | Namespace for cluster wide configuration resources like ClusterFlow and ClusterOutput. This should be a protected namespace from regular users. Resources like fluentbit and fluentd will run in this namespace as well. |
| logging.allowClusterResourcesFromAllNamespaces | bool | `false` | Allow configuration of cluster resources from any namespace. Mutually exclusive with ControlNamespace restriction of Cluster resources |
| logging.nodeAgents | object | `{}` | NodeAgent Configuration |
| logging.enableRecreateWorkloadOnImmutableFieldChange | bool | `false` | EnableRecreateWorkloadOnImmutableFieldChange enables the operator to recreate the fluentbit daemonset and the fluentd statefulset (and possibly other resource in the future) in case there is a change in an immutable field that otherwise couldn’t be managed with a simple update. |
| logging.clusterFlows | list | `[]` | ClusterFlows to deploy |
| logging.clusterOutputs | list | `[]` | ClusterOutputs to deploy |
| logging.eventTailer | object | `{}` | EventTailer config |
| logging.hostTailer | object | `{}` | HostTailer config |
| testReceiver.enabled | bool | `false` |  |
| testReceiver.image | string | `"fluent/fluent-bit"` |  |
| testReceiver.pullPolicy | string | `"IfNotPresent"` |  |
| testReceiver.port | int | `8080` |  |
| testReceiver.args[0] | string | `"-i"` |  |
| testReceiver.args[1] | string | `"http"` |  |
| testReceiver.args[2] | string | `"-p"` |  |
| testReceiver.args[3] | string | `"port=8080"` |  |
| testReceiver.args[4] | string | `"-o"` |  |
| testReceiver.args[5] | string | `"stdout"` |  |
| testReceiver.resources.limits.cpu | string | `"100m"` |  |
| testReceiver.resources.limits.memory | string | `"50Mi"` |  |
| testReceiver.resources.requests.cpu | string | `"20m"` |  |
| testReceiver.resources.requests.memory | string | `"25Mi"` |  |

## Installing Fluentd and Fluent-bit via logging

The chart does **not** install `logging` resource to deploy Fluentd (or Syslog-ng) and Fluent-bit on the cluster by default, but
it can be enabled by setting the `logging.enabled` value to true.
