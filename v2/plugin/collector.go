/*
The package "plugin" provides interfaces to define custom plugins and Context interface
which allows to perform any collection-related operation.
*/
package plugin

type Collector interface {
	Collect(ctx CollectContext) error
}

type LoadableCollector interface {
	Load(Context) error
	Unload(Context) error
}

type DefinableCollector interface {
	PluginDefinition(CollectorDefinition) error
}

///////////////////////////////////////////////////////////////////////////////

// CollectContext provides metric, state and configuration API to be used by custom code.
type CollectContext interface {
	Context

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

///////////////////////////////////////////////////////////////////////////////

// CollectorDefinition provides API for specifying plugin metadata (supported metrics, descriptions etc)
type CollectorDefinition interface {
	// Define supported metric, its description and indication if metric is default
	DefineMetric(string, string, bool, string)

	// Define description for dynamic element
	DefineGroup(string, string)

	// Define global tags that will be applied to all metrics
	DefineGlobalTags(string, map[string]string)

	// Define example config (which will be presented when example task is printed)
	DefineExampleConfig(cfg string) error
}
