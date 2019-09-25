# Proxy Collector

## Obtaining system information - Psutil library

Go language has numerous 3rd party libraries that we can use, instead of writing our own functionalities.
Retrieving system information is quite a common task for any developer and in golang ecosystem we could easily find library that offers easy-to-use cross platform API.
For our plugin we will use [gopsutil](https://github.com/shirou/gopsutil).

What we need to retrieve is:
- processes names (metric name)
- cpu utilization of each process (value) 
- memory utilization of each process (value)
- process's PID (for tags)
- total cpu utilization (value)
- total memory utilization (value)

### API

First four can be obtained via the following snippet (for now let's ignore errors):
```go
import ("github.com/shirou/gopsutil/process")

func Do() {
    processesData, _ := process.Processes()
    for _, proc := range processesData {
        name, _ := proc.Name()
        cpuPerc, _ := proc.CPUPercent()
        memPerc, _ := proc.MemoryPercent()
        ...
    }
}
```

Total CPU and memory are accessed by: 
```go
cpu.Percent(1 * time.Second, false) // result[0] 
mem.VirtualMemory()                 // result.UsedPercent
```

from other gopsutil modules:
```go
import (
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
)
```

> `cpu.Percent` returns an array. Total value is in the first element of the array.

> `mem.VirtualMemory` return general information structure. Field `.UserPercent` provides measurement we are interested in.

## Implementing `Data` Module

As mention earlier `Data` module will contain common structures. 
We need one to hold simplified information about the process. 
Also, our collector module should not depend on the gopsutil library directly.

So far only file we need is `./collector/data/processinfo.go` containing only simple structure 

```go
package data

type ProcessInfo struct {
    ProcessName string
    CpuUsage    float64
    MemoryUsage float64
    PID         int32
}
```

## Implementing `Proxy` Module

`Proxy` will utilize psutil API described above. Also at this time we will start handling errors.

Our implementation will be located in `./collector/proxy/proxy.go`

Let's start from the headers
```go
package proxy

import (
    "fmt"
    "github.com/librato/snap-plugin-lib-go/tutorial/07-proxy/collector/data"
    "github.com/shirou/gopsutil/cpu"
    "github.com/shirou/gopsutil/mem"
    "github.com/shirou/gopsutil/process"
    "time"
)
```

First argument of `cpu.Percent` is the duration for which measurement is taken. 
Later we will be able to configure it, but now let's define const value

```go
const defaultCPUMeasurementTime = 1 * time.Second
``` 

Now, let's define interface which will be used by the `Collector`.

```go
type Proxy interface {
	ProcessesInfo() ([]data.ProcessInfo, error)
	TotalCpuUsage() (float64, error)
	TotalMemoryUsage() (float64, error)
}
```

`Proxy` interface will be implemented by `proxyCollector`:
```go
type proxyCollector struct{}

func New() Proxy {
	return &proxyCollector{}
}
```

Implementation of first function:
```go
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
```

First a list of all processes is obtained from psutil (`process.Processes()`). 
Then we iterate over each process and read: name, cpu and memory utilization and PID. 
At the end (of iteration) we create simplified structure representing single process.
Function returns list of all such structures. 

Total CPU is even simpler to write:
```go
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
```

We are calling `cpu.Percent` API and in case of a failure returning wrapped error.
Additionally, we are checking that slice returned from gopsutil has only one element (it can have more if the second argument is true).
If everything is correct, total cpu is returned.

Last method:
```go
func (p proxyCollector) TotalMemoryUsage() (float64, error) {
	memoryInfo, err := mem.VirtualMemory()
	if err != nil {
		return 0, fmt.Errorf("can't obtain memory information: %v", err)
	}

	return memoryInfo.UsedPercent, nil
}
```

Here we are calling `mem.VirutalUsage` and retrieve value from field `.UserPercent`. 

## Smoke tests

Let's create a test file to manually validate that functions are working correctly (`./collector/proxy/proxy_test.go`)

> Be aware that tests presented in this chapter are not good candidates for stable regressions.
> They are introduced simply to show result of our initial implementation

Headers:
```go
package proxy

import (
	"fmt"
	"testing"
)
```

Simple test (we are checking if result is different from 0.0 which in real situation may not always be the case)
```go
func TestTotalCPUUsage(t *testing.T) {
	p := proxyCollector{}

	cpu, _ := p.TotalCpuUsage()

	if cpu == 0.0 {
		t.Fail()
	}

	fmt.Printf("CPU=%v%%\n", cpu)
}
```

Adequate test for total memory usage 
```go
func TestTotalMemoryUsage(t *testing.T) {
	p := proxyCollector{}

	memory, _ := p.TotalMemoryUsage()

	if memory == 0.0 {
		t.Fail()
	}

	fmt.Printf("Memory=%v%%\n", memory)
}
```

Last manual test for process list 
```go
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
```

We can execute using the following command line to see passed/failed result of a test
```bash
go test ./...
```

or 
```bash
go test ./... -v
```

to see output by `fmt.Printf`

----

* [Table of contents](/tutorial/README.md)
- Previous Chapter: [Advanced Plugin - Introduction](/tutorial/06-overview/README.md)
- Next Chapter: [Implementing System collector](/tutorial/08-collector/README.md)
