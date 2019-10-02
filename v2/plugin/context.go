package plugin

// CollectContext provides state and configuration API to be used by custom code.
type Context interface {
	// Returns configuration value by providing path (representing its position in JSON tree)
	Config(string) (string, bool)

	// Returns list of allowed configuration paths
	ConfigKeys() []string

	// Return raw configuration (JSON string)
	RawConfig() []byte

	// Store any object between Collect requests using key
	Store(string, interface{})

	// Load stored object between Collect requests using key
	Load(string) (interface{}, bool)
}
