package plugin

type InProcessPlugin interface {
	Name() string
	Version() string
}

type InProcessCollector interface {
	Collector
	InProcessPlugin
}

type InProcessStreamingCollector interface {
	StreamingCollector
	InProcessPlugin
}

type InProcessPublisher interface {
	Publisher
	InProcessPlugin
}
