## Design and manual test document for the Persistent Volume drainer feature

### Context

Scaling out fluentd is trivial as long as the sender can distribute load evenly across fluentd instances.
In Kubernetes this is achieved using a standard service.

However scaling it back is nontrivial because a stopped fluentd pod will not necessarily flush all it's buffers.

### Approach

The approach to fix this was based on https://medium.com/@marko.luksa/graceful-scaledown-of-stateful-apps-in-kubernetes-2205fc556ba9 
but instead of starting a pod with the same identity as the statefulset would, we start a job with the same fluentd config mounting the same volume.

We have two simple rules here:
- Volumes that belong to running statefulset pods don't have to be drained.
- Orphan volumes that are created for the statefulset's PVC template but where the associated statefulset pod does not exist anymore need to be drained.
- Volumes that have already been drained need to be tracked.
- Volumes that may need to be drained again need to be tracked.

To accomplish all the above the following algorithm takes place:
- Loop over all the PVCs and 
  - check if they have an existing pod that they can be associated with, so they are `live`
  - check if they have the special `logging.banzaicloud.io/drain-status` label set to `drained`
  - check if they have a drainer `job` in progress
    - If it's `live` and `drained` **remove the label** because it will need to be drained again after use
    - If it's not `live`, not `drained` and does not have a successfully completed `job`, then **create a drainer job**
    - If it has a `job` that has successfully been completed then add **the `drained` label and delete the `job`**
    - If it has a `job` that has failed, then log the error and skip

### Local test environment

```
kind create cluster --name drainer

docker build -f drain-watch-image/Dockerfile -t fluentd-drain-watch:latest drain-watch-image && kind load docker-image --name drainer fluentd-drain-watch:latest

one-eye log-generator install --update

make install
```

Install log in and start log receiver
```
kubectl run log-target --image node:latest --expose=true --port 8080 --command -- sleep 99999

kubectl exec -ti log-target -- node
var server = http.createServer(function(req, res) { let d = ''; req.on('data', c => {d+=c;}); req.on('end', () => { console.log('got request', d); res.writeHead(200); res.end();}) }).listen(8080, '0.0.0.0', () => console.log('server is running'))
```

```
cat <<EOT | kubectl apply -f-
apiVersion: logging.banzaicloud.io/v1beta1
kind: Logging
metadata:
  name: drainer
spec:
  enableRecreateWorkloadOnImmutableFieldChange: true
  controlNamespace: default
  fluentd:
    scaling:
      replicas: 2
      drainWatch:
        image:
          repository: fluentd-drain-watch
          tag: latest
EOT

cat <<EOT | kubectl apply -f-
apiVersion: logging.banzaicloud.io/v1beta1
kind: Output
metadata:
  name: http
spec:
  http:
    endpoint: http://log-target:8080
    buffer:
      type: file
      tags: time
      timekey: 1m
      timekey_wait: 0s
EOT

cat <<EOT | kubectl apply -f-
apiVersion: logging.banzaicloud.io/v1beta1
kind: Flow
metadata:
  name: drainer-http
spec:
  match:
  - select: {}
  localOutputRefs:
  - http
EOT
```

Reconcile a single round
```
make run
```