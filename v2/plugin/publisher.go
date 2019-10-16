package plugin

type Publisher interface {
	Publish(ctx PublishContext) error
}

type LoadablePublisher interface {
	Load(ctx Context) error
	Unload(ctx Context) error
}

type PublishContext interface {
	Context

	ListAllMetrics() []Metric
	Count() int
}
