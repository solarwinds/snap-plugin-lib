package types

import (
	"fmt"
	"strings"
	"time"
)

const (
	metricSeparator = "."
)

type NamespaceElement struct {
	Name        string
	Value       string
	Description string
}

type Metric struct {
	Namespace   []NamespaceElement
	Value       interface{}
	Tags        map[string]string
	Unit        string
	Timestamp   time.Time
	Description string
}

func (ns *NamespaceElement) String() string {
	if ns.Name == "" {
		return ns.Value
	}

	return fmt.Sprintf("[%s=%s]", ns.Name, ns.Value)
}

func (m *Metric) String() string {
	var sb strings.Builder

	for i, ns := range m.Namespace {
		sb.WriteString(ns.String())

		if i != len(m.Namespace)-1 {
			sb.WriteString(metricSeparator)
		}
	}

	return fmt.Sprintf("%s %v {tags: %v}", sb.String(), m.Value, m.Tags)
}
