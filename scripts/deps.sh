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

set -e
set -u
set -o pipefail

__dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
__proj_dir="$(dirname "$__dir")"

# shellcheck source=scripts/common.sh
. "${__dir}/common.sh"

detect_go_dep() {
  [[ -f "${__proj_dir}/Godeps/Godeps.json" ]] && _dep='godep'
  [[ -f "${__proj_dir}/glide.yaml" ]] && _dep='glide'
  [[ -f "${__proj_dir}/vendor/vendor.json" ]] && _dep='govendor'
  [[ -f "${__proj_dir}/Gopkg.toml" ]] && _dep='dep'
  [[ -f "${__proj_dir}/go.mod" ]] && _dep='mod'
  _info "golang dependency tool: ${_dep}"
  echo "${_dep}"
}

install_go_dep() {
  local _dep=${_dep:=$(_detect_dep)}
  _info "ensuring ${_dep} is available"
  case $_dep in
    godep)
      _go_install github.com/tools/godep
      ;;
    dep)
      _go_install github.com/golang/dep/cmd/dep
      ;;
    mod)
      _info "nothing special to do for go modules"
      ;;
    glide)
      _go_install github.com/Masterminds/glide
      ;;
    govendor)
      _go_install github.com/kardianos/govendor
      ;;
  esac
}

restore_go_dep() {
  local _dep=${_dep:=$(_detect_dep)}
  _info "restoring dependency with ${_dep}"
  case $_dep in
    godep)
      (cd "${__proj_dir}" && godep restore)
      ;;
    glide)
      (cd "${__proj_dir}" && glide install)
      ;;
    govendor)
      (cd "${__proj_dir}" && govendor sync)
      ;;
    dep)
      (cd "${__proj_dir}" && dep ensure -v -vendor-only)
      ;;
    mod)
      (cd "${__proj_dir}" && go mod download)
      ;;
  esac
}

_dep=$(detect_go_dep)

if [ $_dep != "mod" ]; then
  _warning "Deprecated dependencies management tool, please update to go modules"

  if [ $_dep != "dep" ]; then
    _error "No even using dep, you need to update the Go dependencies management tool"
    exit 1
  fi
fi

install_go_dep
restore_go_dep
