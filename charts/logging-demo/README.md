
# Logging Operator DemoApp Application  

[Logging Operator](https://github.com/banzaicloud/logging-operator) is a managed centralized logging component based on fluentd and fluent-bit.
## tl;dr:

```bash
$ helm repo add banzaicloud-stable https://kubernetes-charts.banzaicloud.com/
$ helm repo update
$ helm install banzaicloud-stable/logging-demo
```

## Introduction

This chart demonstrates the use of the  [Logging Operator](https://github.com/banzaicloud/banzai-charts/logging-operator) with an [Log-Generator](https://github.com/banzaicloud/log-generator) deployment on a [Kubernetes](http://kubernetes.io) cluster using the [Helm](https://helm.sh) package manager.

## Prerequisites

- [Logging Operator](https://github.com/banzaicloud/logging-operator/blob/master/docs/deploy/README.md) available on the cluster


## Installing the Chart

To install the chart with the release name `logging-demo`:

```bash
$ helm install --namespace logging --name logging-demo banzaicloud-stable/logging-demo
```
## Uninstalling the Chart

To uninstall/delete the `logging-demo` deployment:

```bash
$ helm delete logging-demo
```

The command removes all the Kubernetes components associated with the chart and deletes the release.

## Configuration

The following tables lists the configurable parameters of the logging-demo chart and their default values.

|                          Parameter                                |                        Description                      |     Default    |
| ------------------------------------------------------------------| ------------------------------------------------------- | -------------- |
| `nameOverride`                                                    | Override name of app                                    | ``             |
| `fullnameOverride`                                                | Override full name of app                               | ``             |
| `loggingOperator.fluentd.metrics.serviceMonitor`                  | Enable to create ServiceMonitor for Prometheus operator | `false`        |
| `loggingOperator.fluentd.metrics.prometheusAnnotations`           | Add prometheus labels to fluent pods.                   | `false`        |
| `loggingOperator.fluentd.metrics.port`                            | Metrics Port.                                           | ``             |
| `loggingOperator.fluentd.metrics.path`                            | Metrics Path                                            | ``             |
| `loggingOperator.fluentd.metrics.timeout`                         | Scrape timeout.                                         | ``             |
| `loggingOperator.fluentd.metrics.interval`                        | Scrape interval.                                        | ``             |
| `loggingOperator.fluentd.logLevel`                                | FluentD loglevel fatal,error,warn,info,debug,trace      | `info`         |
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
| `elasticsearch.enabled`                                           | Enable ElasticSearch logging output                     | `false`        |
| `loki.enabled`                                                    | Enable Grafana Loki logging output                      | `false`        |
| `kafka.enabled`                                                   | Enable Kafka logging output                             | `false`        |
| `minio.enabled`                                                   | Enable Minio logging output and install chart           | `false`        |
| `cloudwatch.enabled`                                              | Enable AWS Cloudwatch logging output                    | `false`        |
| `cloudwatch.aws.secret_key`                                       | AWS Secret Access Key                                   | ``             |
| `cloudwatch.aws.access_key`                                       | AWS Access Key ID                                       | ``             |
| `cloudwatch.aws.region`                                           | AWS CLoudWatch Region                                   | ``             |
| `cloudwatch.aws.log_group_name`                                   | AWS CLoudWatch Log Group                                | ``             |
| `cloudwatch.aws.log_stream_name`                                  | AWS CLoudWatch Log Stream                               | ``             |
| `logdna.enabled`                                                  | Enable LogDNA logging output                            | `false`        |
| `logdna.api_key`                                                  | LogDNA Api key                                          | ``             |
| `logdna.hostname`                                                 | Hostname                                                | ``             |
| `logdna.app`                                                      | Application name                                        | ``             |
| `logGenerator.enabled`                                            | Enable Demo Log-Gen application                         | `true`         |
| `loggingOperator.tls.enabled`                                     | Enabled TLS communication between components            | `true`         |
| `loggingOperator.tls.fluentdSecretName`                           | Specified secret name, which contain tls certs          | This will overwrite automatic Helm certificate generation. |
| `loggingOperator.tls.fluentbitSecretName`                         | Specified secret name, which contain tls certs          | This will overwrite automatic Helm certificate generation. |
| `loggingOperator.tls.sharedKey`                                   | Shared key between nodes (fluentd-fluentbit)            | [autogenerated] |


Alternatively, a YAML file that specifies the values for the parameters can be provided while installing the chart. For example:

```bash
$ helm install --name my-release -f values.yaml banzaicloud-stable/logging-demo
```

> **Tip**: You can use the default [values.yaml](values.yaml)

