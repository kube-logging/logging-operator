<p align="center"><img src="./img/nll.png" width="340"></p>

# Store Nginx Access Logs in Grafana Loki with Logging Operator

<p align="center"><img src="./img/nginx-loki.png" width="900"></p>

---
## Contents
- **Installation**
  - Loki
    - [Deploy with Helm](#add-operator-chart-repository)
  - **Logging Operator**
    - [Deploy with Helm](#install-with-helm)
    - [Deploy with Kubernetes Manifests](./deploy/README.md#deploy-logging-operator-from-kubernetes-manifests)
   - **Demo Application**  
    - [Deploy with Helm](#nginx-app--logging-definition)
    - [Deploy with Kubernetes Manifests](#install-from-manifest)
- **Validation**
    - [Grafana Dashboard](#grafana-dashboard)
---

### Add operator chart repository:
```bash
helm repo add loki https://grafana.github.io/loki/charts
helm repo add banzaicloud-stable https://kubernetes-charts.banzaicloud.com
helm repo update
```

### Install Loki
```bash
helm install --namespace logging --name loki loki/loki
```
> [Grafana Loki Documentation](https://github.com/grafana/loki/tree/master/production/helm)
### Install Grafana
```bash
helm install --namespace logging --name grafana stable/grafana \
 --set "datasources.datasources\\.yaml.apiVersion=1" \
 --set "datasources.datasources\\.yaml.datasources[0].name=Loki" \
 --set "datasources.datasources\\.yaml.datasources[0].type=loki" \
 --set "datasources.datasources\\.yaml.datasources[0].url=http://loki:3100" \
 --set "datasources.datasources\\.yaml.datasources[0].access=proxy"
```


## Install with Helm 
### Logging Operator
```bash
helm install --namespace logging --name logging banzaicloud-stable/logging-operator
```
> You can install `logging` resource via [Helm chart](/charts/logging-operator-logging) with built-in TLS generation.

### Nginx App + Logging Definition
```bash
helm install --namespace logging --name nginx-demo banzaicloud-stable/nginx-logging-loki-demo
```

## Install from manifest

#### Create `logging` resource
```bash
cat <<EOF | kubectl -n logging apply -f -
apiVersion: logging.banzaicloud.io/v1beta1
kind: Logging
metadata:
  name: default-logging-simple
spec:
  fluentd: {}
  fluentbit: {}
  controlNamespace: logging
EOF
```
> Note: `ClusterOutput` and `ClusterFlow` resource will only be accepted in the `controlNamespace` 


#### Create an Loki output definition 
```bash
cat <<EOF | kubectl -n logging apply -f -
apiVersion: logging.banzaicloud.io/v1beta1
kind: Output
metadata:
  name: loki-output
spec:
  loki:
    url: http://loki:3100
    configure_kubernetes_labels: true
    buffer:
      timekey: 1m
      timekey_wait: 30s
      timekey_use_utc: true
EOF
```
> Note: For production set-up we recommend using longer `timekey` interval to avoid generating too many object.

#### Create `flow` resource
```bash
cat <<EOF | kubectl -n logging apply -f -
apiVersion: logging.banzaicloud.io/v1beta1
kind: Flow
metadata:
  name: loki-flow
spec:
  filters:
    - parser:
        key_name: message
        remove_key_name_field: true
        reserve_data: true
        parsers:
          - type: nginx
  selectors:
    app: nginx
  outputRefs:
    - loki-output
EOF
```

#### Install nginx deployment
```bash
cat <<EOF | kubectl -n logging apply -f -
apiVersion: apps/v1 
kind: Deployment
metadata:
  name: nginx-deployment
spec:
  selector:
    matchLabels:
      app: nginx
  replicas: 1
  template:
    metadata:
      labels:
        app: nginx
    spec:
      containers:
      - name: nginx
        image: banzaicloud/loggen:latest
EOF
```

### Grafana Dashboard

#### Get Grafana login credantials
```bash
kubectl -n logging get secrets grafana -o json | jq '.data | map_values(@base64d)'
```

#### Forward Grafana Service
```bash
kubectl -n logging port-forward svc/grafana 3000:80
```
[Gradana Dashboard: http://localhost:3000](http://localhost:3000)
<p align="center"><img src="./img/loki1.png" width="660"></p>


