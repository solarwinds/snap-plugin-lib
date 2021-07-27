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

/*
Package proxy :
1) Manages context for different task created for the same plugin
2) Serves as an entry point for any "controller" (like. Rpc)
*/
package proxy

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"runtime/debug"
	"sort"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	commonProxy "github.com/solarwinds/snap-plugin-lib/v2/internal/plugins/common/proxy"
	"github.com/solarwinds/snap-plugin-lib/v2/internal/plugins/common/stats"
	"github.com/solarwinds/snap-plugin-lib/v2/internal/util/log"
	"github.com/solarwinds/snap-plugin-lib/v2/internal/util/metrictree"
	"github.com/solarwinds/snap-plugin-lib/v2/internal/util/types"
	"github.com/solarwinds/snap-plugin-lib/v2/plugin"
)

var moduleFields = logrus.Fields{"layer": "lib", "module": "collector-proxy"}

const (
	RequestAllMetricsFilter = "/*"

	unloadMaxRetries    = 3
	unloadRetryInterval = 1 * time.Second

	streamingCheckInterval = 1 * time.Second
)

type Collector interface {
	RequestCollect(id string) <-chan types.CollectChunk
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
	ctx context.Context

	collector  types.Collector // reference to custom plugin code
	contextMap sync.Map        // (synced map[int]*PluginContext) map of contexts associated with taskIDs

	metricsDefinition *metrictree.TreeValidator // metrics defined by plugin (code)

	metricsMetadata   map[string]metricMetadata // metadata associated with each metric (is default?, description, unit)
	groupsDescription map[string]string         // description associated with each group (dynamic element)

	statsController stats.Controller // reference to statistics controller
}

func NewContextManager(ctx context.Context, collector types.Collector, statsController stats.Controller) *ContextManager {
	cm := &ContextManager{
		ContextManager: commonProxy.NewContextManager(),
		ctx:            ctx,

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

func (cm *ContextManager) Context() context.Context {
	return cm.ctx
}

func (cm *ContextManager) RequestCollect(id string) <-chan types.CollectChunk {
	chunkCh := make(chan types.CollectChunk)
	go cm.requestCollect(id, chunkCh)
	return chunkCh
}

func (cm *ContextManager) requestCollect(id string, chunkCh chan<- types.CollectChunk) {
	if !cm.AcquireTask(id) {
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

	pContext := contextIf.(*PluginContext)

	pContext.AttachContext(cm.TaskContext(id))
	pContext.ClearCollectorSession()
	pContext.ResetWarnings()

	switch cm.collector.Type() {
	case types.PluginTypeCollector:
		cm.collect(id, pContext, chunkCh)
	case types.PluginTypeStreamingCollector:
		cm.streamingCollect(id, pContext, chunkCh)
	}

	cm.MarkTaskAsCompleted(id)
	pContext.ReleaseContext()
}

func (cm *ContextManager) collect(id string, context *PluginContext, chunkCh chan<- types.CollectChunk) {
	logF := cm.logger()
	taskCtx := cm.TaskContext(id)

	var mts []*types.Metric
	var warnings []types.Warning
	var err error

	go func() {
		defer func() {
			// catch panics (since it's running in it's own goroutine)
			if r := recover(); r != nil {
				logF.WithError(fmt.Errorf("%v", r)).Error("user-defined function has ended with panic")
				logF.WithField("block", "recover").Trace(string(debug.Stack()))
				err = fmt.Errorf("user-defined function has ended with panic: %v", r)
			}

			cm.ReleaseTask(id)
		}()

		startTime := time.Now()
		err = cm.collector.Collect(context) // calling to user defined code
		endTime := time.Now()

		if !context.Context.IsDone() {
			mts = context.Metrics(false)
			warnings = context.Warnings(false)

			cm.statsController.UpdateExecutionStat(id, len(mts), err != nil, startTime, endTime)

			if err != nil {
				err = fmt.Errorf("user-defined Collect method ended with error: %v", err)
			}

			logF.WithFields(logrus.Fields{
				"elapsed":      endTime.Sub(startTime).String(),
				"metrics-num":  len(mts),
				"warnings-num": len(warnings),
			}).Debug("Collect completed")
		} else {
			logF.WithFields(logrus.Fields{
				"elapsed": endTime.Sub(startTime).String(),
			}).Info("Collect completed after task had been canceled")
		}
	}()

	<-taskCtx.Done()

	chunkCh <- types.CollectChunk{
		Metrics:  mts,
		Warnings: warnings,
		Err:      err,
	}

	close(chunkCh)
}

func (cm *ContextManager) streamingCollect(id string, context *PluginContext, chunkCh chan<- types.CollectChunk) {
	logF := cm.logger()
	var err error

	startTime := time.Now()

	taskCtx := cm.TaskContext(id)

	go func() {
		defer func() {
			// catch panics (since it's running in it's own goroutine)
			if r := recover(); r != nil {
				logF.WithError(fmt.Errorf("%v", r)).Error("user-defined function has ended with panic")
				logF.WithField("block", "recover").Trace(string(debug.Stack()))
				err = fmt.Errorf("user-defined function has ended with panic: %v", r)
			}

			cm.ReleaseTask(id)
		}()

		err = cm.collector.StreamingCollect(context)
	}()

	for {
		select {
		case <-taskCtx.Done():
			cm.handleChunk(id, err, context, chunkCh, startTime)
			close(chunkCh)
			return
		case <-time.After(streamingCheckInterval):
			cm.handleChunk(id, err, context, chunkCh, startTime)
		}
	}
}

func (cm *ContextManager) handleChunk(id string, err error, context *PluginContext, chunkCh chan<- types.CollectChunk, startTime time.Time) {
	mts := context.Metrics(true)
	warnings := context.Warnings(true)

	if len(mts) > 0 || len(warnings) > 0 || err != nil {
		lastUpdate := time.Now()

		chunkCh <- types.CollectChunk{
			Metrics:  mts,
			Warnings: warnings,
			Err:      err,
		}

		cm.statsController.UpdateStreamingStat(id, len(mts), startTime, lastUpdate)
	}
}

func (cm *ContextManager) LoadTask(id string, rawConfig []byte, mtsFilter []string) error {
	if !cm.AcquireTask(id) {
		return fmt.Errorf("can't process load request, other request for the same id (%s) is in progress", id)
	}
	defer cm.MarkTaskAsCompleted(id)

	if _, ok := cm.contextMap.Load(id); ok {
		return errors.New("context with given id was already defined")
	}
	newCtx, err := NewPluginContext(cm, id, rawConfig)
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
			return fmt.Errorf("wrong filtering rule (%v): %v", mtFilter, err)
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
	logF := cm.logger()

	// Unload may be called when Collect (especially stream) is in progress. If so, try to cancel it.
	for retry := 1; retry <= unloadMaxRetries; retry++ {
		ok := cm.AcquireTask(id)
		if !ok {
			if retry == unloadMaxRetries {
				return fmt.Errorf("can't process unload request, unable to cancel other task with the same ID")
			}

			logF.WithFields(moduleFields).WithField("task-id", id).Trace("other action is active, requesting stop")

			cm.ReleaseTask(id)
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

	pluginCtx := contextI.(*PluginContext)
	if unloadable, ok := cm.collector.Unwrap().(plugin.UnloadableCollector); ok {
		err := unloadable.Unload(pluginCtx)
		if err != nil {
			return fmt.Errorf("error occured when trying to unload a task (%s): %v", id, err)
		}
	}

	cm.contextMap.Delete(id)
	cm.statsController.UpdateUnloadStat(id)

	return nil
}

func (cm *ContextManager) CustomInfo(id string) ([]byte, error) {
	// Do not call cm.AcquireTask as above methods. CustomInfo is read-only

	contextI, ok := cm.contextMap.Load(id)
	if !ok {
		return nil, errors.New("context with given id is not defined")
	}
	pluginCtx := contextI.(*PluginContext)

	if collectorWithCustomInfo, ok := cm.collector.Unwrap().(plugin.CustomizableInfoCollector); ok {
		infoObj := collectorWithCustomInfo.CustomInfo(pluginCtx)

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
	logF := cm.logger()

	err := cm.metricsDefinition.AddRule(ns)
	if err != nil {
		logF.WithFields(moduleFields).WithError(err).WithFields(logrus.Fields{"namespace": ns}).Fatal("Wrong metric definition")
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

func (cm *ContextManager) SetAllowDynamicLastElement() {
	cm.metricsDefinition.AllowDynamicLastElement()
}

func (cm *ContextManager) SetAllowAddingUndefinedMetrics() {
	cm.metricsDefinition.AllowAddingUndefinedMetrics()
}

func (cm *ContextManager) SetAllowValuesAtAnyNamespaceLevel() {
	cm.metricsDefinition.AllowValuesAtAnyNamespaceLevel()
}

///////////////////////////////////////////////////////////////////////////////

func (cm *ContextManager) RequestPluginDefinition() {
	logF := cm.logger()

	if definable, ok := cm.collector.Unwrap().(plugin.DefinableCollector); ok {
		err := definable.PluginDefinition(cm)
		if err != nil {
			logF.WithFields(moduleFields).WithError(err).Errorf("Error occurred during plugin definition")
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

func (cm *ContextManager) logger() logrus.FieldLogger {
	return log.WithCtx(cm.ctx).WithFields(moduleFields).WithField("service", "manager")
}
