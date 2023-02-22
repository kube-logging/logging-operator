## Design and manual test document for the Persistent Volume drainer feature

### Context

Scaling out fluentd is trivial as long as the sender can distribute load evenly across fluentd instances.
In Kubernetes this is achieved using a standard service.

However scaling it back is nontrivial because a stopped fluentd pod will not necessarily flush all its buffers.

### Approach

The approach to fix this was based on https://medium.com/@marko.luksa/graceful-scaledown-of-stateful-apps-in-kubernetes-2205fc556ba9
but in addition to starting a pod with the same identity as the statefulset would, we start a job with the same fluentd config mounting the same volume.

We have four simple rules here:
- Volumes that belong to running statefulset pods don't have to be drained.
- Orphan volumes that are created for the statefulset's PVC template but where the associated statefulset pod does not exist anymore need to be drained.
- Volumes that have already been drained need to be tracked.
- Volumes that may need to be drained again need to be tracked.

To accomplish all the above the following algorithm takes place:
- Loop over all the PVCs and
  - check if they have an existing pod that they can be associated with, so they are *in use*
  - check if they have the special `logging.banzaicloud.io/drain-status` label set to `drained`
  - check if they have a *drainer job* in progress
  - take one of the following actions:
    - if it's *in use* and *drained*, then remove the label because it will need to be drained again after use
    - if it's not *in use*, not *drained* and does not have a successfully completed *job*, then create a placeholder pod and a drainer job for it
    - if it has a *job* that has successfully been completed, then add the `drained` label, delete the *job* and the placeholder pod**
    - if it has a *job* that has failed, then log the error and skip

Additionally, if you want to exclude certain PVCs from draining you can do so by marking them with the special `logging.banzaicloud.io/drain: no` label.

### Local test environment

Create a new cluster
```sh
kind create cluster --name drainer
```

Build and upload fresh drain-watch sidecar image
```sh
docker build drain-watch-image -t ghcr.io/kube-logging/fluentd-drain-watch:testing && kind load docker-image --name drainer ghcr.io/kube-logging/fluentd-drain-watch:testing
```

Install log-generator
```sh
one-eye log-generator install --update
```

Install CRDs into cluster
```sh
make install
```

Start log receiver
```sh
kubectl run log-target --image node:latest --expose --port 8080 --command -- node -e "var log=console.log;var ms=http.createServer((rq,rs)=>{var b='';rq.on('data',c=>{b+=c;});rq.on('end',()=>{log('got request',b);rs.writeHead(200);rs.end()})});var on=(cb)=>{ms.listening?(cb?cb():null):ms.listen(8080,'0.0.0.0',()=>{log('main server is listening on 8080');if(cb)cb()})};var off=(cb)=>{ms.listening?ms.close(()=>{log('main server stopped listening');if(cb)cb()}):(cb?cb():null)};http.createServer((rq,rs)=>{rq.url==='/on'?on(()=>{rs.writeHead(200);rs.end()}):rq.url==='/off'?off(()=>{rs.writeHead(200);rs.end()}):(log('invalid path',rq.url),rs.writeHead(404),rs.end())}).listen(8081,'0.0.0.0',()=>log('side server is listening on 8081'));on()"
```

Create Logging resource with 2 replicas, an output to the log receiver created previously, and flow to route every log to the output
```sh
kubectl apply -f- <<EOT
apiVersion: logging.banzaicloud.io/v1beta1
kind: Logging
metadata:
  name: drainer
spec:
  enableRecreateWorkloadOnImmutableFieldChange: true
  controlNamespace: default
  fluentbit: {}
  fluentd:
    scaling:
      replicas: 2
      drain:
        enabled: true
        image:
          tag: testing
EOT

kubectl apply -f- <<EOT
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

kubectl apply -f- <<EOT
apiVersion: logging.banzaicloud.io/v1beta1
kind: Flow
metadata:
  name: drainer-http
spec:
  match:
  - select:
      labels:
        app.kubernetes.io/name: log-generator
  localOutputRefs:
  - http
EOT
```

Start the operator
```sh
make run
```

> **NOTE**: This will keep running, so execute subsequent commands in a new terminal.

Check that logs are arriving at the receiver (might have to wait a few minutes)
```sh
kubectl logs -f log-target
```

### Primary test cases

Stop receiving logs
```sh
kubectl exec log-target -- curl -sS http://localhost:8081/off
```

Wait for logs to accumulate in buffers
```sh
kubectl exec -it drainer-fluentd-1 -- watch ls -lht /buffers # Press Ctrl+C when enough logs have accumulated
```

**Repeat the above steps before each test case.**

#### Simple downscale

Downscale fluentd
```sh
kubectl patch logging drainer --type merge -p '{"spec":{"fluentd":{"scaling":{"replicas":1}}}}'
```

Wait for the logging operator to reconcile and start the drain job.

In a new terminal, start a watch to examine how buffers change
```sh
kubectl exec -it $(kubectl get pod -l job-name=drainer-fluentd-1-drainer -o custom-columns=name:metadata.name --no-headers | head -1) -- watch ls -lht /buffers
```

Start receiving logs again
```sh
kubectl exec log-target -- curl -sS http://localhost:8081/on
```

Buffers should start disappearing after a while.
When all buffers are gone, the drainer job should be deleted along with its pod(s) and the placeholder pod, and the PVC `drainer-fluentd-buffer-drainer-fluentd-1` should be marked as drained.

#### Scale up while draining

Downscale fluentd
```sh
kubectl patch logging drainer --type merge -p '{"spec":{"fluentd":{"scaling":{"replicas":1}}}}'
```

Wait for the logging operator to reconcile and start the drain job.
A placeholder pod should also appear with the same name as the deleted statefulset pod (`drainer-fluentd-1`).

Upscale fluentd
```sh
kubectl patch logging drainer --type merge -p '{"spec":{"fluentd":{"scaling":{"replicas":3}}}}'
```

The drainer job and the placeholder pod should be terminated to allow the statefulset to create its required pods and giving back control over the buffers to the statefulset pods.
The previously drained PVC should **not** be marked as *drained*.

### Clean it up

```
one-eye log-generator delete --update
k delete logging-all --all
k delete po log-target
k delete svc log-target

kind delete cluster --name drainer
```