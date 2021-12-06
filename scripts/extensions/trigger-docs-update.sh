#!/bin/bash

set -euf

PROJECT_SLUG='gh/banzaicloud/one-eye-docs'
RELEASE_TAG="$1"

function main()
{
    curl \
        -u "${CIRCLE_TOKEN}:" \
        -X POST \
        --header "Content-Type: application/json" \
        -d "{
            \"branch\": \"master\",
            \"parameters\": {
                \"remote-trigger\": true,
                \"project\": \"logging-extensions\",
                \"release-tag\": \"${RELEASE_TAG}\",
                \"build-dir\": \"cmd/build/\"
            }
        }" "https://circleci.com/api/v2/project/${PROJECT_SLUG}/pipeline"
}

main "$@"
