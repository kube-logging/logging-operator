# TC-routing

## Setup

In this setup customer-a sends logs over the traditional path (fluentbit -> fluentd) while
customer-b sends over the telemetry-controller path (opentelemetry collector -> fluentd).

In both cases, flows and outputs should work and behave the same.

```shell
make kind-cluster
make docker-build
kind load docker-image controller:local

helm dependency update charts/logging-operator
helm upgrade --install \
    --wait \
    --create-namespace \
    --namespace logging \
    logging-operator ./charts/logging-operator/ \
    --set image.repository=controller \
    --set image.tag=local \
    --set extraArgs='{"-enable-leader-election=true","-enable-telemetry-controller-route"}' \
    --set telemetry-controller.install=true \
    --set testReceiver.enabled=true

kubectl apply -f config/samples/telemetry-controller-routing/

helm upgrade --install --namespace customer-a log-generator oci://ghcr.io/kube-logging/helm-charts/log-generator
helm upgrade --install --namespace customer-b log-generator oci://ghcr.io/kube-logging/helm-charts/log-generator
helm upgrade --install --namespace infra log-generator oci://ghcr.io/kube-logging/helm-charts/log-generator
```

## Switching between logging-routes or TC routing

To use only logging-routes (This is the default behaviour!):

- Don't set anything.

```yaml
apiVersion: logging.banzaicloud.io/v1beta1
kind: Logging
spec:
  routeConfig:
    enableTelemetryControllerRoute: false
    disableLoggingRoute: false
```

To use logging-routes and TC routing in parallel:

```yaml
apiVersion: logging.banzaicloud.io/v1beta1
kind: Logging
spec:
  routeConfig:
    enableTelemetryControllerRoute: true
    disableLoggingRoute: false
```

To use only TC routing:

```yaml
apiVersion: logging.banzaicloud.io/v1beta1
kind: Logging
spec:
  routeConfig:
    enableTelemetryControllerRoute: true
    disableLoggingRoute: true
```

NOTE: You can change these setting on the fly, the only requirement is to have the necessary components deployed.
