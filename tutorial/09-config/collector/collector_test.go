package collector

import (
	"testing"

	"github.com/librato/snap-plugin-lib-go/tutorial/09-config/collector/data"
	pluginMock "github.com/librato/snap-plugin-lib-go/v2/mock"
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
	// Arrange
	proxy := &mockProxy{}
	ctx := &pluginMock.Context{}

	proxy.On("TotalMemoryUsage").
		Return(15.0, nil).Once()

	ctx.On("AddMetric", mock.Anything, mock.Anything).
		Return(nil).Once()

	c := systemCollector{
		proxyCollector: proxy,
	}

	// Act
	err := c.collectTotalMemory(ctx)

	// Assert
	require.Nil(t, err)
}
