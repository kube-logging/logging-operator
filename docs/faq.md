---
title: Logging Operator frequently asked questions
shorttitle: FAQ
weight: 900
---

{{< contents >}}

## How can I run the unreleased master version?

1. Clone the logging-operator repo.

    ```bash
    git clone git@github.com:banzaicloud/logging-operator.git
    ```

1. Navigate to the `logging-operator` folder.

    ```bash
    cd logging-operator
    ```

1. Install with helm

    - Helm v2

        ```bash
        helm install --namespace logging --name logging ./charts/logging-operator --set image.tag=master
        ```

    - Helm v3

        ```bash
        helm install --namespace logging --name logging ./charts/logging-operator --set createCustomResource=false --set image.tag=master

## How can I support the project?

- Give a star to this repository :star:
- Add your company to the [adopters](https://github.com/banzaicloud/logging-operator/blob/master/ADOPTERS.md) list
[![IMAGE ALT TEXT HERE](http://img.youtube.com/vi/2iaK8adpwfk/0.jpg)](http://www.youtube.com/watch?v=2iaK8adpwfk)
