/*
 Copyright 2016 Intel Corporation

    Licensed under the Apache License, Version 2.0 (the "License");
    you may not use this file except in compliance with the License.
    You may obtain a copy of the License at

      http://www.apache.org/licenses/LICENSE-2.0

    Unless required by applicable law or agreed to in writing, software
    distributed under the License is distributed on an "AS IS" BASIS,
    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
    See the License for the specific language governing permissions and
    limitations under the License.

 Copyright (c) 2020 SolarWinds Worldwide, LLC

    Licensed under the Apache License, Version 2.0 (the "License");
    you may not use this file except in compliance with the License.
    You may obtain a copy of the License at

      http://www.apache.org/licenses/LICENSE-2.0

    Unless required by applicable law or agreed to in writing, software
    distributed under the License is distributed on an "AS IS" BASIS,
    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
    See the License for the specific language governing permissions and
    limitations under the License.

*/

package main

import (
	"github.com/librato/snap-plugin-lib-go/examples/v1/snap-plugin-collector-rand/rand"
	"github.com/librato/snap-plugin-lib-go/v1/plugin"
	"google.golang.org/grpc"
)

const (
	pluginName      = "test-rand-collector"
	pluginVersion   = 1
	maxGRPCMesgSize = 2 * 1024
)

func main() {
	plugin.StartCollector(rand.RandCollector{}, pluginName, pluginVersion, plugin.GRPCServerOptions(grpc.MaxRecvMsgSize(maxGRPCMesgSize)))
}
