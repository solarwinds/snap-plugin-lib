//go:build tutorial
// +build tutorial

/*
 Copyright (c) 2022 SolarWinds Worldwide, LLC

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
	"testing"
)

func TestTotalCPUUsage(t *testing.T) {
	p := proxyCollector{}

	cpu, _ := p.TotalCpuUsage()

	if cpu == 0.0 {
		t.Fail()
	}

	fmt.Printf("CPU=%v%%\n", cpu)
}

func TestTotalMemoryUsage(t *testing.T) {
	p := proxyCollector{}

	memory, _ := p.TotalMemoryUsage()

	if memory == 0.0 {
		t.Fail()
	}

	fmt.Printf("Memory=%v%%\n", memory)
}

func TestProcessesInfo(t *testing.T) {
	p := proxyCollector{}

	procInfo, _ := p.ProcessesInfo()
	if len(procInfo) == 0 {
		t.Fail()
	}

	for _, proc := range procInfo {
		fmt.Printf("%s(%d) cpu=%f, memory=%f\n", proc.ProcessName, proc.PID, proc.CpuUsage, proc.MemoryUsage)
	}
}
