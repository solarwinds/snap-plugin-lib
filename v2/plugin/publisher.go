package plugin

type Publisher interface {
	Publish(ctx PublishContext) error
}

type LoadablePublisher interface {
	Load(ctx Context) error
}

type UnloadablePublisher interface {
	Unload(ctx Context) error
}

type DefinablePublisher interface {
	PluginDefinition(def PublisherDefinition) error
}

type PublishContext interface {
	Context

	ListAllMetrics() []Metric
	Count() int
}

// PublisherDefinition provides API for specifying plugin (publisher) metadata (supported metrics, descriptions etc)
type PublisherDefinition interface {
	Definition
}
