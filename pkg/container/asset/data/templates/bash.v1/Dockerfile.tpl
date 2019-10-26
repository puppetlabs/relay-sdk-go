{{- $FilePath := printf "/nebula/step-%s.sh" .Name -}}
FROM {{ .Settings.Image }}
RUN apk --no-cache add bash ca-certificates curl git jq openssh && update-ca-certificates
{{- if .Settings.AdditionalPackages }}
RUN apk --no-cache add{{ range .Settings.AdditionalPackages }} {{ . }}{{ end }}
{{- end }}
{{- range .Settings.AdditionalCommands }}
RUN ["/bin/bash", "-c", {{ . | mustToJson }}]
{{- end }}
COPY "./{{ .Settings.CommandPath }}" "{{ $FilePath }}"
CMD ["/bin/bash", "{{ $FilePath }}"]
