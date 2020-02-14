<p align="center"><img src="./img/lo-pro.png" width="340"></p>

# Monitor your logging pipeline with Prometheus Operator

<p align="center"><img src="./img/monitor.png" width="900"></p>

---
## Contents
- **Installation**
  - Prometheus Operator
    - [Deploy with Helm](#install-prometheus-operator-with-helm)
  - **Logging Operator**
    - [Deploy with Helm](#install-with-helm)
    - [Deploy with Kubernetes Manifests](./deploy/README.md#deploy-logging-operator-from-kubernetes-manifests)
  - **Minio**  
    - [Deploy with Kubernetes Manifests](#install-minio)
   - **Demo Application**  
    - [Deploy with Helm](#deploy-demo-nginx-app--logging-definition-with-metrics)
    - [Deploy with Kubernetes Manifests](#install-from-manifest)
- **Validation**
    - [Prometheus Dashboard](#prometheus)
    - [Minio Dashboard](#minio)
    - [Grafana Dashboard](#grafana)
---


### Metrics Variables
| Variable Name | Type | Required | Default | Description |
|---|---|---|---|---|
| interval | string | No | "15s" | Scrape Interval |
| timeout | string | No | "5s" |  Scrape Timeout |
| port | int | No | - | Metrics Port. |
| path | int | No | - | Metrics Path. |
| serviceMonitor | bool | No | false | Enable to create ServiceMonitor for Prometheus operator |
| prometheusAnnotations | bool | No | false | Add prometheus labes to fluent pods. |


## Install Prometheus Operator with Helm 

### Create `logging` namespace
```bash
kubectl create namespace logging
```
### Install Prometheus Operator
```bash
helm install --namespace logging --name monitor stable/prometheus-operator \
    --set "grafana.dashboardProviders.dashboardproviders\\.yaml.apiVersion=1" \
    --set "grafana.dashboardProviders.dashboardproviders\\.yaml.providers[0].orgId=1" \
    --set "grafana.dashboardProviders.dashboardproviders\\.yaml.providers[0].type=file" \
    --set "grafana.dashboardProviders.dashboardproviders\\.yaml.providers[0].disableDeletion=false" \
    --set "grafana.dashboardProviders.dashboardproviders\\.yaml.providers[0].options.path=/var/lib/grafana/dashboards/default" \
    --set "grafana.dashboards.default.logging.gnetId=7752" \
    --set "grafana.dashboards.default.logging.revision=3" \
    --set "grafana.dashboards.default.logging.datasource=Prometheus" \
    --set "prometheus.prometheusSpec.serviceMonitorSelectorNilUsesHelmValues=False"
```
> [Prometheus Operator Documentation](https://github.com/coreos/prometheus-operator)

> The prometheus-operator install may take a few more minutes. ***Please be patient.*** 
> The logging-operator metrics function depends on the prometheus-operator's resources.
> If those do not exist in the cluster it may cause the logging-operator's malfunction.

---
## Install with Helm 

### Add operator chart repository:
```bash
helm repo add banzaicloud-stable https://kubernetes-charts.banzaicloud.com
```

### Logging Operator
```bash
helm install --namespace logging --name logging banzaicloud-stable/logging-operator
```
> You also can install logging-operator from manifest [guideline is here](./deploy/README.md#deploy-logging-operator-from-kubernetes-manifests)

### Deploy Demo App + Logging Definition with metrics
```bash
helm install --namespace logging --name logging-demo banzaicloud-stable/logging-demo \
    --set "minio.enabled=True" \
    --set=loggingOperator.fluentd.metrics.serviceMonitor=True \
    --set=loggingOperator.fluentbit.metrics.serviceMonitor=True
```
---
## Install from manifest

### Install Minio

#### Create Minio Credential Secret
```bash
kubectl -n logging create secret generic logging-s3 --from-literal=accesskey='AKIAIOSFODNN7EXAMPLE' --from-literal=secretkey='wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY'
```
#### Deploy Minio
```bash
kubectl -n logging apply -f - <<"EOF" 
apiVersion: apps/v1
kind: Deployment
metadata:
  name: minio-deployment
  namespace: logging
spec:
  strategy:
    type: Recreate
  template:
    metadata:
      labels:
        app: minio
    spec:
      containers:
      - name: minio
        image: minio/minio
        args:
        - server
        - /storage
        readinessProbe:
          httpGet:
            path: /minio/health/ready
            port: 9000
          initialDelaySeconds: 10
          periodSeconds: 5
        env:
        - name: MINIO_REGION
          value: 'test_region'
        - name: MINIO_ACCESS_KEY
          valueFrom:
            secretKeyRef:
              name: logging-s3
              key: accesskey
        - name: MINIO_SECRET_KEY
          valueFrom:
            secretKeyRef:
              name: logging-s3
              key: secretkey
        ports:
        - containerPort: 9000
      volumes:
        - name: logging-s3
          secret:
            secretName: logging-s3
---
kind: Service
apiVersion: v1
metadata:
  name: nginx-demo-minio
  namespace: logging
spec:
  selector:
    app: minio
  ports:
  - protocol: TCP
    port: 9000
    targetPort: 9000

EOF
```

#### Create `logging` resource
```bash
kubectl -n logging apply -f - <<"EOF" 
apiVersion: logging.banzaicloud.io/v1beta1
kind: Logging
metadata:
  name: default-logging-simple
spec:
  fluentd:
    metrics:
      serviceMonitor: true
  fluentbit:
    metrics:
      serviceMonitor: true
  controlNamespace: logging
EOF
```
> Note: `ClusterOutput` and `ClusterFlow` resource will only be accepted in the `controlNamespace` 

#### Create Minio output definition 
```bash
kubectl -n logging apply -f - <<"EOF" 
apiVersion: logging.banzaicloud.io/v1beta1
kind: Output
metadata:
  name: demo-output
spec:
  s3:
    aws_key_id:
      valueFrom:
        secretKeyRef:
          key: accesskey
          name: logging-s3
    aws_sec_key:
      valueFrom:
        secretKeyRef:
          key: secretkey
          name: logging-s3
    buffer:
      timekey: 10s
      timekey_use_utc: true
      timekey_wait: 0s
    force_path_style: "true"
    path: logs/${tag}/%Y/%m/%d/
    s3_bucket: demo
    s3_endpoint: http://nginx-demo-minio.logging.svc.cluster.local:9000
    s3_region: test_region
EOF
```
> Note: For production set-up we recommend using longer `timekey` interval to avoid generating too many object.

#### Create `flow` resource
```bash
kubectl -n logging apply -f - <<"EOF" 
apiVersion: logging.banzaicloud.io/v1beta1
kind: Flow
metadata:
  name: demo-flow
spec:
  filters:
    - tag_normaliser: {}
    - parser:
        remove_key_name_field: true
        reserve_data: true
        parse:
          type: nginx
  selectors:
    app.kubernetes.io/instance: nginx-demo
    app.kubernetes.io/name: nginx-logging-demo
  outputRefs:
    - demo-output
EOF
```

#### Install test deployment
```bash
kubectl -n logging apply -f - <<"EOF" 
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

## Validation

### Minio
#### Get Minio login credantials
```bash
kubectl -n logging get secrets logging-s3 -o json | jq '.data | map_values(@base64d)'
```
#### Forward Service
```bash
kubectl -n logging port-forward svc/logging-demo-minio 9000
```
[Minio Dashboard: http://localhost:9000](http://localhost:9000)
<p align="center"><img src="./img/servicemonitor_minio.png" width="660"></p>


### Prometheus
#### Forward Service
```bash
kubectl port-forward svc/monitor-prometheus-operato-prometheus 9090
```
[Prometheus Dashboard: http://localhost:9090](http://localhost:9090)
<p align="center"><img src="./img/servicemonitor_prometheus.png" width="660"></p>


### Grafana 
#### Get Grafana login credantials
```bash
kubectl get secret --namespace logging monitor-grafana -o jsonpath="{.data.admin-password}" | base64 --decode ; echo
```
> Default username: `admin`

#### Forward Service
```bash
kubectl -n logging port-forward svc/monitor-grafana 3000:80
```
[Gradana Dashboard: http://localhost:3000](http://localhost:3000)
<p align="center"><img src="./img/servicemonitor_grafana.png" width="660"></p>


