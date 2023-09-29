# Mutli-Worker Fluentd Setup

## Necessity

In specific scenarios, a fluentd with a single worker instance cannot process and forward the high amount of logs produced on clusters. This can lead to fluentd Pods not accepting additional traffic from fluent-bits and fluent-bits suffering under Backpressure. In the end, both fluentd and fluent-bit Pods might run into their memory limits and get restarted by Kubernetes. Enabling multiple worker processes per fluentd Pod will increase the performance of this component, so it is recommended to use a multi-worker approach in environments with high log volume. Additionally the official [fluentd documentation](https://docs.fluentd.org/deployment/multi-process-workers) might be helpful.

## Recommended implementation

When enabling the multi-worker setup, it is recommended to ensure the following things:
- Place the fluentds on separate Nodes
- Increase the compute and memory resources 
- Do not use specific filter plugins
  - detectExceptions [is not working](https://github.com/kube-logging/logging-operator/issues/1490)
- Also scale [horizontally](https://kube-logging.dev/docs/logging-infrastructure/fluentd/#autoscaling)

To ensure that the fluentd Pods have enough resources, a common approach is to use specific Nodes for the fluentds and to reserve enough computing and memory resources. A new nodePool should be created with a specific label and a taint. Ideally, the nodeType is compute-optimized. It could look like the following:
```yaml
apiVersion: v1
kind: Node
metadata:
  labels:    
    type: cpu
  name: node1
spec:
  taints:
  - effect: NoSchedule
    key: type
    value: cpu
```

The corresponding setting in the FluentdSpec looks like follows:
```yaml
nodeSelector:
  type: cpu
tolerations:
- effect: NoSchedule
  key: type
  operator: Equal
  value: cpu
```

Additionaly we will have to increase the resources that are requested by the fluentd Pods. In the default setting they use following requests and limits:
```yaml
resources:
  limits:
    cpu: 1
    memory: 400M
  requests:
    cpu: 500m
    memory: 100M
```

In this short walkthrough, we will increase the fluentd workers from `1` to `5`. Therefore, we will multiply the requests and limits with factor 5 to ensure enough resources are reserved. Additionally, we will set requests and limits to the same values to ensure that the fluentd Pods are not affected by other workloads on the Node. This is, in general, a good practice. It is necessary to set the following settings in the FluentdSpec:
```yaml
resources:
  limits:
    cpu: 5
    memory: 2G
  requests:
    cpu: 5
    memory: 2G
```

Lastly we can increase the number of fluentd-workers that are used per Pod and set the rootDir field. It is important that those two settings are changed together otherwise the fluentd process will not work correctly:
```yaml
workers: 5
rootDir: /buffers
```

The full configuration of the Logging resource looks like follows:
```yaml
fluentd:
  nodeSelector:
    type: cpu
  tolerations:
  - effect: NoSchedule
    key: type
    operator: Equal
    value: cpu
  resources:
    limits:
      cpu: 5
      memory: 2G
    requests:
      cpu: 5
      memory: 2G
  workers: 5
  rootDir: /buffers
```