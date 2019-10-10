package collector

import (
	"testing"

	pluginMock "github.com/librato/snap-plugin-lib-go/v2/mock"
	"github.com/stretchr/testify/require"
)

func TestConfig_AllFields(t *testing.T) {
	// Arrange
	ctx := &pluginMock.Context{}

	expectedCfg := &config{
		Processes: configProcesses{
			MinCPUUsage:    0.45,
			MinMemoryUsage: 0.15,
		},
		TotalCPUMeasureDuration: "5s",
	}

	ctx.On("RawConfig").
		Once().Return([]byte(`{"processes": {"minCpuUsage": 0.45, "minMemoryUsage": 0.15}, "totalCpuMeasureDuration": "5s"}`))

	ctx.On("Store", "config", expectedCfg).
		Once().Return()

	// Act
	err := handleConfig(ctx)

	// Assert
	require.Nil(t, err)
}

func TestConfig_EmptyJSON(t *testing.T) {
	// Arrange
	ctx := &pluginMock.Context{}

	expectedCfg := defaultConfig()

	ctx.On("RawConfig").
		Once().Return([]byte(`{}`))

	ctx.On("Store", "config", &expectedCfg).
		Once().Return()

	// Act
	err := handleConfig(ctx)

	// Assert
	require.Nil(t, err)
}

func TestConfig_ValidationError(t *testing.T) {
	// Arrange
	ctx := &pluginMock.Context{}

	expectedCfg := defaultConfig()

	ctx.On("RawConfig").
		Once().Return([]byte(`{"processes": {"minCpuUsage": 110, "minMemoryUsage": 15}, "totalCpuMeasureDuration": "5s"}`))

	ctx.On("Store", "config", expectedCfg).
		Once().Return()

	// Act
	err := handleConfig(ctx)

	// Assert
	require.Error(t, err)
}
