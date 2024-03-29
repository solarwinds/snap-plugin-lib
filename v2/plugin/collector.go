/*
The package "plugin" provides interfaces to define custom plugins and Context interface
which allows to perform any collection-related operation.
*/

/*
 Copyright (c) 2021 SolarWinds Worldwide, LLC

    Licensed under the Apache License, Version 2.0 (the "License");
    you may not use this file except in compliance with the License.
    You may obtain a copy of the License at

      http://www.apache.org/licenses/LICENSE-2.0

    Unless required by applicable law or agreed to in writing, software
    distributed under the License is distributed on an "AS IS" BASIS,
    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
    See the License for the specific language governing permissions and
    limitations under the License.
*/

package plugin

type Collector interface {
	Collect(ctx CollectContext) error
}

type StreamingCollector interface {
	StreamingCollect(ctx CollectContext) error
}

type LoadableCollector interface {
	Load(ctx Context) error
}

type UnloadableCollector interface {
	Unload(ctx Context) error
}

type DefinableCollector interface {
	PluginDefinition(def CollectorDefinition) error
}

type CustomizableInfoCollector interface {
	CustomInfo(ctx Context) interface{}
}

///////////////////////////////////////////////////////////////////////////////

// CollectContext provides metric, state and configuration API to be used by custom code.
type CollectContext interface {
	Context

	// Add concrete metric with calculated value
	AddMetric(namespace string, value interface{}, modifier ...MetricModifier) error

	// Always apply specific modifier(s) for a metrics matching namespace selector
	// Returns object which may be used to dismiss modifiers (make them no-active)
	AlwaysApply(namespaceSelector string, modifier ...MetricModifier) (Dismisser, error)

	// Dismisses all modifiers created by calling AlwaysApply
	DismissAllModifiers()

	// Provide information whether metric or metric group is reasonable to process (won't be filtered).
	ShouldProcess(namespace string) bool

	// List all requested metrics (filter).
	// WARNING: library automatically filters metrics based on provided list. You should use this function
	// in scenarios when output metrics namespaces are constructed based on input list (ie. snmp metrics based on OIDs)
	RequestedMetrics() []string
}

///////////////////////////////////////////////////////////////////////////////

// CollectorDefinition provides API for specifying plugin metadata (supported metrics, descriptions etc)
type CollectorDefinition interface {
	Definition

	// Define supported metric, its description and indication if metric is default
	DefineMetric(namespace string, unit string, isDefault bool, description string)

	// Define description for dynamic element
	DefineGroup(name string, description string)

	// Define example config (which will be presented when example task is printed)
	DefineExampleConfig(cfg string) error

	// Allow submitting metrics with namespace not being explicitly defined earlier
	// The only requirement here is that metrics should have matching root namespace element
	// This allows implementing DefineMetric/DefineGroup thus having dynamic metrics but
	// allows some a priori unknown metrics namespaces at the same
	AllowUndefinedMetrics()

	// Allow metrics values not only on leaves but at any namespace level
	AllowValuesAtAnyNamespaceLevel()

	// Set global namespace prefix for all metrics
	// This is used by "swi" collector (plugins bundle) for metrics have different root names
	// - prefix should contain separator, ie "/swi", "#swi"
	// - removePrefixFromOutput set to true removes global prefix from all namespaces when result is to be sent to agent.
	// !! Should be set before any call to DefineMetric
	SetGlobalMetricPrefix(prefix string, removePrefixFromOutput bool) error
}
