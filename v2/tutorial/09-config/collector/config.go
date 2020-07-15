/*
 Copyright (c) 2020 SolarWinds Worldwide, LLC

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
	"encoding/json"
	"fmt"
	"time"

	"github.com/librato/snap-plugin-lib-go/v2/plugin"
)

const (
	defaultMinCPUUsage             = 0.05
	defaultMinMemoryUsage          = 0.01
	defaultTotalCPUMeasureDuration = "1s"

	configObjectKey = "config"
)

type config struct {
	Processes               configProcesses
	TotalCPUMeasureDuration string
}

type configProcesses struct {
	MinCPUUsage    float64
	MinMemoryUsage float64
}

func defaultConfig() config {
	return config{
		Processes: configProcesses{
			MinCPUUsage:    defaultMinCPUUsage,
			MinMemoryUsage: defaultMinMemoryUsage,
		},
		TotalCPUMeasureDuration: defaultTotalCPUMeasureDuration,
	}
}

func handleConfig(ctx plugin.Context) error {
	cfg := defaultConfig()

	err := json.Unmarshal(ctx.RawConfig(), &cfg)
	if err != nil {
		return fmt.Errorf("invalid config: %v", err)
	}

	_, err = time.ParseDuration(cfg.TotalCPUMeasureDuration)
	if err != nil {
		return fmt.Errorf("invalid value for totalCpuMeasureDuration: %v", err)
	}

	if cfg.Processes.MinCPUUsage < 0 || cfg.Processes.MinCPUUsage > 100 {
		return fmt.Errorf("invalid value for minCpuUsage: %v", err)
	}

	if cfg.Processes.MinMemoryUsage < 0 || cfg.Processes.MinMemoryUsage > 100 {
		return fmt.Errorf("invalid value for minMemoryUsage: %v", err)
	}

	ctx.Store(configObjectKey, &cfg)

	return nil
}

func getConfig(ctx plugin.Context) config {
	obj, ok := ctx.Load(configObjectKey)
	if !ok {
		return defaultConfig()
	}
	return *(obj.(*config))
}
