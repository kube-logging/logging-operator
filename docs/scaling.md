# Scaling Fluentd

## Autoscaling with HPA

First, the [aggregation layer has to be configured](https://kubernetes.io/docs/tasks/extend-kubernetes/configure-aggregation-layer/).
Many providers already have this configured, so does `kind`.

Next, Prometheus and the [Prometheus Adapter](https://github.com/kubernetes-sigs/prometheus-adapter) have to be installed.
The default Prometheus address values will probably have to be adjusted (set `prometheus.url`, `prometheus.port` and `prometheus.path` to the appropriate values).

> You might also want to install [`metrics-server`](https://github.com/kubernetes-sigs/metrics-server) to access basic metrics.
>
> PRO TIP: If the `metrics-server` pod's readiness fails with HTTP 500, try adding the `--kubelet-insecure-tls` flag to the container.

If the necessary custom metric is not available in Prometheus, a Prometheus recording rule should be defined:
```yaml
groups:
- name: my-logging-hpa.rules
  rules:
  - expr: (node_filesystem_size_bytes{container="buffer-metrics-sidecar",mountpoint="/buffers"}-node_filesystem_free_bytes{container="buffer-metrics-sidecar",mountpoint="/buffers"})/node_filesystem_size_bytes{container="buffer-metrics-sidecar",mountpoint="/buffers"}
    record: buffer_space_usage_ratio
```

> Alternatively, the derived metric could be defined as a config rule in the Prometheus Adapter's config map.

Next, if it's not already installed, install the logging-operator and configure a logging with at least one flow.
Make sure that the logging resource has buffer volume metrics monitoring enabled under `spec.fluentd`:
```yaml
#spec:
#  fluentd:
    bufferVolumeMetrics:
      serviceMonitor: true
```

You can check that the custom metric is available with:
```sh
kubectl get --raw '/apis/custom.metrics.k8s.io/v1beta1/namespaces/default/pods/*/buffer_space_usage_ratio'
```

The replica count of the stateful set is enforced by the logging-operator only if it is explicitly set in the logging resource's replica count.

If it's not set explicitly in the logging resource, then the logging operator allows managing it from the outside.

To allow for HPA to control the replica count of the statefulset simply avoid setting the replica count through the logging resource.

Finally, the HPA resource must be created.
```yaml
apiVersion: autoscaling/v2beta2
kind: HorizontalPodAutoscaler
metadata:
  name: one-eye-fluentd
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: StatefulSet
    name: one-eye-fluentd
  minReplicas: 1
  maxReplicas: 10
  metrics:
  - type: Pods
    pods:
      metric:
        name: buffer_space_usage_ratio
      target:
        type: AverageValue
        averageValue: 800m
```
This example tries to keep the average buffer volume usage of Fluentd instances at 80%.
