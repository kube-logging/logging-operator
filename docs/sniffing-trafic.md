# Sniffing fluentd traffic from your local machine

This tutorial shows an easy way to test your Flows and Outputs by tailing them on your local machine.

# Prerequisites

- Docker
- Kurun


## Install Kurun

```bash 
brew install banzaicloud/tap/kurun
```

## Set up local fluentd

Create a *fluentd.conf*.
```yaml
<source>
  @type http
  @id http_input
  port 8888
  <parse>
    @type json
  </parse>
</source>

<match **>
  @type stdout
</match>
```

Start local fluentd
```bash
docker run -p 8888:8888 -v $PWD/fluentd.conf:/fluent/fluent.conf --rm ghcr.io/kube-logging/fluentd:v1.11 -c /fluent/fluent.conf
```

Setup kurun service for port-forward
```bash
kurun port-forward --servicename kurun --serviceport 8888 localhost:8888
```

> Note: Kurun should run next to the fluentd instances or you need to specify explicit service names.

Setup Output/ClusterOutput

```yaml
apiVersion: logging.banzaicloud.io/v1beta1
kind: ClusterOutput
metadata:
  name: sniffer
spec:
  http:
    endpoint: http://kurun:8888
    json_array: true
    buffer:
      flush_interval: 10s
      timekey: 5s
      timekey_wait: 1s
      flush_mode: interval
    format:
      type: json
```

Set up a debug daemonset

```yaml
apiVersion: apps/v1
kind: DaemonSet
metadata:
  name: pod-info-logger
  labels:
    app: pod-info-logger
spec:
  selector:
    matchLabels:
      name: pod-info-logger
  template:
    metadata:
      labels:
        name: pod-info-logger
        app: pod-info-logger
    spec:
      containers:
      - name: test-container
        image: k8s.gcr.io/busybox
        command: [ "sh", "-c"]
        args:
        - while true; do
            echo $(date) $MY_NODE_NAME $MY_POD_NAMESPACE/$MY_POD_NAME;
            sleep 10;
          done;
        env:
          - name: MY_NODE_NAME
            valueFrom:
              fieldRef:
                fieldPath: spec.nodeName
          - name: MY_POD_NAME
            valueFrom:
              fieldRef:
                fieldPath: metadata.name
          - name: MY_POD_NAMESPACE
            valueFrom:
              fieldRef:
                fieldPath: metadata.namespace
          - name: MY_POD_IP
            valueFrom:
              fieldRef:
                fieldPath: status.podIP
          - name: MY_POD_SERVICE_ACCOUNT
            valueFrom:
              fieldRef:
                fieldPath: spec.serviceAccountName
```

ClusterFlow
```yaml
apiVersion: logging.banzaicloud.io/v1beta1
kind: ClusterFlow
metadata:
  name: sniffing-demo-containers
spec:
  filters:
  selectors:
    match:
      labels:
        app: pod-info-logger
  outputRefs:
    - sniffer
```