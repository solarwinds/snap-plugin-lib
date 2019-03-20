package plugin

// Type representing metric tags (additional information associated with metric)
type Tags map[string]string

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
	AddMetric(string, interface{}) error

	// Add concrete metric with calculated value and tags
	AddMetricWithTags(string, interface{}, Tags) error

	// Add tags to specific metric
	ApplyTagsByPath(string, Tags) error

	// Add tags to all metrics matching regular expression
	ApplyTagsByRegExp(string, Tags) error
}

// CollectorDefinition provides API for specifying plugin metadata (supported metrics, descriptions etc)
type CollectorDefinition interface {
	// Define supported metric, its description and indication if metric is default
	DefineMetric(string, bool, string)

	// Define description for dynamic element
	DefineGroup(string, string)

	// Define global tags that will be applied to all metrics
	DefineGlobalTags(string, Tags)
}
