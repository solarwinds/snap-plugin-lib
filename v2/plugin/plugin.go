package plugin

import "time"

///////////////////////////////////////////////////////////////////////////////

type Collector interface {
	Collect(ctx Context) error
}

type LoadableCollector interface {
	Load(Context) error
	Unload(Context) error
}

type DefinableCollector interface {
	DefineMetrics(CollectorContext) error
}

///////////////////////////////////////////////////////////////////////////////

type Tags map[string]string

type Namespace string

type MetricValue interface{}

type Metric struct {
	Namespace Namespace
	Value     MetricValue
	Tags      Tags
	Timestamp time.Time
}
