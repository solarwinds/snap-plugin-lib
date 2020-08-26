#!/bin/bash

protoc_ver="3.13.0"

licence='/*
Copyright (c) '"$(date +"%Y")"' SolarWinds Worldwide, LLC

    Licensed under the Apache License, Version 2.0 (the "License");
    you may not use this file except in compliance with the License.
    You may obtain a copy of the License at

      http://www.apache.org/licenses/LICENSE-2.0

    Unless required by applicable law or agreed to in writing, software
    distributed under the License is distributed on an "AS IS" BASIS,
    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
    See the License for the specific language governing permissions and
    limitations under the License.
*/'

if ! which protoca &> /dev/null ; then
    if [[ $OS == "Windows_NT" ]]; then
        protoc_os="win64"
        protoc_ext=".exe"
    else
        protoc_os="linux-x86_64"
        protoc_ext=""
    fi

    protoc_path="https://github.com/protocolbuffers/protobuf/releases/download/v${protoc_ver}/protoc-${protoc_ver}-${protoc_os}.zip"
    
    echo "Installing protoc from ${protoc_path}"
    curl -o protoc.zip -L $protoc_path -s
    unzip -p protoc.zip "bin/protoc${protoc_ext}" > "protoc${protoc_ext}"
    mv "protoc${protoc_ext}" $GOPATH/bin
fi


