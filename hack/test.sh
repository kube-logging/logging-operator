#!/bin/bash

set -eufo pipefail
set -x

SCRIPT_PATH="hack"

function main()
{
    kubectl get namespace logging || kubectl create namespace logging
    kubectl config set-context --current --namespace=logging

    load_images
    helm_deploy_logging_operator

    test_pod=$(get_test_pod_name)

    wait_for_logs "${test_pod}" 300
}

function load_images()
{
    local images=( "controller:local")
    for image in ${images[@]}; do
        kind load docker-image "${image}"
    done
}

function helm_deploy_logging_operator()
{
    helm upgrade --install \
        --debug \
        --wait \
        --create-namespace \
        -f hack/values.yaml \
        logging-operator \
        "charts/logging-operator"
}

function get_test_pod_name()
{
    kubectl get pod -l app.kubernetes.io/name=e2e-test-receiver -o 'jsonpath={.items[0].metadata.name}'
}

function wait_for_logs()
{
    local test_pod="$1"
    local deadline="$(( $(date +%s) + $2 ))"

    echo 'Waiting for log files...'
    while [ $(date +%s) -lt ${deadline} ]; do
        if [[ $(kubectl logs -l app.kubernetes.io/name=e2e-test-receiver | grep e2e.tag | wc -l) -gt 0 ]]; then
            return
        fi
        sleep 2
    done

    echo 'Cannot find any log files within timeout'
    kubectl get pod,svc
    kubectl exec -it e2e-fluentd-0 cat /fluentd/log/out
    exit 1
}

main "$@"
