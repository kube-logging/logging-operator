{{/* vim: set filetype=mustache: */}}
{{/*
Expand the name of the chart.
*/}}
{{- define "logging-operator.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "logging-operator.fullname" -}}
{{- if .Values.fullnameOverride -}}
{{- .Values.fullnameOverride | trunc 63 | trimSuffix "-" -}}
{{- else -}}
{{- $name := default .Chart.Name .Values.nameOverride -}}
{{- if contains $name .Release.Name -}}
{{- .Release.Name | trunc 63 | trimSuffix "-" -}}
{{- else -}}
{{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" -}}
{{- end -}}
{{- end -}}
{{- end -}}

{{/*
Expand the name of the release.
*/}}
{{- define "logging-operator.releasename" -}}
{{- default .Release.Name .Values.nameOverride | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{/*
Provides the namespace the chart will be installed in using the builtin .Release.Namespace,
or, if provided, a manually overwritten namespace value.
*/}}
{{- define "logging-operator.namespace" -}}
{{- if .Values.namespaceOverride -}}
{{ .Values.namespaceOverride -}}
{{- else -}}
{{ .Release.Namespace }}
{{- end -}}
{{- end -}}


{{/*
Create chart name and version as used by the chart label.
*/}}
{{- define "logging-operator.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{/*
Common labels
*/}}
{{- define "logging-operator.labels" -}}
app.kubernetes.io/name: {{ include "logging-operator.name" . }}
helm.sh/chart: {{ include "logging-operator.chart" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end -}}

{{/*
Formats the cluster domain as a suffix, e.g.:
.Values.logging.clusterDomain == "", returns ""
.Values.logging.clusterDomain == "cluster.local.", returns ".cluster.local."
*/}}
{{- define "logging-operator.clusterDomainAsSuffix" -}}
{{- if .Values.logging.clusterDomain -}}
{{- printf ".%s" .Values.logging.clusterDomain -}}
{{- end -}}
{{- end -}}

{{/*
Replace the registry of an image with global.registry when it is set.
An image without a registry, like fluent/fluent-bit, gets the registry added in front.
Usage: include "logging-operator.image" (dict "registry" .Values.global.registry "image" .Values.image.repository)
*/}}
{{- define "logging-operator.image" -}}
{{- $registry := .registry | default "" | trimSuffix "/" -}}
{{- $image := .image -}}
{{- if $registry -}}
{{- $parts := splitList "/" $image -}}
{{- $host := first $parts -}}
{{- if and (gt (len $parts) 1) (or (contains "." $host) (contains ":" $host) (eq $host "localhost")) -}}
{{- $image = rest $parts | join "/" -}}
{{- end -}}
{{- printf "%s/%s" $registry $image -}}
{{- else -}}
{{- $image -}}
{{- end -}}
{{- end -}}
