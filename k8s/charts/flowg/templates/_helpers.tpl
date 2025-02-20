{{/*
Expand the name of the chart.
*/}}
{{- define "flowg.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "flowg.fullname" -}}
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

{{- define "flowg.fqdn" -}}
{{- $fullname := include "flowg.fullname" . -}}
{{- printf "%s.%s.svc.cluster.local" $fullname .Release.Namespace -}}
{{- end }}

{{/*
Create chart name and version as used by the chart label.
*/}}
{{- define "flowg.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Common labels
*/}}
{{- define "flowg.labels" -}}
helm.sh/chart: {{ include "flowg.chart" . }}
{{ include "flowg.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{/*
Selector labels
*/}}
{{- define "flowg.selectorLabels" -}}
app.kubernetes.io/name: {{ include "flowg.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{/*
Create the name of the service account to use for FlowG
*/}}
{{- define "flowg.serviceAccountName" -}}
{{- if .Values.flowg.serviceAccount.create }}
{{- default (include "flowg.fullname" .) .Values.flowg.serviceAccount.name }}
{{- else }}
{{- default "default" .Values.flowg.serviceAccount.name }}
{{- end }}
{{- end }}

{{/*
Lookup the Kubernetes node for the FlowG deployment
*/}}
{{- define "flowg.nodeName" -}}
{{- if or (not .Values.flowg.nodeName) (eq .Values.flowg.nodeName "") -}}
{{- fail "Missing value 'flowg.nodeName', it is required because FlowG does not support clustering yet" -}}
{{- else -}}
{{- .Values.flowg.nodeName -}}
{{- end -}}
{{- end -}}

{{/*
Fluentd component variables
*/}}
{{- define "fluentd.name" -}}
{{- printf "%s-fluentd" (include "flowg.name" .) -}}
{{- end -}}

{{- define "fluentd.fullname" -}}
{{- printf "%s-fluentd" (include "flowg.fullname" .) -}}
{{- end -}}

{{- define "fluentd.labels" -}}
helm.sh/chart: {{ include "flowg.chart" . }}
{{ include "fluentd.selectorLabels" . }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{- define "fluentd.selectorLabels" -}}
app.kubernetes.io/name: {{ include "fluentd.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end -}}

{{- define "fluentd.serviceAccountName" -}}
{{- if .Values.fluentd.serviceAccount.create }}
{{- default (include "fluentd.fullname" .) .Values.fluentd.serviceAccount.name }}
{{- else }}
{{- default "default" .Values.fluentd.serviceAccount.name }}
{{- end }}
{{- end -}}
