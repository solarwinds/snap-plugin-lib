package proxy

import (
	"sync"

	"github.com/librato/snap-plugin-lib-go/v2/plugin"
)

type ContextManager struct {
	activeTasksMutex sync.RWMutex        // mutex associated with activeTasks
	activeTasks      map[string]struct{} // map of active tasks (tasks for which Collect RPC request is progressing)

	TasksLimit     int
	InstancesLimit int
}

func NewContextManager() *ContextManager {
	return &ContextManager{
		activeTasks:    map[string]struct{}{},
		TasksLimit:     plugin.NoLimit,
		InstancesLimit: plugin.NoLimit,
	}
}

func (cm *ContextManager) ActivateTask(id string) bool {
	cm.activeTasksMutex.Lock()
	defer cm.activeTasksMutex.Unlock()

	if _, ok := cm.activeTasks[id]; ok {
		return false
	}

	cm.activeTasks[id] = struct{}{}
	return true
}

func (cm *ContextManager) MarkTaskAsCompleted(id string) {
	cm.activeTasksMutex.Lock()
	defer cm.activeTasksMutex.Unlock()

	delete(cm.activeTasks, id)
}

func (cm *ContextManager) DefineTasksPerInstanceLimit(limit int) {
	if limit < -1 {
		panic("invalid tasks limit")
	}

	cm.TasksLimit = limit
}

func (cm *ContextManager) DefineInstancesLimit(limit int) {
	if limit < -1 {
		panic("invalid instances limit")
	}

	cm.InstancesLimit = limit
}
