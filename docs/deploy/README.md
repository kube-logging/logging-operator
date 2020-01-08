<p align="center"><img src="../img/lo.svg" width="260"></p>
<p align="center">

# Requirements

- Logging operator requires Kubernetes v1.14.x or later.
- For the [Helm base installation](#deploy-logging-operator-with-helm) you need Helm v2.16.0 or later.

# Deploy the Logging operator from Kubernetes Manifests

Complete the following steps to deploy the Logging operator using Kubernetes manifests. Alternatively, you can also [install the operator using Helm](#deploy-logging-operator-with-helm).

1. Clone the logging-operator repo.
    ```bash
    git clone git@github.com:banzaicloud/logging-operator.git
    ```
1. Navigate to the `logging-operator` folder.
    ```bash
    cd logging-operator
    ```
1. Create a controlNamespace called “logging”.
    ```bash
    kubectl create ns logging
    ```
1. Create a ServiceAccount and install cluster roles.
    ```bash
    kubectl -n logging create -f ./docs/deploy/manifests/rbac.yaml
    ```
1. Apply the ClusterResources.
    ```bash
    kubectl -n logging create -f ./config/crd/bases/logging.banzaicloud.io_clusterflows.yaml
    kubectl -n logging create -f ./config/crd/bases/logging.banzaicloud.io_clusteroutputs.yaml
    kubectl -n logging create -f ./config/crd/bases/logging.banzaicloud.io_flows.yaml
    kubectl -n logging create -f ./config/crd/bases/logging.banzaicloud.io_loggings.yaml
    kubectl -n logging create -f ./config/crd/bases/logging.banzaicloud.io_outputs.yaml
    ```
1. Deploy the Logging operator.
    ```bash
    kubectl -n logging create -f ./docs/deploy/manifests/deployment.yaml
    ```
---

# Deploy Logging operator with Helm

<p align="center"><img src="../img/helm.svg" width="150"></p>
<p align="center">

Complete the following steps to deploy the Logging operator using Helm. Alternatively, you can also [install the operator using Kubernetes manifests](./Readme.md).
> Note: For the [Helm base installation](#deploy-logging-operator-with-helm) you need Helm v2.16.0 or later.

1. Add operator chart repository.
    ```bash
    helm repo add banzaicloud-stable https://kubernetes-charts.banzaicloud.com
    helm repo update
    ```
2. Install the Logging Operator
    ```bash
    helm install --namespace logging --name logging banzaicloud-stable/logging-operator
    ```
    > You can install the `logging` resource with built-in TLS generation using a [Helm chart](/charts/logging-operator-logging).

---

# Check the Logging operator deployment

To verify that the installation was successful, complete the following steps.

1. Check the status of the pods. You should see a new logging-operator pod.
    ```bash
    $ kubectl -n logging get pods
    NAME                                        READY   STATUS    RESTARTS   AGE
    logging-logging-operator-599c9cf846-5nw2n   1/1     Running   0          52s
    ```
1. Check the CRDs. You should see the following five new CRDs.
    ```bash
    $  kubectl get crd
    NAME                                    CREATED AT
    clusterflows.logging.banzaicloud.io     2019-11-01T21:30:18Z
    clusteroutputs.logging.banzaicloud.io   2019-11-01T21:30:18Z
    flows.logging.banzaicloud.io            2019-11-01T21:30:18Z
    loggings.logging.banzaicloud.io         2019-11-01T21:30:18Z
    outputs.logging.banzaicloud.io          2019-11-01T21:30:18Z
    ```