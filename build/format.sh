#!/usr/bin/env bash

# Copyright 2020 The Kubernetes Authors.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# This script checks whether the source codes need to be formatted or not by
# `gofmt`. We should run `hack/update-gofmt.sh` if actually formats them.
# Usage: `hack/verify-gofmt.sh`.
# Note: GoFmt apparently is changing @ head...

set -o errexit
set -o nounset
set -o pipefail

SCRIPT_ROOT=$(dirname "${BASH_SOURCE[0]}")/..

cd "${SCRIPT_ROOT}"

find_files() {
  find . -not \( \
    \( \
    -wholename './output' \
    -o -wholename './.git' \
    -o -wholename './_output' \
    -o -wholename './release' \
    -o -wholename './target' \
    -o -wholename '*/third_party/*' \
    -o -wholename '*/vendor/*' \
    -o -wholename './staging/src/k8s.io/client-go/*vendor/*' \
    -o -wholename '*/bindata.go' \
    \) -prune \
    \) -name '*.go' -type f
}

# gofmt exits with non-zero exit code if it finds a problem unrelated to
# formatting (e.g., a file does not parse correctly). Without "|| true" this
# would have led to no useful error message from gofmt, because the script would
# have failed before getting to the "echo" in the block below.

find_files | xargs goimports -w -local github.com/kom0055/git-mirror
find_files | xargs gofmt -s -d -w
golangci-lint run -c "${SCRIPT_ROOT}"/.golangci.yaml "${SCRIPT_ROOT}"/...
find_files | xargs cat | wc -l
