<p align="center"><img src="../img/lo.svg" width="260"></p>
<p align="center">

# Deploy logging-operator from Kubernetes Manifests

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
1. Deploy the Operator
    ```bash
    kubectl -n logging create -f ./docs/deploy/manifests/deployment.yaml
    ```
---

# Deploy logging-operator with Helm

<p align="center"><img src="../img/helm.svg" width="150"></p>
<p align="center">

1. Add operator chart repository.
    ```bash
    helm repo add banzaicloud-stable https://kubernetes-charts.banzaicloud.com
    helm repo update
    ```
1. Install the Logging Operator
    ```bash
    helm install --namespace logging --name logging banzaicloud-stable/logging-operator
    ```
    > You can install the `logging` resource with built-in TLS generation using a [Helm chart](/charts/logging-operator-logging).

---

# Check logging-operator deployment

### Check pods status

```bash
$ kubectl -n logging get pods
NAME                                        READY   STATUS    RESTARTS   AGE
logging-logging-operator-599c9cf846-5nw2n   1/1     Running   0          52s
```

### Check CRD 
```bash
$  kubectl get crd
NAME                                    CREATED AT
clusterflows.logging.banzaicloud.io     2019-11-01T21:30:18Z
clusteroutputs.logging.banzaicloud.io   2019-11-01T21:30:18Z
flows.logging.banzaicloud.io            2019-11-01T21:30:18Z
loggings.logging.banzaicloud.io         2019-11-01T21:30:18Z
outputs.logging.banzaicloud.io          2019-11-01T21:30:18Z
```