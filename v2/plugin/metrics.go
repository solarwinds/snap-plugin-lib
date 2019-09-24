package plugin

import "time"

// todo: document
type Metric interface {
	Namespace() []NamespaceElement
	Value() interface{}
	Tags() map[string]string
	Timestamp() time.Time

	HasTagWithKey(key string) bool
	HasTagWithValue(key string) bool
	HasTag(key string, value string) bool

	HasNsElement(el string) bool
	HasNsElementOn(el string, pos int) bool

	NamespaceText() string
}

type NamespaceElement interface {
	Name() string
	Value() string
	Description() string

	IsDynamic() bool
}
