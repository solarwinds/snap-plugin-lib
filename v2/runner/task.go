package runner

import (
	"fmt"

	"github.com/go-yaml/yaml"
	"github.com/librato/snap-plugin-lib-go/v2/internal/proxy"
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
	Config  map[string]interface{}
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

	aoName := fmt.Sprintf("ao-%scollector", pluginName)
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
				Config: map[string]interface{}{
					"config_variable": "value",
				},
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

	b, err := yaml.Marshal(taskExample)
	if err != nil {
		fmt.Printf("Error: can't print task information (%v)", err)
	}

	fmt.Printf("---\n%s\n", string(b))
}
