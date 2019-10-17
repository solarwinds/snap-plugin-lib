/*
The package "plugin" provides interfaces to define custom plugins and Context interface
which allows to perform any collection-related operation.
*/
package plugin

type Collector interface {
	Collect(ctx CollectContext) error
}

type LoadableCollector interface {
	Load(ctx Context) error
	Unload(ctx Context) error
}

type DefinableCollector interface {
	PluginDefinition(def CollectorDefinition) error
}

///////////////////////////////////////////////////////////////////////////////

// CollectContext provides metric, state and configuration API to be used by custom code.
type CollectContext interface {
	Context

	// Add concrete metric with calculated value
	AddMetric(namespace string, value interface{}) error

	// Add concrete metric with calculated value and tags
	AddMetricWithTags(namespace string, value interface{}, tags map[string]string) error

	// Add tags to specific metric
	ApplyTagsByPath(namespace string, tags map[string]string) error

	// Add tags to all metrics matching regular expression
	ApplyTagsByRegExp(namespaceSelector string, tags map[string]string) error

	// Provide information whether metric or metric group is reasonable to process (won't be filtered).
	ShouldProcess(namespace string) bool
}

///////////////////////////////////////////////////////////////////////////////

// CollectorDefinition provides API for specifying plugin metadata (supported metrics, descriptions etc)
type CollectorDefinition interface {
	Definition

	// Define supported metric, its description and indication if metric is default
	DefineMetric(namespace string, unit string, isDefault bool, description string)

	// Define description for dynamic element
	DefineGroup(name string, description string)

	// Define global tags that will be applied to all metrics
	DefineGlobalTags(namespaceSelector string, tags map[string]string)

	// Define example config (which will be presented when example task is printed)
	DefineExampleConfig(cfg string) error
}
