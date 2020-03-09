/*
Package proxy:
1) Manages context for different task created for the same plugin
2) Serves as an entry point for any "controller" (like. Rpc)
*/
package proxy

import (
	"encoding/json"
	"errors"
	"fmt"
	"runtime/debug"
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

	unloadMaxRetries    = 3
	unloadRetryInterval = 1 * time.Second

	streamingCheckInterval = 1 * time.Second
)

type Collector interface {
	RequestCollect(id string) chan types.CollectChunk
	LoadTask(id string, config []byte, selectors []string) error
	UnloadTask(id string) error
	CustomInfo(id string) ([]byte, error)
}

type metricMetadata struct {
	isDefault   bool
	description string
	unit        string
}

type ContextManager struct {
	*commonProxy.ContextManager

	collector  types.Collector // reference to custom plugin code
	contextMap sync.Map        // (synced map[int]*pluginContext) map of contexts associated with taskIDs

	metricsDefinition *metrictree.TreeValidator // metrics defined by plugin (code)

	metricsMetadata   map[string]metricMetadata // metadata associated with each metric (is default?, description, unit)
	groupsDescription map[string]string         // description associated with each group (dynamic element)

	statsController stats.Controller // reference to statistics controller

	ExampleConfig yaml.Node // example config
}

func NewContextManager(collector types.Collector, statsController stats.Controller) *ContextManager {
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

func (cm *ContextManager) RequestCollect(id string) chan types.CollectChunk {
	chunkCh := make(chan types.CollectChunk)
	go cm.requestCollect(id, chunkCh)
	return chunkCh
}

func (cm *ContextManager) requestCollect(id string, chunkCh chan types.CollectChunk) {
	if !cm.ActivateTask(id) {
		chunkCh <- types.CollectChunk{
			Err: fmt.Errorf("can't process collect request, other request for the same id (%s) is in progress", id),
		}
		close(chunkCh)
		return
	}

	contextIf, ok := cm.contextMap.Load(id)
	if !ok {
		chunkCh <- types.CollectChunk{
			Err: fmt.Errorf("can't find a context for a given id: %s", id),
		}
		close(chunkCh)
		return
	}

	pContext := contextIf.(*pluginContext)

	pContext.AttachContext(cm.TaskContext(id))
	pContext.sessionMts = []*types.Metric{}
	pContext.ResetWarnings()

	switch cm.collector.Type() {
	case types.PluginTypeCollector:
		go func() {
			cm.collect(id, pContext, chunkCh)
			cm.MarkTaskAsCompleted(id)
		}()
	case types.PluginTypeStreamingCollector:
		go func() {
			cm.streamingCollect(id, pContext, chunkCh)
			cm.MarkTaskAsCompleted(id)
		}()
	}
}

func (cm *ContextManager) collect(id string, context *pluginContext, chunkCh chan types.CollectChunk) {
	startTime := time.Now()
	err := cm.collector.Collect(context) // calling to user defined code
	endTime := time.Now()

	cm.statsController.UpdateExecutionStat(id, len(context.sessionMts), err != nil, startTime, endTime)

	if err != nil {
		err = fmt.Errorf("user-defined Collect method ended with error: %v", err)
	}

	chunkCh <- types.CollectChunk{
		Metrics:  context.sessionMts,
		Warnings: context.Warnings(),
		Err:      err,
	}

	close(chunkCh)

	log.WithFields(logrus.Fields{
		"elapsed":      endTime.Sub(startTime).String(),
		"metrics-num":  len(context.sessionMts),
		"warnings-num": len(context.Warnings()),
	}).Debug("Collect completed")
}

func (cm *ContextManager) streamingCollect(id string, context *pluginContext, chunkCh chan types.CollectChunk) {
	startTime := time.Now()

	wg := sync.WaitGroup{}

	taskCtx := cm.TaskContext(id)

	wg.Add(1)
	go func() {
		defer func() {
			wg.Done()

			// panic may be caused by two reasons:
			// - user-defined function performed invalid operation
			// - user-defined function called context API after it had been marked as dead
			if r := recover(); r != nil {
				log.WithError(fmt.Errorf("%v", r)).Warn("user-defined function has been ended")
				log.Trace(string(debug.Stack()))
			}
		}()
		for {
			select {
			case <-taskCtx.Done():
				return
			default:
				cm.collector.StreamingCollect(context)
			}
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-taskCtx.Done():
				close(chunkCh)
				return
			case <-time.After(streamingCheckInterval):
				mts := context.sessionMts
				warnings := context.Warnings()

				if len(mts) > 0 || len(warnings) > 0 {
					lastUpdate := time.Now()

					chunkCh <- types.CollectChunk{
						Metrics:  context.sessionMts,
						Warnings: context.Warnings(),
					}

					cm.statsController.UpdateExecutionStat(id, len(context.sessionMts), true, startTime, lastUpdate)
				}

				context.sessionMts = nil
				context.ResetWarnings()

				// todo: adamik: synchro
			}
		}
	}()

	wg.Wait()
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

	if loadable, ok := cm.collector.Unwrap().(plugin.LoadableCollector); ok {
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
	// Unload may be called when Collect (especially stream) is in progress. If so, try to cancel it.
	for retry := 1; retry <= unloadMaxRetries; retry++ {
		ok := cm.ActivateTask(id)
		if !ok {
			if retry == unloadMaxRetries {
				return fmt.Errorf("can't process unload request, unable to cancel other task with the same ID")
			}

			log.WithField("taskID", id).Trace("other action is active, requesting stop")
			cm.CancelTask(id)
			time.Sleep(unloadRetryInterval)

			continue
		}

		defer cm.MarkTaskAsCompleted(id)
		break
	}

	contextI, ok := cm.contextMap.Load(id)
	if !ok {
		return errors.New("context with given id is not defined")
	}

	context := contextI.(*pluginContext)
	if unloadable, ok := cm.collector.Unwrap().(plugin.UnloadableCollector); ok {
		err := unloadable.Unload(context)
		if err != nil {
			return fmt.Errorf("error occured when trying to unload a task (%s): %v", id, err)
		}
	}

	cm.contextMap.Delete(id)
	cm.statsController.UpdateUnloadStat(id)

	return nil
}

func (cm *ContextManager) CustomInfo(id string) ([]byte, error) {
	// Do not call cm.ActivateTask as above methods. CustomInfo is read-only

	contextI, ok := cm.contextMap.Load(id)
	if !ok {
		return nil, errors.New("context with given id is not defined")
	}
	context := contextI.(*pluginContext)

	if collectorWithCustomInfo, ok := cm.collector.Unwrap().(plugin.CustomizableInfoCollector); ok {
		infoObj := collectorWithCustomInfo.CustomInfo(context)

		infoJSON, err := json.Marshal(infoObj)
		if err != nil {
			return nil, fmt.Errorf("can't unmarshal custom info to JSON: %v", err)
		}

		return infoJSON, nil
	}

	return []byte{}, nil
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
	if definable, ok := cm.collector.Unwrap().(plugin.DefinableCollector); ok {
		err := definable.PluginDefinition(cm)
		if err != nil {
			log.WithError(err).Errorf("Error occurred during plugin definition")
		}
	}
}

///////////////////////////////////////////////////////////////////////////////

func (cm *ContextManager) ListDefaultMetrics() []string {
	var result []string
	for mt, meta := range cm.metricsMetadata {
		if meta.isDefault {
			result = append(result, mt)
		}
	}

	sort.Strings(result)

	return result
}
