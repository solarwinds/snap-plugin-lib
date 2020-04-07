package types

import (
	"fmt"
	"time"

	"github.com/librato/snap-plugin-lib-go/v2/plugin"
)

const (
	metricSeparator = "/"
)

type Metric struct {
	Namespace_   []NamespaceElement
	Value_       interface{}
	Tags_        map[string]string
	Unit_        string
	Timestamp_   time.Time
	Description_ string
}

func (m Metric) Namespace() plugin.Namespace {
	ns := make(Namespace, 0, len(m.Namespace_))

	for i := range m.Namespace_ {
		ns = append(ns, m.Namespace_[i])
	}

	return ns
}

func (m Metric) Value() interface{} {
	return m.Value_
}

func (m Metric) Tags() map[string]string {
	return m.Tags_
}

func (m Metric) Unit() string {
	return m.Unit_
}

func (m Metric) Description() string {
	return m.Description_
}

func (m Metric) Timestamp() time.Time {
	return m.Timestamp_
}

func (m Metric) String() string {
	return fmt.Sprintf("%s %v {%v}", m.Namespace().String(), m.Value_, m.Tags_)
}

func (m *Metric) AddTags(tags map[string]string) {
	if m.Tags_ == nil { // lazy initialization
		m.Tags_ = map[string]string{}
	}

	for k, v := range tags {
		m.Tags_[k] = v
	}
}

func (m *Metric) SetDescription(description string) {
	m.Description_ = description
}

func (m *Metric) SetUnit(unit string) {
	m.Unit_ = unit
}

func (m *Metric) SetTimestamp(timestamp time.Time) {
	m.Timestamp_ = timestamp
}
