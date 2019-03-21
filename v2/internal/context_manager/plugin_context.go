package context_manager

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

func (pc *pluginContext) AddMetricWithTags(string, interface{}, map[string]string) error {
	panic("implement me")
}

func (pc *pluginContext) ApplyTagsByPath(string, map[string]string) error {
	panic("implement me")
}

func (pc *pluginContext) ApplyTagsByRegExp(string, map[string]string) error {
	panic("implement me")
}
