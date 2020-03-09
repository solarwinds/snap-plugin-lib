package proxy

import (
	"context"
	"fmt"
	"reflect"
	"sync"
	"time"

	"github.com/librato/snap-plugin-lib-go/v2/internal/util/simpleconfig"
	"github.com/librato/snap-plugin-lib-go/v2/internal/util/types"
	"github.com/sirupsen/logrus"
)

const (
	maxWarningMsgSize = 256 // maximum length of a single warning message
	maxNoOfWarnings   = 40  // maximum number of warnings added during one collect/publish operation
)

var (
	log = logrus.WithFields(logrus.Fields{"layer": "lib", "module": "common-proxy"})
)

type Context struct {
	rawConfig          []byte
	flattenedConfig    map[string]string
	storedObjects      map[string]interface{}
	storedObjectsMutex sync.RWMutex
	sessionWarnings    []types.Warning
	warningsMutex      sync.RWMutex

	ctx context.Context
}

type Warning struct {
	Message   string
	Timestamp time.Time
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
		ctx:             context.Background(),
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

func (c *Context) LoadTo(key string, dest interface{}) error {
	c.storedObjectsMutex.RLock()
	defer c.storedObjectsMutex.RUnlock()

	obj, ok := c.storedObjects[key]
	if !ok {
		return fmt.Errorf("couldn't find object with a given key (%s)", key)
	}

	vDest := reflect.ValueOf(dest)
	if vDest.Kind() != reflect.Ptr || vDest.IsNil() {
		return fmt.Errorf("passed variable should be a non-nill pointer")
	}
	if reflect.TypeOf(dest).Elem() != reflect.TypeOf(obj) {
		return fmt.Errorf("type of destination variable don't match to type of stored value")
	}

	vDest.Elem().Set(reflect.ValueOf(obj))

	return nil
}

func (c *Context) AddWarning(msg string) {
	c.warningsMutex.Lock()
	defer c.warningsMutex.Unlock()

	if len(c.sessionWarnings) >= maxNoOfWarnings {
		log.Warning("Maximum number of warnings logged. New warning has been ignored")
		return
	}

	if len(msg) > maxWarningMsgSize {
		log.Info("Warning message size exceeds maximum allowed value and will be cut off")
		msg = msg[:maxWarningMsgSize]
	}

	c.sessionWarnings = append(c.sessionWarnings, types.Warning{
		Message:   msg,
		Timestamp: time.Now(),
	})
}

func (c *Context) Warnings(clear bool) []types.Warning {
	c.warningsMutex.RLock()
	defer c.warningsMutex.RUnlock()

	warnings := c.sessionWarnings
	if clear {
		warnings = []types.Warning{}
	}
	return warnings
}

func (c *Context) ResetWarnings() {
	c.warningsMutex.RLock()
	defer c.warningsMutex.RUnlock()

	c.sessionWarnings = []types.Warning{}
}

func (c *Context) AttachContext(parentCtx context.Context) {
	c.ctx, _ = context.WithCancel(parentCtx)
}

func (c *Context) IsDone() bool {
	return c.ctx.Err() != nil
}

func (c *Context) Done() <-chan struct{} {
	return c.ctx.Done()
}
