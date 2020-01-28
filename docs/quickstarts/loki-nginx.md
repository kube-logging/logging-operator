<p align="center"><img src="../img/nll.png" width="340"></p>

# Store Nginx Access Logs in Grafana Loki with Logging operator

This guide describes how to collect application and container logs in Kubernetes using the Logging operator, and how to send them to Grafana Loki.

The following figure gives you an overview about how the system works. The Logging operator collects the logs from the application, selects which logs to forward to the output, and sends the selected log messages to the output (in this case, to Loki). For more details about the Logging operator, see the [Logging operator overview](../Readme.md).

<p align="center"><img src="../img/nginx-loki.png" width="900"></p>

---
## Contents
- **Installation**
  - **Loki**
    - [Deploy with Helm](#deploy-loki-and-grafana)
  - **Logging operator and Demo application**
    - [Deploy with Helm](#deploy-the-logging-operator-with-helm)
    - [Deploy with Kubernetes Manifests](#deploy-the-logging-operator-with-kubernetes-manifests)
- **Validation**
    - [Grafana Dashboard](#grafana-dashboard)
---

<br />

## Deploy Loki and Grafana

1 Add loki chart repository:
```bash
helm repo add loki https://grafana.github.io/loki/charts
helm repo update
```

1 Install Loki
```bash
helm install --namespace logging --name loki loki/loki
```
> [Grafana Loki Documentation](https://github.com/grafana/loki/tree/master/production/helm)

1 Install Grafana
```bash
helm install --namespace logging --name grafana stable/grafana \
 --set "datasources.datasources\\.yaml.apiVersion=1" \
 --set "datasources.datasources\\.yaml.datasources[0].name=Loki" \
 --set "datasources.datasources\\.yaml.datasources[0].type=loki" \
 --set "datasources.datasources\\.yaml.datasources[0].url=http://loki:3100" \
 --set "datasources.datasources\\.yaml.datasources[0].access=proxy"
```
<br />



## Deploy the Logging operator and a demo Application

Next, install the Logging operator and a demo application to provide sample log messages.

### Deploy the Logging operator with Helm

To install the Logging operator using Helm, complete these steps. If you want to install the Logging operator using Kubernetes manifests, see [Deploy the Logging operator with Kubernetes manifests](../deploy/README.md#deploy-the-logging-operator-from-kubernetes-manifests).

1. Add the chart repository of the Logging operator using the following commands:
    ```bash
    helm repo add banzaicloud-stable https://kubernetes-charts.banzaicloud.com
    helm repo update
    ```
1. Install the Logging operator. For details, see [How to install Logging-operator with Helm](../deploy/README.md#deploy-logging-operator-with-helm)
1. Install the demo application and its logging definition.
    ```bash
    helm install --namespace logging --name logging-demo banzaicloud-stable/logging-demo \
      --set "loki.enabled=True"
    ```

### Deploy the Logging operator with Kubernetes manifests

To deploy the Logging operator using Kubernetes manifests, complete these steps. If you want to install the Logging operator using Helm, see [Deploy the Logging operator with Helm](#deploy-the-logging-operator-with-helm).   

1. Install the Logging operator. For details, see [How to install Logging-operator from manifests](../deploy/README.md#deploy-the-logging-operator-from-kubernetes-manifests)
1. Create the `logging` resource.
     ```bash
     kubectl -n logging apply -f - <<"EOF" 
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
     > Note: You can use the `ClusterOutput` and `ClusterFlow` resources only in the `controlNamespace`.
1. Create an Elasticsearch `output` definition.
     ```bash
    kubectl -n logging apply -f - <<"EOF" 
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
     > Note: In production environment, use a longer `timekey` interval to avoid generating too many objects.
1. Create a `flow` resource.
     ```bash
     kubectl -n logging apply -f - <<"EOF" 
     apiVersion: logging.banzaicloud.io/v1beta1
     kind: Flow
     metadata:
       name: loki-flow
     spec:
       filters:
         - tag_normaliser: {}
         - parser:
             remove_key_name_field: true
             reserve_data: true
             parse:
               type: nginx
       selectors:
         app.kubernetes.io/name: log-generator
       outputRefs:
         - loki-output
     EOF
     ```
1. Install the demo application.
     ```bash
    kubectl -n logging apply -f - <<"EOF" 
    apiVersion: apps/v1 
    kind: Deployment
    metadata:
      name: log-generator
    spec:
      selector:
        matchLabels:
          app.kubernetes.io/name: log-generator
      replicas: 1
      template:
        metadata:
          labels:   
            app.kubernetes.io/name: log-generator
        spec:
          containers:
          - name: nginx
            image: banzaicloud/log-generator:0.3.2
    EOF
     ```


## Deployment Validation

### Grafana Dashboard

#### Get Grafana login credantials
```bash
kubectl -n logging get secrets grafana -o json | jq '.data | map_values(@base64d)'
```

#### Forward Grafana Service
```bash
kubectl -n logging port-forward svc/grafana 3000:80
```
Gradana Dashboard: [http://localhost:3000](http://localhost:3000)
<p align="center"><img src="../img/loki1.png" width="660"></p>

> If you don't get the expected result you can find help in the [troubleshooting-guideline](../troubleshooting.md).
