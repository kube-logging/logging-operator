#!/bin/sh

CHECK_INTERVAL="${CHECK_INTERVAL:-60}"
RPC_ADDRESS="${RPC_ADDRESS:-127.0.0.1:24444}"
CUSTOM_RUNNER_ADDRESS="${CUSTOM_RUNNER_ADDRESS:-127.0.0.1:7357}"

[ -z "$BUFFER_PATH" ] && exit 2

# this loop will not go on indefinitely because the fluentd RPC endpoint should
# come up eventually and won't terminate without a signal from outside (barring errors)
echo '['$(date)']' 'waiting for fluentd RPC endpoint to become available'
until netstat -tln | grep "$RPC_ADDRESS" >/dev/null
do
  [ -z "$DEBUG" ] && echo '['$(date)']' 'fluentd RPC endpoint not available, waiting'
  sleep 1
done

# this loop will not go on indefinitely because the custom-runner's HTTP endpoint should
# come up eventually and won't terminate without a signal from outside (barring errors)
echo '['$(date)']' 'waiting for custom-runner HTTP endpoint to become available'
until curl -so /dev/null ${CUSTOM_RUNNER_ADDRESS}
do 
  [ -z "$DEBUG" ] && echo '['$(date)']' 'custom-runner HTTP endpoint not available, waiting'
  sleep 1
done

echo '['$(date)']' 'waiting for fluentd to exit' # i.e. stop listening on the RPC address
while netstat -tln | grep "$RPC_ADDRESS" >/dev/null
do
  [ -z "$DEBUG" ] && echo '['$(date)']' 'RPC endpoint still listening'

  if [ "$(find $BUFFER_PATH -iname '*.buffer' -or -iname '*.buffer.meta' | wc -l)" = 0 ]
  then
    echo '['$(date)']' 'exiting node exporter custom runner:' "$(curl --silent --show-error http://$CUSTOM_RUNNER_ADDRESS/exit)"
    echo '['$(date)']' 'no buffers left, terminating workers:' "$(curl --silent --show-error http://$RPC_ADDRESS/api/processes.killWorkers)"
    exit 0
  fi

  sleep "$CHECK_INTERVAL"
done

echo '['$(date)']' 'checking for remaining buffers'
[ "$(find $BUFFER_PATH -iname '*.buffer' -or -iname '*.buffer.meta' | wc -l)" -gt 0 ] && exit 1

exit 0