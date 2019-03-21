package proxy

import "github.com/librato/snap-plugin-lib-go/v2/plugin"

type LoadableCollector interface {
	Load(plugin.Context) error
	Unload(plugin.Context) error
}

type DefinableCollector interface {
	DefineMetrics(plugin.CollectorDefinition) error
}
