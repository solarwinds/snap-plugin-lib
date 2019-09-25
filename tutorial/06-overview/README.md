# Advanced Plugin - Introduction

In the previous chapters you've learnt how to write a very simple plugin and utilize functionality that the library provided. 
Now, we will teach you how to build advanced, practical collector which will gather basic information about monitored system:
- cpu usage (percentage) for each running process 
- memory usage (percentage) for each running process
- total cpu usage (percentage)
- total memory usage (percentage)

## Metrics

## Static metrics

In our simple example, a metric could hold the following information
- metric name (namespace), ie. `/example/time/hour`
- value, ie. `11`
- additional text information (tags), ie. `weekday: Monday`
- unit, ie. `second` (could be defined in `DefinePlugin` method)
- measurement time (added by a library)

Total CPU and memory usage can be represented by two following metrics:
- `/minisystem/usage/cpu`
- `/minisystem/usage/memory`

### Dynamic metrics

Sometimes it's very useful to define metric in a way that specific element (or elements) of its name is not constant.
When gathering cpu or memory utilization for each process, we don't know the list of processes running in minisystem.
It would be convenient to define some kind of a template name, based on which concrete metrics could be generated.

That functionality is offered by dynamic metrics. We can define it in `PluginDefinition` using special form of metric name:
- `/minisystem/processes/[processName]/cpu`
- `/minisystem/processes/[processName]/memory`

Then, when we would iterate over a list of processes, we can replace `processName` with a concrete name, ie.
- `/minisystem/processes/mysql/cpu`
- `/minisystem/processes/mysql/memory`
- `/minisystem/processes/chrome/cpu`
- `/minisystem/processes/chrome/memory`
- `/minisystem/processes/vbox/cpu`
- `/minisystem/processes/vbox/memory`

Later we will see that we can control (via configuration) which dynamic metrics are gathered (ie. we may want to gather only `chrome` and `mysql` utilization).

To sum up, our plugin will define 4 metrics:
- `/minisystem/usage/cpu`
- `/minisystem/usage/memory`
- `/minisystem/processes/[processName]/cpu`
- `/minisystem/processes/[processName]/memory`

We will write related code in [Chapter 8](/tutorial/08-collector/README.md).

### Tags

As previously stated, tags offer additional information associated with measurement. 
System collector will add information about PID to processes related metrics.

## Configuration 

In a typical system with many running processes, only specific thresholds might be of interest (ie. CPU > 5%).
Also, accuracy of some measurements may be related to the duration that the measurement takes.

Giving that requirements, our example configuration may be written as:

```json
{
    "processes": {
        "minCpuUsage": 0.05,
        "minMemoryUsage": 0.01
    },
    "totalCpuMeasureDuration": "1s"
}
```

Configuration handling will be covered in [Chapter 8](/tutorial/08-collector/README.md).

## Code structure

Since created plugin is more complicated and will require a lot of coding, we will structure go files into separate submodules.

Collector structure will look as follows:
- main.go
- collector/
  - data/
  - proxy/

`Collector` module will contain code related to plugin (implementation of `Load`, `Unload`, `Collect` and `DefineMetrics`).
`Proxy` module will retrieve system measurement using psutil library.
Separating `Collector` and `Proxy` will enable us to easily write unit tests. 
`Data` module will contain structures shared between `Collector` and `Proxy`.

----

* [Table of contents](/tutorial/README.md)
- Previous Chapter: [Useful tools](/tutorial/05-tools/README.md)
- Next Chapter: [Gathering data (Proxy Collector)](/tutorial/07-proxy/README.md)

