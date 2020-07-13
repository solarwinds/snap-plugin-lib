package plugin

// CollectContext provides state and configuration API to be used by custom code.
type Context interface {
	// Returns configuration value by providing path (representing its position in JSON tree)
	Config(key string) (string, bool)

	// Returns list of allowed configuration paths
	ConfigKeys() []string

	// Return raw configuration (JSON string)
	RawConfig() []byte

	// Store any object using key to have access from different Collect requests
	Store(key string, value interface{})

	// Load any object using key from different Collect requests (returns an interface{} which need to be casted to concrete type)
	Load(key string) (interface{}, bool)

	// Load any object using key from different Collect requests (passing it to provided reference).
	// Will throw error when dest type doesn't match to type of stored value or object with a given key wasn't found.
	LoadTo(key string, dest interface{}) error

	// Add warning information to current collect / process operation.
	AddWarning(msg string)

	// Check if task is completed
	IsDone() bool

	// Check if task is completed (via listening on a channel)
	Done() <-chan struct{}
}
