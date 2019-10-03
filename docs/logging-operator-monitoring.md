<p align="center"><img src="./img/nll.png" width="340"></p>

# Monitor your logging pipeline with Prometheus Operator

<p align="center"><img src="./img/nginx-loki.png" width="900"></p>

### Create `logging` namespace
```bash
kubectl create namespace logging
```


> [Prometheus Operator Documentation](https://github.com/coreos/prometheus-operator)
### Install Prometheus Operator with Helm
```bash
helm install --namespace logging stable/prometheus-operator 
```


## Install with Helm 

### Add operator chart repository:
```bash
helm repo add banzaicloud-stable https://kubernetes-charts.banzaicloud.com
```

### Logging Operator
```bash
helm install --namespace logging --name logging banzaicloud-stable/logging-operator
```
> You can install `logging` resource via [Helm chart](/charts/logging-operator-logging) with built-in TLS generation.

### Demo Nginx App + Logging Definition with metrics
```bash
helm install --namespace logging --name nginx-demo banzaicloud-stable/nginx-logging-demo \
    --set=loggingOperator.fluentd.metrics.enabled=True \
    --set=loggingOperator.fluentbit.metrics.enabled=True
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
  fluentd:
    metics: {}
  fluentbit:
    metics: {}
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
    buffer:
      path: /tmp/buffer
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
    - tag_normaliser: {}
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
        image: nginx:latest
        ports:
        - containerPort: 80
          name: http
          protocol: TCP
        livenessProbe:
          failureThreshold: 3
          httpGet:
            path: /
            port: http
            scheme: HTTP
          periodSeconds: 10
          successThreshold: 1
          timeoutSeconds: 1
        readinessProbe:
          failureThreshold: 3
          httpGet:
            path: /
            port: http
            scheme: HTTP
          periodSeconds: 10
          successThreshold: 1
          timeoutSeconds: 1
EOF
```


#### Forward Prometheus Service
```bash
kubectl -n logging port-forward svc/grafana 3000:80
```
[Gradana Dashboard: http://localhost:3000](http://localhost:3000)
<p align="center"><img src="./img/loki1.png" width="660"></p>



### Grafana Dashboard

#### Get Minio login credantials
```bash
kubectl -n logging get secrets logging-s3 -o json | jq '.data | map_values(@base64d)'
```

#### Forward Minio Service
```bash
kubectl -n logging port-forward svc/grafana 3000:80
```
[Gradana Dashboard: http://localhost:3000](http://localhost:3000)
<p align="center"><img src="./img/loki1.png" width="660"></p>


