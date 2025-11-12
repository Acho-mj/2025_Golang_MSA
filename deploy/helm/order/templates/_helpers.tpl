{{- define "order-service.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{- define "order-service.fullname" -}}
{{- $name := default .Chart.Name .Values.nameOverride -}}
{{- if .Values.fullnameOverride -}}
{{- .Values.fullnameOverride | trunc 63 | trimSuffix "-" -}}
{{- else -}}
{{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" -}}
{{- end -}}
{{- end -}}

{{- define "order-service.labels" -}}
app.kubernetes.io/name: {{ include "order-service.name" . }}
helm.sh/chart: {{ .Chart.Name }}-{{ .Chart.Version }}
app.kubernetes.io/instance: {{ .Release.Name }}
app.kubernetes.io/version: {{ .Chart.AppVersion }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end -}}

{{- define "order-service.selectorLabels" -}}
app.kubernetes.io/name: {{ include "order-service.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end -}}

{{- define "order-service.serviceAccountName" -}}
{{- if .Values.serviceAccount.create -}}
  {{- if .Values.serviceAccount.name -}}
    {{ .Values.serviceAccount.name }}
  {{- else -}}
    {{ include "order-service.fullname" . }}
  {{- end -}}
{{- else -}}
  {{- if .Values.serviceAccount.name -}}
    {{ .Values.serviceAccount.name }}
  {{- else -}}
    default
  {{- end -}}
{{- end -}}
{{- end -}}

