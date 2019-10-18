package plugin

// Common plugin definition interface
type Definition interface {
	// Define maximum number of tasks that one instance of plugin should handle
	DefineTasksPerInstanceLimit(limit int) error

	// Define maximum number of instances
	DefineInstancesLimit(limit int) error
}
