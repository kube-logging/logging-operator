#!/usr/bin/env bash

stderr() {
    printf "$*" >&2
}

cmd="$1"
wait=${2-2}
timeout=${3-30}
startat=$(date +%s)
stderr "polling '$cmd'"
while true; do
    now=$(date +%s)
    stderr "."
    eval "$cmd" > /dev/null
    if [[ $? -eq 0 ]]; then
        break
    elif [[ $now -ge $((startat+timeout)) ]]; then
        stderr "timed out: (wait:$wait timeout:$timeout)\n"
        exit 1
    else
        sleep $wait
    fi
done

stderr "OK\n"
exit 0