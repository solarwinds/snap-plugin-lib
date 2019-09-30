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




----

* [Table of contents](/tutorial/README.md)
- Previous Chapter: [Implementing System collector](/tutorial/08-collector/README.md)
- Next Chapter: [FAQ](/tutorial/faq/README.md)
