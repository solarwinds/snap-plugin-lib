package types

import (
	"fmt"
	"strings"
)

type Namespace []NamespaceElement

func (ns Namespace) HasElement(el string) bool {
	for _, nsElem := range ns {
		if el == nsElem.String() {
			return true
		}
	}

	return false
}

func (ns Namespace) HasElementOn(el string, pos int) bool {
	if pos < len(ns) && pos >= 0 {
		if el == ns[pos].String() {
			return true
		}
	}

	return false
}

func (ns Namespace) String() string {
	var sb strings.Builder

	sb.WriteString(metricSeparator)

	for i, nsElem := range ns {
		sb.WriteString(nsElem.String())

		if i != len(ns)-1 {
			sb.WriteString(metricSeparator)
		}
	}

	return sb.String()
}

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
