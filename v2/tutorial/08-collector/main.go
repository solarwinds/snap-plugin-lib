/*
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
	"context"

	"github.com/librato/snap-plugin-lib-go/v2/runner"
	"github.com/librato/snap-plugin-lib-go/v2/tutorial/08-collector/collector"
	"github.com/librato/snap-plugin-lib-go/v2/tutorial/08-collector/collector/proxy"
)

/*****************************************************************************/

const pluginName = "system-collector"
const pluginVersion = "1.0.0"

func main() {
	runner.StartCollectorWithContext(context.Background(), collector.New(proxy.New()), pluginName, pluginVersion)
}
