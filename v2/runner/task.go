/*
 Copyright (c) 2021 SolarWinds Worldwide, LLC

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

package runner

import (
	"fmt"

	"github.com/solarwinds/snap-plugin-lib/v2/internal/util/types"
	"gopkg.in/yaml.v3"
)

const template = `
# THIS IS AN AUTOMATICALLY GENERATED TASK TEMPLATE EXAMPLE

version: 2
schedule:
  type: cron
  interval: "0 * * * * *"
plugins:
  - plugin_name: %s
    # plugin_binary:

    # config:

    # metrics:

    publish:
      - plugin_name: %s
`

func printExampleTask(exampleConfig yaml.Node, pluginName string, pluginType types.PluginType) {
	var b []byte
	var err error

	if len(exampleConfig.Content) != 0 {
		b, err = yaml.Marshal(&exampleConfig)
	} else {
		var filledTemplate string

		switch pluginType {
		case types.PluginTypeCollector:
			filledTemplate = fmt.Sprintf(template, pluginName, "publisher-appoptics")
		case types.PluginTypeStreamingCollector:
			filledTemplate = fmt.Sprintf(template, pluginName, "publisher-appoptics")
		case types.PluginTypePublisher:
			filledTemplate = fmt.Sprintf(template, "collector-name", pluginName)
		default:
			err = fmt.Errorf("invalid plugin type")
		}

		b = []byte(filledTemplate)
	}

	if err != nil {
		fmt.Printf("Error: can't print task information (%v)", err)
	}

	fmt.Printf("---\n%s\n", string(b))
}
