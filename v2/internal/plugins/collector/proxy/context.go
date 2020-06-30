package proxy

import (
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/librato/snap-plugin-lib-go/v2/plugin"

	commonProxy "github.com/librato/snap-plugin-lib-go/v2/internal/plugins/common/proxy"
	"github.com/librato/snap-plugin-lib-go/v2/internal/util/metrictree"
	"github.com/librato/snap-plugin-lib-go/v2/internal/util/types"
)

///////////////////////////////////////////////////////////////////////////////

type modifiersMetadata struct {
	nsSelector string
	modifiers  []plugin.MetricModifier
	validator  *metrictree.TreeValidator
	active     bool
}

func (m *modifiersMetadata) Dismiss() {
	m.active = false
}

///////////////////////////////////////////////////////////////////////////////

type pluginContext struct {
	*commonProxy.Context

	metricsFilters  *metrictree.TreeValidator // metric filters defined by task (yaml)
	sessionMtsMutex sync.RWMutex
	sessionMts      []*types.Metric
	modifiersTable  []*modifiersMetadata
	ctxManager      *ContextManager // back-reference to context manager
}

func NewPluginContext(ctxManager *ContextManager, rawConfig []byte) (*pluginContext, error) {
	if ctxManager == nil {
		return nil, errors.New("can't create context without valid context manager")
	}

	baseContext, err := commonProxy.NewContext(rawConfig)
	if err != nil {
		return nil, err
	}

	pc := &pluginContext{
		Context:        baseContext,
		metricsFilters: metrictree.NewMetricFilter(ctxManager.metricsDefinition),
		ctxManager:     ctxManager,
		sessionMts:     nil,
	}

	return pc, nil
}

func (pc *pluginContext) AddMetric(ns string, v interface{}, modifiers ...plugin.MetricModifier) error {
	if pc.IsDone() {
		return fmt.Errorf("task has been canceled")
	}

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
		return fmt.Errorf("couldn't match metric with plugin definition: %v", ns)
	}

	if !matchFilters {
		if logrus.IsLevelEnabled(logrus.TraceLevel) {
			log.WithField("ns", ns).Trace("couldn't match metrics with plugin filters")
		}
		return nil // don't throw error when metric is just filtered
	}

	var mtNamespace []types.NamespaceElement
	nsDefFormat, nsSeparator, err := metrictree.SplitNamespace(ns)
	if err != nil {
		return err
	}
	nsDefFormat = nsDefFormat[1:]

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

	// if performance would suffer at some point in future proposed solution (indefinite chan) may be introduced
	// https://github.com/librato/snap-plugin-lib-go/pull/49/files#r390325795
	// https://medium.com/capital-one-tech/building-an-unbounded-channel-in-go-789e175cd2cd
	pc.sessionMtsMutex.Lock()
	defer pc.sessionMtsMutex.Unlock()

	mt := &types.Metric{
		Namespace_:   mtNamespace,
		Value_:       v,
		Unit_:        mtMeta.unit,
		Timestamp_:   time.Now(),
		Description_: mtMeta.description,
	}

	// modifiers related to AddMetric
	for _, m := range modifiers {
		m.UpdateMetric(mt)
	}

	// modifiers list defined by AlwaysApply calls
	for _, modElement := range pc.modifiersTable {
		if !modElement.active {
			continue
		}

		isValid, _ := modElement.validator.IsValid(mt.Namespace().String())
		if isValid {
			for _, modifier := range modElement.modifiers {
				modifier.UpdateMetric(mt)
			}
		}
	}

	pc.sessionMts = append(pc.sessionMts, mt)

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

func (pc *pluginContext) AlwaysApply(namespaceSelector string, modifiers ...plugin.MetricModifier) (plugin.Dismisser, error) {
	pc.sessionMtsMutex.Lock()
	defer pc.sessionMtsMutex.Unlock()

	validator := metrictree.NewMetricFilter(metrictree.NewMetricDefinition())
	err := validator.AddRule(namespaceSelector)
	if err != nil {
		return nil, fmt.Errorf("can't apply modifiers: %v", err)
	}

	modifierMeta := &modifiersMetadata{
		nsSelector: namespaceSelector,
		modifiers:  modifiers,
		validator:  validator,
		active:     true,
	}

	pc.modifiersTable = append(pc.modifiersTable, modifierMeta)
	return modifierMeta, nil
}

func (pc *pluginContext) DismissAllModifiers() {
	pc.sessionMtsMutex.Lock()
	defer pc.sessionMtsMutex.Unlock()

	for _, m := range pc.modifiersTable {
		m.Dismiss()
	}
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

func (pc *pluginContext) ClearCollectorSession() {
	pc.sessionMtsMutex.Lock()
	defer pc.sessionMtsMutex.Unlock()

	pc.sessionMts = nil
	pc.modifiersTable = nil
}

func (pc *pluginContext) Metrics(clear bool) []*types.Metric {
	pc.sessionMtsMutex.Lock()
	defer pc.sessionMtsMutex.Unlock()

	mts := pc.sessionMts
	if clear {
		pc.sessionMts = nil
	}
	return mts
}

func (pc *pluginContext) RequestedMetrics() []string {
	return pc.metricsFilters.ListRules()
}
