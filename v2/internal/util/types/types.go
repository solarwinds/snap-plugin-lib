package types

import "time"

type NamespaceElement struct {
	Name        string
	Value       string
	Description string
}

type Metric struct {
	Namespace   []NamespaceElement
	Value       interface{}
	Tags        map[string]string
	Timestamp   time.Time
	Description string
}
