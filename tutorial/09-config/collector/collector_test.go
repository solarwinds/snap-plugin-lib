package collector

import (
	"github.com/librato/snap-plugin-lib-go/tutorial/09-config/collector/data"
	"github.com/smartystreets/assertions"
	"github.com/stretchr/testify/mock"
	"testing"
)

///////////////////////////////////////////////////////////////////////////////

type mockContext struct {
	mock.Mock
}

func (m *mockContext) Config(key string) (string, bool) {
	args := m.Called(key)
	return args.String(0), args.Bool(1)
}

func (m *mockContext) ConfigKeys() []string {
	args := m.Called()
	return args.Get(0).([]string)
}

func (m *mockContext) RawConfig() []byte {
	args := m.Called()
	return args.Get(0).([]byte)
}

func (m *mockContext) Store(key string, value interface{}) {
	m.Called(key, value)
}

func (m *mockContext) Load(key string) (interface{}, bool) {
	args := m.Called(key)
	return args.Get(0), args.Bool(1)
}

func (m *mockContext) AddMetric(ns string, value interface{}) error {
	args := m.Called(ns, value)
	return args.Error(0)
}

func (m *mockContext) AddMetricWithTags(ns string, value interface{}, tags map[string]string) error {
	args := m.Called(ns, value, tags)
	return args.Error(0)
}

func (m *mockContext) ApplyTagsByPath(ns string, tags map[string]string) error {
	args := m.Called(ns, tags)
	return args.Error(0)
}

func (m *mockContext) ApplyTagsByRegExp(ns string, tags map[string]string) error {
	args := m.Called(ns, tags)
	return args.Error(0)
}

func (m *mockContext) ShouldProcess(ns string) bool {
	args := m.Called(ns)
	return args.Bool(0)
}

///////////////////////////////////////////////////////////////////////////////

type mockProxy struct {
	mock.Mock
}

func (m *mockProxy) ProcessesInfo() ([]data.ProcessInfo, error) {
	args := m.Called()
	return args.Get(0).([]data.ProcessInfo), args.Error(1)
}

func (m *mockProxy) TotalCpuUsage() (float64, error) {
	args := m.Called()
	return args.Get(0).(float64), args.Error(1)
}

func (m *mockProxy) TotalMemoryUsage() (float64, error) {
	args := m.Called()
	return args.Get(0).(float64), args.Error(1)
}

///////////////////////////////////////////////////////////////////////////////

func TestCollectTotalMemory(t *testing.T) {
	proxy := &mockProxy{}
	ctx := &mockContext{}

	proxy.On("TotalMemoryUsage").
		Return(15.0, nil).Once()

	ctx.On("AddMetric", mock.Anything, mock.Anything).
		Return(nil).Once()

	c := systemCollector{
		proxyCollector: proxy,
	}

	err := c.collectTotalMemory(ctx)
	assertions.ShouldBeNil(err)
}
