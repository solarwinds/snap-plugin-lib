/*
 Copyright (c) 2020 SolarWinds Worldwide, LLC

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
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"

	commonProxy "github.com/librato/snap-plugin-lib-go/v2/internal/plugins/common/proxy"
	"github.com/librato/snap-plugin-lib-go/v2/internal/plugins/common/stats"
	"github.com/librato/snap-plugin-lib-go/v2/internal/util/types"
	"github.com/librato/snap-plugin-lib-go/v2/plugin"
	"github.com/sirupsen/logrus"
)

var log *logrus.Entry

func init() {
	log = logrus.WithFields(logrus.Fields{"layer": "lib", "module": "publisher-proxy"})
}

type Publisher interface {
	RequestPublish(id string, mts []*types.Metric) types.ProcessingStatus
	LoadTask(id string, config []byte) error
	UnloadTask(id string) error
	CustomInfo(id string) ([]byte, error)
}

type ContextManager struct {
	*commonProxy.ContextManager

	publisher  plugin.Publisher
	contextMap sync.Map

	statsController stats.Controller // reference to statistics controller
}

func NewContextManager(publisher plugin.Publisher, statsController stats.Controller) *ContextManager {
	cm := &ContextManager{
		ContextManager: commonProxy.NewContextManager(),
		publisher:      publisher,
		contextMap:     sync.Map{},

		statsController: statsController,
	}

	cm.RequestPluginDefinition()

	return cm
}

///////////////////////////////////////////////////////////////////////////////
// proxy.Publisher related methods

func (cm *ContextManager) RequestPublish(id string, mts []*types.Metric) types.ProcessingStatus {
	if !cm.AcquireTask(id) {
		return types.ProcessingStatus{
			Error: fmt.Errorf("can't process publish request, other request for the same id (%s) is in progress", id),
		}
	}
	defer cm.MarkTaskAsCompleted(id)

	contextIf, ok := cm.contextMap.Load(id)
	if !ok {
		return types.ProcessingStatus{
			Error: fmt.Errorf("can't find a context for a given id: %s", id),
		}
	}
	context := contextIf.(*PluginContext)

	context.sessionMts = mts // metrics to publish are set within context
	context.ResetWarnings()

	startTime := time.Now()
	err := cm.publisher.Publish(context) // calling to user defined code
	warnings := context.Warnings(false)
	endTime := time.Now()

	cm.statsController.UpdateExecutionStat(id, len(context.sessionMts), err != nil, startTime, endTime)

	if err != nil {
		return types.ProcessingStatus{
			Error:    fmt.Errorf("user-defined Publish method ended with error: %v", err),
			Warnings: warnings,
		}
	}

	log.WithFields(logrus.Fields{
		"elapsed":      endTime.Sub(startTime).String(),
		"metrics-num":  len(mts),
		"warnings-num": len(warnings),
	}).Debug("Publish completed")

	return types.ProcessingStatus{
		Warnings: warnings,
	}
}

func (cm *ContextManager) LoadTask(id string, config []byte) error {
	if !cm.AcquireTask(id) {
		return fmt.Errorf("can't process load request, other request for the same id (%s) is in progress", id)
	}
	defer cm.MarkTaskAsCompleted(id)

	if _, ok := cm.contextMap.Load(id); ok {
		return errors.New("context with given id was already defined")
	}

	newCtx, err := NewPluginContext(cm,id, config)
	if err != nil {
		return fmt.Errorf("can't load task: %v", err)
	}

	if loadable, ok := cm.publisher.(plugin.LoadablePublisher); ok {
		err := loadable.Load(newCtx)
		if err != nil {
			return fmt.Errorf("can't load task due to errors returned from user-defined function: %s", err)
		}
	}

	cm.contextMap.Store(id, newCtx)
	cm.statsController.UpdateLoadStat(id, string(config), nil)

	return nil
}

func (cm *ContextManager) UnloadTask(id string) error {
	if !cm.AcquireTask(id) {
		return fmt.Errorf("can't process unload request, other request for the same id (%s) is in progress", id)
	}
	defer cm.MarkTaskAsCompleted(id)

	contextI, ok := cm.contextMap.Load(id)
	if !ok {
		return errors.New("context with given id is not defined")
	}

	context := contextI.(*PluginContext)
	if unloadable, ok := cm.publisher.(plugin.UnloadablePublisher); ok {
		err := unloadable.Unload(context)
		if err != nil {
			return fmt.Errorf("error occured when trying to unload a publisher task (%s): %v", id, err)
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
	context := contextI.(*PluginContext)

	if publisherWithCustomInfo, ok := cm.publisher.(plugin.CustomizableInfoPublisher); ok {
		infoObj := publisherWithCustomInfo.CustomInfo(context)

		infoJSON, err := json.Marshal(infoObj)
		if err != nil {
			return nil, fmt.Errorf("can't unmarshal custom info to JSON: %v", err)
		}

		return infoJSON, nil
	}

	return []byte{}, nil
}

func (cm *ContextManager) RequestPluginDefinition() {
	if definable, ok := cm.publisher.(plugin.DefinablePublisher); ok {
		err := definable.PluginDefinition(cm)
		if err != nil {
			log.WithError(err).Errorf("Error occurred during plugin definition")
		}
	}
}
