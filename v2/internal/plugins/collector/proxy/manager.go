/*
Package proxy:
1) Manages context for different task created for the same plugin
2) Serves as an entry point for any "controller" (like. Rpc)
*/
package proxy

import (
	"errors"
	"fmt"
	"sort"
	"sync"
	"time"

	commonProxy "github.com/librato/snap-plugin-lib-go/v2/internal/plugins/common/proxy"
	"github.com/librato/snap-plugin-lib-go/v2/internal/plugins/common/stats"
	"github.com/librato/snap-plugin-lib-go/v2/internal/util/metrictree"
	"github.com/librato/snap-plugin-lib-go/v2/internal/util/types"
	"github.com/librato/snap-plugin-lib-go/v2/plugin"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

var log = logrus.WithFields(logrus.Fields{"layer": "lib", "module": "collector-proxy"})

const (
	RequestAllMetricsFilter = "/*"
)

type Collector interface {
	RequestCollect(id string) ([]*types.Metric, types.ProcessingStatus)
	LoadTask(id string, config []byte, selectors []string) error
	UnloadTask(id string) error
}

type metricMetadata struct {
	isDefault   bool
	description string
	unit        string
}

type ContextManager struct {
	*commonProxy.ContextManager

	collector  plugin.Collector // reference to custom plugin code
	contextMap sync.Map         // (synced map[int]*pluginContext) map of contexts associated with taskIDs

	metricsDefinition *metrictree.TreeValidator // metrics defined by plugin (code)

	metricsMetadata   map[string]metricMetadata // metadata associated with each metric (is default?, description, unit)
	groupsDescription map[string]string         // description associated with each group (dynamic element)

	statsController stats.Controller // reference to statistics controller

	ExampleConfig yaml.Node // example config
}

func NewContextManager(collector plugin.Collector, statsController stats.Controller) *ContextManager {
	cm := &ContextManager{
		ContextManager: commonProxy.NewContextManager(),

		collector:  collector,
		contextMap: sync.Map{},

		metricsDefinition: metrictree.NewMetricDefinition(),

		metricsMetadata:   map[string]metricMetadata{},
		groupsDescription: map[string]string{},

		statsController: statsController,
	}

	cm.RequestPluginDefinition()

	return cm
}

///////////////////////////////////////////////////////////////////////////////
// proxy.Collector related methods

func (cm *ContextManager) RequestCollect(id string) ([]*types.Metric, types.ProcessingStatus) {
	if !cm.ActivateTask(id) {
		return nil, types.ProcessingStatus{
			Error: fmt.Errorf("can't process collect request, other request for the same id (%s) is in progress", id),
		}
	}
	defer cm.MarkTaskAsCompleted(id)

	contextIf, ok := cm.contextMap.Load(id)
	if !ok {
		return nil, types.ProcessingStatus{
			Error: fmt.Errorf("can't find a context for a given id: %s", id),
		}
	}
	context := contextIf.(*pluginContext)

	context.sessionMts = []*types.Metric{}
	context.ResetWarnings()

	startTime := time.Now()
	err := cm.collector.Collect(context) // calling to user defined code
	endTime := time.Now()

	cm.statsController.UpdateExecutionStat(id, len(context.sessionMts), err != nil, startTime, endTime)

	if err != nil {
		return nil, types.ProcessingStatus{
			Error:    fmt.Errorf("user-defined Collect method ended with error: %v", err),
			Warnings: context.Warnings(),
		}
	}

	log.WithFields(logrus.Fields{
		"elapsed":      endTime.Sub(startTime).String(),
		"metrics-num":  len(context.sessionMts),
		"warnings-num": len(context.Warnings()),
	}).Debug("Collect completed")

	return context.sessionMts, types.ProcessingStatus{
		Warnings: context.Warnings(),
	}
}

func (cm *ContextManager) LoadTask(id string, rawConfig []byte, mtsFilter []string) error {
	if !cm.ActivateTask(id) {
		return fmt.Errorf("can't process load request, other request for the same id (%s) is in progress", id)
	}
	defer cm.MarkTaskAsCompleted(id)

	if _, ok := cm.contextMap.Load(id); ok {
		return errors.New("context with given id was already defined")
	}

	newCtx, err := NewPluginContext(cm, rawConfig)
	if err != nil {
		return fmt.Errorf("can't load task: %v", err)
	}

	for _, mtFilter := range mtsFilter {
		// If requested metrics are not provided in config, snap sends requested metric in a form of "/*"
		if mtFilter == RequestAllMetricsFilter {
			continue
		}

		err := newCtx.metricsFilters.AddRule(mtFilter)
		if err != nil {
			log.WithError(err).WithField("rule", mtFilter).Warn("can't add filtering rule, it will be ignored")
		}
	}

	if loadable, ok := cm.collector.(plugin.LoadableCollector); ok {
		err := loadable.Load(newCtx)
		if err != nil {
			return fmt.Errorf("can't load task due to errors returned from user-defined function: %s", err)
		}
	}

	cm.contextMap.Store(id, newCtx)
	cm.statsController.UpdateLoadStat(id, string(rawConfig), mtsFilter)

	return nil
}

func (cm *ContextManager) UnloadTask(id string) error {
	if !cm.ActivateTask(id) {
		return fmt.Errorf("can't process unload request, other request for the same id (%s) is in progress", id)
	}
	defer cm.MarkTaskAsCompleted(id)

	contextI, ok := cm.contextMap.Load(id)
	if !ok {
		return errors.New("context with given id is not defined")
	}

	context := contextI.(*pluginContext)
	if unloadable, ok := cm.collector.(plugin.UnloadableCollector); ok {
		err := unloadable.Unload(context)
		if err != nil {
			return fmt.Errorf("error occured when trying to unload a task (%s): %v", id, err)
		}
	}

	cm.contextMap.Delete(id)
	cm.statsController.UpdateUnloadStat(id)

	return nil
}

///////////////////////////////////////////////////////////////////////////////
// plugin.CollectorDefinition related methods

func (cm *ContextManager) DefineMetric(ns string, unit string, isDefault bool, description string) {
	err := cm.metricsDefinition.AddRule(ns)
	if err != nil {
		log.WithError(err).WithFields(logrus.Fields{"namespace": ns}).Errorf("Wrong metric definition")
	}

	cm.metricsMetadata[ns] = metricMetadata{
		isDefault:   isDefault,
		description: description,
		unit:        unit,
	}
}

// Define description for dynamic element
func (cm *ContextManager) DefineGroup(name string, description string) {
	cm.groupsDescription[name] = description
}

// Define global tags that will be applied to all metrics
func (cm *ContextManager) DefineGlobalTags(string, map[string]string) {
	panic("implement")
}

func (cm *ContextManager) DefineExampleConfig(cfg string) error {
	err := yaml.Unmarshal([]byte(cfg), &cm.ExampleConfig)
	if err != nil {
		return fmt.Errorf("invalid YAML provided by user: %v", err)
	}

	return nil
}

///////////////////////////////////////////////////////////////////////////////

func (cm *ContextManager) RequestPluginDefinition() {
	if definable, ok := cm.collector.(plugin.DefinableCollector); ok {
		err := definable.PluginDefinition(cm)
		if err != nil {
			log.WithError(err).Errorf("Error occurred during plugin definition")
		}
	}
}

///////////////////////////////////////////////////////////////////////////////

func (cm *ContextManager) ListDefaultMetrics() []string {
	result := []string{}
	for mt, meta := range cm.metricsMetadata {
		if meta.isDefault {
			result = append(result, mt)
		}
	}

	sort.Strings(result)

	return result
}
