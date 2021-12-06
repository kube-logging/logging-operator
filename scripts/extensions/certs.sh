#!/usr/bin/env bash

# tmpdir=$TMPDIR/k8s-webhook-server/serving-certs/
tmpdir=/tmp/k8s-webhook-server/serving-certs/

mkdir -p ${tmpdir}

pushd ${tmpdir}
mkcert sample-tailer-webhook sample-tailer-webhook.default sample-tailer-webhook.default.svc localhost
mv sample-tailer-webhook+3.pem tls.crt
mv sample-tailer-webhook+3-key.pem tls.key
popd

ls -l ${tmpdir}