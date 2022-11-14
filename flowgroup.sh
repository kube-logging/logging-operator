#!/usr/bin/env bash

kubectl apply -f -<<EOF
apiVersion: logging.banzaicloud.io/v1beta1
kind: Output
metadata:
  name: output-sample
spec:
  nullout: {}
EOF

kubectl apply -f -<<EOF
apiVersion: logging.banzaicloud.io/v1beta1
kind: Flow
metadata:
  name: flow-sample
  namespace: default
spec:
  filters:
    - parser:
        remove_key_name_field: true
        parse:
          type: nginx
    - tag_normaliser:
        format: ${namespace_name}.${pod_name}.${container_name}
  localOutputRefs:
    - output-sample
  match:
    - select:
        labels:
          app: nginx
EOF
