#!/bin/bash

set -eufo pipefail

SCRIPT_PATH="$(dirname "$(readlink -f "$0")")"
BUCKET='minio/logs'

function main()
{
    helm_deploy_logging_operator

    apply_s3_output
    mc_pod="$(get_mc_pod_name)"
    wait_for_log_files "${mc_pod}" 300
    print_logs "${mc_pod}"
}

function helm_deploy_logging_operator()
{
    helm install \
        --wait \
        --name logging-operator \
        --set image.tag='local' \
        banzaicloud-stable/logging-operator
}

function apply_s3_output()
{
    kubectl apply -f "${SCRIPT_PATH}/test-s3-output.yaml"
}

function get_mc_pod_name()
{
    kubectl get pod -l app=minio-mc -o 'jsonpath={.items[0].metadata.name}'
}

function wait_for_log_files()
{
    local mc_pod="$1"
    local deadline="$(( $(date +%s) + $2 ))"

    echo 'Waiting for log files...'
    while [ $(date +%s) -lt ${deadline} ]; do
        if [ $(count_log_files "${mc_pod}") -gt 0 ]; then
            return
        fi
        sleep 5
    done

    echo 'Cannot find any log files within timeout'
    exit 1
}

function count_log_files()
{
    local mc_pod="$1"

    get_log_files "${mc_pod}" |  wc -l
}

function get_log_files()
{
    local mc_pod="$1"

    kubectl exec "${mc_pod}" -- mc find "${BUCKET}"  --name '*.gz'
}

function print_logs()
{
    local mc_pod="$1"

    kubectl exec "${mc_pod}" -- mc find "${BUCKET}" --name '*.gz' -exec 'mc cat {}' | gzip -d
}

main "$@"
