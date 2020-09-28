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

package proxy

import (
	"fmt"
	"time"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/process"
	"github.com/solarwinds/snap-plugin-lib/v2/tutorial/07-proxy/collector/data"
)

const defaultCPUMeasurementTime = 1 * time.Second

type Proxy interface {
	ProcessesInfo() ([]data.ProcessInfo, error)
	TotalCpuUsage() (float64, error)
	TotalMemoryUsage() (float64, error)
}

type proxyCollector struct{}

func New() Proxy {
	return &proxyCollector{}
}

func (p proxyCollector) ProcessesInfo() ([]data.ProcessInfo, error) {
	procInfo := []data.ProcessInfo{}

	processesData, err := process.Processes()
	if err != nil {
		return procInfo, fmt.Errorf("can't obtain list of processes: %v", err)
	}

	for _, proc := range processesData {
		name, err := proc.Name()
		if err != nil {
			continue
		}

		cpuPerc, err := proc.CPUPercent()
		if err != nil {
			continue
		}

		memPerc, err := proc.MemoryPercent()
		if err != nil {
			continue
		}

		procInfo = append(procInfo, data.ProcessInfo{
			ProcessName: name,
			CpuUsage:    cpuPerc,
			MemoryUsage: float64(memPerc),
			PID:         proc.Pid,
		})
	}

	return procInfo, nil
}

func (p proxyCollector) TotalCpuUsage() (float64, error) {
	totalCpu, err := cpu.Percent(defaultCPUMeasurementTime, false)
	if err != nil {
		return 0, fmt.Errorf("can't obtain cpu information: %v", err)
	}
	if len(totalCpu) == 0 {
		return 0, fmt.Errorf("unexpected cpu information: %v", err)
	}

	return totalCpu[0], nil
}

func (p proxyCollector) TotalMemoryUsage() (float64, error) {
	memoryInfo, err := mem.VirtualMemory()
	if err != nil {
		return 0, fmt.Errorf("can't obtain memory information: %v", err)
	}

	return memoryInfo.UsedPercent, nil
}
