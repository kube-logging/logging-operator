#!/usr/bin/env bash

# helm install -f charts/logging-operator/values.yaml ./charts/logging-operator --generate-name

# helm install -f charts/logging-operator/values.yaml foobar ./charts/logging-operator --dry-run

helm install -f charts/logging-operator/values.yaml operator ./charts/logging-operator

helm install -f charts/logging-operator-logging/values.yaml foobar ./charts/logging-operator-logging --dry-run --debug

