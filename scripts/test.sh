#!/usr/bin/env bash

#
# Copyright 2016 Intel Corporation
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
#
#
# Copyright (c) 2020 SolarWinds Worldwide, LLC
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
#

# Support travis.ci environment matrix:
TEST_TYPE="${TEST_TYPE:-$1}"
UNIT_TEST="${UNIT_TEST:-"gofmt goimports go_vet go_test go_cover"}"
TEST_K8S="${TEST_K8S:-0}"
DEBUG="${DEBUG:-}"

set -e
set -u
set -o pipefail

__dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
__proj_dir="$(dirname "$__dir")"

# shellcheck source=scripts/common.sh
. "${__dir}/common.sh"

_debug "script directory ${__dir}"
_debug "project directory ${__proj_dir}"
_info "skipping go test in following directories: ${NO_GO_TEST}"

[[ "$TEST_TYPE" =~ ^(lint|small|medium|large|deploy|build)$ ]] || _error "invalid TEST_TYPE (value must be 'lint', 'small', 'medium', 'large', 'deploy', or 'build' recieved:${TEST_TYPE}"

lint() {
  # This function runs various linters and static analysis tools
  # 1. gofmt         (http://golang.org/cmd/gofmt/)
  # 2. goimports     (https://github.com/bradfitz/goimports)
  # 3. golint        (https://github.com/golang/lint)
  # 4. go vet        (http://golang.org/cmd/vet)
  local go_analyzers
  go_analyzers=(gofmt goimports golint go_vet copyrights)

  ((n_elements=${#go_analyzers[@]}, max=n_elements - 1))
  for ((i = 0; i <= max; i++)); do
    _info "running ${go_analyzers[i]}"
    _"${go_analyzers[i]}"
  done
}

test_unit() {
  # This function runs unit tests with code coverage analysis
  local go_tests
  go_tests=(go_test go_cover go_race)

  _debug "available unit tests: ${go_tests[*]}"
  _debug "user specified tests: ${UNIT_TEST}"

  ((n_elements=${#go_tests[@]}, max=n_elements - 1))
  for ((i = 0; i <= max; i++)); do
    if [[ "${UNIT_TEST}" =~ (^| )"${go_tests[i]}"( |$) ]]; then
      _info "running ${go_tests[i]}"
      _"${go_tests[i]}"
    else
      _debug "skipping ${go_tests[i]}"
    fi
  done
}

if [[ $TEST_TYPE == "lint" ]]; then
  lint
elif [[ $TEST_TYPE == "small" ]]; then
  if [[ -f "${__dir}/small.sh" ]]; then
    . "${__dir}/small.sh"
  else
    echo "mode: count" > profile.cov
    test_unit
  fi
elif [[ $TEST_TYPE == "medium" ]]; then
  if [[ -f "${__dir}/setup_medium.sh" ]]; then
    _info "running the medium test setup script. . ."
    . "${__dir}/setup_medium.sh"
  fi
  if [[ -f "${__dir}/medium.sh" ]]; then
    . "${__dir}/medium.sh"
  else
    UNIT_TEST="go_test go_cover"
    echo "mode: count" > profile.cov
    test_unit
  fi
fi
