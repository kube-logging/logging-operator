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

## Feature highlights

- [x] Namespace isolation 
- [x] Native Kubernetes label selectors
- [x] Secure communication (TLS)
- [x] Configuration validation
- [x] Multiple flow support (multiply logs for different transformations)
- [x] Multiple [output](docs/plugins/outputs) support (store the same logs in multiple storage: S3, GCS, ES, Loki and more...)
- [x] Multiple logging system support (multiple fluentd, fluent-bit deployment on the same cluster)

## Motivation

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

## QuickStart

This short movie shows how to get a complete logging solution on Kubernetes using Logging-Operator in less than 2 minutes:

### Install



## Documentation

 You can find the complete documentation of Logging Operator v2 at [here](./docs/Readme.md) :blue_book: <br>
:construction: The **master branch** is under heavy development. Please use [releases](https://github.com/banzaicloud/logging-operator/releases) instead of the master branch to get stable software.

If you encounter any problems that the documentation does not address, please [file an issue](https://github.com/banzaicloud/logging-operator/issues) or talk to us on the Banzai Cloud Slack channel [#logging-operator](https://slack.banzaicloud.io/).


## Contributing

If you find this project useful here's how you can help:

- Send a pull request with your new features and bug fixes :rocket: 
- Help new users with issues they may encounter :muscle:
- Support the development of this project and star this repo! :star:
- If you use the operator, we would like to kindly ask you to add yourself to the list of production [adopters](https://github.com/banzaicloud/logging-operator/blob/master/ADOPTERS.md).:metal: <br> 

*For more information please read the [developer documentation](./docs/developers.md)*

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
