<p align="center"><img src="../img/nle.png" width="340"></p>

# Store NGINX access logs in Elasticsearch with Logging operator

This guide describes how to collect application and container logs in Kubernetes using the Logging operator, and how to send them to Elasticsearch.

The following figure gives you an overview about how the system works. The Logging operator collects the logs from the application, selects which logs to forward to the output, and sends the selected log messages to the output (in this case, to Elasticsearch). For more details about the Logging operator, see the [Logging operator overview](../Readme.md).

<p align="center"><img src="../img/nginx-elastic.png" width="900"></p>

---

## Contents

- **Installation**
  - **Elasticsearch operator**
    - [Deploy with Kubernetes Manifests](#deploy-elasticsearch)
  - **Logging operator and Demo application**
    - [Deploy with Helm](#deploy-the-logging-operator-with-helm)
    - [Deploy with Kubernetes Manifests](#deploy-the-logging-operator-with-kubernetes-manifests)
- **Validation**
    - [Kibana Dashboard](#validate-the-deployment)

---

## Deploy Elasticsearch

First, deploy Elasticsearch in your Kubernetes cluster. The following procedure is based on the [Elastic Cloud on Kubernetes quickstart](https://www.elastic.co/guide/en/cloud-on-k8s/current/k8s-quickstart.html)

1. Install the Elasticsearch operator.
    ```yaml
    kubectl apply -f https://download.elastic.co/downloads/eck/1.0.0-beta1/all-in-one.yaml
    ```
1. Create the `logging` Namespace.
    ```bash
    kubectl create ns logging
    ```
1. Install the Elasticsearch cluster.
    ```yaml
    cat <<EOF | kubectl apply -n logging -f -
    apiVersion: elasticsearch.k8s.elastic.co/v1beta1
    kind: Elasticsearch
    metadata:
      name: quickstart
    spec:
      version: 7.5.0
      nodeSets:
      - name: default
        count: 1
        config:
          node.master: true
          node.data: true
          node.ingest: true
          node.store.allow_mmap: false
    EOF
    ```
1. Install Kibana.
    ```yaml
    cat <<EOF | kubectl apply -n logging -f -
    apiVersion: kibana.k8s.elastic.co/v1beta1
    kind: Kibana
    metadata:
      name: quickstart
    spec:
      version: 7.5.0
      count: 1
      elasticsearchRef:
        name: quickstart
    EOF
    ```

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
      --set "elasticsearch.enabled=True"
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
       name: es-output
     spec:
       elasticsearch:
         host: quickstart-es-http.logging.svc.cluster.local
         port: 9200
         scheme: https
         ssl_verify: false
         ssl_version: TLSv1_2
         user: elastic
         password:
           valueFrom:
             secretKeyRef:
               name: quickstart-es-elastic-user
               key: elastic
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
       name: es-flow
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
         - es-output
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

## Validate the deployment

To validate that the deployment was successful, complete the following steps.

1. Use the following command to retrieve the password of the `elastic` user:
    ```bash
    kubectl -n logging get secret quickstart-es-elastic-user -o=jsonpath='{.data.elastic}' | base64 --decode
    ```
1. Enable port forwarding to the Kibana Dashboard Service.
    ```bash
    kubectl -n logging port-forward svc/quickstart-kb-http 5601
    ```
1. Open the Kibana dashboard in your browser: [https://localhost:5601](https://localhost:5601). You should see the dashboard and some sample log messages from the demo application.

<p align="center"><img src="../img/es_kibana.png" width="660"></p>
> If you don't get the expected result you can find help in the [troubleshooting-guideline](../troubleshooting.md).
