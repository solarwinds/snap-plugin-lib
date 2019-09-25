package plugin

import "time"

// Representation of AppOptics measurement
type Metric interface {
	// Name of metric, ie: [system, cpu, percentage]
	Namespace() []NamespaceElement

	// Value associated with measurement
	Value() interface{}

	// Text-like object associated with measurement
	Tags() map[string]string

	// Time, when measurement was taken
	Timestamp() time.Time

	// True, when metric contains tag with specific key
	HasTagWithKey(key string) bool

	// True, when metric contains tag with specific value
	HasTagWithValue(key string) bool

	// True, when metric contains specific
	HasTag(key string, value string) bool

	// True, when metric name contains given element
	HasNsElement(el string) bool

	// True, when metric name contains element on a given position
	HasNsElementOn(el string, pos int) bool

	// Name of metric, ie: /system/cpu/percentage
	NamespaceText() string
}

// Representation of part of AppOptics measurement name
type NamespaceElement interface {
	// Name of element (not empty in case element is dynamic)
	Name() string

	// Value of element (not empty in case element is dynamic)
	Value() string

	// Description associated with element (not empty in case element is dynamic)
	Description() string

	// True, if element is dynamic
	IsDynamic() bool
}
