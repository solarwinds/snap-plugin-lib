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

package runner

import (
	"fmt"

	"github.com/solarwinds/snap-plugin-lib/v2/internal/plugins/collector/proxy"
	"gopkg.in/yaml.v3"
)

func printExampleTask(ctxMan *proxy.ContextManager, pluginName string) {
	var b []byte
	var err error

	if len(ctxMan.ExampleConfig.Content) != 0 {
		b, err = yaml.Marshal(&ctxMan.ExampleConfig)
	} else {
		template := fmt.Sprintf(`
# THIS IS GENERIC EXAMPLE TASK TEMPLATE
---
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
      - plugin_name: publisher-appoptics
        binary_name: snap-plugin-publisher-appoptics

`, pluginName)
		b = []byte(template)
	}

	if err != nil {
		fmt.Printf("Error: can't print task information (%v)", err)
	}

	fmt.Printf("---\n%s\n", string(b))
}
