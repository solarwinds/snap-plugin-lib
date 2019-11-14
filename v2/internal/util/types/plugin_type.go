package types

type PluginType int

const (
	PluginTypeCollector PluginType = iota
	PluginTypeProcessor
	PluginTypePublisher
	PluginTypeStreamingCollector
)

func (pt PluginType) String() string {
	switch pt {
	case PluginTypeCollector:
		return "Collector"
	case PluginTypeProcessor:
		return "Publisher"
	case PluginTypePublisher:
		return "Publisher"
	case PluginTypeStreamingCollector:
		return "Streaming Collector"
	default:
		return "Unknown plugin type"
	}
}
