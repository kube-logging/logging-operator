```
# start two nodes, e.g.
minikube start --nodes=2

# label the nodes for tenants
kubectl label node minikube tenant=tenant-a
kubectl label node minikube-m02 tenant=tenant-b

make codegen manifests install
kubectl apply -f config/samples/multitenant-hard/logging

helm upgrade --install --namespace a --create-namespace --set "nodeSelector.tenant=tenant-a" log-generator oci://ghcr.io/kube-logging/helm-charts/log-generator
helm upgrade --install --namespace b --create-namespace --set "nodeSelector.tenant=tenant-b" log-generator oci://ghcr.io/kube-logging/helm-charts/log-generator
# in a separate shell
make run

# expected output
NAMESPACE     NAME                               READY   STATUS    RESTARTS      AGE     IP             NODE           NOMINATED NODE   READINESS GATES
a-control     a-fluentbit-2997s                  1/1     Running   0             9m15s   10.244.0.5     minikube       <none>           <none>
a-control     a-fluentd-0                        2/2     Running   0             9m15s   10.244.0.6     minikube       <none>           <none>
a             log-generator-6cfb45c684-kbzk4     1/1     Running   0             11m     10.244.0.3     minikube       <none>           <none>
b-control     b-fluentbit-9bvbn                  1/1     Running   0             7m30s   10.244.1.7     minikube-m02   <none>           <none>
b-control     b-fluentd-0                        2/2     Running   0             7m29s   10.244.1.8     minikube-m02   <none>           <none>
b             log-generator-7b95b6fdc5-62bnr     1/1     Running   0             11m     10.244.1.3     minikube-m02   <none>           <none>
```
