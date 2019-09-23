package plugin

import "time"

type Metric interface {
	Namespace() []NamespaceElement
	Data() interface{}
	Tags() map[string]string
	Timestamp() time.Time

	HasTagWithKey(key string)
	HasTagWithValue(key string)
	HasTag(key string, value string)

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
