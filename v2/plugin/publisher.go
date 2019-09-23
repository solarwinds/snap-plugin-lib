package plugin

type Publisher interface {
	Publish(ctx PublishContext) error
}

type PublishContext interface {
	Context

	ListMetrics()
	ListAllMetrics()
	HasMetric(ns string)
	Count() int
}

type MetricFilter func(Metric) bool
