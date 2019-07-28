Create the namespaces
```
kubectl create ns source
kubectl create ns target
```

To test local changes we can install the operator requirements through the chart but we will start the operator 
locally watching different namespaces.
```
helm upgrade --install logging-operator charts/logging-operator --set replicaCount=0
```

Setup the `target` namespace
```
kubens target

# create the fluentd resource
helm upgrade --install logging-operator-target charts/logging-operator-fluent \
  --set fluentbit.enabled=false

# send everything to stdout for for this simple demonstration
kubectl apply -f example/stdout.yaml

# start the operator to reconcile the desired state on the target namespace
# stop it once it created all resources successfully
WATCH_NAMESPACE=target go run cmd/manager/main.go

kubectl rollout status deployment fluentd
```

Setup the `source` namespace to collect logs and forward to `target` 
```
kubens source

# create the fluentd resource
helm upgrade --install logging-operator-source charts/logging-operator-fluent

# install the demo app that writes logs
helm upgrade --install nginx-logging-demo charts/nginx-logging-demo \
  --set forwarding.enabled=true \
  --set forwarding.targetHost=fluentd.target.svc \
  --set forwarding.targetPort=24240

# start the operator to reconcile the desired state on the source namespace
# stop it once it created all resources successfully
WATCH_NAMESPACE=source go run cmd/manager/main.go

# both fluent-bit and fluentd should be successfully rolled out
kubectl rollout status daemonset fluent-bit-daemon
kubectl rollout status deployment fluentd
```

Watch the logs as they arrive to the target cluster
```
kubetail -n target
```