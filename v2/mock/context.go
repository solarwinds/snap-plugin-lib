package mock

import (
	"context"
	"github.com/stretchr/testify/mock"
)

type Context struct {
	mock.Mock
}

func (m *Context) Config(key string) (string, bool) {
	args := m.Called(key)
	return args.String(0), args.Bool(1)
}

func (m *Context) ConfigKeys() []string {
	args := m.Called()
	return args.Get(0).([]string)
}

func (m *Context) RawConfig() []byte {
	args := m.Called()
	return args.Get(0).([]byte)
}

func (m *Context) Store(key string, value interface{}) {
	m.Called(key, value)
}

func (m *Context) Load(key string) (interface{}, bool) {
	args := m.Called(key)
	return args.Get(0), args.Bool(1)
}

func (m *Context) LoadTo(key string, dest interface{}) error {
	args := m.Called(key, dest)
	return args.Error(0)
}

func (m *Context) AddWarning(msg string) {
	m.Called(msg)
}

func (m *Context) AddMetric(ns string, value interface{}) error {
	args := m.Called(ns, value)
	return args.Error(0)
}

func (m *Context) AddMetricWithTags(ns string, value interface{}, tags map[string]string) error {
	args := m.Called(ns, value, tags)
	return args.Error(0)
}

func (m *Context) ApplyTagsByPath(ns string, tags map[string]string) error {
	args := m.Called(ns, tags)
	return args.Error(0)
}

func (m *Context) ApplyTagsByRegExp(ns string, tags map[string]string) error {
	args := m.Called(ns, tags)
	return args.Error(0)
}

func (m *Context) ShouldProcess(ns string) bool {
	args := m.Called(ns)
	return args.Bool(0)
}

func (m *Context) AttachContext(parentCtx context.Context) {
	m.Called(parentCtx)
	return
}

func (m *Context) IsDone() bool {
	args := m.Called()
	return args.Bool(0)
}

func (m *Context) Done() <-chan struct{} {
	args := m.Called()
	return args.Get(0).(<-chan struct{})
}
