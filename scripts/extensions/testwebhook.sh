#!/usr/bin/env bash

helm repo add jetstack https://charts.jetstack.io
helm repo update
kubectl apply -f https://github.com/jetstack/cert-manager/releases/download/v1.6.1/cert-manager.crds.yaml

# lsof -nP | grep LISTEN
# curl -v --cacert "$(mkcert -CAROOT)/rootCA.pem" 'https://localhost:9443/mutate-v1-pod' -H "Content-Type: application/json" --request POST --data '{"foo":"bar"}'

FLAGFILE=runflag
# export KUBECONFIG="$(k3d get-kubeconfig --name='k3s-default')"
tmpdir=$TMPDIR/k8s-webhook-server/serving-certs/

# function cleanup {
#     echo "Cleanup..."
#     kubectl exec test-pod -- sh -c "rm /test/${FLAGFILE}"
#     pkill -f urandom
# }
# trap cleanup EXIT

# kubectl apply --validate=false -f https://github.com/jetstack/cert-manager/releases/download/v0.14.0/cert-manager.yaml

# kubectl apply -f config/samples/logging-extensions_v1alpha1_webhook.yaml
kubectl apply -f - <<EOF
apiVersion: admissionregistration.k8s.io/v1
kind: MutatingWebhookConfiguration
metadata:
  name: sample-webhook-cfg
  labels:
    app: sample-webhook
  annotations:
    certmanager.k8s.io/inject-ca-from: "namespace/certificate-name"
webhooks:
  - name: sample-webhook.banzaicloud.com
    clientConfig:
      service:
        name: sample-tailer-webhook
        namespace: default
        path: "/mutate-v1-pod"
      caBundle: $(cat "$(mkcert --CAROOT)/rootCA.pem" | base64)
    rules:
      - operations: [ "CREATE" ]
        apiGroups: [""]
        apiVersions: ["v1"]
        resources: ["pods"]
        scope: "*"
    sideEffects: None
    admissionReviewVersions:
      - v1
      - v1alpha1
EOF

# inject certs
kubectl apply -f - <<EOF
apiVersion: v1
kind: Secret
metadata:
  name: ca-key-pair
  namespace: default
data:
  tls.crt: $(cat "$(mkcert --CAROOT)/rootCA.pem" | base64)
  tls.key: $(cat "$(mkcert --CAROOT)/rootCA-key.pem" | base64)
EOF

./scripts/pollstatus.sh "kubectl get svc -l app=webhook -n cert-manager -o=name | grep cert-manager" 2 30
./scripts/pollstatus.sh "kubectl get svc -l app=cert-manager -n cert-manager -o=name | grep cert-manager" 2 30
./scripts/pollstatus.sh "kubectl get pod -l app=cert-manager -n cert-manager | grep -i 1/1" 2 30
./scripts/pollstatus.sh "kubectl get pod -l app=cainjector -n cert-manager | grep -i 1/1" 2 30
./scripts/pollstatus.sh "kubectl get pod -l app=webhook -n cert-manager | grep -i 1/1" 2 30

kubectl apply -f - <<EOF
apiVersion: cert-manager.io/v1
kind: Issuer
metadata:
  name: ca-issuer
  namespace: default
spec:
  ca:
    secretName: ca-key-pair
EOF

kubectl apply -f - <<EOF
apiVersion: cert-manager.io/v1
kind: Certificate
metadata:
  name: dev-webhook
  namespace: default
spec:
  secretName: dev-webhook-tls
  duration: 2160h # 90d
  renewBefore: 360h # 15d
  commonName: example.com
  isCA: false
  usages:
    - server auth
    - client auth
  dnsNames:
  - sample-tailer-webhook
  - sample-tailer-webhook.default
  - sample-tailer-webhook.default.svc
  - localhost
  ipAddresses:
  - 192.168.0.5
  issuerRef:
    name: ca-issuer
    kind: Issuer
    group: cert-manager.io
EOF

kurun port-forward --namespace default --servicename sample-tailer-webhook https://localhost:9443 --tlssecret dev-webhook-tls --serviceport 443

# kubectl apply -f config/samples/testpod_annotated.yaml
# kubectl exec test-pod -- sh -c "cd /test && touch ${FLAGFILE} && mkdir -p /var/foo && while [ -f /test/${FLAGFILE} ]; do head -c 64 /dev/urandom | sha1sum | tee -a /var/foo/bar | tee -a /var/zzz; sleep 0.5; done" &
