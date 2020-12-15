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

version=${CIRCLE_TAG:-${TRAVIS_TAG:-250.0.0}}
versionarr=($(echo $version | tr "." "\n"))
commit=${CIRCLE_SHA1:-${TRAVIS_COMMIT:-0}}
ubuntu_release=(lsb_release -rs)

# shellcheck source=scripts/common.sh
. "${__dir}/common.sh"

# Build examples for v1 and v2
for f in $(find ./examples/ -name "main.go"); do (cd $(dirname $f) && echo "Building $PWD" && go build) ; done

pushd "${__proj_dir}/v2/bindings"
# Build CGO for linux
_info "Building swisnap-plugin-lib.so"
(export GOOS=linux && _go_build "--buildmode=c-shared" "swisnap-plugin-lib.so")

# Build CGO for windows 
_info "Building swisnap-plugin-lib.dll"
_go_get github.com/josephspurrier/goversioninfo/cmd/goversioninfo
cp versioninfo.template versioninfo.json
sed -i 's/{major_version}/'${versionarr[0]}'/g' versioninfo.json
sed -i 's/{minor_version}/'${versionarr[1]}'/g' versioninfo.json
sed -i 's/{patch_version}/'${versionarr[2]}'/g' versioninfo.json
sed -i 's/{version}/'${version}'/g' versioninfo.json
sed -i 's/{commit}/'${commit}'/g' versioninfo.json
go generate
(export GOOS=windows && export CGO_ENABLED=1 && export CXX=x86_64-w64-mingw32-g++ && export CC=x86_64-w64-mingw32-gcc  \
 && _go_build "--buildmode=c-shared" "swisnap-plugin-lib.dll")

# Building .Net 
pushd csharp/SnapPluginLibCs
dotnet build SnapPluginLibCs.sln -c Release -p:Version=${version}

popd
popd
