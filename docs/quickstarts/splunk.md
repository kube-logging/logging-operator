---
title: Demo Splunk operator with Logging operator
shorttitle: Splunk HEC
weight: 300
---

{{< contents >}}

<p align="center"><img src="../../img/splunk.png" width="340"></p>

This guide describes how to collect application and container logs in Kubernetes using the Logging operator, and how to send them to Splunk.

The following figure gives you an overview about how the system works. The Logging operator collects the logs from the application, selects which logs to forward to the output, and sends the selected log messages to the output (in this case, to Splunk). For more details about the Logging operator, see the [Logging operator overview]({{< relref "docs/one-eye/logging-operator/_index.md">}}).

## Deploy Splunk

First, deploy Splunk Standalone in your Kubernetes cluster. The following procedure is based on the [Splunk on Kubernetes quickstart](https://www.splunk.com/en_us/blog/it/an-insider-s-guide-to-splunk-on-containers-and-kubernetes.html).

1. Create the `logging` Namespace.

    ```bash
    kubectl create ns logging
    ```

1. Install the Splunk operator.

    ```yaml
    kubectl apply -n logging -f https://tiny.cc/splunk-operator-install
    ```
  
1. Install the Splunk cluster.

    ```yaml
    kubectl apply -n logging -f - <<"EOF" 
    apiVersion: enterprise.splunk.com/v1alpha2
    kind: Standalone
    metadata:
      name: single
    spec:
      config:
        splunkPassword: helloworld456
        splunkStartArgs: --accept-license
      topology:
        standalones: 1
    EOF
    ```


## Deploy the Logging operator and a demo Application

Install the Logging operator and a demo application to provide sample log messages.

### Deploy the Logging operator with Kubernetes manifests

To deploy the Logging operator using Kubernetes manifests, complete these steps. If you want to install the Logging operator using Helm, see [Deploy the Logging operator with Helm](#deploy-the-logging-operator-with-helm).

1. Install the Logging operator. For details, see [How to install Logging-operator from manifests]({{< relref "docs/one-eye/logging-operator/deploy/_index.md#deploy-the-logging-operator-from-kubernetes-manifests" >}}.
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

1. Get a Splunk HEC Token.

     ```bash
     HEC_TOKEN=$(kubectl get secret -n logging  splunk-single-standalone-secrets -o jsonpath='{.data.hec_token}' | base64 --decode)
     ```

1. Create a Splunk output secret from the token.
     ```bash
     kubectl  create secret generic splunk-token -n logging --from-literal "SplunkHecToken=${HEC_TOKEN}"
     ```


1. Define a Splunk `output`.

    ```bash
    kubectl -n logging apply -f - <<"EOF"
    apiVersion: logging.banzaicloud.io/v1beta1
    kind: Output
    metadata:
     name: splunk-output
    spec:
     splunkHec:
        hec_host: splunk-single-standalone-headless
        insecure_ssl: true
        hec_port: 8088
        hec_token: 
            valueFrom:
               secretKeyRef:
                  name:  splunk-token
                  key: SplunkHecToken
        index: main
        format:
          type: json 
    EOF
    ```


1. Create a `flow` resource.

    ```bash
    kubectl -n logging apply -f - <<"EOF"
    apiVersion: logging.banzaicloud.io/v1beta1
    kind: Flow
    metadata:
      name: splunk-flow
    spec:
      filters:
        - tag_normaliser: {}
        - parser:
            remove_key_name_field: true
            reserve_data: true
            parse:
              type: nginx
      match:
        - select:
            labels:
              app.kubernetes.io/name: log-generator
      outputRefs:
        - splunk-output
    EOF
    ```

1. Install the demo application.

     ```bash
    kubectl -n logging apply -f - <<"EOF" 
    apiVersion: apps/v1
    kind: Deployment
    metadata:
      name: log-generator
      labels:
        app.kubernetes.io/name: log-generator
    spec:
      replicas: 3
      selector:
        matchLabels:
          app.kubernetes.io/name: log-generator
      template:
        metadata:
          labels:
            app.kubernetes.io/name: log-generator
        spec:
          containers:
          - name: log-generator
            image: banzaicloud/log-generator:0.3.2
    EOF
     ```

## Validate the deployment

To validate that the deployment was successful, complete the following steps.

1. Use the following command to retrieve the password of the `admin` user:

    ```bash
    kubectl -n logging get secret splunk-single-standalone-secrets -o jsonpath='{.data.password}' | base64 --decode
    ```

1. Enable port forwarding to the Splunk Dashboard Service.

    ```bash
    kubectl -n logging port-forward svc/splunk-single-standalone-headless 8000
    ```

1. Open the Splunk dashboard in your browser: [http://localhost:8000](http://localhost:8000). You should see the dashboard and some sample log messages from the demo application.

<p align="center"><img src="../../img/splunk_dash.png" width="660"></p>

> If you don't get the expected result you can find help in the [troubleshooting section]({{< relref "docs/one-eye/logging-operator/troubleshooting.md">}}).
