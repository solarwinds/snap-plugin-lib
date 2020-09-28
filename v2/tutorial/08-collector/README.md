# System collector

In the previous chapter we have written code responsible for gathering system information, yet unable to work in snap environment.
Now we will complete a missing part of the collector.  

## Implementing `Collector` Module

Let's start from creating a structure representing our collector, holding reference to `Proxy` interface (in `./collector/collector.go`)

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

Next step is to provide implementations of `Collect` and `PluginDefinition` methods.
Due to the fact that we will be collecting dynamic metrics, the second method is required.

> `PluginDefinition` is not required when plugin collects only static metrics. 
> However, it's a good practice to provide as much information about plugin as possible to help others understand its purpose (and produce better results, ie. when running a plugin with `-print-example-task` flag) 

Let's start from implementing the second one:
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
At the beginning we are defining a dynamic element (often referred to a group) binding its name with a description.
Then we are defining 4 metrics, 2 of which are dynamic. 

> Dynamic element is always surrounded by `[]` in definition.

> Metrics can contain more than 1 dynamic element, but there are restrictions:
> - first and last element have to be static (ie. `/minisystem/devices/[type]/[producer]/mem_usage)
> - you can't define static and dynamic element at the same position when they have common prefix. For example: 2nd elements of given metrics `/minisystem/[processName]/cpu`, `/minisystem/usage/cpu`: `[processName]` and `usage` have the same prefix; it's not allowed)

Before we introduce `Collect` implementation let's add some helper methods which will convert measurements (of type defined in `data` module) into metrics:

```go
func (s systemCollector) collectTotalCPU(ctx plugin.CollectContext) error {
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
If there were no errors, metric is added to results by calling `ctx.AddMetric()`.

> You can ignore error value returned from `ctx.AddMetric()`. It's should be rather used only during debugging (see [explanation](/v2/tutorial/faq#should-i-handle-error-value-from-ctxaddmetric))    

`collectTotalMemory` is similar to method we have just implemented:
```go
func (s systemCollector) collectTotalMemory(ctx plugin.CollectContext) error {
	memory, err := s.proxyCollector.TotalCpuUsage()
	if err != nil {
		return fmt.Errorf("can't create metric for total memory utilization: %v", err)
	}

	_ = ctx.AddMetric("/minisystem/usage/memory", memory)
	return nil
}
```

The last helper method will be a little bit more complicated:
```go
func (s systemCollector) collectProcessesInfo(ctx plugin.CollectContext) error {
    procsInfo, err := s.proxyCollector.ProcessesInfo()
    if err != nil {
        return fmt.Errorf("can't create metrics associated with processes")
    }

    for _, p := range procsInfo {
        pName := s.sanitizeName(p.ProcessName)

        cpuMetricNs := fmt.Sprintf("/minisystem/processes/[processName=%s]/cpu", pName)
        _ = ctx.AddMetric(cpuMetricNs, p.CpuUsage, plugin.MetricTag("PID", fmt.Sprintf("%d", p.PID)))
        
        memMetricNs := fmt.Sprintf("/minisystem/processes/[processName=%s]/memory", pName)
        _ = ctx.AddMetric(memMetricNs, p.MemoryUsage, plugin.MetricTag("PID", fmt.Sprintf("%d", p.PID)))
    }

    return nil
}
```
At the begging we are getting list of processes by calling `ProcessesInfo()` on `Proxy` interface.
After handling error, we are iterating over each element from results and create two metrics: one related to cpu, one related to memory.

Be aware that we are using special format for dynamic element (`/minisystem/processes/[processName=mysql]/memory`).
Generally syntax `/minisystem/processes/mysql/memory` is also allowed, but it's much more readable for other developers to use '[]' when working with dynamic elements.

Last but not least, pay attention that we needed to sanitize process name (since it will become part of metric name).

> Valid character for namespace elements are letters, numbers and `_`

```go
var sanitizeRegex = regexp.MustCompile(`[()[\]{}<> ,./?;':"|^!\-_+\\]+`)

func (s systemCollector) sanitizeName(n string) string {
	return strings.ToLower(sanitizeRegex.ReplaceAllString(n, "_"))
}
```

When all helpers are finished, we can finally implement `Collect`.
```go
func (s systemCollector) Collect(ctx plugin.CollectContext) error {
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

Method simply calls all helpers and finishes with error if any problem arise. 

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

* [Table of contents](/v2/README.md)
- Previous Chapter: [Gathering data (Proxy Collector)](/v2/tutorial/07-proxy/README.md)
- Next Chapter: [Handle configuration](/v2/tutorial/09-config/README.md)
