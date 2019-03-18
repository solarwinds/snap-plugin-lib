package proxy

import (
	"github.com/librato/snap-plugin-lib-go/v2/plugin"
)

type ContextManager struct {
	collector plugin.Collector
}

func NewContextManager(collector plugin.Collector, pluginName string, version string) Collector {
	return &ContextManager{
		collector: collector,
	}
}

func (cp *ContextManager) RequestCollect(id int) []plugin.Metric {
	return nil
}

func (cp *ContextManager) LoadTask(id int, config string, selectors []plugin.Namespace) error {
	return nil
}

func (cp *ContextManager) UnloadTask(id int) error {
	return nil
}

func (cp *ContextManager) RequestInfo() plugin.Info {
	return plugin.Info{}
}
