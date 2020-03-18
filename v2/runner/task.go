package runner

import (
	"fmt"

	"github.com/librato/snap-plugin-lib-go/v2/internal/plugins/collector/proxy"
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
