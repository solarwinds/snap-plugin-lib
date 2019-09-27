package collector

import (
	"github.com/librato/snap-plugin-lib-go/tutorial/08-collector/collector/data"
	"github.com/stretchr/testify/mock"
	"testing"
)

///////////////////////////////////////////////////////////////////////////////

type mockContext struct {
}

func (m *mockContext) Config(string) (string, bool) {
	panic("implement me")
}

func (m *mockContext) ConfigKeys() []string {
	panic("implement me")
}

func (m *mockContext) RawConfig() []byte {
	panic("implement me")
}

func (m *mockContext) Store(string, interface{}) {
	panic("implement me")
}

func (m *mockContext) Load(string) (interface{}, bool) {
	panic("implement me")
}

func (m *mockContext) AddMetric(string, interface{}) error {
	panic("implement me")
}

func (m *mockContext) AddMetricWithTags(string, interface{}, map[string]string) error {
	panic("implement me")
}

func (m *mockContext) ApplyTagsByPath(string, map[string]string) error {
	panic("implement me")
}

func (m *mockContext) ApplyTagsByRegExp(string, map[string]string) error {
	panic("implement me")
}

func (m *mockContext) ShouldProcess(string) bool {
	panic("implement me")
}

///////////////////////////////////////////////////////////////////////////////

type mockProxy struct {
	mock.Mock
}

func (m *mockProxy) ProcessesInfo() ([]data.ProcessInfo, error) {
	panic("implement me")
}

func (m *mockProxy) TotalCpuUsage() (float64, error) {
	panic("implement me")
}

func (m *mockProxy) TotalMemoryUsage() (float64, error) {
	panic("implement me")
}

///////////////////////////////////////////////////////////////////////////////

func TestCollectTotalCpu(t *testing.T) {
	c := systemCollector{
		proxyCollector: &mockProxy{},
	}

	_ = c.collectTotalCPU()
}
