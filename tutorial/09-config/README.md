# Handle configuration

## Configuration

### Overview

A plugin written in [Chapter 8](/tutorial/08-collector/README.md) already provided us quite useful functionality.
Yet, when we look at metrics associated with processes there is a lot of metric produced - also for processes which cpu and memory utilization is infinitesimal.
We might want to track only processes which resource utilization is above specific threshold. 
Configuration is perfect way to dynamically provide it: in one system administator may be interested in tracking processes which use above 50% of total memory, in other those values may differ.

In [Overview](/tutorial/06-overview/README.md) we've already mentioned that config will be given as a JSON:
```json
{
    "processes": {
        "minCpuUsage": 0.05,
        "minMemoryUsage": 0.01
    },
    "totalCpuMeasureDuration": "1s"
}
```

### Implementation

Now, let's write code associated with configuration:
- basic validation (ie. `minCpuUsage` should be within range <0;1>)
- providing default values when JSON is not completed
- accessing configuration values from `Collect` and `Load`

Majority of the code will be put into new file (`./collector/config.go`)

We will start from defining default values
```go
const (
	defaultMinCpuUsage             = 0.05
	defaultMinMemoryUsage          = 0.01
	defaultTotalCpuMeasureDuration = "1s"

	configObjectKey = "config"
)
```
First three values are associated with JSON fields. 
Processing configuration will be done only during `Load` stage and stored in plugin `Context` using `configObjectKey`. 
When user would like to have access to processed configuration fields in `Collect`, he will simply call `ctx.Load(configObjectKey)` instead of `ctx.RawConfig()`.

Next step is to create structure that represent configuration.
```go
type config struct {
	Processes               configProcesses
	TotalCpuMeasureDuration string
}

type configProcesses struct {
	MinCpuUsage    float64
	MinMemoryUsage float64
}
```
Go language offers very simple API to convert (unmarchal) bytes into native structures.
You may notice that our `config` and `configProcesses` fields are the same (ignoring word case) to expected from JSON.

Now, we can implement first function (factory method), which will return default configuration.
```go
func defaultConfig() *config {
	return &config{
		Processes: configProcesses{
			MinCpuUsage:    defaultMinCpuUsage,
			MinMemoryUsage: defaultMinMemoryUsage,
		},
		TotalCpuMeasureDuration: defaultTotalCpuMeasureDuration,
	}
}
```

After that we are able to implement handleConfig partially: 
```go
func handleConfig(ctx plugin.Context) error {
    // (...)
	return nil
}
```

First step is to create structure representing default configuration. 
Then we should marshal JSON configuration received from snap (we can access it via `ctx.RawConfig()`)
If JSON is not complete (or even empty) it's not a problem, since we are operating on structure partially filled in by default.

```go
    // (...)
    cfg := defaultConfig()

	err := json.Unmarshal(ctx.RawConfig(), cfg)
	if err != nil {
		return fmt.Errorf("invalid config: %v", err)
	}
    // (...)
```

> We will validate how our plugin react on passing different JSON configurations in unit tests.
> Till now, you can take a look at test code in `./collector/config_test.go`.

We have now access to passed configuration via cfg variable.
What we can do next is to validate values of configuration fields, especially: 
- `totalCpuMeasureDuration` should represent string based on which `time.Duration` object can be build later.
- `minCpuUsage` and `minMemoryUsage` should be in range <0;100> 

Code responsible for validation is given below:
```go
    // (...)
	_, err = time.ParseDuration(cfg.TotalCpuMeasureDuration)
	if err != nil {
		return fmt.Errorf("invalid value for totalCpuMeasureDuration: %v", err)
	}

	if cfg.Processes.MinCpuUsage < 0 || cfg.Processes.MinCpuUsage > 100 {
		return fmt.Errorf("invalid value for minCpuUsage: %v", err)
	}

	if cfg.Processes.MinMemoryUsage < 0 || cfg.Processes.MinMemoryUsage > 100 {
		return fmt.Errorf("invalid value for minMemoryUsage: %v", err)
	}
    // (...)
```

When unmarshalling and validation ended without error we can store final structure and access it from `Collect`.
```go
    // (...)
    ctx.Store(configObjectKey, cfg)
    return nil
```

Complete function:
```go
func handleConfig(ctx plugin.Context) error {
	cfg := defaultConfig()

	err := json.Unmarshal(ctx.RawConfig(), cfg)
	if err != nil {
		return fmt.Errorf("invalid config: %v", err)
	}

	_, err = time.ParseDuration(cfg.TotalCpuMeasureDuration)
	if err != nil {
		return fmt.Errorf("invalid value for totalCpuMeasureDuration: %v", err)
	}

	if cfg.Processes.MinCpuUsage < 0 || cfg.Processes.MinCpuUsage > 100 {
		return fmt.Errorf("invalid value for minCpuUsage: %v", err)
	}

	if cfg.Processes.MinMemoryUsage < 0 || cfg.Processes.MinMemoryUsage > 100 {
		return fmt.Errorf("invalid value for minMemoryUsage: %v", err)
	}

	ctx.Store(configObjectKey, cfg)

	return nil
}
```

One more thing that left is helper method to will give access to remember structure.
```go
func getConfig(ctx plugin.Context) *config {
	obj, ok := ctx.Load(configObjectKey)
	if !ok {
		return defaultConfig()
	}
	return obj.(*config)
}
```

What's is done here is simply calling `ctx.Load()` with cast to appropriate type.
Also we added safety net: if someone will call `getConfig()` before `handleConfig()` we would get default config.
Other solutions would be to return error (as a additional parameter) or throw panic (since it's generally developer mistake to call `getConfig()` earlier).

### `Collect`

All configuration helpers are in place. 
No we can change Collector code.

At first let's process configuration in `Load` stage:
```go
func (s systemCollector) Load(ctx plugin.Context) error {
	return handleConfig(ctx)
}
```

collectTotalCPU at some point calls blocking operation of gopsutil. 
Based on entry from configuration how long blocking should take (how precise result will be).  
```go
func (s systemCollector) collectTotalCPU(ctx plugin.Context) error {
	cfg := getConfig(ctx)
	measurementDur, _ := time.ParseDuration(cfg.TotalCpuMeasureDuration)

	cpu, err := s.proxyCollector.TotalCpuUsage(measurementDur)
	if err != nil {
		return fmt.Errorf("can't create metric for total cpu utilization: %v", err)
	}

	_ = ctx.AddMetric("/minisystem/usage/cpu", cpu)
	return nil
}
``` 

Notice, that we needed to change proxy API: TotalCpuUsage is now taking one parameter - duration.
```go
type Proxy interface {
	ProcessesInfo() ([]data.ProcessInfo, error)
	TotalCpuUsage(time.Duration) (float64, error)
	TotalMemoryUsage() (float64, error)
}
```

The change in function `TotalCpuUsage` is simple (using passed parameter instead default value).
Remaining code should be left unchanged.  
```go
func (p proxyCollector) TotalCpuUsage(d time.Duration) (float64, error) {
	totalCpu, err := cpu.Percent(d, false)
	...
}
```

Other function which we should modified a bit is `collectProcessesInfo`. 
In short: when processes uses cpu below given limit, metric shouldn't be created.
New function is given below:
```go
func (s systemCollector) collectProcessesInfo(ctx plugin.Context) error {
	procsInfo, err := s.proxyCollector.ProcessesInfo()
	if err != nil {
		return fmt.Errorf("can't create metrics associated with processes")
	}

	cfg := getConfig(ctx)

	for _, p := range procsInfo {
		pName := s.sanitizeName(p.ProcessName)

		if p.CpuUsage >= cfg.Processes.MinCpuUsage {
			cpuMetricNs := fmt.Sprintf("/minisystem/processes/[processName=%s]/cpu", pName)
			_ = ctx.AddMetricWithTags(cpuMetricNs, p.CpuUsage, map[string]string{"PID": fmt.Sprintf("%d", p.PID)})
		}

		if p.MemoryUsage >= cfg.Processes.MinMemoryUsage {
			memMetricNs := fmt.Sprintf("/minisystem/processes/[processName=%s]/memory", pName)
			_ = ctx.AddMetricWithTags(memMetricNs, p.MemoryUsage, map[string]string{"PID": fmt.Sprintf("%d", p.PID)})
		}
	}

	return nil
}
```

After retrieving process list we are calling `getConfig(ctx)` which returns processed configuration.
Then, in the loop, we are checking if cpu and memory values are greater that given thresholds.
If so, metrics are created (so only the most "meaningful" resources are gathered).

> You can take a look at example unit test in `./collector/collector_test.go` which validates using limits.

----

* [Table of contents](/tutorial/README.md)
- Previous Chapter: [Implementing System collector](/tutorial/08-collector/README.md)
- Next Chapter: [FAQ](/tutorial/faq/README.md)
