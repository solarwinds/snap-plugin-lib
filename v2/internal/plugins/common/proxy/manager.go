package proxy

import "sync"

type ContextManager struct {
	activeTasksMutex sync.RWMutex        // mutex associated with activeTasks
	activeTasks      map[string]struct{} // map of active tasks (tasks for which Collect RPC request is progressing)
}

func NewContextManager() *ContextManager {
	return &ContextManager{
		activeTasks: map[string]struct{}{},
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
