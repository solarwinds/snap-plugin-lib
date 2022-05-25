//go:build tutorial
// +build tutorial

/*
 Copyright (c) 2021 SolarWinds Worldwide, LLC

    Licensed under the Apache License, Version 2.0 (the "License");
    you may not use this file except in compliance with the License.
    You may obtain a copy of the License at

      http://www.apache.org/licenses/LICENSE-2.0

    Unless required by applicable law or agreed to in writing, software
    distributed under the License is distributed on an "AS IS" BASIS,
    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
    See the License for the specific language governing permissions and
    limitations under the License.
*/

package collector

import (
	"testing"

	pluginMock "github.com/solarwinds/snap-plugin-lib/v2/mock"
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
