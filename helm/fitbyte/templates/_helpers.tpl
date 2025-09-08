{{/*
Expand the name of the chart.
*/}}
{{- define "fitbyte.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "fitbyte.fullname" -}}
{{- if .Values.fullnameOverride }}
{{- .Values.fullnameOverride | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- $name := default .Chart.Name .Values.nameOverride }}
{{- if contains $name .Release.Name }}
{{- .Release.Name | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" }}
{{- end }}
{{- end }}
{{- end }}

{{/*
Create chart name and version as used by the chart label.
*/}}
{{- define "fitbyte.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Common labels
*/}}
{{- define "fitbyte.labels" -}}
helm.sh/chart: {{ include "fitbyte.chart" . }}
{{ include "fitbyte.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{/*
Selector labels
*/}}
{{- define "fitbyte.selectorLabels" -}}
app.kubernetes.io/name: {{ include "fitbyte.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{/*
Create the name of the service account to use
*/}}
{{- define "fitbyte.serviceAccountName" -}}
{{- if .Values.serviceAccount.create }}
{{- default (include "fitbyte.fullname" .) .Values.serviceAccount.name }}
{{- else }}
{{- default "default" .Values.serviceAccount.name }}
{{- end }}
{{- end }}

{{/*
PostgreSQL fullname
*/}}
{{- define "fitbyte.postgresql.fullname" -}}
{{- if .Values.database.postgresql.enabled }}
{{- printf "%s-postgresql" (include "fitbyte.fullname" .) }}
{{- else }}
{{- .Values.database.postgresql.external.host }}
{{- end }}
{{- end }}

{{/*
MinIO fullname
*/}}
{{- define "fitbyte.minio.fullname" -}}
{{- if .Values.minio.enabled }}
{{- printf "%s-minio" (include "fitbyte.fullname" .) }}
{{- else }}
{{- .Values.minio.external.host }}
{{- end }}
{{- end }}

{{/*
Prometheus fullname
*/}}
{{- define "fitbyte.prometheus.fullname" -}}
{{- if .Values.prometheus.enabled }}
{{- printf "%s-prometheus" (include "fitbyte.fullname" .) }}
{{- else }}
{{- .Values.prometheus.external.host }}
{{- end }}
{{- end }}

{{/*
Grafana fullname
*/}}
{{- define "fitbyte.grafana.fullname" -}}
{{- if .Values.grafana.enabled }}
{{- printf "%s-grafana" (include "fitbyte.fullname" .) }}
{{- else }}
{{- .Values.grafana.external.host }}
{{- end }}
{{- end }}
