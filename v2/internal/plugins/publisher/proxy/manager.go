package proxy

import (
	"errors"
	"fmt"
	"sync"
	"time"

	commonProxy "github.com/librato/snap-plugin-lib-go/v2/internal/plugins/common/proxy"
	"github.com/librato/snap-plugin-lib-go/v2/internal/util/types"
	"github.com/librato/snap-plugin-lib-go/v2/plugin"
	"github.com/sirupsen/logrus"
)

var log *logrus.Entry

func init() {
	log = logrus.WithFields(logrus.Fields{"layer": "lib", "module": "publisher-proxy"})
}

type Publisher interface {
	RequestPublish(id string, mts []*types.Metric) error
	LoadTask(id string, config []byte) error
	UnloadTask(id string) error
}

type ContextManager struct {
	*commonProxy.ContextManager

	publisher  plugin.Publisher
	contextMap sync.Map
}

func NewContextManager(publisher plugin.Publisher) *ContextManager {
	cm := &ContextManager{
		ContextManager: commonProxy.NewContextManager(),
		publisher:      publisher,
		contextMap:     sync.Map{},
	}

	return cm
}

///////////////////////////////////////////////////////////////////////////////
// proxy.Publisher related methods

func (cm *ContextManager) RequestPublish(id string, mts []*types.Metric) error {
	if !cm.ActivateTask(id) {
		return fmt.Errorf("can't process publish request, other request for the same id (%s) is in progress", id)
	}
	defer cm.MarkTaskAsCompleted(id)

	contextIf, ok := cm.contextMap.Load(id)
	if !ok {
		return fmt.Errorf("can't find a context for a given id: %s", id)
	}
	context := contextIf.(*pluginContext)

	context.sessionMts = mts // metrics to publish are set withing context

	startTime := time.Now()
	err := cm.publisher.Publish(context) // calling to user defined code
	endTime := time.Now()

	// todo: update statistics https://swicloud.atlassian.net/browse/AO-14142

	if err != nil {
		return fmt.Errorf("user-defined Publish method ended with error: %v", err)
	}

	log.WithField("elapsed", endTime.Sub(startTime).String()).Debug("Publish completed")

	return nil
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

	// todo: update statistics https://swicloud.atlassian.net/browse/AO-14142

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
	if loadable, ok := cm.publisher.(plugin.LoadablePublisher); ok {
		err := loadable.Unload(context)
		if err != nil {
			return fmt.Errorf("error occured when trying to unload a publisher task (%s): %v", id, err)
		}
	}

	cm.contextMap.Delete(id)

	// todo: update statistics https://swicloud.atlassian.net/browse/AO-14142

	return nil
}