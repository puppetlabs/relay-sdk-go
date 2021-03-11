# Relay Go SDK

This repository contains APIs and utility commands used to interact with Relay
from inside a step container.

## Ni

Ni is a command provided in every container running in Relay at the location
`/var/lib/puppet/relay/ni`. Usually, it is automatically added to the container `$PATH`.

### Getting Ni

If you want to bundle Ni with your container or test it outside of a
containerized environment, you may download a binary copy of Ni:

| SDK Version | Platform | Binary Archive | SHA-256 Checksum |
|-------------|----------|----------------|------------------|
| v1 | macOS/amd64 | [ni-v1-darwin-amd64.tar.xz](https://packages.nebula.puppet.net/sdk/ni/v1/ni-v1-darwin-amd64.tar.xz) | [ni-v1-darwin-amd64.tar.xz.sha256](https://packages.nebula.puppet.net/sdk/ni/v1/ni-v1-darwin-amd64.tar.xz.sha256) |
| | Linux/amd64 | [ni-v1-linux-amd64.tar.xz](https://packages.nebula.puppet.net/sdk/ni/v1/ni-v1-linux-amd64.tar.xz) | [ni-v1-linux-amd64.tar.xz.sha256](https://packages.nebula.puppet.net/sdk/ni/v1/ni-v1-linux-amd64.tar.xz.sha256) |
| | Linux/386 | [ni-v1-linux-386.tar.xz](https://packages.nebula.puppet.net/sdk/ni/v1/ni-v1-linux-386.tar.xz) | [ni-v1-linux-386.tar.xz.sha256](https://packages.nebula.puppet.net/sdk/ni/v1/ni-v1-linux-386.tar.xz.sha256) |
| | Linux/arm64 | [ni-v1-linux-arm64.tar.xz](https://packages.nebula.puppet.net/sdk/ni/v1/ni-v1-linux-arm64.tar.xz) | [ni-v1-linux-arm64.tar.xz.sha256](https://packages.nebula.puppet.net/sdk/ni/v1/ni-v1-linux-arm64.tar.xz.sha256) |
| | Linux/ppc64le | [ni-v1-linux-ppc64le.tar.xz](https://packages.nebula.puppet.net/sdk/ni/v1/ni-v1-linux-ppc64le.tar.xz) | [ni-v1-linux-ppc64le.tar.xz.sha256](https://packages.nebula.puppet.net/sdk/ni/v1/ni-v1-linux-ppc64le.tar.xz.sha256) |
| | Linux/s390x | [ni-v1-linux-s390x.tar.xz](https://packages.nebula.puppet.net/sdk/ni/v1/ni-v1-linux-s390x.tar.xz) | [ni-v1-linux-s390x.tar.xz.sha256](https://packages.nebula.puppet.net/sdk/ni/v1/ni-v1-linux-s390x.tar.xz.sha256) |
| | Windows/amd64 | [ni-v1-windows-amd64.zip](https://packages.nebula.puppet.net/sdk/ni/v1/ni-v1-windows-amd64.zip) | [ni-v1-windows-amd64.zip.sha256](https://packages.nebula.puppet.net/sdk/ni/v1/ni-v1-windows-amd64.zip.sha256) |
