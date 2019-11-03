<p align="center"><img src="../img/lo.svg" width="260"></p>
<p align="center">

# Deploy logging-operator from Kubernetes Manifests

### Clone the logging-operator repo
```bash
git clone git@github.com:banzaicloud/logging-operator.git
```

### Navigate to the `logging-operator` folder 
```bash
cd logging-operator
```

### Create a controlNamespace named “logging”: 
```bash
kubectl create ns logging
```

### Create ServiceAccount and install cluster roles
```bash
kubectl -n logging create -f ./docs/deploy/manifests/rbac.yaml
```

### Apply the ClusterResources
```bash
kubectl -n logging create -f ./config/crd/bases/logging.banzaicloud.io_clusterflows.yaml
kubectl -n logging create -f ./config/crd/bases/logging.banzaicloud.io_clusteroutputs.yaml
kubectl -n logging create -f ./config/crd/bases/logging.banzaicloud.io_flows.yaml
kubectl -n logging create -f ./config/crd/bases/logging.banzaicloud.io_loggings.yaml
kubectl -n logging create -f ./config/crd/bases/logging.banzaicloud.io_outputs.yaml
```

### Deploy the Operator
```bash
kubectl -n logging create -f ./docs/deploy/manifests/deployment.yaml
```

---
<br />
<p align="center"><img src="../img/helm.svg" width="150"></p>
<p align="center">

# Deploy logging-operator with Helm

### Add operator chart repository:
```bash
helm repo add banzaicloud-stable https://kubernetes-charts.banzaicloud.com
helm repo update
```

### Logging Operator
```bash
helm install --namespace logging --name logging banzaicloud-stable/logging-operator
```
> You can install `logging` resource via [Helm chart](/charts/logging-operator-logging) with built-in TLS generation.


---
<br />

# Check logging-operator deployment

### Pods Status

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

