<p align="center"><img src="https://raw.githubusercontent.com/cncf/landscape/master/hosted_logos/logging-operator.svg" width="260"></p>
<p align="center">

  <a href="https://goreportcard.com/badge/github.com/kube-logging/logging-operator">
    <img src="https://goreportcard.com/badge/github.com/kube-logging/logging-operator" alt="Go Report Card">
  </a>

  <a href="https://github.com/kube-logging/logging-operator/">
    <img src="https://img.shields.io/badge/license-Apache%20v2-orange.svg" alt="license">
  </a>

  <a href="https://www.bestpractices.dev/projects/7941">
    <img src="https://www.bestpractices.dev/projects/7941/badge?cache=20231206">
  </a>

</p>

# Logging operator

The Logging Operator is now a [CNCF Sandbox](https://www.cncf.io/sandbox-projects/) project.

The Logging operator solves your logging-related problems in Kubernetes environments by automating the deployment and configuration of a Kubernetes logging pipeline.

1. The operator deploys and configures a log collector (currently a Fluent Bit DaemonSet) on every node to collect container and application logs from the node file system.
1. Fluent Bit queries the Kubernetes API and enriches the logs with metadata about the pods, and transfers both the logs and the metadata to a log forwarder instance.
1. The log forwarder instance receives, filters, and transforms the incoming the logs, and transfers them to one or more destination outputs. The Logging operator supports Fluentd and syslog-ng as log forwarders.

Your logs are always transferred on authenticated and encrypted channels.

This operator helps you bundle logging information with your applications: you can describe the behavior of your application in its charts, the Logging operator does the rest.

## What is this operator for?

This operator helps you bundle logging information with your applications: you can describe the behavior of your application in its charts, the Logging operator does the rest.

<p align="center"><img src="https://kube-logging.github.io/docs/img/logging_operator_flow.png" ></p>

## Feature highlights

- [x] Namespace isolation
- [x] Native Kubernetes label selectors
- [x] Secure communication (TLS)
- [x] Configuration validation
- [x] Multiple flow support (multiply logs for different transformations)
- [x] Multiple output support (store the same logs in multiple storage: S3, GCS, ES, Loki and more...)
- [x] Multiple logging system support (multiple fluentd, fluent-bit deployment on the same cluster)

## Architecture

The Logging operator manages the log collectors and log forwarders of your logging infrastructure, and the routing rules that specify where you want to send your different log messages.

The **log collectors** are endpoint agents that collect the logs of your Kubernetes nodes and send them to the log forwarders. Logging operator currently uses Fluent Bit as log collector agents.

The **log forwarder** instance receives, filters, and transforms the incoming the logs, and transfers them to one or more destination outputs. The Logging operator supports Fluentd and syslog-ng as log forwarders. Which log forwarder is best for you depends on your logging requirements.

You can filter and process the incoming log messages using the **flow** custom resource of the log forwarder to route them to the appropriate **output**. The outputs are the destinations where you want to send your log messages, for example, Elasticsearch, or an Amazon S3 bucket. You can also define cluster-wide outputs and flows, for example, to use a centralized output that namespaced users can reference but cannot modify. Note that flows and outputs are specific to the type of log forwarder you use (Fluentd or syslog-ng).

You can configure the Logging operator using the following Custom Resource Definitions.

- [Logging](https://kube-logging.github.io/docs/logging-infrastructure/logging/) - The `Logging` resource defines the logging infrastructure (the log collectors and forwarders) for your cluster that collects and transports your log messages. It also contains configurations for Fluent Bit, Fluentd, and syslog-ng.
- CRDs for Fluentd:
    - [Output](https://kube-logging.github.io/docs/configuration/output/) - Defines a Fluentd Output for a logging flow, where the log messages are sent using Fluentd. This is a namespaced resource. See also `ClusterOutput`. To configure syslog-ng outputs, see `SyslogNGOutput`.
    - [Flow](https://kube-logging.github.io/docs/configuration/flow/) - Defines a Fluentd logging flow using `filters` and `outputs`. Basically, the flow routes the selected log messages to the specified outputs. This is a namespaced resource. See also `ClusterFlow`. To configure syslog-ng flows, see `SyslogNGFlow`.
    - [ClusterOutput](https://kube-logging.github.io/docs/configuration/output/) - Defines a Fluentd output that is available from all flows and clusterflows. The operator evaluates clusteroutputs in the `controlNamespace` only unless `allowClusterResourcesFromAllNamespaces` is set to true.
    - [ClusterFlow](https://kube-logging.github.io/docs/configuration/flow/) - Defines a Fluentd logging flow that collects logs from all namespaces by default. The operator evaluates clusterflows in the `controlNamespace` only unless `allowClusterResourcesFromAllNamespaces` is set to true. To configure syslog-ng clusterflows, see `SyslogNGClusterFlow`.
- CRDs for syslog-ng (these resources like their Fluentd counterparts, but are tailored to features available via syslog-ng):
    - [SyslogNGOutput](https://kube-logging.github.io/docs/configuration/output/#syslogngoutput) - Defines a syslog-ng Output for a logging flow, where the log messages are sent using Fluentd. This is a namespaced resource. See also `SyslogNGClusterOutput`. To configure Fluentd outputs, see `output`.
    - [SyslogNGFlow](https://kube-logging.github.io/docs/configuration/flow/#syslogngflow) - Defines a syslog-ng logging flow using `filters` and `outputs`. Basically, the flow routes the selected log messages to the specified outputs. This is a namespaced resource. See also `SyslogNGClusterFlow`. To configure Fluentd flows, see `flow`.
    - [SyslogNGClusterOutput](https://kube-logging.github.io/docs/configuration/output/#syslogngoutput) - Defines a syslog-ng output that is available from all flows and clusterflows. The operator evaluates clusteroutputs in the `controlNamespace` only unless `allowClusterResourcesFromAllNamespaces` is set to true.
    - [SyslogNGClusterFlow](https://kube-logging.github.io/docs/configuration/flow/#syslogngflow) - Defines a syslog-ng logging flow that collects logs from all namespaces by default. The operator evaluates clusterflows in the `controlNamespace` only unless `allowClusterResourcesFromAllNamespaces` is set to true. To configure Fluentd clusterflows, see `clusterflow`.

See the [detailed CRDs documentation](https://kube-logging.github.io/docs/configuration/crds/).

<p align="center"><img src="https://kube-logging.github.io/docs/img/logging-operator-v2-architecture.png" ></p>

## Quickstart

[![asciicast](https://asciinema.org/a/315998.svg)](https://asciinema.org/a/315998)

Follow these [quickstart guides](https://kube-logging.github.io/docs/quickstarts/) to try out the Logging operator!

### Install

Deploy Logging Operator with our [Helm chart](https://kube-logging.github.io/docs/install/#deploy-logging-operator-with-helm).

> Caution: The **master branch** is under heavy development. Use [releases](https://github.com/kube-logging/logging-operator/releases) instead of the master branch to get stable software.

## Support

If you encounter problems while using the Logging operator the documentation does not address, [open an issue](https://github.com/kube-logging/logging-operator/issues) or talk to us on the [#logging-operator Discord channel](https://discord.gg/eAcqmAVU2u).

## Documentation

 You can find the complete documentation on the [Logging operator documentation page](https://kube-logging.github.io/docs/) :blue_book: <br>

## Contributing

If you find this project useful, help us:

- Support the development of this project and star this repo! :star:
- If you use the Logging operator in a production environment, add yourself to the list of production [adopters](https://github.com/kube-logging/logging-operator/blob/master/ADOPTERS.md).:metal: <br> 
- Help new users with issues they may encounter :muscle:
- Send a pull request with your new features and bug fixes :rocket: 

Please read the [Organisation's Code of Conduct](https://github.com/kube-logging/.github/blob/main/CODE_OF_CONDUCT.md)!

*For more information, read the [developer documentation](https://kube-logging.github.io/docs/developers)*.

## License

Copyright (c) 2021-2023 [Cisco Systems, Inc. and its affiliates](https://cisco.com)
Copyright (c) 2017-2020 [Banzai Cloud, Inc.](https://banzaicloud.com)

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

[http://www.apache.org/licenses/LICENSE-2.0](http://www.apache.org/licenses/LICENSE-2.0)

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
