<p align="center"><img src="./img/lo.svg" width="260"></p>

# Routing your logs with match directive

---

The first step to process your logs is to select what logs goes to where.
The logging operator uses Kubernetes labels and namespaces to separate
different log flows.

#### Match statement

To select or exclude logs you can use the `match` statement. Match is a collection
of `select` and `exclude` expressions. In both expression you can use the `labels`
attribute to filter for pod's labels. Moreover in Cluster flow you can use `namespaces`
as a selecting or excluding criteria.

The list of `select` and `exclude` statements are evaluated **in order**!

Flow
```
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

ClusterFlow
```
  kind: ClusterFlow
  metadata:
    name: flow-sample
  spec:
    match:
      - exclude:
          labels:
            exclude-this: label
          namespaces: developer 
      - select:
          labels:
            app: nginx
            label/xxx: example
          namespaces: production,beta
```

## Examples

#### 1. Select logs with `app: nginx` labels from the namespace

  ```
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

#### 2. Exclude logs with `app: nginx` labels from the namespace
  ```
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

#### 3. Exclude logs with `env: dev` labels but select `app: nginx` labels from the namespace
  ```
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

#### 4. Exclude cluster logs from  `dev`, `sandbox` namespaces and select `app: nginx` from all namespaces
  ```
  apiVersion: logging.banzaicloud.io/v1beta1
  kind: ClusterFlow
  metadata:
    name: clusterflow-sample
  spec:
    outputRefs:
      - forward-output-sample
    match:
      - exclude:
          namespaces: dev,sandbox
      - select:
          labels:
            app: nginx
  ```


#### 5. Exclude cluster logs from  `dev`, `sandbox` namespaces and select `app: nginx` from all `prod` and `infra` namespaces
  ```
  apiVersion: logging.banzaicloud.io/v1beta1
  kind: ClusterFlow
  metadata:
    name: clusterflow-sample
  spec:
    outputRefs:
      - forward-output-sample
    match:
      - exclude:
          namespaces: dev,sandbox
      - select:
          labels:
            app: nginx
          namespaces: prod,infra
  ```