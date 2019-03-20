package context_manager

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

func (pc *pluginContext) AddMetric(string, interface{}) error {
	panic("implement me")
}

func (pc *pluginContext) AddMetricWithTags(string, interface{}, plugin.Tags) error {
	panic("implement me")
}

func (pc *pluginContext) ApplyTagsByPath(string, plugin.Tags) error {
	panic("implement me")
}

func (pc *pluginContext) ApplyTagsByRegExp(string, plugin.Tags) error {
	panic("implement me")
}
