{{- $PackageFileName := printf "ni-%s-linux-amd64.tar.xz" .SDKVersion -}}
{{- $PackageSHA256FileName := printf "%s.sha256.asc" $PackageFileName -}}
{{- $PackageRepoURL := printf "https://packages.nebula.puppet.net/ni/%s" .SDKVersion }}
FROM {{ .Containers.Base.Ref }}
RUN set -eux ; \
    mkdir -p /tmp/ni && \
    cd /tmp/ni && \
    wget {{ printf "%s/%s" $PackageRepoURL $PackageFileName }} && \
    wget {{ printf "%s/%s" $PackageRepoURL $PackageSHA256FileName }} && \
    sha256sum -c $PackageSHA256FileName && \
    tar -xvJf $PackageFileName && \
    mv ni-{{ .SDKVersion }}*-linux-amd64/ni /usr/local/bin/ni && \
    cd - && \
    rm -fr /tmp/ni
