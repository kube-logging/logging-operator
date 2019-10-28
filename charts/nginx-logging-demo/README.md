
# Logging Operator Nginx DemoApp with Minio (local S3-compatible object store) output  

[Logging Operator](https://github.com/banzaicloud/logging-operator) is a managed centralized logging component based on fluentd and fluent-bit.
## tl;dr:

```bash
$ helm repo add banzaicloud-stable https://kubernetes-charts.banzaicloud.com/
$ helm repo update
$ helm install banzaicloud-stable/nginx-logging-demo
```

## Introduction

This chart demonstrates the use of the  [Logging Operator](https://github.com/banzaicloud/banzai-charts/logging-operator) with an Nginx deployment on a [Kubernetes](http://kubernetes.io) cluster using the [Helm](https://helm.sh) package manager.

## Prerequisites

- [Logging Operator](https://github.com/banzaicloud/logging-operator) available on the cluster


## Installing the Chart

To install the chart with the release name `log-test-nginx`:

```bash
$ helm install --name log-test-nginx banzaicloud-stable/nginx-logging-demo
```
## Uninstalling the Chart

To uninstall/delete the `log-test-nginx` deployment:

```bash
$ helm delete log-test-nginx
```

The command removes all the Kubernetes components associated with the chart and deletes the release.

## Configuration

The following tables lists the configurable parameters of the nginx-logging-demo chart and their default values.

|                          Parameter                                |                        Description                      |     Default    |
| ---------------------------------------------------------         | ------------------------------------------------------- | -------------- |
| `image.repository`                                                | Container image repository                              | `nginx`        |
| `image.tag`                                                       | Container image tag                                     | `stable`       |
| `image.pullPolicy`                                                | Container pull policy                                   | `IfNotPresent` |
| `nameOverride`                                                    | Override name of app                                    | ``             |
| `fullnameOverride`                                                | Override full name of app                               | ``             |
| `affinity`                                                        | Node Affinity                                           | `{}`           |
| `resources`                                                       | CPU/Memory resource requests/limits                     | `{}`           |
| `tolerations`                                                     | Node Toleration                                         | `[]`           |
| `nodeSelector`                                                    | Define which Nodes the Pods are scheduled on.           | `{}`           |
| `loggingOperator.fluentd.metrics.serviceMonitor`                  | Enable to create ServiceMonitor for Prometheus operator | `false`        |
| `loggingOperator.fluentd.metrics.prometheusAnnotations`           | Add prometheus labels to fluent pods.                   | `false`        |
| `loggingOperator.fluentd.metrics.port`                            | Metrics Port.                                           | ``             |
| `loggingOperator.fluentd.metrics.path`                            | Metrics Path                                            | ``             |
| `loggingOperator.fluentd.metrics.timeout`                         | Scrape timeout.                                         | ``             |
| `loggingOperator.fluentd.metrics.interval`                        | Scrape interval.                                        | ``             |
| `loggingOperator.fluentbit.metrics.serviceMonitor`                | Enable to create ServiceMonitor for Prometheus operator | `false`        |
| `loggingOperator.fluentbit.metrics.prometheusAnnotations`         | Add prometheus labels to fluent pods.                   | `false`        |
| `loggingOperator.fluentbit.metrics.port`                          | Metrics Port.                                           | ``             |
| `loggingOperator.fluentbit.metrics.path`                          | Metrics Path                                            | ``             |
| `loggingOperator.fluentbit.metrics.timeout`                       | Scrape timeout.                                         | ``             |
| `loggingOperator.fluentbit.metrics.interval`                      | Scrape interval.                                        | ``             |
| `loggingOperator.fluentd.security.roleBasedAccessControlCreate`   | Enable fluentd RBAC                                     | `true`         |
| `loggingOperator.fluentd.security.podSecurityPolicyCreate`        | Enable fluentd PSP                                      | `true`         |
| `loggingOperator.fluentd.security.serviceAccount`                 | Set fluentd Service Account                             | ``             |
| `loggingOperator.fluentbit.security.roleBasedAccessControlCreate` | Enable fluentbit RBAC                                   | `true`         |
| `loggingOperator.fluentbit.security.podSecurityPolicyCreate`      | Enable fluentbit PSP                                    | `true`         |
| `loggingOperator.fluentbit.security.serviceAccount`               | Set fluentbit Service Account                           | ``             |


Alternatively, a YAML file that specifies the values for the parameters can be provided while installing the chart. For example:

```bash
$ helm install --name my-release -f values.yaml banzaicloud-stable/nginx-logging-demo
```

> **Tip**: You can use the default [values.yaml](values.yaml)

