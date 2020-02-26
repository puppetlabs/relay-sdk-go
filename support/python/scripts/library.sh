#!/usr/bin/env bash

#
# Commands
#

FIND="${FIND:-find}"
PYTHON="${PYTHON:-python}"
SED="${SED:-sed}"

#
#
#

. ../../scripts/library.sh

nebula::sdk::support_python::package_name() {
  $PYTHON setup.py --name
}

nebula::sdk::support_python::package_name_wheel() {
  # See https://www.python.org/dev/peps/pep-0427/
  nebula::sdk::support_python::package_name | $SED -Ee 's/[^A-Za-z0-9_.]+/_/g'
}

nebula::sdk::support_python::package_artifacts() {
  if [[ "$#" -ne 1 ]]; then
    echo "usage: ${FUNCNAME[0]} <directory>" >&2
    return 1
  fi

  local PACKAGE_VERSION
  PACKAGE_VERSION="$( $PYTHON setup.py --version )"

  $FIND "$1" \
    -type f \
    -name "$( nebula::sdk::support_python::package_name )-${PACKAGE_VERSION}"'*.tar.gz' \
    -or -name "$( nebula::sdk::support_python::package_name_wheel )-${PACKAGE_VERSION}"'*.whl'
}

# If we're in a CI context, we may need to force the version.
if [ -n "${GIT_TAG_OVERRIDE:-}" ]; then
  export SETUPTOOLS_SCM_PRETEND_VERSION="${GIT_TAG_OVERRIDE##v}"
fi
