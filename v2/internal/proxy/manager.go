/*
Package proxy:
1) Manages context for different task created for the same plugin
2) Serves as an entry point for any "controller" (like. Rpc)
 */
package proxy

import (
	"errors"

	"github.com/librato/snap-plugin-lib-go/v2/plugin"
)

type Collector interface {
	RequestCollect(id int) ([]plugin.Metric, error)
	LoadTask(id int, config string, selectors []string) error
	UnloadTask(id int) error
	RequestInfo()
}

type ContextManager struct {
	collector  plugin.Collector
	contextMap map[int]*pluginContext
}

func NewContextManager(collector plugin.Collector, pluginName string, version string) Collector {
	return &ContextManager{
		collector:  collector,
		contextMap: map[int]*pluginContext{},
	}
}

///////////////////////////////////////////////////////////////////////////////
// proxy.Collector related methods

func (cm *ContextManager) RequestCollect(id int) ([]plugin.Metric, error) {
	if context, ok := cm.contextMap[id]; ok {
		cm.collector.Collect(context)
	}

	return nil, nil
}

func (cm *ContextManager) LoadTask(id int, config string, selectors []string) error {
	if _, ok := cm.contextMap[id]; ok {
		return errors.New("context with given id was already defined")
	}

	cm.contextMap[id] = &pluginContext{}

	if loadable, ok := cm.collector.(LoadableCollector); ok {
		loadable.Load(cm.contextMap[id])
	}

	return nil
}

func (cm *ContextManager) UnloadTask(id int) error {
	if _, ok := cm.contextMap[id]; !ok {
		return errors.New("context with given id is not defined")
	}

	if loadable, ok := cm.collector.(LoadableCollector); ok {
		loadable.Unload(cm.contextMap[id])
	}

	delete(cm.contextMap, id)
	return nil
}

func (cm *ContextManager) RequestInfo() {
	return
}
