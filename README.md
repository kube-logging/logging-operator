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

# Logging operator

Logging operator for Kubernetes based on Fluentd and Fluent-bit.

The Logging operator automates the deployment and configuration of a Kubernetes logging pipeline. The operator configures a fluent-bit daemonset for collecting container logs from the node file system. Fluent-bit enriches the logs with Kubernetes metadata and transfers them to fluentd. Fluentd receives, filters, and transfer logs to multiple outputs. Your logs will always be transferred on authenticated and encrypted channels.

## What is this operator for?

This operator helps you bundle logging information with your applications: you can describe the behavior of your application in its charts, the Logging operator does the rest.

<p align="center"><img src="docs/img/logging_operator_flow.png" ></p>

## Feature highlights

- [x] Namespace isolation
- [x] Native Kubernetes label selectors
- [x] Secure communication (TLS)
- [x] Configuration validation
- [x] Multiple flow support (multiply logs for different transformations)
- [x] Multiple [output](docs/plugins/outputs) support (store the same logs in multiple storage: S3, GCS, ES, Loki and more...)
- [x] Multiple logging system support (multiple fluentd, fluent-bit deployment on the same cluster)

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

Follow these [quickstart guides](docs/quickstarts) to try out Logging Operator!

### Install

Deploy Logging Operator with [Kubernetes manifests](/docs/deploy/README.md) or [Helm chart](/docs/deploy/README.md#Deploy-logging-operator-with-Helm).

> Caution: The **master branch** is under heavy development. Use [releases](https://github.com/banzaicloud/logging-operator/releases) instead of the master branch to get stable software.

If you encounter any problems that the documentation does not address, [open an issue](https://github.com/banzaicloud/logging-operator/issues) or talk to us on the Banzai Cloud Slack channel [#logging-operator](https://slack.banzaicloud.io/).

## Documentation

 You can find the complete documentation of Logging operator v2 [here](./docs/Readme.md) :blue_book: <br>

## Contributing

If you find this project useful, help us:

- Support the development of this project and star this repo! :star:
- If you use the Logging operator in a production environment, add yourself to the list of production [adopters](https://github.com/banzaicloud/logging-operator/blob/master/ADOPTERS.md).:metal: <br> 
- Help new users with issues they may encounter :muscle:
- Send a pull request with your new features and bug fixes :rocket: 

*For more information, read the [developer documentation](./docs/developers.md)*.
