package runner

import (
	"fmt"

	"github.com/librato/snap-plugin-lib-go/v2/internal/plugins/collector/proxy"
	"gopkg.in/yaml.v3"
)

type taskEntry struct {
	Version  int
	Schedule scheduleEntry
	Plugins  []pluginEntry
}

type scheduleEntry struct {
	SType    string `yaml:"type"`
	Interval string
}

type pluginEntry struct {
	Name    string
	Metrics []string
	Config  *yaml.Node `yaml:",omitempty"`
	Tags    map[string]interface{}
	Publish publishEntry
}

type publishEntry struct {
	Config publishSubEntry
}

type publishSubEntry struct {
	Period      int
	FloorSecond int `yaml:"floor_second"`
}

func printExampleTask(ctxMan *proxy.ContextManager, pluginName string) {
	mtsList := ctxMan.ListDefaultMetrics()

	aoName := fmt.Sprintf("%scollector", pluginName)
	tagKey := fmt.Sprintf("/%s", pluginName)

	taskExample := taskEntry{
		Version: 2,
		Schedule: scheduleEntry{
			SType:    "simple",
			Interval: "60s",
		},
		Plugins: []pluginEntry{
			{
				Name:    aoName,
				Metrics: mtsList,
				Tags: map[string]interface{}{
					tagKey: map[string]interface{}{
						"plugin_tag": "tag",
					},
				},
				Publish: publishEntry{
					Config: publishSubEntry{
						Period:      60,
						FloorSecond: 60,
					},
				},
			},
		},
	}

	if len(ctxMan.ExampleConfig.Content) != 0 {
		taskExample.Plugins[0].Config = ctxMan.ExampleConfig.Content[0]
	}

	b, err := yaml.Marshal(&taskExample)
	if err != nil {
		fmt.Printf("Error: can't print task information (%v)", err)
	}

	fmt.Printf("---\n%s\n", string(b))
}
