package proxy

import (
	"fmt"
	"sync"

	"github.com/librato/snap-plugin-lib-go/v2/internal/util/simpleconfig"
)

type Context struct {
	rawConfig          []byte
	flattenedConfig    map[string]string
	storedObjects      map[string]interface{}
	storedObjectsMutex sync.RWMutex
}

func NewContext(rawConfig []byte) (*Context, error) {
	flattenedConfig, err := simpleconfig.JSONToFlatMap(rawConfig)
	if err != nil {
		return nil, fmt.Errorf("can't create context due to invalid json: %v", err)
	}

	return &Context{
		rawConfig:       rawConfig,
		flattenedConfig: flattenedConfig,
		storedObjects:   map[string]interface{}{},
	}, nil
}

func (c *Context) Config(key string) (string, bool) {
	v, ok := c.flattenedConfig[key]
	return v, ok
}

func (c *Context) ConfigKeys() []string {
	var keysList []string
	for k := range c.flattenedConfig {
		keysList = append(keysList, k)
	}
	return keysList
}

func (c *Context) RawConfig() []byte {
	return c.rawConfig
}

func (c *Context) Store(key string, obj interface{}) {
	c.storedObjectsMutex.Lock()
	defer c.storedObjectsMutex.Unlock()

	c.storedObjects[key] = obj
}

func (c *Context) Load(key string) (interface{}, bool) {
	c.storedObjectsMutex.RLock()
	defer c.storedObjectsMutex.RUnlock()

	obj, ok := c.storedObjects[key]
	return obj, ok
}
