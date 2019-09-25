package types

import (
	"fmt"
	"strings"
	"time"

	"github.com/librato/snap-plugin-lib-go/v2/plugin"
)

const (
	metricSeparator = "/"
)

///////////////////////////////////////////////////////////////////////////////

type NamespaceElement struct {
	Name_        string
	Value_       string
	Description_ string
}

func (ns *NamespaceElement) Name() string {
	return ns.Name_
}

func (ns *NamespaceElement) Value() string {
	return ns.Value_
}

func (ns *NamespaceElement) Description() string {
	return ns.Description_
}

func (ns *NamespaceElement) IsDynamic() bool {
	return ns.Name_ != ""
}

func (ns *NamespaceElement) String() string {
	if ns.Name_ == "" {
		return ns.Value_
	}

	return fmt.Sprintf("[%s=%s]", ns.Name_, ns.Value_)
}

///////////////////////////////////////////////////////////////////////////////

type Metric struct {
	Namespace_   []NamespaceElement
	Value_       interface{}
	Tags_        map[string]string
	Unit_        string
	Timestamp_   time.Time
	Description_ string
}

///////////////////////////////////////////////////////////////////////////////

func (m *Metric) Namespace() []plugin.NamespaceElement {
	ns := make([]plugin.NamespaceElement, 0, len(m.Namespace_))

	for _, nsElem := range m.Namespace_ {
		ns = append(ns, &nsElem)
	}

	return ns
}

func (m *Metric) Value() interface{} {
	return m.Value_
}

func (m *Metric) Tags() map[string]string {
	return m.Tags_
}

func (m *Metric) Timestamp() time.Time {
	return m.Timestamp_
}

func (m *Metric) HasTagWithKey(key string) bool {
	_, ok := m.Tags_[key]
	return ok
}

func (m *Metric) HasTagWithValue(value string) bool {
	for _, v := range m.Tags_ {
		if v == value {
			return true
		}
	}

	return false
}

func (m *Metric) HasTag(key string, value string) bool {
	for k, v := range m.Tags_ {
		if k == key && v == value {
			return true
		}
	}

	return false
}

func (m *Metric) HasNsElement(el string) bool {
	for _, nsElem := range m.Namespace_ {
		if el == nsElem.String() {
			return true
		}
	}

	return false
}

func (m *Metric) HasNsElementOn(el string, pos int) bool {
	if pos < len(m.Namespace_) {
		if el == m.Namespace_[pos].String() {
			return true
		}
	}

	return false
}

func (m *Metric) NamespaceText() string {
	var sb strings.Builder

	sb.WriteString(metricSeparator)

	for i, ns := range m.Namespace_ {
		sb.WriteString(ns.String())

		if i != len(m.Namespace_)-1 {
			sb.WriteString(metricSeparator)
		}
	}

	return sb.String()
}

func (m *Metric) String() string {
	return fmt.Sprintf("%s %v {%v}", m.NamespaceText(), m.Value_, m.Tags_)
}
