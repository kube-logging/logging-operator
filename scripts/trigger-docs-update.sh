#!/bin/bash

set -euf

PROJECT_SLUG='gh/banzaicloud/logging-operator-docs'
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
                \"generated-docs-update\": true,
                \"release-tag\": \"${RELEASE_TAG}\"
            }
        }" "https://circleci.com/api/v2/project/${PROJECT_SLUG}/pipeline"
}

main "$@"
