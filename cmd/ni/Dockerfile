FROM golang:1.14-alpine AS builder
ENV GO111MODULE on
ENV CGO_ENABLED 0
WORKDIR /build
COPY . .
RUN go build -a -installsuffix cgo -o /usr/local/bin/ni ./cmd/ni

FROM alpine:latest
COPY --from=builder /usr/local/bin/ni /usr/local/bin/ni
RUN apk --no-cache add bash ca-certificates curl git jq openssh && update-ca-certificates
RUN ["/bin/bash", "-c", "set -eux\nmkdir -p /tmp/kubectl \u0026\u0026 \\\n  cd /tmp/kubectl \u0026\u0026 \\\n  curl -LO https://storage.googleapis.com/kubernetes-release/release/v1.16.0/bin/linux/amd64/kubectl \u0026\u0026 \\\n  echo \"$( curl -sL https://storage.googleapis.com/kubernetes-release/release/v1.16.0/bin/linux/amd64/kubectl.sha256 )  kubectl\" | sha256sum -c - \u0026\u0026 \\\n  chmod +x kubectl \u0026\u0026 \\\n  mv kubectl /usr/local/bin/kubectl \u0026\u0026 \\\n  cd - \u0026\u0026 \\\n  rm -fr /tmp/kubectl\n"]
CMD ["/usr/local/bin/ni"]
