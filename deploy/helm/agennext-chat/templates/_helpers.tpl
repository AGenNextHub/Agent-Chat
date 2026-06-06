{{- define "agennext-chat.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{- define "agennext-chat.fullname" -}}
{{- printf "%s-%s" .Release.Name (include "agennext-chat.name" .) | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{- define "agennext-chat.labels" -}}
app.kubernetes.io/name: {{ include "agennext-chat.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
helm.sh/chart: {{ printf "%s-%s" .Chart.Name .Chart.Version }}
{{- end -}}

{{- define "agennext-chat.selectorLabels" -}}
app.kubernetes.io/name: {{ include "agennext-chat.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end -}}

{{- define "agennext-chat.serviceAccountName" -}}
{{- if .Values.serviceAccount.create -}}
{{- default (include "agennext-chat.fullname" .) .Values.serviceAccount.name -}}
{{- else -}}
{{- default "default" .Values.serviceAccount.name -}}
{{- end -}}
{{- end -}}
