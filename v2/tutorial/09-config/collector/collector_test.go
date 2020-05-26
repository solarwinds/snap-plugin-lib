package collector

import (
	"testing"
	"time"

	"github.com/librato/snap-plugin-lib-go/v2/plugin"

	pluginMock "github.com/librato/snap-plugin-lib-go/v2/mock"
	"github.com/librato/snap-plugin-lib-go/v2/tutorial/09-config/collector/data"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

///////////////////////////////////////////////////////////////////////////////

type mockProxy struct {
	mock.Mock
}

func (m *mockProxy) ProcessesInfo() ([]data.ProcessInfo, error) {
	args := m.Called()
	return args.Get(0).([]data.ProcessInfo), args.Error(1)
}

func (m *mockProxy) TotalCpuUsage(d time.Duration) (float64, error) {
	args := m.Called(d)
	return args.Get(0).(float64), args.Error(1)
}

func (m *mockProxy) TotalMemoryUsage() (float64, error) {
	args := m.Called()
	return args.Get(0).(float64), args.Error(1)
}

///////////////////////////////////////////////////////////////////////////////

func TestCollectProcessMetrics(t *testing.T) {
	// Arrange
	proxy := &mockProxy{}
	ctx := &pluginMock.Context{}

	processList := []data.ProcessInfo{
		{ProcessName: "mysql", CpuUsage: 0.3, MemoryUsage: 0.3, PID: 1232},
		{ProcessName: "rabbit", CpuUsage: 0.1, MemoryUsage: 0.2, PID: 4514},
		{ProcessName: "chrome", CpuUsage: 0.5, MemoryUsage: 0.4, PID: 2012},
	}

	pluginConfig := &config{
		Processes: configProcesses{
			MinCPUUsage:    0.3,
			MinMemoryUsage: 0.4,
		},
	}

	ctx.On("Load", "config").
		Once().Return(pluginConfig, true)

	proxy.On("ProcessesInfo").
		Once().Return(processList, nil)

	ctx.On("AddMetric", "/minisystem/processes/[processName=mysql]/cpu", 0.3, []plugin.MetricModifier{plugin.MetricTag("PID", "1232")}).
		Once().Return(nil)

	ctx.On("AddMetric", "/minisystem/processes/[processName=chrome]/cpu", 0.5, []plugin.MetricModifier{plugin.MetricTag("PID", "2012")}).
		Once().Return(nil)

	ctx.On("AddMetric", "/minisystem/processes/[processName=chrome]/memory", 0.4, []plugin.MetricModifier{plugin.MetricTag("PID", "2012")}).
		Once().Return(nil)

	c := systemCollector{
		proxyCollector: proxy,
	}

	// Act
	err := c.collectProcessesInfo(ctx)

	// Assert
	require.Nil(t, err)

	require.True(t, ctx.AssertExpectations(t))
}
