
# Logging Operator Chart

[Logging Operator](https://github.com/banzaicloud/logging-operator) Managed centralized logging component fluentd and fluent-bit instance on cluster.
## tl;dr:

```bash
$ helm repo add banzaicloud-stable https://kubernetes-charts.banzaicloud.com
$ helm repo update
$ helm install banzaicloud-stable/logging-operator
```

## Introduction

This chart bootstraps an [Logging Operator](https://github.com/banzaicloud/banzai-charts/logging-operator) deployment on a [Kubernetes](http://kubernetes.io) cluster using the [Helm](https://helm.sh) package manager.

## Prerequisites

- Kubernetes 1.8+ with Beta APIs enabled

## Installing the Chart

To install the chart with the release name `my-release`:

```bash
$ helm install --name my-release banzaicloud-stable/logging-operator
```

The command deploys **logging-operator** on the Kubernetes cluster in the default configuration. The [configuration](#configuration) section lists the parameters that can be configured during installation.

## Uninstalling the Chart

To uninstall/delete the `my-release` deployment:

```bash
$ helm delete my-release
```

The command removes all the Kubernetes components associated with the chart and deletes the release.

## Configuration

The following tables lists the configurable parameters of the logging-operator chart and their default values.

|                      Parameter                      |                        Description                     |             Default            |
| --------------------------------------------------- | ------------------------------------------------------ | ------------------------------ |
| `image.repository`                                  | Container image repository                             | `banzaicloud/logging-operator` |
| `image.tag`                                         | Container image tag                                    | `0.1.7`                        |
| `image.pullPolicy`                                  | Container pull policy                                  | `IfNotPresent`                 |
| `nameOverride`                                      | Override name of app                                   | ``                             |
| `fullnameOverride`                                  | Override full name of app                              | ``                             |
| `watchNamespace`                                    | Namespace to watch fot LoggingOperator CRD             | ``                             |
| `grafana.dashboard.enabled`                         | Install grafana logging-operator dashboard             | `true`                         |
| `rbac.enabled`                                      | Create rbac service account and roles                  | `true`                         |
| `rbac.psp.enabled`                                  | Must be used with `rbac.enabled` true. If true, creates & uses RBAC resources required in the cluster with [Pod Security Policies](https://kubernetes.io/docs/concepts/policy/pod-security-policy/) enabled.              | `false`                        |
| `affinity`                                          | Node Affinity                                          | `{}`                           |
| `resources`                                         | CPU/Memory resource requests/limits                    | `{}`                           |
| `tolerations`                                       | Node Tolerations                                       | `[]`                           |
| `nodeSelector`                                      | Define which Nodes the Pods are scheduled on.          | `{}`                           |
| `podSecurityContext`                                | Pod SecurityContext for Logging operator. [More info](https://kubernetes.io/docs/concepts/policy/security-context/)                                                                                             | `{"runAsNonRoot": true, "runAsUser": 1000, "fsGroup": 2000}` |
| `securityContext`                                   | Container SecurityContext for Logging operator. [More info](https://kubernetes.io/docs/concepts/policy/security-context/)                                                                                             | `{"allowPrivilegeEscalation": false, "readOnlyRootFilesystem": true}` |

Alternatively, a YAML file that specifies the values for the parameters can be provided while installing the chart. For example:

```bash
$ helm install --name my-release -f values.yaml banzaicloud-stable/logging-operator
```

> **Tip**: You can use the default [values.yaml](values.yaml)

## Installing Fluentd and Fluent-bit

The previous chart does **not** install Fluentd or Fluent-bit custom resource. To install them please use the [Logging Operator Fluent](https://github.com/banzaicloud/banzai-charts/logging-operator-fluent) chart.

## tl;dr:

```bash
$ helm repo add banzaicloud-stable https://kubernetes-charts.banzaicloud.com
$ helm repo update
$ helm install banzaicloud-stable/logging-operator-fluent
```

|                      Parameter                      |                        Description                     |             Default            |
| --------------------------------------------------- | ------------------------------------------------------ | ------------------------------ |
| `tls.enabled`                                       | Enabled TLS communication between components           | true                           |
| `tls.secretName`                                    | Specified secret name, which contain tls certs         | This will overwrite automatic Helm certificate generation. |
| `tls.sharedKey`                                     | Shared key between nodes (fluentd-fluentbit)           | [autogenerated]                |
| `fluentbit.enabled`                                 | Install fluent-bit                                     | true                           |
| `fluentbit.namespace`                               | Specified fluentbit installation namespace             | same as operator namespace     |
| `fluentbit.image.tag`                               | Fluentbit container image tag                          | `latest`                       |
| `fluentbit.image.repository`                        | Fluentbit container image repository                   | `fluent/fluent-bit`            |
| `fluentbit.image.pullPolicy`                        | Fluentbit container pull policy                        | `IfNotPresent`                 |
| `fluentd.enabled`                                   | Install fluentd                                        | true                           |
| `fluentd.namespace`                                 | Specified fluentd installation namespace               | same as operator namespace     |
| `fluentd.image.tag`                                 | Fluentd container image tag                            | `v1.1.4`                       |
| `fluentd.image.repository`                          | Fluentd container image repository                     | `banzaicloud/fluentd`          |
| `fluentd.image.pullPolicy`                          | Fluentd container pull policy                          | `IfNotPresent`                 |
| `fluentd.volumeModImage.tag`                        | Fluentd volumeModImage container image tag             | `latest`                       |
| `fluentd.volumeModImage.repository`                 | Fluentd volumeModImage container image repository      | `busybox`                      |
| `fluentd.volumeModImage.pullPolicy`                 | Fluentd volumeModImage container pull policy           | `IfNotPresent`                 |
| `fluentd.configReloaderImage.tag`                   | Fluentd configReloaderImage container image tag        | `v0.2.2`                       |
| `fluentd.configReloaderImage.repository`            | Fluentd configReloaderImage container image repository | `jimmidyson/configmap-reload`  |
| `fluentd.configReloaderImage.pullPolicy`            | Fluentd configReloaderImage container pull policy      | `IfNotPresent`                 |
| `fluentd.fluentdPvcSpec.accessModes`                | Fluentd persistence volume access modes                | `[ReadWriteOnce]`              |
| `fluentd.fluentdPvcSpec.resources.requests.storage` | Fluentd persistence volume size                        | `21Gi`                         |
