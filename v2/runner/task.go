package runner

import (
	"fmt"

	"github.com/librato/snap-plugin-lib-go/v2/internal/proxy"
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
	Config  *yaml.Node
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

	temp := yaml.Node{}

	err := yaml.Unmarshal([]byte(ctxMan.ExampleConfig__), &temp)
	if err != nil {
		fmt.Printf("Error: can't do (%v)", err)
	}

	m := map[string]interface{}{}
	err = temp.Decode(m)
	if err != nil {
		fmt.Printf("Error: 1 (%v)", err)
	}

	bbbbbb, err := yaml.Marshal(&temp)
	fmt.Printf("bbbbbb=%s\n", string(bbbbbb))
	fmt.Printf("err=%#v\n", err)

	taskExample.Plugins[0].Config = temp.Content[0]

	b, err := yaml.Marshal(&taskExample)
	if err != nil {
		fmt.Printf("Error: can't print task information (%v)", err)
	}

	fmt.Printf("---\n%s\n", string(b))
}
