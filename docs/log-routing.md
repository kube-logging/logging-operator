#

The first step to process your logs is to select what logs goes to where.

#### Flow

#### Cluster flow

Select in

> Note: Flows are namespace scoped

1. Select logs with `app: nginx` labels from the namespace

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

2. Exclude logs with `app: nginx` labels from the namespace
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

3. Exclude logs with `env: dev` labels but select `app: nginx` labels from the namespace
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

4. Exclude cluster logs from  `dev`, `sandbox` namespaces and select `app: nginx` from all namespaces
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


5. Exclude cluster logs from  `dev`, `sandbox` namespaces and select `app: nginx` from all `prod` and `infra` namespaces
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