package proxy

import (
	"github.com/librato/snap-plugin-lib-go/v2/plugin/types"
)

type ContextManager struct {
	collector types.Collector
}

func NewContextManager(collector types.Collector, pluginName string, version string) Collector {
	return &ContextManager{
		collector: collector,
	}
}

func (cp *ContextManager) RequestCollect(id int) []types.Metric {
	return nil
}

func (cp *ContextManager) LoadTask(id int, config string, selectors []types.Namespace) error {
	return nil
}

func (cp *ContextManager) UnloadTask(id int) error {
	return nil
}

func (cp *ContextManager) RequestInfo() types.Info {
	return types.Info{}
}
