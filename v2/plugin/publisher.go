package plugin

type Publisher interface {
	Publish(ctx PublishContext) error
}

type PublishContext interface {
	Context

	ListMetrics(MetricFilter) []Metric
	ListAllMetrics() []Metric
	HasMetric(ns string)
	Count() int
}

type MetricFilter func(Metric) bool
