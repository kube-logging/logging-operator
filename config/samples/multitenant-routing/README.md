```
make generate install
kubectl apply -f config/samples/multitenant-routing/logging
helm upgrade --install --namespace customer-a log-generator oci://ghcr.io/kube-logging/helm-charts/log-generator
helm upgrade --install --namespace customer-b log-generator oci://ghcr.io/kube-logging/helm-charts/log-generator
# in a separate shell
make run
```
