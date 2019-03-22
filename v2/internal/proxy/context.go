package proxy

import (
	"fmt"

	"github.com/librato/snap-plugin-lib-go/v2/internal/utils"
)

type pluginContext struct {
	rawConfig     string
	flattenConfig map[string]string
	mtsSelectors  []string
	storedObjects map[string]interface{}
}

func NewPluginContext(config string, selectors []string) (*pluginContext, error) {
	flattenConfig, err := utils.JSONToFlatMap(config)
	if err != nil {
		return nil, fmt.Errorf("can't create context due to invalid json: %v", err)
	}

	return &pluginContext{
		rawConfig:     config,
		flattenConfig: flattenConfig,
		mtsSelectors:  selectors,
	}, nil
}

func (pc *pluginContext) Config(key string) (string, bool) {
	v, ok := pc.flattenConfig[key]
	return v, ok
}

func (pc *pluginContext) ConfigKeys() []string {
	keysList := []string{}
	for k := range pc.flattenConfig {
		keysList = append(keysList, k)
	}
	return keysList
}

func (pc *pluginContext) RawConfig() string {
	return pc.rawConfig
}

func (pc *pluginContext) Store(key string, obj interface{}) {
	pc.storedObjects[key] = obj
}

func (pc *pluginContext) Load(key string) (interface{}, bool) {
	obj, ok := pc.storedObjects[key]
	return obj, ok
}

func (pc *pluginContext) AddMetric(string, interface{}) error {
	panic("implement me")
}

func (pc *pluginContext) AddMetricWithTags(string, interface{}, map[string]string) error {
	panic("implement me")
}

func (pc *pluginContext) ApplyTagsByPath(string, map[string]string) error {
	panic("implement me")
}

func (pc *pluginContext) ApplyTagsByRegExp(string, map[string]string) error {
	panic("implement me")
}
