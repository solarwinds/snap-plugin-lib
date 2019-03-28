package proxy

import (
	"fmt"
	"sync"

	"github.com/librato/snap-plugin-lib-go/v2/internal/util/simpleconfig"
)

type pluginContext struct {
	rawConfig          []byte
	flattenedConfig    map[string]string
	mtsSelectors       []string
	storedObjects      map[string]interface{}
	storedObjectsMutex sync.RWMutex
}

func NewPluginContext(rawConfig []byte, mtsSelectors []string) (*pluginContext, error) {
	flattenConfig, err := simpleconfig.JSONToFlatMap(rawConfig)
	if err != nil {
		return nil, fmt.Errorf("can't create context due to invalid json: %v", err)
	}

	return &pluginContext{
		rawConfig:       []byte(rawConfig),
		flattenedConfig: flattenConfig,
		mtsSelectors:    mtsSelectors,
		storedObjects:   map[string]interface{}{},
	}, nil
}

func (pc *pluginContext) Config(key string) (string, bool) {
	v, ok := pc.flattenedConfig[key]
	return v, ok
}

func (pc *pluginContext) ConfigKeys() []string {
	keysList := []string{}
	for k := range pc.flattenedConfig {
		keysList = append(keysList, k)
	}
	return keysList
}

func (pc *pluginContext) RawConfig() []byte {
	return pc.rawConfig
}

func (pc *pluginContext) Store(key string, obj interface{}) {
	pc.storedObjectsMutex.Lock()
	defer pc.storedObjectsMutex.Unlock()

	pc.storedObjects[key] = obj
}

func (pc *pluginContext) Load(key string) (interface{}, bool) {
	pc.storedObjectsMutex.RLock()
	defer pc.storedObjectsMutex.RUnlock()

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
