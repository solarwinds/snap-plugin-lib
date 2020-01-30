package proxy

import (
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
	if !cm.ActivateTask(id) {
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
	context := contextIf.(*pluginContext)

	context.sessionMts = mts // metrics to publish are set within context
	context.ResetWarnings()

	startTime := time.Now()
	err := cm.publisher.Publish(context) // calling to user defined code
	endTime := time.Now()

	cm.statsController.UpdateExecutionStat(id, len(context.sessionMts), err != nil, startTime, endTime)

	if err != nil {
		return types.ProcessingStatus{
			Error:    fmt.Errorf("user-defined Publish method ended with error: %v", err),
			Warnings: context.Warnings(),
		}
	}

	log.WithFields(logrus.Fields{
		"elapsed": endTime.Sub(startTime).String(),
		"metrics": len(mts),
	}).Debug("Publish completed")

	return types.ProcessingStatus{
		Warnings: context.Warnings(),
	}
}

func (cm *ContextManager) LoadTask(id string, config []byte) error {
	if !cm.ActivateTask(id) {
		return fmt.Errorf("can't process load request, other request for the same id (%s) is in progress", id)
	}
	defer cm.MarkTaskAsCompleted(id)

	if _, ok := cm.contextMap.Load(id); ok {
		return errors.New("context with given id was already defined")
	}

	newCtx, err := NewPluginContext(cm, config)
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
	if !cm.ActivateTask(id) {
		return fmt.Errorf("can't process unload request, other request for the same id (%s) is in progress", id)
	}
	defer cm.MarkTaskAsCompleted(id)

	contextI, ok := cm.contextMap.Load(id)
	if !ok {
		return errors.New("context with given id is not defined")
	}

	context := contextI.(*pluginContext)
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

func (cm *ContextManager) RequestPluginDefinition() {
	if definable, ok := cm.publisher.(plugin.DefinablePublisher); ok {
		err := definable.PluginDefinition(cm)
		if err != nil {
			log.WithError(err).Errorf("Error occurred during plugin definition")
		}
	}
}
