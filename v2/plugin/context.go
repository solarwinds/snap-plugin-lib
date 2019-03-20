package plugin

// Context provides metric and configuration API to be used by custom code.
type Context interface {
	// Returns configuration value by providing path (representing its position in JSON tree)
	Config(string) (string, bool)

	// Returns list of allowed configuration paths
	ConfigKeys() []string

	// Return raw configuration (JSON string)
	RawConfig() string

	// Store any object between Collect requests using key
	Store(string, interface{})

	// Load stored object between Collect requests using key
	Load(string) (interface{}, bool)

	// Add concrete metric with calculated value
	AddMetric(Namespace, MetricValue) error

	// Add concrete metric with calculated value and tags
	AddMetricWithTags(Namespace, MetricValue, Tags) error

	// Add tags to specific metric
	ApplyTagsByPath(Namespace, Tags) error

	// Add tags to all metrics matching regular expression
	ApplyTagsByRegExp(Namespace, Tags) error
}

type CollectorContext interface {
	DefineMetric(Namespace, bool, string)
	DefineGroup(string, string)
	DefineGlobalTags(Namespace, Tags)
}
