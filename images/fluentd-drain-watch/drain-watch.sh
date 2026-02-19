#!/bin/sh

CHECK_INTERVAL="${CHECK_INTERVAL:-60}"
RPC_ADDRESS="${RPC_ADDRESS:-127.0.0.1:24444}"
CUSTOM_RUNNER_ADDRESS="${CUSTOM_RUNNER_ADDRESS:-127.0.0.1:7357}"
CUSTOM_RUNNER_TIMEOUT="${CUSTOM_RUNNER_TIMEOUT:-30}"
KILL_TIMEOUT="${KILL_TIMEOUT:-300}"

[ -z "$BUFFER_PATH" ] && exit 2

# this loop will not go on indefinitely because the fluentd RPC endpoint should
# come up eventually and won't terminate without a signal from outside (barring errors)
echo '['$(date)']' 'waiting for fluentd RPC endpoint to become available'
until netstat -tln | grep "$RPC_ADDRESS" >/dev/null
do
  [ -z "$DEBUG" ] && echo '['$(date)']' 'fluentd RPC endpoint not available, waiting'
  sleep 1
done

# Wait for custom-runner with a timeout since it may not be deployed
# (it's only present when buffer volume metrics are enabled)
echo '['$(date)']' 'waiting for custom-runner HTTP endpoint to become available (timeout: '${CUSTOM_RUNNER_TIMEOUT}'s)'
CUSTOM_RUNNER_AVAILABLE=false
elapsed=0
while [ $elapsed -lt $CUSTOM_RUNNER_TIMEOUT ]
do
  if curl -so /dev/null ${CUSTOM_RUNNER_ADDRESS}
  then
    CUSTOM_RUNNER_AVAILABLE=true
    echo '['$(date)']' 'custom-runner HTTP endpoint is available'
    break
  fi
  [ -z "$DEBUG" ] && echo '['$(date)']' 'custom-runner HTTP endpoint not available, waiting ('$elapsed's/'${CUSTOM_RUNNER_TIMEOUT}'s)'
  sleep 1
  elapsed=$((elapsed + 1))
done

if [ "$CUSTOM_RUNNER_AVAILABLE" = "false" ]
then
  echo '['$(date)']' 'custom-runner HTTP endpoint not available after '${CUSTOM_RUNNER_TIMEOUT}' seconds, assuming it is not deployed (buffer metrics sidecar disabled)'
fi

echo '['$(date)']' 'waiting for fluentd to exit' # i.e. stop listening on the RPC address
WORKERS_KILLED=false
kill_elapsed=0

while netstat -tln | grep "$RPC_ADDRESS" >/dev/null
do
  [ -z "$DEBUG" ] && echo '['$(date)']' 'RPC endpoint still listening'

  if [ "$WORKERS_KILLED" = "true" ]
  then
    kill_elapsed=$((kill_elapsed + CHECK_INTERVAL))
    if [ "$kill_elapsed" -ge "$KILL_TIMEOUT" ]
    then
      echo '['$(date)']' 'ERROR: fluentd did not shut down within '${KILL_TIMEOUT}'s after killing workers'
      exit 1
    fi
  fi

  if [ "$(find $BUFFER_PATH -iname '*.buffer' -or -iname '*.buffer.meta' | wc -l)" = 0 ]
  then
    if [ "$WORKERS_KILLED" = "false" ]
    then
      if [ "$CUSTOM_RUNNER_AVAILABLE" = "true" ]
      then
        echo '['$(date)']' 'exiting node exporter custom runner:' "$(curl --silent --show-error http://$CUSTOM_RUNNER_ADDRESS/exit)"
      fi
      echo '['$(date)']' 'no buffers left, terminating workers:' "$(curl --silent --show-error http://$RPC_ADDRESS/api/processes.killWorkers)"
      WORKERS_KILLED=true
      echo '['$(date)']' 'waiting for fluentd to shut down gracefully'
    fi
  fi

  sleep "$CHECK_INTERVAL"
done

echo '['$(date)']' 'fluentd has stopped listening on RPC endpoint'
echo '['$(date)']' 'checking for remaining buffers'
[ "$(find $BUFFER_PATH -iname '*.buffer' -or -iname '*.buffer.meta' | wc -l)" -gt 0 ] && exit 1

exit 0
