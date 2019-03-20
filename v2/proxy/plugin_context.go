package proxy

import "github.com/librato/snap-plugin-lib-go/v2/plugin"

type pluginContext struct {
}

func (pc *pluginContext) Config(string) (string, bool) {
	panic("implement me")
}

func (pc *pluginContext) ConfigKeys() []string {
	panic("implement me")
}

func (pc *pluginContext) RawConfig() string {
	panic("implement me")
}

func (pc *pluginContext) Store(string, interface{}) {
	panic("implement me")
}

func (pc *pluginContext) Load(string) (interface{}, bool) {
	panic("implement me")
}

func (pc *pluginContext) AddMetric(plugin.Namespace, plugin.MetricValue) error {
	panic("implement me")
}

func (pc *pluginContext) AddMetricWithTags(plugin.Namespace, plugin.MetricValue, plugin.Tags) error {
	panic("implement me")
}

func (pc *pluginContext) ApplyTagsByPath(plugin.Namespace, plugin.Tags) error {
	panic("implement me")
}

func (pc *pluginContext) ApplyTagsByRegExp(plugin.Namespace, plugin.Tags) error {
	panic("implement me")
}
