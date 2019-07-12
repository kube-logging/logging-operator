#!/bin/bash -x

if ! which cfssl; then
  go get -tags nopkcs11 github.com/cloudflare/cfssl/cmd/cfssl
  go get -tags nopkcs11 github.com/cloudflare/cfssl/cmd/cfssljson
fi

cfssl gencert -initca cfssl-ca.json | cfssljson -bare ca
cfssl gencert -ca ca.pem -ca-key ca-key.pem -config cfssl-ca.json -profile server cfssl-csr.json | cfssljson -bare server
cfssl gencert -ca ca.pem -ca-key ca-key.pem -config cfssl-ca.json -profile client cfssl-csr.json | cfssljson -bare client

FILE_ARGS=()

for i in ca server client; do
  FILE_ARGS+=(--from-file "${i}Cert=${i}.pem" --from-file "${i}Key=${i}-key.pem")
done

kubectl create secret generic fluentd-tls -n target "${FILE_ARGS[@]}"
kubectl create secret generic fluentd-tls -n source "${FILE_ARGS[@]}"
