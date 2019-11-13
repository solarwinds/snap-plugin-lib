package types

type PluginType int

const (
	PluginTypeCollector PluginType = iota
	PluginTypeProcessor
	PluginTypePublisher
	PluginTypeStreamingCollector
)
