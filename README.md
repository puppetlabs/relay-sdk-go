# nebula-sdk

This repository contains APIs and utility commands used to interact with Nebula
from inside a step container.

## Ni

Ni, short for Nebula Interface, is a command provided in every container running
in Nebula at the location `/nebula/bin/ni`. Usually, it is automatically added
to the container `$PATH`.

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

## Spindle

Spindle is a utility that generates Nebula-compatible Dockerfiles from YAML
configuration. Spindle uses an inheritance model, where configuration builds on
zero or more parent templates, to produce a unified configuration model. Once
the configuration is resolved, Spindle generates one or more Dockerfiles and a
Bash script that builds and tags images.

**Note:** It is not necessary to use Spindle to build Nebula-compatible Docker
images; in fact, you can use Nebula with any Docker image. If you are building
Nebula-oriented images, Spindle might make your content authoring experience
more straightforward.

### Getting Spindle

If your development process already uses [Go](https://golang.org/), it's easy to
embed Spindle into it. Simply run the command:

```
$ go run github.com/puppetlabs/nebula-sdk/cmd/spindle
# ...
```

You can automatically run Spindle using `go generate`. For an example, see the
[`hack/generate` directory in our step content repository](https://github.com/puppetlabs/nebula-steps/tree/master/hack/generate).

We also provide Spindle in binary form:

| SDK Version | Platform | Binary Archive | SHA-256 Checksum |
|-------------|----------|----------------|------------------|
| v1 | macOS/amd64 | [spindle-v1-darwin-amd64.tar.xz](https://packages.nebula.puppet.net/sdk/spindle/v1/spindle-v1-darwin-amd64.tar.xz) | [spindle-v1-darwin-amd64.tar.xz.sha256](https://packages.nebula.puppet.net/sdk/spindle/v1/spindle-v1-darwin-amd64.tar.xz.sha256) |
| | Linux/amd64 | [spindle-v1-linux-amd64.tar.xz](https://packages.nebula.puppet.net/sdk/spindle/v1/spindle-v1-linux-amd64.tar.xz) | [spindle-v1-linux-amd64.tar.xz.sha256](https://packages.nebula.puppet.net/sdk/spindle/v1/spindle-v1-linux-amd64.tar.xz.sha256) |
| | Linux/386 | [spindle-v1-linux-386.tar.xz](https://packages.nebula.puppet.net/sdk/spindle/v1/spindle-v1-linux-386.tar.xz) | [spindle-v1-linux-386.tar.xz.sha256](https://packages.nebula.puppet.net/sdk/spindle/v1/spindle-v1-linux-386.tar.xz.sha256) |
| | Linux/arm64 | [spindle-v1-linux-arm64.tar.xz](https://packages.nebula.puppet.net/sdk/spindle/v1/spindle-v1-linux-arm64.tar.xz) | [spindle-v1-linux-arm64.tar.xz.sha256](https://packages.nebula.puppet.net/sdk/spindle/v1/spindle-v1-linux-arm64.tar.xz.sha256) |
| | Linux/ppc64le | [spindle-v1-linux-ppc64le.tar.xz](https://packages.nebula.puppet.net/sdk/spindle/v1/spindle-v1-linux-ppc64le.tar.xz) | [spindle-v1-linux-ppc64le.tar.xz.sha256](https://packages.nebula.puppet.net/sdk/spindle/v1/spindle-v1-linux-ppc64le.tar.xz.sha256) |
| | Linux/s390x | [spindle-v1-linux-s390x.tar.xz](https://packages.nebula.puppet.net/sdk/spindle/v1/spindle-v1-linux-s390x.tar.xz) | [spindle-v1-linux-s390x.tar.xz.sha256](https://packages.nebula.puppet.net/sdk/spindle/v1/spindle-v1-linux-s390x.tar.xz.sha256) |
| | Windows/amd64 | [spindle-v1-windows-amd64.zip](https://packages.nebula.puppet.net/sdk/spindle/v1/spindle-v1-windows-amd64.zip) | [spindle-v1-windows-amd64.zip.sha256](https://packages.nebula.puppet.net/sdk/spindle/v1/spindle-v1-windows-amd64.zip.sha256) |

### Generating Dockerfiles with Spindle

Spindle comes with several templates that help to automatically generate
Dockerfiles:

| Name | File Reference | Description |
|------|----------------|-------------|
| `bash.v1` | `{from: sdk, name: bash.v1}` | Alpine Linux-based image configured to run a Bash script |
| `go.v1` | `{from: sdk, name: go.v1}` | Alpine Linux-based image configured to run a compiled Go binary |
| `python.v1` | `{from: sdk, name: python.v1}` | Alpine Linux-based image configured to run a Python script |

You can also write your own templates, or write templates that build from these.
For examples of customizing these templates, check out the
[`hack/defs` directory of our step content repository](https://github.com/puppetlabs/nebula-steps/tree/master/hack/defs).

To learn how to configure a template, you can use the `spindle desc template`
command:

```console
$ spindle desc template --from=sdk bash.v1
Images:
 NAME  TEMPLATE        DEPENDENCIES
 base  Dockerfile.tpl

Settings:
 NAME                DESCRIPTION                          DEFAULT VALUE
 Image               The Alpine Linux-based image to use  "alpine:3"
 AdditionalPackages  Additional APK packages to install   []
 AdditionalCommands  Additional Bash commands to run      []
 CommandPath         The path to the shell script to run  "step.sh"
$ spindle desc template ./hack/defs/aws
Images:
 NAME  TEMPLATE        DEPENDENCIES
 base  Dockerfile.tpl

Settings:
 NAME                DESCRIPTION                          DEFAULT VALUE
 Image               The Alpine Linux-based image to use  "python:3-alpine"
 AdditionalCommands  Additional Bash commands to run      ["pip install awscli"]
 CommandPath         The path to the shell script to run  "step.sh"
 AdditionalPackages  Additional APK packages to install   []
```

Once you have selected a template to drive from, you simply write a
`container.yaml` file that references it and overrides any settings. The Spindle
command knows how to select YAML files based on the content of the `apiVersion`
and `kind` keys. A minimal `container.yaml`, when paired with a Bash script called `step.sh`, is:

```yaml
apiVersion: container/v1
kind: StepContainer
inherit: {from: sdk, name: bash.v1}
title: My step
description: |
  This task does some very complex work.
```

### Example

Let's create a container that prints a fortune for us:

```console
$ mkdir fortune
$ cd fortune
$ cat >container.yaml <<'EOF'
apiVersion: container/v1
kind: StepContainer
inherit: {from: sdk, name: bash.v1}
title: Fortune
description: The fortune task prints a wholesome fortune to the log.
settings:
  AdditionalPackages: [fortune]
EOF
$ cat >step.sh <<'EOF'
#!/bin/bash
fortune
EOF
$ chmod +x step.sh
$ spindle gen --repo-name-base=example.com/nebula --write --verbose
scripts/build-container
Dockerfile
$ scripts/build-container -t latest
# [...]
# Tagged example.com/nebula/fortune:latest
$ docker run --rm example.com/nebula/fortune:latest
Q:  How many DEC repairmen does it take to fix a flat?
A:  Five; four to hold the car up and one to swap tires.
```

For more examples, take a look at the steps defined in our
[step content repository](https://github.com/puppetlabs/nebula-steps).

## Client APIs

This repository contains client APIs for interacting with the Nebula metadata
service in steps:

| SDK Version | Language | Source Code | Archive | SHA-256 Checksum |
|-------------|----------|-------------|---------|------------------|
| v1 | Python | [`support/python`](support/python) | [nebula-sdk-1.tar.gz](https://packages.nebula.puppet.net/sdk/support/python/v1/nebula-sdk-1.tar.gz) | [nebula-sdk-1.tar.gz.sha256](https://packages.nebula.puppet.net/sdk/support/python/v1/nebula-sdk-1.tar.gz.sha256) |
|||| [nebula_sdk-1-py3-none-any.whl](https://packages.nebula.puppet.net/sdk/support/python/v1/nebula_sdk-1-py3-none-any.whl) | [nebula_sdk-1-py3-none-any.whl.sha256](https://packages.nebula.puppet.net/sdk/support/python/v1/nebula_sdk-1-py3-none-any.whl.sha256) |
