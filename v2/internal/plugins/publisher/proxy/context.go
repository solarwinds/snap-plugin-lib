package proxy

import (
	"errors"
	"fmt"
	"sync"

	"github.com/librato/snap-plugin-lib-go/v2/internal/util/simpleconfig"
	"github.com/librato/snap-plugin-lib-go/v2/internal/util/types"
	"github.com/librato/snap-plugin-lib-go/v2/plugin"
)

type pluginContext struct {
	rawConfig          []byte // todo: make common context
	flattenedConfig    map[string]string
	storedObjects      map[string]interface{}
	storedObjectsMutex sync.RWMutex

	sessionMts []*types.Metric
}

func NewPluginContext(ctxManager *ContextManager, rawConfig []byte) (*pluginContext, error) {
	if ctxManager == nil {
		return nil, errors.New("can't create context without valid context manager")
	}

	flattenedConfig, err := simpleconfig.JSONToFlatMap(rawConfig)
	if err != nil {
		return nil, fmt.Errorf("can't create context due to invalid json: %v", err)
	}

	pc := &pluginContext{
		rawConfig:       rawConfig,
		flattenedConfig: flattenedConfig,
		storedObjects:   map[string]interface{}{},
	}

	return pc, nil
}

func (pc *pluginContext) Config(key string) (string, bool) {
	v, ok := pc.flattenedConfig[key]
	return v, ok
}

func (pc *pluginContext) ConfigKeys() []string {
	var keysList []string
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

func (pc *pluginContext) ListAllMetrics() []plugin.Metric {
	mts := make([]plugin.Metric, 0, len(pc.sessionMts))

	for _, mt := range pc.sessionMts {
		mts = append(mts, mt)
	}

	return mts
}

func (pc *pluginContext) Count() int {
	return len(pc.sessionMts)
}
