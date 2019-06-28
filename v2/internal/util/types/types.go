package types

import (
	"fmt"
	"strings"
	"time"
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
	nsStr := []string{}
	for _, ns := range m.Namespace {
		nsStr = append(nsStr, fmt.Sprintf("%s", ns.String()))
	}

	return fmt.Sprintf("%s %v {tags: %v}", strings.Join(nsStr, "."), m.Value, m.Tags)
}
