#!/usr/bin/env bash

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

version=${TRAVIS_TAG:-250.0.1}
readarray -d . -t versionarr < <(printf '%s' ${version})
commit=${TRAVIS_COMMIT:-0}

# shellcheck source=scripts/common.sh
. "${__dir}/common.sh"

# Build CGO for linux
_info "Building swisnap-plugin-lib.so"
(export GOOS=linux && cd "${__proj_dir}/v2/bindings" && _go_build "--buildmode=c-shared" "swisnap-plugin-lib.so")

# Build CGO for windows 
_info "Building swisnap-plugin-lib.dll"
_go_get github.com/josephspurrier/goversioninfo/cmd/goversioninfo
(cd "${__proj_dir}/v2/bindings" && cp versioninfo.template versioninfo.json)
(cd "${__proj_dir}/v2/bindings" && sed -i 's/{major_version}/'${versionarr[0]}'/g' versioninfo.json)
(cd "${__proj_dir}/v2/bindings" && sed -i 's/{minor_version}/'${versionarr[1]}'/g' versioninfo.json)
(cd "${__proj_dir}/v2/bindings" && sed -i 's/{patch_version}/'${versionarr[2]}'/g' versioninfo.json)
(cd "${__proj_dir}/v2/bindings" && sed -i 's/{version}/'${version}'/g' versioninfo.json)
(cd "${__proj_dir}/v2/bindings" && sed -i 's/$TRAVIS_COMMIT/'${commit}'/g' versioninfo.json)
(cd "${__proj_dir}/v2/bindings" && go generate)
(export GOOS=windows && export CGO_ENABLED=1 && export CXX=x86_64-w64-mingw32-g++ && export CC=x86_64-w64-mingw32-gcc && cd "${__proj_dir}/v2/bindings" && _go_build "--buildmode=c-shared" "swisnap-plugin-lib.dll")
