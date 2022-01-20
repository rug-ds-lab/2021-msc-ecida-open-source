package modulegen

import (
	"ecida/pkg/meta"
	"fmt"
	"os"
	"strings"

	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/chartutil"
)

func NewEcidaChart(name string) *chart.Metadata {
	metadata := meta.NewEcidaMetadata()

	return &chart.Metadata{
		Name:        name,
		Description: fmt.Sprintf("The %s ECiDA helmchart chart", name),
		APIVersion:  "v2",
		Version:     "0.1.0",
		Type:        "application",
		Annotations: metadata.ToChartAnnotations(),
	}
}

const DeploymentStub = `apiVersion: apps/v1
kind: Deployment
metadata:
  name: {{ include "__MODULE_NAME.fullname" . }}
  labels:
    {{- include "__MODULE_NAME.labels" . | nindent 4 }}
spec:
  selector:
    matchLabels:
      {{- include "__MODULE_NAME.selectorLabels" . | nindent 6 }}
  template:
    metadata:
      labels:
        {{- include "__MODULE_NAME.selectorLabels" . | nindent 8 }}
    spec:
      containers:
        - name: {{ .Chart.Name }}
          image: "{{ .Values.image.repository }}:{{ .Values.image.tag }}"
          ports:
            - name: http
              containerPort: {{ .Values.image.port }}
              protocol: TCP
`
const ServiceStub = `apiVersion: v1
kind: Service
metadata:
  name: {{ .Values.service.name }}
  labels:
    {{- include "__MODULE_NAME.labels" . | nindent 4 }}
spec:
  ports:
    - port: {{ .Values.service.port }}
      targetPort: http
      protocol: TCP
      name: http
  selector:
      {{- include "__MODULE_NAME.selectorLabels" . | nindent 4 }}
`

const ValuesStub = `service:
  name: "__MODULE_NAME"
  port: "5000"

# Configuration on the image associated with this ECiDA module
image:
  port: "3000"
  repository: "fillmein/myimagemodule"
  tag: "latest"
`

const TemplateStub = `{{/*
Expand the name of the chart.
*/}}
{{- define "__MODULE_NAME.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "__MODULE_NAME.fullname" -}}
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
{{- define "__MODULE_NAME.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Common labels
*/}}
{{- define "__MODULE_NAME.labels" -}}
helm.sh/chart: {{ include "__MODULE_NAME.chart" . }}
{{ include "__MODULE_NAME.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{/*
Selector labels
*/}}
{{- define "__MODULE_NAME.selectorLabels" -}}
app.kubernetes.io/name: {{ include "__MODULE_NAME.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}
`

func RenderStub(stub string, name string, path string) error {
	contents := strings.ReplaceAll(stub, "__MODULE_NAME", name)

	outfile, err := os.Create(path)

	if err != nil {
		return err
	}

	defer outfile.Close()

	outfile.WriteString(contents)

	return nil
}

func GenerateModule(name string, location string) error {
	if err := os.Mkdir(fmt.Sprintf("%s/templates", location), os.ModePerm); err != nil {
		return err
	}

	if err := RenderStub(DeploymentStub, name, fmt.Sprintf("%s/templates/deployment.yaml", location)); err != nil {
		return err
	}

	if err := RenderStub(ServiceStub, name, fmt.Sprintf("%s/templates/service.yaml", location)); err != nil {
		return err
	}

	if err := RenderStub(TemplateStub, name, fmt.Sprintf("%s/templates/_helpers.tpl", location)); err != nil {
		return err
	}

	if err := RenderStub(ValuesStub, name, fmt.Sprintf("%s/values.yaml", location)); err != nil {
		return err
	}

	chartFile := NewEcidaChart(name)
	chartutil.SaveChartfile(fmt.Sprintf("%s/Chart.yaml", location), chartFile)

	return nil
}
