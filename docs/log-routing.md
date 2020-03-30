---
title: Routing your logs with match directive
shorttitle: Log routing
weight: 700
---

{{< contents >}}

The first step to process your logs is to select what logs goes to where.
The logging operator uses Kubernetes labels and namespaces to separate
different log flows.

## Match statement

To select or exclude logs you can use the `match` statement. Match is a collection
of `select` and `exclude` expressions. In both expression you can use the `labels`
attribute to filter for pod's labels. Moreover in Cluster flow you can use `namespaces`
as a selecting or excluding criteria.

The list of `select` and `exclude` statements are evaluated **in order**!

Flow:

```yaml
  kind: Flow
  metadata:
    name: flow-sample
  spec:
    match:
      - exclude:
          labels:
            exclude-this: label
      - select:
          labels:
            app: nginx
            label/xxx: example
```

ClusterFlow:

```yaml
  kind: ClusterFlow
  metadata:
    name: flow-sample
  spec:
    match:
      - exclude:
          labels:
            exclude-this: label
          namespaces:
            - developer 
      - select:
          labels:
            app: nginx
            label/xxx: example
          namespaces:
            - production
            - beta
```

## Examples

### Example 1. Select logs by label

Select logs with `app: nginx` labels from the namespace:

  ```yaml
  apiVersion: logging.banzaicloud.io/v1beta1
  kind: Flow
  metadata:
    name: flow-sample
    namespace: default
  spec:
    outputRefs:
      - forward-output-sample
    match:
      - select:
          labels:
            app: nginx
  ```

### Example 2. Exclude logs by label

Exclude logs with `app: nginx` labels from the namespace

  ```yaml
  apiVersion: logging.banzaicloud.io/v1beta1
  kind: Flow
  metadata:
    name: flow-sample
    namespace: default
  spec:
    outputRefs:
      - forward-output-sample
    match:
      - exclude:
          labels:
            app: nginx
  ```

### Example 3. Exclude and select logs by label

Exclude logs with `env: dev` labels but select `app: nginx` labels from the namespace

  ```yaml
  apiVersion: logging.banzaicloud.io/v1beta1
  kind: Flow
  metadata:
    name: flow-sample
    namespace: default
  spec:
    outputRefs:
      - forward-output-sample
    match:
      - exclude:
          labels:
            env: dev
      - select:
          labels:
            app: nginx
  ```

### Example 4. Exclude cluster logs by namespace

Exclude cluster logs from  `dev`, `sandbox` namespaces and select `app: nginx` from all namespaces

  ```yaml
  apiVersion: logging.banzaicloud.io/v1beta1
  kind: ClusterFlow
  metadata:
    name: clusterflow-sample
  spec:
    outputRefs:
      - forward-output-sample
    match:
      - exclude:
          namespaces:
            - dev
            - sandbox
      - select:
          labels:
            app: nginx
  ```

### Example 5. Exclude and select cluster logs by namespace

Exclude cluster logs from  `dev`, `sandbox` namespaces and select `app: nginx` from all `prod` and `infra` namespaces

  ```yaml
  apiVersion: logging.banzaicloud.io/v1beta1
  kind: ClusterFlow
  metadata:
    name: clusterflow-sample
  spec:
    outputRefs:
      - forward-output-sample
    match:
      - exclude:
          namespaces:
            - dev
            - sandbox
      - select:
          labels:
            app: nginx
          namespaces:
            - prod
            - infra
  ```
