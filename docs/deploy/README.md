<p align="center"><img src="../img/lo.svg" width="260"></p>
<p align="center">

# Deploy logging-operator from Kubernetes Manifests

### Clone the logging-operator repo
```bash
git clone git@github.com:banzaicloud/logging-operator.git
```

### Navigate to the nginx-data-www folder 
```bash
cd logging-operator
```

### Now, let’s create a Namespace called “logging” to work in: 
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

