package proxy

import (
	"sync"

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
	publisher  plugin.Publisher
	contextMap sync.Map // todo: is it possible to make those 3 common

	activeTasksMutex sync.RWMutex        // mutex associated with activeTasks
	activeTasks      map[string]struct{} // map of active tasks (tasks for which Collect RPC request is progressing)
}

func NewContextManager(publisher plugin.Publisher) *ContextManager {
	cm := &ContextManager{
		publisher:   publisher,
		contextMap:  sync.Map{},
		activeTasks: nil,
	}

	return cm
}

///////////////////////////////////////////////////////////////////////////////
// proxy.Publisher related methods

func (cm *ContextManager) RequestPublish(id string, mts []*types.Metric) error {
	return nil
}

func (cm *ContextManager) LoadTask(id string, config []byte) error {
	return nil
}

func (cm *ContextManager) UnloadTask(id string) error {
	return nil
}

///////////////////////////////////////////////////////////////////////////////

func (cm *ContextManager) activateTask(id string) bool {
	cm.activeTasksMutex.Lock()
	defer cm.activeTasksMutex.Unlock()

	if _, ok := cm.activeTasks[id]; ok {
		return false
	}

	cm.activeTasks[id] = struct{}{}
	return true
}

func (cm *ContextManager) markTaskAsCompleted(id string) {
	cm.activeTasksMutex.Lock()
	defer cm.activeTasksMutex.Unlock()

	delete(cm.activeTasks, id)
}
