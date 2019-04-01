package proxy

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/librato/snap-plugin-lib-go/v2/internal/util/metrictree"
	"github.com/librato/snap-plugin-lib-go/v2/internal/util/simpleconfig"
	"github.com/librato/snap-plugin-lib-go/v2/plugin"
)

type pluginContext struct {
	rawConfig          []byte
	flattenedConfig    map[string]string
	storedObjects      map[string]interface{}
	storedObjectsMutex sync.RWMutex
	metricsDefinition  metricValidator // metrics defined by plugin (code)
	metricsFilters     metricValidator // metric filters defined by task (yaml)

	sessionMts []*plugin.Metric
}

func NewPluginContext(mtsDefinition metricValidator, rawConfig []byte, mtsSelectors []string) (*pluginContext, error) {
	flattenConfig, err := simpleconfig.JSONToFlatMap(rawConfig)
	if err != nil {
		return nil, fmt.Errorf("can't create context due to invalid json: %v", err)
	}

	return &pluginContext{
		rawConfig:         []byte(rawConfig),
		flattenedConfig:   flattenConfig,
		storedObjects:     map[string]interface{}{},
		metricsDefinition: mtsDefinition,
		metricsFilters:    metrictree.NewMetricFilter(),
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

func (pc *pluginContext) AddMetric(ns string, v interface{}) error {
	return pc.AddMetricWithTags(ns, v, nil)
}

func (pc *pluginContext) AddMetricWithTags(ns string, v interface{}, tags map[string]string) error {
	matchDefinition := pc.metricsDefinition.IsValid(ns)
	matchFilters := pc.metricsFilters.IsValid(ns)

	if !matchDefinition {
		return errors.New("couldn't match metric with plugin definition")
	}

	if !matchFilters {
		return errors.New("couldn't match metrics with plugin filters")
	}

	pc.sessionMts = append(pc.sessionMts, &plugin.Metric{
		Namespace: ns,
		Value:     v,
		Tags:      tags,
		Timestamp: time.Now(),
	})

	return nil
}

func (pc *pluginContext) ApplyTagsByPath(string, map[string]string) error {
	panic("implement me")
}

func (pc *pluginContext) ApplyTagsByRegExp(string, map[string]string) error {
	panic("implement me")
}
