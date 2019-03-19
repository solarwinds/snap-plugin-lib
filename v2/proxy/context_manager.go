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

///////////////////////////////////////////////////////////////////////////////
// proxy.Collector related methods

func (cp *ContextManager) RequestCollect(id int) ([]plugin.Metric, error) {
	cp.collector.Collect(cp)
	return nil, nil
}

func (cp *ContextManager) LoadTask(id int, config string, selectors []plugin.Namespace) error {
	if loadable, ok := cp.collector.(plugin.LoadableCollector); ok {
		loadable.Load(cp)
	}
	return nil
}

func (cp *ContextManager) UnloadTask(id int) error {
	if loadable, ok := cp.collector.(plugin.LoadableCollector); ok {
		loadable.Unload(cp)
	}
	return nil
}

func (cp *ContextManager) RequestInfo() plugin.Info {
	return plugin.Info{}
}

///////////////////////////////////////////////////////////////////////////////
// Context related methods

func (cp *ContextManager) Config(string) (string, bool) {
	panic("implement me")
}

func (cp *ContextManager) ConfigKeys() []string {
	panic("implement me")
}

func (cp *ContextManager) RawConfig() string {
	panic("implement me")
}

func (cp *ContextManager) Store(string, interface{}) {
	panic("implement me")
}

func (cp *ContextManager) Load(string) (interface{}, bool) {
	panic("implement me")
}

func (cp *ContextManager) AddMetric(plugin.Namespace, plugin.MetricValue) error {
	panic("implement me")
}

func (cp *ContextManager) AddMetricWithTags(plugin.Namespace, plugin.MetricValue, plugin.Tags) error {
	panic("implement me")
}

func (cp *ContextManager) ApplyTagsByPath(plugin.Namespace, plugin.Tags) error {
	panic("implement me")
}

func (cp *ContextManager) ApplyTagsByRegExp(plugin.Namespace, plugin.Tags) error {
	panic("implement me")
}
