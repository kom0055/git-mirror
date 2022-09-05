#!/usr/bin/env bash

set -o errexit
set -o pipefail
set -o nounset
#set -x

PACKAGE_NAME="gclone"
BASE_DIR=$(cd $(dirname $0)/.. && pwd)
OUTPUT_PATH=${BASE_DIR}/_output
mkdir -p ${OUTPUT_PATH}/bin

source "${BASE_DIR}/hack/version.sh"

go build -gcflags=all="-N -l" -a -o "${OUTPUT_PATH}"/bin/${PACKAGE_NAME} \
  -ldflags "$(api::version::ldflags)" "${OUTPUT_PATH}"/cmd/
