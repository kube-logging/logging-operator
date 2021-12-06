#!/usr/bin/env bash

kubectl apply -f config/samples/logging-extensions_v1alpha1_promtail.yaml

helm install generator banzaicloud-stable/log-generator --wait

helm upgrade --install loki loki/loki-stack --set promtail.enabled=false --wait

kubectl run -i --tty logtest --image=grafana/logcli:master-0254070-amd64 --restart=Never -- query '{container="log-generator"}' --addr="http://loki:3100"
