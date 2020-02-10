<p align="center"><img src="../img/kafka_logo.png" width="340"></p>

# Transport Nginx Access Logs into Kafka with Logging operator

This guide describes how to collect application and container logs in Kubernetes using the Logging operator, and how to send them to Kafka.

The following figure gives you an overview about how the system works. The Logging operator collects the logs from the application, selects which logs to forward to the output, and sends the selected log messages to the output (in this case, to Kafka). For more details about the Logging operator, see the [Logging operator overview](../Readme.md).

<p align="center"><img src="../img/nignx-kafka.png" width="900"></p>

---
## Contents
- **Installation**
  - **Kafka** 
    - [Deploy with Helm](#deploy-kafka)
  - **Logging operator and Demo application**
    - [Deploy with Helm](#deploy-the-logging-operator-with-helm)
    - [Deploy with Kubernetes Manifests](#deploy-the-logging-operator-with-kubernetes-manifests)
- **Validation**
    - [Kafkacat](#test-your-deployment-with-kafkacat)
---
<br />

## Deploy Kafka
>In this demo we are using our kafka operator.
> [Easy Way Installing with Helm](https://github.com/banzaicloud/kafka-operator#easy-way-installing-with-helm)
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
      --set "kafka.enabled=True"
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
      name: kafka-output
    spec:
      kafka:
        brokers: kafka-headless.kafka.svc.cluster.local:29092
        default_topic: topic
        format: 
          type: json    
        buffer:
          tags: topic
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
       name: kafka-flow
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
         - kafka-output
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

## Test Your Deployment with kafkacat
#### Exec Kafaka test pod
```bash
kubectl -n kafka exec -it kafka-test-c sh
```

#### Run kafkacat
```bash
kafkacat -C -b kafka-0.kafka-headless.kafka.svc.cluster.local:29092 -t topic
```

[![asciicast](https://asciinema.org/a/273236.svg)](https://asciinema.org/a/273236)

> If you don't get the expected result you can find help in the [troubleshooting-guideline](../troubleshooting.md).
