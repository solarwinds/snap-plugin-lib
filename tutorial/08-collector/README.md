# System collector

In the previous chapter we have written code responsible for gathering system information, yet enable to work in snap environment.
Now will write a missing part.  

## Implementing `Collector` Module

Let's start from creating a structure representing our collector which hold reference to `Proxy` interface (in `./collector/collector.go`)

```go
type systemCollector struct {
	proxyCollector proxy.Proxy
}

func New(proxy proxy.Proxy) plugin.Collector {
	return systemCollector{
		proxyCollector: proxy,
	}
}
```

Next, important step is to provide implementation of `Collect` and `PluginDefinition`.
Due to the fact that we will be collecting dynamic metrics the second method is required.

> When plugin collects only static metrics (as in the case of first plugin), `PluginDefinition` is not required. 
> However, it's a good practice to provide as much information about plugin as possible by plugin's creator. 

Let's start from the second one:
```go
func (s systemCollector) PluginDefinition(def plugin.CollectorDefinition) error {
	def.DefineGroup("processName", "process name")

	def.DefineMetric("/minisystem/processes/[processName]/cpu", "%", true, "CPU Utilization by current process")
	def.DefineMetric("/minisystem/processes/[processName]/memory", "%", true, "Memory Utilization by current process")
	def.DefineMetric("/minisystem/usage/cpu", "%", true, "Total CPU Utilization")
	def.DefineMetric("/minisystem/usage/memory", "%", true, "Total memory Utilization")

	return nil
}
```
At the beginning we are defining a dynamic element (ofter referred as a group) binding its name with a description.
Then we are defining 4 metrics, 2 of which are dynamic. 

> Dynamic element is always surrounded by `[]` in definition.

> Metrics can contain more than 1 dynamic elements, but there are restrictions:
> - first and last elements must be static (ie. `/minisystem/devices/[type]/[producer]/mem_usage)
> - you can't define static and dynamic element at the same position when they have common prefix (ie. for `/minisystem/[processName]/cpu`, `/minisystem/usage/cpu`: `[processName]` and `usage` have the same prefix - it's not allowed)

Before we introduce `Collect` implementation let's implement method which convert measurements (of type defined in `data` module) into metrics:

```
func (s systemCollector) collectTotalCPU(ctx plugin.Context) error {
	cpu, err := s.proxyCollector.TotalCpuUsage()
	if err != nil {
		return fmt.Errorf("can't create metric for total cpu utilization: %v", err)
	}

	_ = ctx.AddMetric("/minisystem/usage/cpu", cpu)
	return nil
}
```

In the first line we are calling `TotalCpuUsage`, which is part of `Proxy` interface. 
Then we handle possible errors by wrapping it into new one.
If there were no error, metric is added to results by calling `ctx.AddMetric()`.
(refer to **TODO** why not to handle errors)

`collectTotalMemory` is almost the same:
```go
func (s systemCollector) collectTotalMemory(ctx plugin.Context) error {
	memory, err := s.proxyCollector.TotalCpuUsage()
	if err != nil {
		return fmt.Errorf("can't create metric for total memory utilization: %v", err)
	}

	_ = ctx.AddMetric("/minisystem/usage/memory", memory)
	return nil
}
```

The last helper method is a little bit more complicated:
```go
func (s systemCollector) collectProcessesInfo(ctx plugin.Context) error {
	procsInfo, err := s.proxyCollector.ProcessesInfo()
	if err != nil {
		return fmt.Errorf("can't create metrics associated with processes")
	}

	for _, p := range procsInfo {
		pName := s.sanitizeName(p.ProcessName)

		cpuMetricNs := fmt.Sprintf("/minisystem/processes/[processName=%s]/cpu", pName)
		_ = ctx.AddMetricWithTags(cpuMetricNs, p.CpuUsage, map[string]string{"PID": fmt.Sprintf("%d", p.PID)})

		memMetricNs := fmt.Sprintf("/minisystem/processes/[processName=%s]/memory", pName)
		_ = ctx.AddMetricWithTags(memMetricNs, p.MemoryUsage, map[string]string{"PID": fmt.Sprintf("%d", p.PID)})
	}

	return nil
}
```
At the begging we are getting list of processes by calling `ProcessesInfo()` on `Proxy` interface.
After handling error, we are iterating over each element from results and create two metrics: one related to cpu, one related to memory.
Please take a look that we are using special format for dynamic element (ie. `/minisystem/processes/\[processName=mysql\]/memory")
Generally syntax `/minisystem/processes/mysql/memory` would be also allowed, but it's much more readable for other developers to use '[]' when adding dynamic elements.
Last but not least, pay attention that we needed to sanitize process name (since it will become part of metric name).

> Valid character for namespace elements are letters, numbers and `_`

```go
var sanitizeRegex = regexp.MustCompile(`[()[\]{}<> ,./?;':"|^!\-_+\\]+`)

func (s systemCollector) sanitizeName(n string) string {
	return strings.ToLower(sanitizeRegex.ReplaceAllString(n, "_"))
}
```

When all helpers are finished we can finally implement `Collect`.
```go
func (s systemCollector) Collect(ctx plugin.Context) error {
	err := s.collectTotalCPU(ctx)
	if err != nil {
		return err
	}

	err = s.collectTotalMemory(ctx)
	if err != nil {
		return err
	}

	err = s.collectProcessesInfo(ctx)
	if err != nil {
		return err
	}

	return nil
}
```

Method simply calls all helpers and end with errors if any problem arise. 

> Instead of returning error from `Collect` which is signal to framework that measurement went wrong, we could just log errors and continue gathering other metrics.
> It's up to developer to decide if it's satisfactory to return only partial measurement. 

## Manual validation

Let's execute plugin in debug mode with some filters defined (to limit output)
```bash
./08-collector -debug-mode -debug-collect-interval=1s -debug-collect-counts=4 -plugin-filter="/minisystem/usage/*;/minisystem/processes/08_collector_exe/memory"
```

Example output:
```
Gathered metrics (length=3):
minisystem.usage.cpu 0.9803921568627451 {map[]}
minisystem.usage.memory 25 {map[]}
minisystem.processes.[processName=08_collector_exe].memory 0.03228510916233063 {map[PID:14796]}

Gathered metrics (length=3):
minisystem.usage.cpu 0.3913894324853229 {map[]}
minisystem.usage.memory 25 {map[]}
minisystem.processes.[processName=08_collector_exe].memory 0.03609180822968483 {map[PID:14796]}

Gathered metrics (length=3):
minisystem.usage.cpu 0 {map[]}
minisystem.usage.memory 25 {map[]}
minisystem.processes.[processName=08_collector_exe].memory 0.040461134165525436 {map[PID:14796]}

Gathered metrics (length=3):
minisystem.usage.cpu 0.1949317738791423 {map[]}
minisystem.usage.memory 25 {map[]}
minisystem.processes.[processName=08_collector_exe].memory 0.04088011011481285 {map[PID:14796]}
```

----

* [Table of contents](/tutorial/README.md)
- Previous Chapter: [Gathering data (Proxy Collector)](/tutorial/07-proxy/README.md)
- Next Chapter: [Handle configuration](/tutorial/09-config/README.md)
