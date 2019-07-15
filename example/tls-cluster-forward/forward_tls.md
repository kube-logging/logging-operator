This is an example to demonstrate fluentd event forwarding. 
For the sake of simplicity it is demonstrated between namespaces.

### Create the TLS certificate

In this example we will create a single TLS bundle with the following contents:
- CA cert
- Server cert + key for the `source` and the `target` fluentd instance
- Client cert + key for the `source` fluentbit and the `source` fluentd forwarder

Caveats:
 - certs on the source and target side must use the same CA
 - the fluentd forwarder will look for the client cert in the same bundle that is used by the fluentd server
 
Enough said, let's create the namespaces and the cert
```
kubectl create ns target
kubectl create ns source
(cd example/tls-cluster-forward; ./gencert.sh)
```

To test local changes we can install the operator requirements through the chart but we will start the operator 
locally watching different namespaces.
```
helm upgrade --install logging-operator charts/logging-operator --set replicaCount=0
```

> use https://github.com/ahmetb/kubectx to switch namespaces with the `kubens` command
> use https://github.com/johanhaleby/kubetail to tail logs with the `kubetail` command

### Create and setup fluentd in the target namespace
```
kubens target

# create the fluentd resource
helm upgrade --install logging-operator-target charts/logging-operator-fluent \
  --set fluentbit.enabled=false \
  --set tls.enabled=true \
  --set tls.secretName=fluentd-tls \
  --set tls.sharedKey=example

# send everything to stdout for checking the forwarded logs from the `source` cluster
kubectl apply -f example/stdout.yaml

# start the operator to reconcile the desired state on the target namespace
# stop it once it created all resources successfully
WATCH_NAMESPACE=target go run cmd/manager/main.go

kubectl rollout status deployment fluentd
```

### Setup the `source` namespace to collect logs and forward to `target` 
```
kubens source

# create the fluentd resource
helm upgrade --install logging-operator-source charts/logging-operator-fluent \
  --set tls.enabled=true \
  --set tls.secretName=fluentd-tls \
  --set tls.sharedKey=example

# setup the tls forwarder with a shared key that matches with the existing key set up by the chart
cat <<EOT | kubectl apply -f-
apiVersion: logging.banzaicloud.com/v1alpha1
kind: Plugin
metadata:
  name: forward
spec:
  input:
    label:
      app: "*"
  output:
    - type: forward
      name: forward
      parameters:
        - name: host
          value: "fluentd.target.svc"
        - name: port
          value: "24240"
        - name: name
          value: target
        - name: tlsSharedKey
          value: ZXhhbXBsZQ==
EOT

# start the operator to reconcile the desired state on the source namespace
# stop it once it created all resources successfully
WATCH_NAMESPACE=source go run cmd/manager/main.go


# both fluent-bit and fluentd should be successfully rolled out
kubectl rollout status daemonset fluent-bit-daemon
kubectl rollout status deployment fluentd

# install an app that writes logs
helm template charts/nginx-logging-demo -x templates/deployment.yaml | kubectl apply -f- 
```


Go back to the target and check the logs
```
# check the logs by tailing our fluentd instance
kubens target
kubetail fluentd
```