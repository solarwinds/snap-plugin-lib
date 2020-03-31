package plugin

import "time"

// Representation of AppOptics measurement
type Metric interface {
	// Name of metric, ie: [system, cpu, percentage]
	Namespace() Namespace

	// Value associated with measurement
	Value() interface{}

	// Text-like object associated with measurement
	Tags() map[string]string

	// Description of measurement
	Description() string

	// Unit of measurement value
	Unit() string

	// Time, when measurement was taken
	Timestamp() time.Time
}

// Interface for setting custom metric metadata
type SettableMetric interface {
	// Add custom text-like object associated with measurement
	AddTags(map[string]string)

	// Set custom Description of measurement
	SetDescription(string)

	// Set custom unit of measurement value
	SetUnit(string)

	// Set custom timestamp
	SetTimestamp(time.Time)
}

// Representation of AppOptics measurement name
type Namespace interface {
	// Return namespace element at the given position
	At(pos int) NamespaceElement

	// Return length of the element
	Len() int

	// True, when metric name contains given element
	HasElement(el string) bool

	// True, when metric name contains element on a given position
	HasElementOn(el string, pos int) bool

	// Name of metric, ie: /system/cpu/percentage
	String() string
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
