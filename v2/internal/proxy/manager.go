/*
Package proxy:
1) Manages context for different task created for the same plugin
2) Serves as an entry point for any "controller" (like. Rpc)
*/
package proxy

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/librato/snap-plugin-lib-go/v2/internal/stats"
	"github.com/librato/snap-plugin-lib-go/v2/internal/util/metrictree"
	"github.com/librato/snap-plugin-lib-go/v2/internal/util/types"
	"github.com/librato/snap-plugin-lib-go/v2/plugin"
	"github.com/sirupsen/logrus"
)

var log *logrus.Entry

func init() {
	log = logrus.WithFields(logrus.Fields{"layer": "lib", "module": "plugin-proxy"})
}

type Collector interface {
	RequestCollect(id int) ([]*types.Metric, error)
	LoadTask(id int, config []byte, selectors []string) error
	UnloadTask(id int) error
}

type metricMetadata struct {
	isDefault   bool
	description string
	unit        string
}

type ContextManager struct {
	collector  plugin.Collector // reference to custom plugin code
	contextMap sync.Map         // (synced map[int]*pluginContext) map of contexts associated with taskIDs

	activeTasksMutex sync.RWMutex     // mutex associated with activeTasks
	activeTasks      map[int]struct{} // map of active tasks (tasks for which Collect RPC request is progressing)

	metricsDefinition *metrictree.TreeValidator // metrics defined by plugin (code)

	metricsMetadata   map[string]metricMetadata // metadata associated with each metric (is default?, description, unit)
	groupsDescription map[string]string         // description associated with each group (dynamic element)

	StatsController *stats.StatsController // reference to statistics controller
}

func NewContextManager(collector plugin.Collector, pluginName string, pluginVersion string) *ContextManager {
	statsController := stats.NewStatsController(pluginName, pluginVersion)
	statsController.Run()

	cm := &ContextManager{
		collector:   collector,
		contextMap:  sync.Map{},
		activeTasks: map[int]struct{}{},

		metricsDefinition: metrictree.NewMetricDefinition(),

		metricsMetadata:   map[string]metricMetadata{},
		groupsDescription: map[string]string{},

		StatsController: statsController,
	}

	cm.RequestPluginDefinition()

	return cm
}

///////////////////////////////////////////////////////////////////////////////
// proxy.Collector related methods

func (cm *ContextManager) RequestCollect(id int) ([]*types.Metric, error) {
	if !cm.activateTask(id) {
		return nil, fmt.Errorf("can't process collect request, other request for the same id (%d) is in progress", id)
	}
	defer cm.markTaskAsCompleted(id)

	contextIf, ok := cm.contextMap.Load(id)
	if !ok {
		return nil, fmt.Errorf("can't find a context for a given id: %d", id)
	}
	context := contextIf.(*pluginContext)

	// collect metrics - user defined code
	context.sessionMts = []*types.Metric{}

	startTime := time.Now()
	err := cm.collector.Collect(context)
	endTime := time.Now()
	cm.StatsController.UpdateCollectStat(id, context.sessionMts, err != nil, startTime, endTime)

	if err != nil {
		return nil, fmt.Errorf("user-defined Collect method ended with error: %v", err)
	}

	return context.sessionMts, nil
}

func (cm *ContextManager) LoadTask(id int, rawConfig []byte, mtsFilter []string) error {
	if !cm.activateTask(id) {
		return fmt.Errorf("can't process load request, other request for the same id (%d) is in progress", id)
	}
	defer cm.markTaskAsCompleted(id)

	if _, ok := cm.contextMap.Load(id); ok {
		return errors.New("context with given id was already defined")
	}

	newCtx, err := NewPluginContext(cm, rawConfig)
	if err != nil {
		return fmt.Errorf("can't load task: %v", err)
	}

	for _, mtFilter := range mtsFilter {
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
	cm.StatsController.UpdateLoadStat(id, string(rawConfig), mtsFilter)

	return nil
}

func (cm *ContextManager) UnloadTask(id int) error {
	if !cm.activateTask(id) {
		return fmt.Errorf("can't process unload request, other request for the same id (%d) is in progress", id)
	}
	defer cm.markTaskAsCompleted(id)

	contextI, ok := cm.contextMap.Load(id)
	if !ok {
		return errors.New("context with given id is not defined")
	}

	context := contextI.(*pluginContext)
	if loadable, ok := cm.collector.(plugin.LoadableCollector); ok {
		err := loadable.Unload(context)
		if err != nil {
			return fmt.Errorf("error occured when trying to unload a task (%d): %v", id, err)
		}
	}

	cm.contextMap.Delete(id)
	cm.StatsController.UpdateUnloadStat(id)

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

///////////////////////////////////////////////////////////////////////////////

func (cm *ContextManager) RequestPluginDefinition() {
	if definable, ok := cm.collector.(plugin.DefinableCollector); ok {
		err := definable.DefineMetrics(cm)
		if err != nil {
			log.WithError(err).Errorf("Error occurred during plugin definition")
		}
	}
}

func (cm *ContextManager) activateTask(id int) bool {
	cm.activeTasksMutex.Lock()
	defer cm.activeTasksMutex.Unlock()

	if _, ok := cm.activeTasks[id]; ok {
		return false
	}

	cm.activeTasks[id] = struct{}{}
	return true
}

func (cm *ContextManager) markTaskAsCompleted(id int) {
	cm.activeTasksMutex.Lock()
	defer cm.activeTasksMutex.Unlock()

	delete(cm.activeTasks, id)
}
