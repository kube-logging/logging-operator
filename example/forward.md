Create the namespaces
```
kubectl create ns source
kubectl create ns target
```

Setup the `source` namespace to collect logs and forward to `target` 
> use https://github.com/ahmetb/kubectx to switch namespaces
```
kubens source

# RBAC
kubectl apply -f deploy/clusterrole.yaml \
  -f deploy/clusterrole_binding.yaml \
  -f deploy/service_account.yaml

# CRDs
kubectl apply -f deploy/crds/logging_v1alpha1_fluentbit_crd.yaml \
  -f deploy/crds/logging_v1alpha1_fluentd_crd.yaml \
  -f deploy/crds/logging_v1alpha1_plugin_crd.yaml

# CRs
kubectl apply -f deploy/crds/logging_v1alpha1_fluentbit_cr.yaml \
  -f deploy/crds/logging_v1alpha1_fluentd_cr.yaml

# setup the forwarder
kubectl apply -f example/cluster_forward.yaml
  
# start the operator locally to reconcile all the resources 
# (in a few seconds it finishes creating all resources then we can stop it)
WATCH_NAMESPACE=source go run cmd/manager/main.go

# both fluent-bit and fluentd should be successfully rolled out
kubectl rollout status daemonset fluent-bit-daemon
kubectl rollout status deployment test-fluentd

# install an app that writes logs
helm template charts/nginx-logging-demo -x templates/deployment.yaml | kubectl apply -f- 
```

Setup the `target` namespace
> use https://github.com/johanhaleby/kubetail to tail logs
```
kubens target

# create the fluentd deployment
kubectl apply -f deploy/crds/logging_v1alpha1_fluentd_cr.yaml

# send everything to stdout for checking the forwarded logs
kubectl apply -f example/stdout.yaml

WATCH_NAMESPACE=target go run cmd/manager/main.go

kubectl rollout status deployment test-fluentd

# check the logs by tailing our fluentd instance
kubetail test-fluentd
```