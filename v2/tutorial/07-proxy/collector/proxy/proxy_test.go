// +build tutorial

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
