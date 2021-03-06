# Handle configuration

## Configuration

### Overview

A plugin written in [Chapter 8](/v2/tutorial/08-collector/README.md) already provides quite useful functionality.
Yet, when we look at a result, there are many metrics produced, numerous of them associated with no-essential information (small utilization of cpu and memory by the majority of processes).
We might want to track only processes which resource utilization is above specific threshold. 
Configuration is a perfect way to dynamically provide it.

In [Overview](/v2/tutorial/06-overview/README.md) we've already mentioned that config will be given as a JSON, ie.
```json
{
	"processes": {
		"minCPUUsage": 0.05,
		"minMemoryUsage": 0.01
	},
	"totalCPUMeasureDuration": "1s"
}
```

### Implementation

Now, let's write code associated with configuration:
- basic validation (ie. `minCPUUsage` and `minMemoryUsage` should be within range <0;100>)
- providing default values when JSON is not completed (or empty)
- accessing configuration values from `Collect`

Majority of the code will be put into new file (`./collector/config.go`)

We will start by defining the default values
```go
const (
	defaultminCPUUsage             = 0.05
	defaultMinMemoryUsage          = 0.01
	defaulttotalCPUMeasureDuration = "1s"

	configObjectKey = "config"
)
```
First three values are associated with JSON fields. 
Processing configuration will be done only during `Load` stage and stored in plugin `Context` using `configObjectKey`. 
When the user wants to have access to processed configuration fields in `Collect`, he can simply call `ctx.Load(configObjectKey)` instead of `ctx.RawConfig()`.

Next step is to create structure that represents configuration.
```go
type config struct {
	Processes               configProcesses
	totalCPUMeasureDuration string
}

type configProcesses struct {
	minCPUUsage    float64
	MinMemoryUsage float64
}
```
Go language offers very simple API to convert (unmarshal) bytes into native structures.
You may notice that our `config` and `configProcesses` contains the same fields as expected from JSON.

Now, we can implement first function (factory method), which will return default configuration.
```go
func defaultConfig() config {
	return config{
		Processes: configProcesses{
			minCPUUsage:    defaultminCPUUsage,
			MinMemoryUsage: defaultMinMemoryUsage,
		},
		totalCPUMeasureDuration: defaulttotalCPUMeasureDuration,
	}
}
```

After that we are able to start implementing `handleConfig`: 
```go
func handleConfig(ctx plugin.Context) error {
	// (...)
	return nil
}
```

First step is to create variable (structure) representing default configuration. 
Then we should unmarshal JSON configuration received from snap (we can access it via `ctx.RawConfig()`)
In case some fields are not set in JSON, the defaults will be preserved.

```go
	// (...)
	cfg := defaultConfig()

	err := json.Unmarshal(ctx.RawConfig(), &cfg)
	if err != nil {
		return fmt.Errorf("invalid config: %v", err)
	}
	// (...)
```

> We will validate how our plugin reacts on passing different JSON configurations in unit tests.
> You can take a look at test code in `./collector/config_test.go`.

The next thing to do is:
- `totalCPUMeasureDuration` should represent string, based on which `time.Duration` can be created later.
- `minCPUUsage` and `minMemoryUsage` should be in range <0;100> 

A sample code responsible for the validation is given below:
```go
	// (...)
	_, err = time.ParseDuration(cfg.totalCPUMeasureDuration)
	if err != nil {
		return fmt.Errorf("invalid value for totalCPUMeasureDuration: %v", err)
	}

	if cfg.Processes.minCPUUsage < 0 || cfg.Processes.minCPUUsage > 100 {
		return fmt.Errorf("invalid value for minCPUUsage: %v", err)
	}

	if cfg.Processes.MinMemoryUsage < 0 || cfg.Processes.MinMemoryUsage > 100 {
		return fmt.Errorf("invalid value for minMemoryUsage: %v", err)
	}
	// (...)
```

If no errors are seen during the processing, we can store configuration structure (to access it later using `Collect`).
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

	_, err = time.ParseDuration(cfg.totalCPUMeasureDuration)
	if err != nil {
		return fmt.Errorf("invalid value for totalCPUMeasureDuration: %v", err)
	}

	if cfg.Processes.minCPUUsage < 0 || cfg.Processes.minCPUUsage > 100 {
		return fmt.Errorf("invalid value for minCPUUsage: %v", err)
	}

	if cfg.Processes.MinMemoryUsage < 0 || cfg.Processes.MinMemoryUsage > 100 {
		return fmt.Errorf("invalid value for minMemoryUsage: %v", err)
	}

	ctx.Store(configObjectKey, &cfg)

	return nil
}
```

The last thing to do is helper method that will give access to remembered configuration structure.
```go
func getConfig(ctx plugin.Context) config {
	obj, ok := ctx.Load(configObjectKey)
	if !ok {
		return defaultConfig()
	}
	return *(obj.(*config))
}
```
We are simply calling `ctx.Load()` with casting to appropriate type.
If `getConfig()` is called before `handleConfig()` default configuration will be returned (other solution would throw error or panic in such a case).

### Implementing `Collect`

Since all configuration helpers are in place, we can implement `Load` and update helpers called from `Collect`.

At first let's process configuration in `Load` stage:
```go
func (s systemCollector) Load(ctx plugin.Context) error {
	return handleConfig(ctx)
}
```

You might remember that `collectTotalCPU` at some point does a blocking call to gopsutil library. 
Having configuration object in place, we can now pass a timeout as an argument to `collectTotalCPU`.
```go
func (s systemCollector) collectTotalCPU(ctx plugin.CollectContext) error {
	cfg := getConfig(ctx)
	measurementDur, _ := time.ParseDuration(cfg.totalCPUMeasureDuration)

	cpu, err := s.proxyCollector.TotalCpuUsage(measurementDur)
	if err != nil {
		return fmt.Errorf("can't create metric for total cpu utilization: %v", err)
	}

	_ = ctx.AddMetric("/minisystem/usage/cpu", cpu)
	return nil
}
``` 

Notice, that we were required to change proxy API: `TotalCpuUsage` is now takes one parameter: the duration.
```go
type Proxy interface {
	ProcessesInfo() ([]data.ProcessInfo, error)
	TotalCpuUsage(time.Duration) (float64, error)
	TotalMemoryUsage() (float64, error)
}
```

The change in function `TotalCpuUsage` is simple (using passed parameter instead of the default value).
Remaining code should be left unchanged.  
```go
func (p proxyCollector) TotalCpuUsage(timeout time.Duration) (float64, error) {
	totalCpu, err := cpu.Percent(timeout, false)
	...
}
```

Other function which we will modify is `collectProcessesInfo`.
When processes uses cpu or memory below given limit, metric shouldn't be created.
```go
func (s systemCollector) collectProcessesInfo(ctx plugin.CollectContext) error {
    procsInfo, err := s.proxyCollector.ProcessesInfo()
    if err != nil {
        return fmt.Errorf("can't create metrics associated with processes")
    }

    cfg := getConfig(ctx)

    for _, p := range procsInfo {
        pName := s.sanitizeName(p.ProcessName)

        if p.CpuUsage >= cfg.Processes.MinCPUUsage {
            cpuMetricNs := fmt.Sprintf("/minisystem/processes/[processName=%s]/cpu", pName)
            _ = ctx.AddMetric(cpuMetricNs, p.CpuUsage, plugin.MetricTag("PID", fmt.Sprintf("%d", p.PID)))
        }

        if p.MemoryUsage >= cfg.Processes.MinMemoryUsage {
            memMetricNs := fmt.Sprintf("/minisystem/processes/[processName=%s]/memory", pName)
            _ = ctx.AddMetric(memMetricNs, p.MemoryUsage, plugin.MetricTag("PID", fmt.Sprintf("%d", p.PID)))
        }
    }

    return nil
}
```

After retrieving process list we are calling `getConfig(ctx)` which returns processed configuration.
Then, in the loop, we are checking if cpu and memory values are greater than given thresholds.
If so, metrics are created (only the most "meaningful" resources are gathered).

> You can take a look at example unit test in `./collector/collector_test.go` which validates usage of limits.

----

* [Table of contents](/v2/README.md)
- Previous Chapter: [Implementing System collector](/v2/tutorial/08-collector/README.md)
- Next Chapter: [Writing plugins in Python and C#](/v2/tutorial/other-languages/README.md)
