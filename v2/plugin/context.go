package plugin

type Context interface {
	Config(string) (string, bool)
	ConfigKeys() []string
	RawConfig() string

	Store(string, interface{})
	Load(string) (interface{}, bool)

	AddMetric(Namespace, MetricValue) error
	AddMetricWithTags(Namespace, MetricValue, Tags) error

	ApplyTagsByPath(Namespace, Tags) error
	ApplyTagsByRegExp(Namespace, Tags) error
}

type CollectorContext interface {
	DefineMetric(Namespace, bool, string)
	DefineGroup(string, string)
	DefineGlobalTags(Namespace, Tags)
}
