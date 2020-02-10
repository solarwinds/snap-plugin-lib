package plugin

type Publisher interface {
	Publish(ctx PublishContext) error
}

type LoadablePublisher interface {
	Publisher
	Load(ctx Context) error
}

type UnloadablePublisher interface {
	Publisher
	Unload(ctx Context) error
}

type DefinablePublisher interface {
	Publisher
	PluginDefinition(def PublisherDefinition) error
}

type CustomizableInfoPublisher interface {
	Publisher
	CustomInfo(ctx Context) interface{}
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
