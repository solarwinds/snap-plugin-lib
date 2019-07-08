package plugin

// Context provides metric and configuration API to be used by custom code.
type Context interface {
	// Returns configuration value by providing path (representing its position in JSON tree)
	Config(string) (string, bool)

	// Returns list of allowed configuration paths
	ConfigKeys() []string

	// Return raw configuration (JSON string)
	RawConfig() []byte

	// Store any object between Collect requests using key
	Store(string, interface{})

	// Load stored object between Collect requests using key
	Load(string) (interface{}, bool)

	// Add concrete metric with calculated value
	AddMetric(string, interface{}) error

	// Add concrete metric with calculated value and tags
	AddMetricWithTags(string, interface{}, map[string]string) error

	// Add tags to specific metric
	ApplyTagsByPath(string, map[string]string) error

	// Add tags to all metrics matching regular expression
	ApplyTagsByRegExp(string, map[string]string) error

	// Provide information whether metric or metric group is reasonable to process (won't be filtered).
	ShouldProcess(string) bool
}

// CollectorDefinition provides API for specifying plugin metadata (supported metrics, descriptions etc)
type CollectorDefinition interface {
	// Define supported metric, its description and indication if metric is default
	DefineMetric(string, string, bool, string)

	// Define description for dynamic element
	DefineGroup(string, string)

	// Define global tags that will be applied to all metrics
	DefineGlobalTags(string, map[string]string)
}
