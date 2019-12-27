#!/bin/bash

set -eufo pipefail

SCRIPT_PATH="$(dirname "$(readlink -f "$0")")"
BUCKET='minio/test'

function main()
{
    local mc_pod="$(get_mc_pod_name)"
    create_test_bucket "${mc_pod}"
    kubectl apply -f "${SCRIPT_PATH}/secret.yaml"
    helm_deploy_logging_operator
    configure_logging

    wait_for_log_files "${mc_pod}" 300
    print_logs "${mc_pod}"
}

function create_test_bucket()
{
    local mc_pod="$1"
    kubectl exec --namespace logging "${mc_pod}" -- \
        mc mb --region 'test_region' "${BUCKET}"
}

function helm_deploy_logging_operator()
{
    helm install \
        --wait \
        --name logging-operator \
        --set image.tag='local' \
        --set image.repository='controller' \
        "${SCRIPT_PATH}/../charts/logging-operator"
}

function configure_logging()
{
    helm install \
        --wait \
        --namespace logging \
        --name 'logging-operator-logging-tls' \
        "${SCRIPT_PATH}/../charts/logging-operator-logging"
    kubectl apply -f "${SCRIPT_PATH}/clusteroutput.yaml"
    kubectl apply -f "${SCRIPT_PATH}/clusterflow.yaml"
}

function get_mc_pod_name()
{
    kubectl get pod --namespace logging -l app=minio-mc -o 'jsonpath={.items[0].metadata.name}'
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
    kubectl get pod,svc --namespace logging
    kubectl exec -it --namespace logging logging-operator-logging-fluentd-0 cat /fluentd/log/out
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

    kubectl exec --namespace logging "${mc_pod}" -- mc find "${BUCKET}"  --name '*.gz'
}

function print_logs()
{
    local mc_pod="$1"

    kubectl exec --namespace logging "${mc_pod}" -- mc find "${BUCKET}" --name '*.gz' -exec 'mc cat {}' | gzip -d
}

main "$@"
