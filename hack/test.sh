#!/bin/bash

set -eufo pipefail
set -x

SCRIPT_PATH="hack"
BUCKET='minio/test'

function main()
{
    load_images
    local mc_pod="$(get_mc_pod_name)"
    remove_test_bucket "${mc_pod}"
    create_test_bucket "${mc_pod}"
    kubectl apply -n logging -f "${SCRIPT_PATH}/secret.yaml"
    helm_deploy_logging_operator
    configure_logging

    wait_for_log_files "${mc_pod}" 300
    print_logs "${mc_pod}"
}

function load_images()
{
    local images=( "controller:local")
    for image in ${images[@]}; do
        kind load docker-image "${image}"
    done
}

function remove_test_bucket()
{
    local mc_pod="$1"
    kubectl exec -n logging "${mc_pod}" -- \
        mc rb "${BUCKET}" || true
}

function create_test_bucket()
{
    local mc_pod="$1"
    kubectl exec -n logging "${mc_pod}" -- \
        mc mb --region 'test_region' "${BUCKET}"
}

function helm_add_repo()
{
    helm repo add kube-logging https://kube-logging.github.io/helm-charts
    helm repo update kube-logging
}

function helm_deploy_logging_operator()
{
    helm upgrade --install \
        --debug \
        --wait \
        --set image.tag='local' \
        --set image.repository='controller' \
        logging-operator \
        "e2e/charts/logging-operator"
}

function configure_logging()
{
    helm upgrade --install \
        --debug \
        --wait \
        --create-namespace \
        -n logging \
        'logging-operator-logging-tls' \
        "e2e/charts/logging-operator-logging"
    kubectl apply -n logging -f "${SCRIPT_PATH}/clusteroutput.yaml"
    kubectl apply -n logging -f "${SCRIPT_PATH}/clusterflow.yaml"
}

function get_mc_pod_name()
{
    kubectl get pod -n logging -l app=minio-mc -o 'jsonpath={.items[0].metadata.name}'
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
    kubectl get pod,svc -n logging
    kubectl exec -it -n logging logging-operator-logging-fluentd-0 cat /fluentd/log/out
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

    kubectl exec -n logging "${mc_pod}" -- mc find "${BUCKET}"  --name '*.gz'
}

function print_logs()
{
    local mc_pod="$1"

    kubectl exec -n logging "${mc_pod}" -- mc find "${BUCKET}" --name '*.gz' -exec 'mc cat {}' | gzip -d
}

main "$@"
