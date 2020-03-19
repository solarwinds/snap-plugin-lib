package mock

import "github.com/stretchr/testify/mock"

type Definition struct {
	mock.Mock
}

func (m *Definition) DefineTasksPerInstanceLimit(limit int) error {
	args := m.Called(limit)
	return args.Error(0)
}

func (m *Definition) DefineInstancesLimit(limit int) error {
	args := m.Called(limit)
	return args.Error(0)
}

type CollectorDefinition struct {
	Definition
}

func (m *CollectorDefinition) DefineMetric(namespace string, unit string, isDefault bool, description string) {
	m.Called(namespace, unit, isDefault, description)
}

func (m *CollectorDefinition) DefineGroup(name string, description string) {
	m.Called(name, description)
}

func (m *CollectorDefinition) DefineGlobalTags(namespaceSelector string, tags map[string]string) {
	m.Called(namespaceSelector, tags)
}

func (m *CollectorDefinition) DefineExampleConfig(cfg string) error {
	args := m.Called(cfg)
	return args.Error(0)
}
