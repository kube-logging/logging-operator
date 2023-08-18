```
make generate manifests install
kubectl apply -f config/samples/fluentbit-ns-routing/logging
helm upgrade --install --namespace a log-generator oci://ghcr.io/kube-logging/helm-charts/log-generator
helm upgrade --install --namespace b1 log-generator oci://ghcr.io/kube-logging/helm-charts/log-generator
helm upgrade --install --namespace b2 log-generator oci://ghcr.io/kube-logging/helm-charts/log-generator
# in a separate shell
make run
```
