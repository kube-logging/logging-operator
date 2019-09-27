<p align="center"><img src="./img/les.png" width="240"></p>

# Save all logs to ElasticSearch


## Deploy ElasticSearch

### Add chart repository:
```bash
helm repo add es-operator https://raw.githubusercontent.com/upmc-enterprises/elasticsearch-operator/master/charts/
helm repo update
```

### Install ElasticSearch with operator
```bash
helm install --namespace logging --name elasticsearch-operator es-operator/elasticsearch-operator --set rbac.enabled=True
helm install --namespace logging --name elasticsearch es-operator/elasticsearch --set kibana.enabled=True --set cerebro.enabled=True
```
> [Elasticsearch Operator Documentation](https://github.com/upmc-enterprises/elasticsearch-operator)

#### Forward cerebro & kibana dashboards
```bash
kubectl -n logging port-forward svc/cerebro-elasticsearch-cluster 9001:80
kubectl -n logging port-forward svc/kibana-elasticsearch-cluster 5601:80
```


## Deploy Logging-Operator

### Install with Helm 
[Install Logging-operator with helm](./deploy/README.md#deploy-logging-operator-with-helm)


### Install from Kubernetes manifests
[Install Logging-operator from manifests](./deploy/README.md#deploy-logging-operator-from-kubernetes-manifests)

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


#### Create an ElasticSearch output definition 
```bash
cat <<EOF | kubectl -n logging apply -f -
apiVersion: logging.banzaicloud.io/v1beta1
kind: ClusterOutput
metadata:
  name: es-output
  namespace: logging
spec:
  elasticsearch:
    host: elasticsearch-elasticsearch-cluster.default.svc.cluster.local
    port: 9200
    scheme: https
    ssl_verify: false
    ssl_version: TLSv1_2
    buffer:
      path: /tmp/buffer
      timekey: 1m
      timekey_wait: 30s
      timekey_use_utc: true
EOF
```

> Note: For production set-up we recommend using longer `timekey` interval to avoid generating too many object.

The following snippet will use [tag_normaliser](./plugins/filters/tagnormaliser.md) to re-tag logs and after push it to ElasticSearch.

#### Create `flow` resource
```bash
cat <<EOF | kubectl -n logging apply -f -
apiVersion: logging.banzaicloud.io/v1beta1
kind: ClusterFlow
metadata:
  name: es-flow
  namespace: logging
spec:
  filters:
    - tag_normaliser: {}
  selectors: {}
  outputRefs:
    - es-output
EOF
```
