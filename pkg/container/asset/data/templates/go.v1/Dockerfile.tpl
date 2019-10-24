{{- $Binary := printf "/usr/local/bin/nebula-%s" .Name -}}
FROM golang:{{ .Settings.GoVersion }}-alpine AS builder
ENV CGO_ENABLED 0
WORKDIR /build
COPY . .
RUN go build -a -o "{{ $Binary }}" "./{{ .Settings.CommandPath }}"

FROM {{ .Settings.Image }}
RUN apk --no-cache add ca-certificates && update-ca-certificates
COPY --from=builder "{{ $Binary }}" "{{ $Binary }}"
CMD ["{{ $Binary }}"]
