<p align="center"><img src="docs/img/lo.svg" width="260"></p>
<p align="center">

  <a href="https://hub.docker.com/r/banzaicloud/logging-operator/">
    <img src="https://img.shields.io/docker/automated/banzaicloud/logging-operator.svg" alt="Docker Automated build">
  </a>

  <a href="https://hub.docker.com/r/banzaicloud/logging-operator/">
    <img src="https://img.shields.io/docker/pulls/banzaicloud/logging-operator.svg?style=shield" alt="Docker Pulls">
  </a>

  <a href="https://circleci.com/gh/banzaicloud/logging-operator">
    <img src="https://circleci.com/gh/banzaicloud/logging-operator.svg?style=shield" alt="CircleCI">
  </a>

  <a href="https://goreportcard.com/badge/github.com/banzaicloud/logging-operator">
    <img src="https://goreportcard.com/badge/github.com/banzaicloud/logging-operator" alt="Go Report Card">
  </a>

  <a href="https://github.com/banzaicloud/logging-operator/">
    <img src="https://img.shields.io/badge/license-Apache%20v2-orange.svg" alt="license">
  </a>

</p>


# logging-operator v2

Logging operator for Kubernetes based on Fluentd and Fluent-bit.


## What is this operator for?

This operator helps you to pack together logging information with your applications. With the help of Custom Resource Definition you can describe the behaviour of your application within its charts. The operator does the rest.

<p align="center"><img src="docs/img/logging_operator_flow.png" ></p>

### Feature highlights

- [x] Namespace isolation 
- [x] Native Kubernetes label selectors
- [x] Secure communication (TLS)
- [x] Configuration validation
- [x] Multiple flow support (multiply logs for different transformations)
- [x] Multiple [output](docs/plugins/outputs) support (store the same logs in multiple storage: S3, GCS, ES, Loki and more...)
- [x] Multiple logging system support (multiple fluentd, fluent-bit deployment on the same cluster)

### Motivation

The logging operator automates the deployment and configuration of a Kubernetes logging pipeline. Under the hood the operator configures a fluent-bit daemonset for collecting container logs from the node file system. Fluent-bit enriches the logs with Kubernetes metadata and transfers them to fluentd. Fluentd receives, filters and transfer logs to multiple outputs. Your logs will always be transferred on authenticated and encrypted channels.

## Architecture

Available custom resources:
- [logging](/docs/crds.md#loggings) - Represents a logging system. Includes `Fluentd` and `Fluent-bit` configuration. Specifies the `controlNamespace`. Fluentd and Fluent-bit will be deployed in the `controlNamespace`
- [output](/docs/crds.md#outputs-clusteroutputs) - Defines an Output for a logging flow. This is a namespaced resource.
- [flow](/docs/crds.md#flows-clusterflows) - Defines a logging flow with `filters` and `outputs`. You can specify `selectors` to filter logs by labels. Outputs can be `output` or `clusteroutput`.  This is a namespaced resource.
- [clusteroutput](/docs/crds.md#outputs-clusteroutputs) - Defines an output without namespace restriction. Only effective in `controlNamespace`.
- [clusterflow](/docs/crds.md#flows-clusterflows) - Defines a logging flow without namespace restriction.

The detailed CRD documentation can be found [here](/docs/crds.md).

<p align="center"><img src="docs/img/logging-operator-v2-architecture.png" ></p>

*connection between custom resources*

### Blogs
  - [Logging-Operator v2](https://banzaicloud.com/blog/logging-operator-v2/)
  - [Eleasticsearch and GeoIP](https://banzaicloud.com/blog/logging-operator-efk/)  

##### Blogs (general logging and operator v1)
  - [Advanced logging on Kubernetes](https://banzaicloud.com/blog/k8s-logging-advanced/)
  - [Secure logging on Kubernetes with Fluentd and Fluent Bit](https://banzaicloud.com/blog/k8s-logging-tls/)
  - [Centralized logging under Kubernetes](https://banzaicloud.com/blog/k8s-logging/)
  - [Centralized logging on Kubernetes automated](https://banzaicloud.com/blog/k8s-logging-operator/)
  - [And more...](https://banzaicloud.com/tags/logging/)


Logging-operator is a core part of the [Pipeline](https://beta.banzaicloud.io) platform, a Cloud Native application and devops platform that natively supports multi- and hybrid-cloud deployments with multiple authentication backends. Check out the developer beta:
 <p align="center">
   <a href="https://beta.banzaicloud.io">
   <img src="https://camo.githubusercontent.com/a487fb3128bcd1ef9fc1bf97ead8d6d6a442049a/68747470733a2f2f62616e7a6169636c6f75642e636f6d2f696d672f7472795f706970656c696e655f627574746f6e2e737667">
   </a>
 </p>

---

## Contents
- **[Installation](./docs/deploy/README.md)**
  - [Deploy with Helm](./docs/deploy/README.md#deploy-logging-operator-with-helm)
  - [Deploy with Kubernetes Manifests](./docs/deploy/README.md#deploy-logging-operator-from-kubernetes-manifests)
- **[Supported Plugins](#supported-plugins)**
- **[Examples](./docs)**
  - [S3 Output](./docs/example-s3.md)
  - [Elasticsearch Output](./docs/example-es-nginx.md)
  - [Loki Output](./docs/example-loki-nginx.md)
  - [Kafka Output](./docs/example-kafka-nginx.md)
  - [Amazon CloudWatch Output](./docs/example-cloudwatch-nginx.md)
  - [And more...](./docs/examples)
- **[Monitoring](./docs/logging-operator-monitoring.md)**
- **[Security](./docs/security/README.md)**
- **[Troubleshooting](#troubleshooting)**
- **[Contributing](#contributing)**
---


## Supported Plugins

For complete list of supported plugins please check the [plugins index](/docs/plugins/Readme.md).

| Name                                                    |  Type  |                                Description               | Status  | Version                                                                                 |
|---------------------------------------------------------|:------:|:--------------------------------------------------------:|---------|-----------------------------------------------------------------------------------------|
| [Alibaba](./docs/plugins/outputs/oss.md)                  | Output | Store logs the Alibaba Cloud Object Storage Service    |    GA   | [0.0.2](https://github.com/aliyun/fluent-plugin-oss)                                    |
| [Amazon S3](./docs/plugins/outputs/s3.md)                 | Output | Store logs in Amazon S3                                |    GA   | [1.2.1](https://github.com/fluent/fluent-plugin-s3/releases/tag/v1.2.1)               |
| [Azure](./docs/plugins/outputs/azurestore.md)             | Output | Store logs in Azure Storega                            |    GA   | [0.1.0](https://github.com/htgc/fluent-plugin-azurestorage/releases/tag/v0.1.0)         |
| [Google Storage](./docs/plugins/outputs/gcs.md)           | Output | Store logs in Google Cloud Storage                     |    GA   | [0.4.0](https://github.com/banzaicloud/fluent-plugin-gcs)                               |
| [Grafana Loki](./docs/plugins/outputs/loki.md)            | Output | Transfer logs to Loki                                  |    GA   | [1.2.2](https://github.com/grafana/loki/tree/master/fluentd/fluent-plugin-grafana-loki)   |
| [ElasticSearch](./docs/plugins/outputs/elasticsearch.md)  | Output | Send your logs to Elasticsearch                        |    GA   | [3.7.0](https://github.com/uken/fluent-plugin-elasticsearch/releases/tag/v3.7.0)        |
| [Sumologic](./docs/plugins/outputs/sumologic.md)          | Output | Send your logs to Sumologic                            |    GA   | [1.6.1](https://github.com/SumoLogic/fluentd-output-sumologic/releases/tag/1.6.1)       |
| [CloudWatch](./docs/plugins/outputs/cloudwatch.md)        | Output | Send your logs to AWS CloudWatch                       |    GA   | [0.7.6](https://github.com/banzaicloud/fluent-plugin-cloudwatch-logs/releases/tag/v0.7.6) |
| [Kafka](./docs/plugins/outputs/kafka.md)                  | Output | Send your logs to Kafka                                |    GA   | [0.12.1](https://github.com/fluent/fluent-plugin-kafka/releases/tag/v0.12.1)            |
| [Exception Detector](./docs/plugins/outputs/detect_exceptions.md)| Filter | Exception detector plugin for fluentd           |    GA   | [0.0.13](https://github.com/GoogleCloudPlatform/fluent-plugin-detect-exceptions/releases/tag/0.0.13) |
| [Tag Normaliser](./docs/plugins/filters/tagnormaliser.md) | Filter | Normalise tags for outputs                             |    GA   |                                                                                         |
| [Parser](./docs/plugins/filters/parser.md)                | Filter | Parse logs with parser plugin                          |    GA   |                                                                                         |
| [Multi Format Parser](./docs/plugins/outputs/detect_exceptions.md)| Filter | Exception detector plugin for fluentd           |    GA   | [1.0.0](https://github.com/repeatedly/fluent-plugin-multi-format-parser/releases/tag/v1.0.0) |

---

## Troubleshooting
:construction: The **master branch** is under heavy development. Please use [releases](https://github.com/banzaicloud/logging-operator/releases) instead of the master branch to get stable software.

If you encounter any problems that the documentation does not address, please [file an issue](https://github.com/banzaicloud/logging-operator/issues) or talk to us on the Banzai Cloud Slack channel [#logging-operator](https://slack.banzaicloud.io/).

## Contributing

If you find this project useful here's how you can help:

:rocket: Send a pull request with your new features and bug fixes
:muscle: Help new users with issues they may encounter
:star: Support the development of this project and star this repo!
:metal: If you use the operator, we would like to kindly ask you to add yourself to the list of production adopters.
For more information please read the [developer documentation](./docs/developers.md)

## License

Copyright (c) 2017-2019 [Banzai Cloud, Inc.](https://banzaicloud.com)

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

[http://www.apache.org/licenses/LICENSE-2.0](http://www.apache.org/licenses/LICENSE-2.0)

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
