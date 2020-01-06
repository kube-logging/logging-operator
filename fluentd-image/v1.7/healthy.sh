#!/bin/sh -x

# Liveness probe is aimed to help in situarions where fluentd
# silently hangs for no apparent reasons until manual restart.
# The idea of this probe is that if fluentd is not queueing or
# flushing chunks for 5 minutes, something is not right. If
# you want to change the fluentd configuration, reducing amount of
# logs fluentd collects, consider changing the threshold or turning
# liveness probe off completely.
# soiurce https://github.com/kubernetes/kubernetes/blob/master/cluster/addons/fluentd-gcp/fluentd-gcp-ds.yaml#L58

BUFFER_PATH=${BUFFER_PATH};
LIVENESS_THRESHOLD_SECONDS=${LIVENESS_THRESHOLD_SECONDS:-300};

if [ ! -e ${BUFFER_PATH} ];
then
  exit 1;
fi;
touch -d "${LIVENESS_THRESHOLD_SECONDS} seconds ago" /tmp/marker-liveness;
if [ -z "$(find ${BUFFER_PATH} -type d -newer /tmp/marker-liveness -print -quit)" ];
then
  exit 1;
fi;