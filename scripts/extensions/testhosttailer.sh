#!/usr/bin/env bash

FLAGFILE=runflag

function cleanup {
    echo "Cleanup..."
    kubectl exec test-pd -- sh -c "rm /test/${FLAGFILE}"
    pkill -f urandom
}
trap cleanup EXIT

kubectl apply -f config/samples/logging-extensions_v1alpha1_hosttailer.yaml
kubectl apply -f config/samples/testpod.yaml

#polling pods to be ready
./scripts/pollstatus.sh "kubectl get pod test-pd | grep -i running" 2 30
./scripts/pollstatus.sh "kubectl get pod -l app.kubernetes.io/name=host-tailer | grep -i running" 2 30

# podname=$(kubectl get pod -l app.kubernetes.io/name=host-tailer -o jsonpath="{.items[0].metadata.name}")
# podname=${$(kubectl get pod -l app.kubernetes.io/name=host-tailer -o name)#*/}
podname=$(kubectl get pod -l app.kubernetes.io/name=host-tailer -o name | sed s/^.\*\\//\/)

echo $podname

# kubectl exec test-pd -- sh -c "cd /test && mkdir -p /test/temp && while true; do head -c 64 /dev/urandom | sha1sum | tee -a shifter | tee -a foobar | tee -a temp/shifter; sleep 0.5; done"
kubectl exec test-pd -- sh -c "cd /test && touch ${FLAGFILE} && mkdir -p /test/temp && while [ -f /test/${FLAGFILE} ]; do head -c 64 /dev/urandom | sha1sum | tee -a temp/foobar | tee -a temp/shifter; sleep 0.5; done" &

stern ${podname}
