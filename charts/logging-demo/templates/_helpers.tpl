{{/* vim: set filetype=mustache: */}}
{{/*
Expand the name of the chart.
*/}}
{{- define "logging-demo.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "logging-demo.fullname" -}}
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

{{- define "flow.fullname" -}}
{{- printf "%s-%s" (include "logging-demo.fullname" .) "flow" | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{- define "output.fullname" -}}
{{- printf "%s-%s" (include "logging-demo.fullname" .) "output" | trunc 63 | trimSuffix "-" -}}
{{- end -}}


{{/*
Create chart name and version as used by the chart label.
*/}}
{{- define "logging-demo.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{/*
Common labels
*/}}
{{- define "logging-demo.labels" -}}
app.kubernetes.io/name: {{ include "logging-demo.name" . }}
helm.sh/chart: {{ include "logging-demo.chart" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end -}}

{{/*
Formats the cluster domain as a suffix, e.g.:
.Values.loggingOperator.clusterDomain == "", returns ""
.Values.loggingOperator.clusterDomain == "cluster.local", returns ".cluster.local"
*/}}
{{- define "logging-demo.clusterDomainAsSuffix" -}}
{{- if .Values.loggingOperator.clusterDomain -}}
{{- printf ".%s" .Values.loggingOperator.clusterDomain -}}
{{- end -}}
{{- end -}}
