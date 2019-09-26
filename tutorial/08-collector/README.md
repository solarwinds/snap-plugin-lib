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




* [Table of contents](/tutorial/README.md)
- Previous Chapter: [Gathering data (Proxy Collector)](/tutorial/07-proxy/README.md)
- Next Chapter: [Handle configuration](/tutorial/09-config/README.md)
