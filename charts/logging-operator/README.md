# logging-operator

![version: 4.2.3](https://img.shields.io/badge/version-4.2.3-informational?style=flat-square) ![type: application](https://img.shields.io/badge/type-application-informational?style=flat-square) ![app version: 4.2.2](https://img.shields.io/badge/app%20version-4.2.2-informational?style=flat-square) ![kube version: >=1.22.0-0](https://img.shields.io/badge/kube%20version->=1.22.0--0-informational?style=flat-square) [![artifact hub](https://img.shields.io/badge/artifact%20hub-logging--operator-informational?style=flat-square)](https://artifacthub.io/packages/helm/kube-logging/logging-operator)

Logging operator for Kubernetes based on Fluentd and Fluentbit.

**Homepage:** <https://kube-logging.github.io>

## TL;DR;

```bash
helm repo add kube-logging https://kube-logging.github.io/helm-charts
helm install --generate-name --wait kube-logging/logging-operator
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
| watchNamespace | string | `""` | Namespace to watch for LoggingOperator Custom Resources. |
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

## Installing Fluentd and Fluent-bit via logging

The previous chart does **not** install `logging` resource to deploy Fluentd and Fluent-bit on cluster.
To install them please use the [Logging Operator Logging](https://github.com/kube-logging/helm-charts/tree/main/charts/logging-operator-logging) chart.
