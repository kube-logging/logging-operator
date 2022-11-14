#!/bin/bash

set -euf

OWNER='banzaicloud'
REPO='logging-operator-docs'
WORKFLOW='generate-docs.yml'
RELEASE_TAG="$1"

function main()
{
    curl \
      -X POST \
      -H "Accept: application/vnd.github+json" \
      -H "Authorization: token ${GITHUB_TOKEN}" \
      "https://api.github.com/repos/${OWNER}/${REPO}/actions/workflows/${WORKFLOW}/dispatches" \
      -d "{\"ref\":\"master\",\"inputs\":{\"release-tag\":\"${RELEASE_TAG}\"}}"
}

main "$@"
