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

const nsSeparator = metrictree.NsSeparator

type pluginContext struct {
	rawConfig          []byte
	flattenedConfig    map[string]string
	storedObjects      map[string]interface{}
	storedObjectsMutex sync.RWMutex
	metricsFilters     *metrictree.TreeValidator // metric filters defined by task (yaml)

	sessionMts []*types.Metric

	ctxManager *ContextManager // back-reference to context manager
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

	pc.metricsFilters = metrictree.NewMetricFilter(ctxManager.metricsDefinition)
	pc.ctxManager = ctxManager

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

func (pc *pluginContext) AddMetric(ns string, v interface{}) error {
	return pc.AddMetricWithTags(ns, v, nil)
}

func (pc *pluginContext) AddMetricWithTags(ns string, v interface{}, tags map[string]string) error {
	parsedNs, err := metrictree.ParseNamespace(ns, false)
	if err != nil {
		return fmt.Errorf("invalid format of namespace: %v", err)
	}
	if !parsedNs.IsUsableForAddition(pc.ctxManager.metricsDefinition.HasRules(), false) {
		return fmt.Errorf("invalid namespace (some elements can't be used when adding metric): %v", err)
	}

	matchDefinition, groupPositions := pc.ctxManager.metricsDefinition.IsValid(ns)
	matchFilters, _ := pc.metricsFilters.IsValid(ns)

	if !matchDefinition {
		return fmt.Errorf("couldn't match metric with plugin definition: %v", err)
	}

	if !matchFilters {
		return fmt.Errorf("couldn't match metrics with plugin filters: %v", err)
	}

	var mtNamespace []types.NamespaceElement
	nsDefFormat := strings.Split(ns, metrictree.NsSeparator)[1:]

	for i, nsElem := range nsDefFormat {
		groupName := groupPositions[i]
		mtNamespace = append(mtNamespace, types.NamespaceElement{
			Name_:        groupName,
			Value_:       pc.extractStaticValue(nsElem),
			Description_: pc.ctxManager.groupsDescription[groupName],
		})

		if groupPositions[i] != "" {
			nsDefFormat[i] = fmt.Sprintf("[%s]", groupPositions[i])
		}
	}

	nsDescKey := nsSeparator + strings.Join(nsDefFormat, nsSeparator)
	mtMeta := pc.metricMeta(nsDescKey)

	pc.sessionMts = append(pc.sessionMts, &types.Metric{
		Namespace_:   mtNamespace,
		Value_:       v,
		Tags_:        tags,
		Unit_:        mtMeta.unit,
		Timestamp_:   time.Now(),
		Description_: mtMeta.description,
	})

	return nil
}

func (pc *pluginContext) ShouldProcess(ns string) bool {
	parsedNs, err := metrictree.ParseNamespace(ns, false)
	if err != nil {
		return false
	}
	if !parsedNs.IsUsableForAddition(pc.ctxManager.metricsDefinition.HasRules(), true) {
		return false
	}

	defValid := pc.ctxManager.metricsDefinition.IsPartiallyValid(ns)
	shouldProcess := defValid && pc.metricsFilters.IsPartiallyValid(ns)

	return shouldProcess
}

func (pc *pluginContext) metricMeta(nsKey string) metricMetadata {
	if mtMeta, ok := pc.ctxManager.metricsMetadata[nsKey]; ok {
		return mtMeta
	}

	// if metric wasn't defined just return structure with empty fields
	return metricMetadata{}
}

func (pc *pluginContext) ApplyTagsByPath(string, map[string]string) error {
	// TODO: https://swicloud.atlassian.net/browse/AO-12232
	panic("implement me")
}

func (pc *pluginContext) ApplyTagsByRegExp(string, map[string]string) error {
	// TODO: https://swicloud.atlassian.net/browse/AO-12232
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
