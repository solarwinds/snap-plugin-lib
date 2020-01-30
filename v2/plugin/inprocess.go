package plugin

type InProcessCollector interface {
	Collector
	Name() string
	Version() string
}

type InProcessPublisher interface {
	Publisher
	Name() string
	Version() string
}
