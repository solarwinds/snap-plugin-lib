package proxy

import (
	"context"
	"fmt"
	"sync"

	"github.com/librato/snap-plugin-lib-go/v2/plugin"
)

type contextHolder struct {
	ctx      context.Context
	cancelFn context.CancelFunc
}

type ContextManager struct {
	activeTasksMutex sync.RWMutex             // mutex associated with activeTasks
	activeTasks      map[string]contextHolder // map of active tasks (tasks for which Collect RPC request is progressing)

	TasksLimit     int
	InstancesLimit int
}

func NewContextManager() *ContextManager {
	return &ContextManager{
		activeTasks:    map[string]contextHolder{},
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

	ctx, cancelFn := context.WithCancel(context.Background())

	cm.activeTasks[id] = contextHolder{
		ctx:      ctx,
		cancelFn: cancelFn,
	}
	return true
}

func (cm *ContextManager) MarkTaskAsCompleted(id string) {
	cm.activeTasksMutex.Lock()
	defer cm.activeTasksMutex.Unlock()

	cm.activeTasks[id].cancelFn()
	delete(cm.activeTasks, id)
}

func (cm *ContextManager) CancelTask(id string) {
	cm.activeTasksMutex.Lock()
	defer cm.activeTasksMutex.Unlock()

	if aTask, ok := cm.activeTasks[id]; ok {
		aTask.cancelFn()
	}
}

func (cm *ContextManager) TaskContext(id string) context.Context {
	cm.activeTasksMutex.Lock()
	defer cm.activeTasksMutex.Unlock()

	return cm.activeTasks[id].ctx
}

func (cm *ContextManager) DefineTasksPerInstanceLimit(limit int) error {
	if limit < -1 {
		return fmt.Errorf("invalid tasks limit")
	}

	cm.TasksLimit = limit
	return nil
}

func (cm *ContextManager) DefineInstancesLimit(limit int) error {
	if limit < -1 {
		return fmt.Errorf("invalid instances limit")
	}

	cm.InstancesLimit = limit
	return nil
}
