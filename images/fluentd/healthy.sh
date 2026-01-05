#!/bin/sh -x

# Liveness probe is aimed to help in situations where fluentd
# silently hangs for no apparent reasons until manual restart.
# The idea of this probe is that if fluentd is not queueing or
# flushing chunks for 5 minutes, something is not right. If
# you want to change the fluentd configuration, reducing amount of
# logs fluentd collects, consider changing the threshold or turning
# liveness probe off completely.
# soiurce https://github.com/kubernetes/kubernetes/blob/master/cluster/addons/fluentd-gcp/fluentd-gcp-ds.yaml#L58

BUFFER_PATH=${BUFFER_PATH:-/buffers};
LIVENESS_THRESHOLD_SECONDS=${LIVENESS_THRESHOLD_SECONDS:-300};

if [ ! -e "${BUFFER_PATH}" ]; then
  exit 1;
fi;
MINUTES=$(( (LIVENESS_THRESHOLD_SECONDS + 59) / 60 ));
if [ -z "$(find "${BUFFER_PATH}" -type d -mmin -"${MINUTES}" -print -quit)" ]; then
  exit 1;
fi;
