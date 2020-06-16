#!/usr/bin/env bash

#
# Commands
#

FIND="${FIND:-find}"
GIT="${GIT:-git}"
GSUTIL="${GSUTIL:-gsutil}"
SHA256SUM="${SHA256SUM:-shasum -a 256}"

#
#
#

relay::sdk::go::default_programs() {
  local DEFAULT_PROGRAMS
  DEFAULT_PROGRAMS=( ni )

  for DEFAULT_PROGRAM in ${DEFAULT_PROGRAMS[@]}; do
    printf "%s\n" "${DEFAULT_PROGRAM}"
  done
}

relay::sdk::go::git_tag() {
  printf "%s\n" "${GIT_TAG_OVERRIDE:-$( $GIT tag --points-at HEAD 'v*.*.*' )}"
}

relay::sdk::go::sha256sum() {
  $SHA256SUM | cut -d' ' -f1
}

relay::sdk::go::escape_shell() {
  printf '%s\n' "'${*//\'/\'\"\'\"\'}'"
}

relay::sdk::go::release_version() {
  local GIT_TAG GIT_CHANGED_FILES
  GIT_TAG="$( relay::sdk::go::git_tag )"
  GIT_CHANGED_FILES="$( $GIT status --short )"

  # Check for releasable version: if we have no tags or any changed files, we
  # can't release.
  if [ -z "${GIT_TAG}" ] || [ -n "${GIT_CHANGED_FILES}" ]; then
    return 1
  fi

  # Arbitrarily pick the first line.
  read GIT_TAG_A <<<"${GIT_TAG}"

  printf "%s\n" "${GIT_TAG_A#v}"
}

relay::sdk::go::release_check() {
  if ! relay::sdk::go::release_version >/dev/null; then
    echo "$0: no release tag (this commit must be tagged with the format vX.Y.Z)" >&2
    return 2
  fi
}

relay::sdk::go::release_vars() {
  RELEASE_VERSION="$( relay::sdk::go::release_version || true )"
  if [ -z "${RELEASE_VERSION}" ]; then
    printf 'RELEASE_VERSION=\n'
    return
  fi

  # Parse the version information.
  IFS='.' read RELEASE_VERSION_MAJOR RELEASE_VERSION_MINOR RELEASE_VERSION_PATCH <<<"${RELEASE_VERSION}"

  printf 'RELEASE_VERSION=%s\n' "$( relay::sdk::go::escape_shell "${RELEASE_VERSION}" )"
  printf 'RELEASE_VERSION_MAJOR=%s\n' "$( relay::sdk::go::escape_shell "${RELEASE_VERSION_MAJOR}" )"
  printf 'RELEASE_VERSION_MINOR=%s\n' "$( relay::sdk::go::escape_shell "${RELEASE_VERSION_MINOR}" )"
  printf 'RELEASE_VERSION_PATCH=%s\n' "$( relay::sdk::go::escape_shell "${RELEASE_VERSION_PATCH}" )"
}

relay::sdk::go::release_vars_local() {
  printf 'local RELEASE_VERSION RELEASE_VERSION_MAJOR RELEASE_VERSION_MINOR RELEASE_VERSION_PATCH\n'
  relay::sdk::go::release_vars "$@"
}

relay::sdk::go::release() {
  if [[ "$#" -lt 3 ]]; then
    echo "usage: ${FUNCNAME[0]} <bucket> <release-name> <filename> [dist-ext [dist-prefix]]" >&2
    return 1
  fi

  relay::sdk::go::release_check
  eval "$( relay::sdk::go::release_vars )"

  local KEY_PREFIX FILENAME DIST_PREFIX DIST_EXT
  KEY_PREFIX="gs://$1/sdk/$2"
  FILENAME="$3"
  DIST_EXT="${4:-}"
  DIST_PREFIX="${5:-"$2-v"}"

  (
    set -x

    local KEY KEY_MAJOR_MINOR KEY_MAJOR
    KEY="${KEY_PREFIX}/v${RELEASE_VERSION}/${DIST_PREFIX}${RELEASE_VERSION}${DIST_EXT}"
    KEY_MAJOR_MINOR="${KEY_PREFIX}/v${RELEASE_VERSION_MAJOR}.${RELEASE_VERSION_MINOR}/${DIST_PREFIX}${RELEASE_VERSION_MAJOR}.${RELEASE_VERSION_MINOR}${DIST_EXT}"
    KEY_MAJOR="${KEY_PREFIX}/v${RELEASE_VERSION_MAJOR}/${DIST_PREFIX}${RELEASE_VERSION_MAJOR}${DIST_EXT}"

    $GSUTIL cp "${FILENAME}" "${KEY}"
    $GSUTIL cp "${KEY}" "${KEY_MAJOR_MINOR}"
    $GSUTIL cp "${KEY}" "${KEY_MAJOR}"
  )
}

relay::sdk::go::version() {
  eval "$( relay::sdk::go::release_vars )"

  if [ -n "${RELEASE_VERSION}" ]; then
    printf "%s\n" "v${RELEASE_VERSION}"
  else
    $GIT describe --tags --always --dirty
  fi
}

relay::sdk::go::cli_vars() {
  if [[ "$#" -ne 1 ]]; then
    echo "usage: ${FUNCNAME[0]} <program>" >&2
    return 1
  fi

  local GO GOOS GOARCH
  GO="${GO:-go}"
  GOOS="$( $GO env GOOS )"
  GOARCH="$( $GO env GOARCH )"

  local EXT=
  [[ "${GOOS}" == "windows" ]] && EXT=.exe

  printf 'CLI_NAME=%s\n' "$( relay::sdk::go::escape_shell "$1" )"
  printf 'CLI_VERSION=%s\n' "$( relay::sdk::go::version )"
  printf 'CLI_FILE_PREFIX="${CLI_NAME}-${CLI_VERSION}"-%s-%s\n' \
    "$( relay::sdk::go::escape_shell "${GOOS}" )" \
    "$( relay::sdk::go::escape_shell "${GOARCH}" )"
  printf 'CLI_FILE_BIN="${CLI_FILE_PREFIX}%s"\n' "${EXT}"
}

relay::sdk::go::cli_vars_local() {
  printf 'local CLI_NAME CLI_FILE_PREFIX CLI_FILE_BIN\n'
  relay::sdk::go::cli_vars "$@"
}

relay::sdk::go::cli_artifacts() {
  if [[ "$#" -ne 2 ]]; then
    echo "usage: ${FUNCNAME[0]} <program> <directory>" >&2
    return 1
  fi

  eval "$( relay::sdk::go::cli_vars_local "$1" )"

  local CLI_MATCH
  CLI_MATCH="${CLI_NAME}-${CLI_VERSION}-"

  $FIND "$2" -type f -name "${CLI_MATCH}"'*.tar.xz' -or -name "${CLI_MATCH}"'*.zip'
}

relay::sdk::go::cli_platform_ext() {
  if [[ "$#" -ne 2 ]]; then
    echo "usage: ${FUNCNAME[0]} <program> <package-file>" >&2
    return 1
  fi

  eval "$( relay::sdk::go::cli_vars_local "$1" )"

  local CLI_FILE
  CLI_FILE="$( basename "$2" )"

  printf "%s\n" "${CLI_FILE##${CLI_NAME}-${CLI_VERSION}-}"
}

relay::sdk::go::usage() {
  echo "usage: $*" >&2
  exit 1
}
