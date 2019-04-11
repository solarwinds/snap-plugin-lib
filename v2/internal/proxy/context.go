package proxy

import (
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/librato/snap-plugin-lib-go/v2/internal/util/metrictree"
	"github.com/librato/snap-plugin-lib-go/v2/internal/util/simpleconfig"
	"github.com/librato/snap-plugin-lib-go/v2/internal/util/types"
)

type pluginContext struct {
	rawConfig          []byte
	flattenedConfig    map[string]string
	storedObjects      map[string]interface{}
	storedObjectsMutex sync.RWMutex
	metricsFilters     metricValidator // metric filters defined by task (yaml)

	sessionMts []*types.Metric

	ctxManager *contextManager // back-reference to context manager
}

func NewPluginContext(ctxManager *contextManager, rawConfig []byte, mtsSelectors []string) (*pluginContext, error) {
	flattenedConfig, err := simpleconfig.JSONToFlatMap(rawConfig)
	if err != nil {
		return nil, fmt.Errorf("can't create context due to invalid json: %v", err)
	}

	pc := &pluginContext{
		rawConfig:       []byte(rawConfig),
		flattenedConfig: flattenedConfig,
		storedObjects:   map[string]interface{}{},
	}

	if ctxManager != nil {
		pc.metricsFilters = metrictree.NewMetricFilter(ctxManager.metricsDefinition.(*metrictree.TreeValidator))
		pc.ctxManager = ctxManager
	}

	return pc, nil
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
	matchDefinition, groupPositions := pc.ctxManager.metricsDefinition.IsValid(ns)
	matchFilters, _ := pc.metricsFilters.IsValid(ns)

	if !matchDefinition {
		return errors.New("couldn't match metric with plugin definition")
	}

	if !matchFilters {
		return errors.New("couldn't match metrics with plugin filters")
	}

	mtNamespace := []types.NamespaceElement{}
	nsDefFormat := strings.Split(ns, "/")[1:]

	for i, nsElem := range nsDefFormat {
		groupName := groupPositions[i]
		mtNamespace = append(mtNamespace, types.NamespaceElement{
			Name:        groupName,
			Value:       pc.extractStaticValue(nsElem),
			Description: pc.ctxManager.groupsDescription[groupName],
		})
	}

	for i := 0; i < len(nsDefFormat); i++ {
		if groupPositions[i] != "" {
			nsDefFormat[i] = fmt.Sprintf("[%s]", groupPositions[i])
		}
	}

	nsDescKey := "/" + strings.Join(nsDefFormat, "/")
	mtMeta := pc.metricMeta(nsDescKey)

	pc.sessionMts = append(pc.sessionMts, &types.Metric{
		Namespace:   mtNamespace,
		Value:       v,
		Tags:        tags,
		Unit:        mtMeta.unit,
		Timestamp:   time.Now(),
		Description: mtMeta.description,
	})

	return nil
}

func (pc *pluginContext) metricMeta(nsKey string) metricMetadata {
	if mtMeta, ok := pc.ctxManager.metricsMetadata[nsKey]; ok {
		return mtMeta
	}

	// if metric wasn't defined just return structure with empty fields
	return metricMetadata{}
}

func (pc *pluginContext) ApplyTagsByPath(string, map[string]string) error {
	panic("implement me")
}

func (pc *pluginContext) ApplyTagsByRegExp(string, map[string]string) error {
	panic("implement me")
}

// extract static value when adding metrics like. /plugin/[grp=id]/m1
// function assumes valid format
func (pc *pluginContext) extractStaticValue(s string) string {
	eqIndex := strings.Index(s, "=")
	if eqIndex != -1 {
		return s[eqIndex+1 : len(s)-1]
	}

	return s
}
