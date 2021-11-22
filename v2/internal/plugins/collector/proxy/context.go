/*
 Copyright (c) 2021 SolarWinds Worldwide, LLC

    Licensed under the Apache License, Version 2.0 (the "License");
    you may not use this file except in compliance with the License.
    You may obtain a copy of the License at

      http://www.apache.org/licenses/LICENSE-2.0

    Unless required by applicable law or agreed to in writing, software
    distributed under the License is distributed on an "AS IS" BASIS,
    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
    See the License for the specific language governing permissions and
    limitations under the License.
*/

package proxy

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/solarwinds/snap-plugin-lib/v2/plugin"

	commonProxy "github.com/solarwinds/snap-plugin-lib/v2/internal/plugins/common/proxy"
	"github.com/solarwinds/snap-plugin-lib/v2/internal/util/log"
	"github.com/solarwinds/snap-plugin-lib/v2/internal/util/metrictree"
	"github.com/solarwinds/snap-plugin-lib/v2/internal/util/types"
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

type PluginContext struct {
	*commonProxy.Context
	ctx context.Context

	taskID          string
	metricsFilters  *metrictree.TreeValidator // metric filters defined by task (yaml)
	sessionMtsMutex sync.RWMutex
	sessionMts      []*types.Metric
	modifiersTable  []*modifiersMetadata
	ctxManager      *ContextManager // back-reference to context manager
}

func NewPluginContext(ctxManager *ContextManager, taskID string, rawConfig []byte) (*PluginContext, error) {
	if ctxManager == nil {
		return nil, errors.New("can't create context without valid context manager")
	}

	baseContext, err := commonProxy.NewContext(rawConfig)
	if err != nil {
		return nil, err
	}

	pc := &PluginContext{
		Context:        baseContext,
		ctx:            ctxManager.ctx,
		taskID:         taskID,
		metricsFilters: metrictree.NewMetricFilter(ctxManager.metricsDefinition),
		ctxManager:     ctxManager,
		sessionMts:     nil,
	}

	return pc, nil
}

func (pc *PluginContext) AddMetric(ns string, v interface{}, modifiers ...plugin.MetricModifier) error {
	logF := log.WithCtx(pc.ctx).WithFields(moduleFields).WithField("service", "metrics")

	if pc.IsDone() {
		return fmt.Errorf("task has been canceled")
	}

	if !pc.isValidValueType(v) {
		return fmt.Errorf("invalid value type (%T) for metric: %s", v, ns)
	}

	if pc.ctxManager.globalPrefix.enabled {
		ns = fmt.Sprintf("/%s%s", pc.ctxManager.globalPrefix.name, ns)
	}

	if err := pc.ctxManager.metricsDefinition.IsUsableForAddition(ns, false); err != nil {
		return fmt.Errorf("invalid namespace (some elements can't be used when adding metric): %w", err)
	}

	matchDefinition, groupPositions := pc.ctxManager.metricsDefinition.IsValid(ns)
	matchFilters, _ := pc.metricsFilters.IsValid(ns)

	if !matchDefinition {
		return fmt.Errorf("couldn't match metric with plugin definition: %v", ns)
	}

	if !matchFilters {
		if logrus.IsLevelEnabled(logrus.TraceLevel) {
			logF.WithField("ns", ns).Trace("couldn't match metrics with plugin filters")
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
	// https://github.com/solarwinds/snap-plugin-lib/pull/49/files#r390325795
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

func (pc *PluginContext) ShouldProcess(ns string) bool {
	logF := log.WithCtx(pc.ctx).WithFields(moduleFields).WithField("service", "metrics")

	if err := pc.ctxManager.metricsDefinition.IsUsableForAddition(ns, true); err != nil {
		if logrus.GetLevel() >= logrus.DebugLevel {
			logF.WithError(err).WithField("namespace", ns).Debug("Should NOT process metric")
		}

		return false
	}

	defValid := pc.ctxManager.metricsDefinition.IsPartiallyValid(ns)
	shouldProcess := defValid && pc.metricsFilters.IsPartiallyValid(ns)

	return shouldProcess
}

func (pc *PluginContext) metricMeta(nsKey string) metricMetadata {
	if mtMeta, ok := pc.ctxManager.metricsMetadata[nsKey]; ok {
		return mtMeta
	}

	// if metric wasn't defined just return structure with empty fields
	return metricMetadata{}
}

func (pc *PluginContext) AlwaysApply(namespaceSelector string, modifiers ...plugin.MetricModifier) (plugin.Dismisser, error) {
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

func (pc *PluginContext) DismissAllModifiers() {
	pc.sessionMtsMutex.Lock()
	defer pc.sessionMtsMutex.Unlock()

	for _, m := range pc.modifiersTable {
		m.Dismiss()
	}
}

// extract static value when adding metrics like. /plugin/[grp=id]/m1
// function assumes valid format
func (pc *PluginContext) extractStaticValue(s string) string {
	eqIndex := strings.Index(s, "=")
	if eqIndex != -1 {
		return s[eqIndex+1 : len(s)-1]
	}

	return s
}

func (pc *PluginContext) ClearCollectorSession() {
	pc.sessionMtsMutex.Lock()
	defer pc.sessionMtsMutex.Unlock()

	pc.sessionMts = nil
	pc.modifiersTable = nil
}

func (pc *PluginContext) Metrics(clear bool) []*types.Metric {
	pc.sessionMtsMutex.Lock()
	defer pc.sessionMtsMutex.Unlock()

	mts := pc.sessionMts
	if clear {
		pc.sessionMts = nil
	}

	globalPrefix := pc.ctxManager.globalPrefix
	if globalPrefix.enabled && globalPrefix.removeFromOutput {
		for i := 0; i < len(mts); i++ {
			mts[i].Namespace_ = mts[i].Namespace_[1:]
		}
	}

	return mts
}

func (pc *PluginContext) RequestedMetrics() []string {
	return pc.metricsFilters.ListRules()
}

func (pc *PluginContext) TaskID() string {
	return pc.taskID
}

func (pc *PluginContext) isValidValueType(value interface{}) bool {
	// when adding new type(s) apply changes also in toGRPCValue() function
	switch value.(type) {
	case string:
	case float64:
	case float32:
	case int32:
	case int:
	case int64:
	case uint32:
	case uint64:
	case uint:
	case []byte:
	case bool:
	case int16:
	case uint16:
	case nil:
	case plugin.Summary:
	case *plugin.Summary:
	case plugin.Histogram:
	case *plugin.Histogram:
	default:
		return false
	}

	return true
}
