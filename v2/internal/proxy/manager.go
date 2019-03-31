/*
Package proxy:
1) Manages context for different task created for the same plugin
2) Serves as an entry point for any "controller" (like. Rpc)
*/
package proxy

import (
	"errors"
	"fmt"

	"github.com/librato/snap-plugin-lib-go/v2/internal/util/metrictree"
	"github.com/librato/snap-plugin-lib-go/v2/plugin"
	"github.com/sirupsen/logrus"
)

var log = logrus.WithFields(logrus.Fields{"module": "plugin-proxy"})

type Collector interface {
	RequestCollect(id int) ([]plugin.Metric, error)
	LoadTask(id int, config []byte, selectors []string) error
	UnloadTask(id int) error
	RequestInfo()
}

type metricValidator interface {
	AddRule(string) error
	IsValid(string) bool
	ListRules() []string
}

type ContextManager struct {
	collector  plugin.Collector       // reference to custom plugin code
	contextMap map[int]*pluginContext // map of contexts associated with taskIDs

	metricsDefinition metricValidator // metrics defined by plugin (code)
	metricsFilters    metricValidator // metric filters defined by task (yaml)
}

func NewContextManager(collector plugin.Collector, pluginName string, version string) Collector {
	cm := &ContextManager{
		collector:  collector,
		contextMap: map[int]*pluginContext{},

		metricsDefinition: metrictree.NewMetricDefinition(),
		metricsFilters:    metrictree.NewMetricFilter(),
	}

	return cm
}

///////////////////////////////////////////////////////////////////////////////
// proxy.Collector related methods

func (cm *ContextManager) RequestCollect(id int) ([]plugin.Metric, error) {
	context, ok := cm.contextMap[id]
	if !ok {
		return nil, fmt.Errorf("can't find a context for a given id: %d", id)
	}

	cm.collector.Collect(context)
	return nil, nil
}

func (cm *ContextManager) LoadTask(id int, rawConfig []byte, mtsSelectors []string) error {
	if _, ok := cm.contextMap[id]; ok {
		return errors.New("context with given id was already defined")
	}

	newCtx, err := NewPluginContext(rawConfig, mtsSelectors)
	if err != nil {
		return fmt.Errorf("can't load task: %v", err)
	}
	cm.contextMap[id] = newCtx

	if loadable, ok := cm.collector.(plugin.LoadableCollector); ok {
		loadable.Load(cm.contextMap[id])
	}

	return nil
}

func (cm *ContextManager) UnloadTask(id int) error {
	if _, ok := cm.contextMap[id]; !ok {
		return errors.New("context with given id is not defined")
	}

	if loadable, ok := cm.collector.(plugin.LoadableCollector); ok {
		loadable.Unload(cm.contextMap[id])
	}

	delete(cm.contextMap, id)
	return nil
}

func (cm *ContextManager) RequestInfo() {
	return
}

///////////////////////////////////////////////////////////////////////////////
// plugin.CollectorDefinition related methods

func (cm *ContextManager) DefineMetric(ns string, isDefault bool, description string) {
	err := cm.metricsDefinition.AddRule(ns)
	if err != nil {
		log.WithError(err).WithFields(logrus.Fields{"namespace": ns}).Errorf("Wrong metric definition")
	}
}

// Define description for dynamic element
func (cm *ContextManager) DefineGroup(string, string) {
	panic("implement")
}

// Define global tags that will be applied to all metrics
func (cm *ContextManager) DefineGlobalTags(string, map[string]string) {
	panic("implement")
}

///////////////////////////////////////////////////////////////////////////////

func (cm *ContextManager) RequestPluginDefinition() {
	if definable, ok := cm.collector.(plugin.DefinableCollector); ok {
		definable.DefineMetrics(cm)
	}
}
