package proxy

import (
	"errors"
	"fmt"
	"strings"
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
	metricsDefinition  metricValidator // metrics defined by plugin (code) // todo: remove
	metricsFilters     metricValidator // metric filters defined by task (yaml)

	sessionMts []*plugin.Metric

	cxtManager *ContextManager // back-reference to context manager
}

func NewPluginContext(ctxManager *ContextManager, mtsDefinition metricValidator, rawConfig []byte, mtsSelectors []string) (*pluginContext, error) {
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

		cxtManager: ctxManager,
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
	matchDefinition, groupPositions := pc.metricsDefinition.IsValid(ns)
	matchFilters, _ := pc.metricsFilters.IsValid(ns)

	if !matchDefinition {
		return errors.New("couldn't match metric with plugin definition")
	}

	if !matchFilters {
		return errors.New("couldn't match metrics with plugin filters")
	}

	mtNamespace := []plugin.NamespaceElement{}
	for i, nsElem := range strings.Split(ns, "/")[1:] {
		groupName := groupPositions[i]
		mtNamespace = append(mtNamespace, plugin.NamespaceElement{
			Name:        groupName,
			Value:       nsElem, // todo: extract only value when someone add /plugin/[group=df]/metr1
			Description: pc.cxtManager.groupsDescription[groupName],
		})
	}

	nsDefFormat := strings.Split(ns, "/")[1:]
	for i := 0; i < len(nsDefFormat); i++ {
		if groupPositions[i] != "" {
			nsDefFormat[i] = fmt.Sprintf("[%s]", groupPositions[i])
		}
	}

	nsDescKey := "/" + strings.Join(nsDefFormat, "/")
	pc.sessionMts = append(pc.sessionMts, &plugin.Metric{
		Namespace:   mtNamespace,
		Value:       v,
		Tags:        tags,
		Timestamp:   time.Now(),
		Description: pc.cxtManager.metricsDescription[nsDescKey],
	})

	return nil
}

func (pc *pluginContext) ApplyTagsByPath(string, map[string]string) error {
	panic("implement me")
}

func (pc *pluginContext) ApplyTagsByRegExp(string, map[string]string) error {
	panic("implement me")
}
