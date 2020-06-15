# Basic concepts

In [Chapter 1](/v2/tutorial/01-simple/README.md) and [Chapter 2](/v2/tutorial/02-testing/README.md) you've learned how to write and validate a simple collector. 
Now we will introduce more advanced concepts that version 2 of plugin-lib-go introduced, this will help developer build their plugins easier and faster.

## Tasks

In a typical situation plugins are controlled and managed by snap.
In snap v3 (along with plugin-lib-go v2) user configures one or several tasks that will be requested from a single or several plugins.
A single task contains information about configuration and requested metrics.
When several tasks are requested from the same plugin, by default only one instance of plugin's binary will be run by snap. 
Plugin-lib-go provides facilities for maintaining different tasks.

### Context 

Previously, when we defined a plugin algorithm we had to provide implementation of `Collect` method:  

```go
func (s simpleCollector) Collect(ctx plugin.CollectContext) error {
    ...
	_ = ctx.AddMetric("/example/date/day", t.Day())
    ...
}
```

`Collect` takes one argument, which is a context - an object associated with current tasks, which allows:
- adding metrics (measurements) during current collection request,
- access configuration values, 
- to maintain state between collections for the same task,
- to optimize metrics calculation

We will slightly modify our previous example in order to check those features.

### Load(), Unload()

When task is defined for a plugin, snap will send a `Load()` request to plugin containing:
- task identifier - unique value maintained by snap
- JSON-like object with configuration fields
- list of metrics that user wants to gather - we can request only a subset of "measurements"

When a task handling is no longer needed snap sends an `Unload()` request.
As for `Load()`, we can provide some custom code which will be executed when the task is finished.

> Keep in mind, however, this doesn't necessarily mean the plugin is not needed as some other tasks that rely on this plugin may still be running.

Let's introduce empty custom implementation of `Load()` method:
```go
func (s simpleCollector) Load(ctx plugin.Context) error {
	return nil
}
```

We could also add custom implementation for `Unload` but it's not required in this example.
```go
func (s simpleCollector) Unload(ctx plugin.Context) error {
	return nil
}
```

> Custom implementation of `Unload()` method should be provided when plugin is storing some object (ie. http client) that needs to be manually released (ie. via `obj.Close()`) to avoid memory or resource leaks.

#### Configuration

Example plugin defined 5 metrics - one of them gives information about current hour (0-23).
We may dynamically request different format (short: 0-11, long:0-23) by defining configuration for a task.
Plugin will expect configuration in JSON format. In our case it may be simply:
```json
{
  "format": "short"
}
```

We can access configuration fields in two ways.
- by accessing method `ctx.Config` which implements simplified access to the configuration values
- by accessing method `ctx.RawConfig` which returns JSON object (bytes) that needs to be manually unmarshaled.

First method will be introduced in this chapter. Alternative will be presented in [Chapter 9](/v2/tutorial/09-config/README.md)

Let's create a helper method which reads `format` field:
```go
func (s simpleCollector) format(ctx plugin.Context) string {
	fm, _ := ctx.Config("format")
	if fm == "short" {
		return fm
	}
	return "long"
}
```
`ctx.Config` returns two values: field associated string and a bool flag indicating that the field was present in configuration.
If `format` field had value `short` we will return it, otherwise default `long` is returned in other situations.

> Notice, that `ctx.Config()` always return a string even if different data type was provided in JSON, for example int or bool.
> If you want strict type control use `ctx.RawConfig()`.

Modified version of `Collect` method:

```go
func (s simpleCollector) Collect(ctx plugin.CollectContext) error {
	// Collect data
	t := time.Now()

	// Handle configuration
	hour := t.Hour()
	if s.format(ctx) == "short" {
		hour %= 12
	}

	// Convert data to metric form
	_ = ctx.AddMetric("/example/date/day", t.Day())
	_ = ctx.AddMetric("/example/date/month", int(t.Month()))
	_ = ctx.AddMetric("/example/time/hour", hour)
	_ = ctx.AddMetric("/example/time/minute", t.Minute())
	_ = ctx.AddMetric("/example/time/second", t.Second())

	return nil;
}
```

Now we can test it using command line (requesting short format):
```bash
./03-concepts -debug-mode=1 -debug-collect-counts=1 -debug-collect-interval=5s -plugin-config='{"format": "short"}'
```

Output:
```
Gathered metrics (length=5):
example.date.day 8 {map[]}
example.date.month 9 {map[]}
example.time.hour 8 {map[]}
example.time.minute 34 {map[]}
example.time.second 15 {map[]}
```

Analogically, long format requested by any of the following commands:
```bash
./03-concepts -debug-mode=1 -debug-collect-counts=1 -debug-collect-interval=5s -plugin-config='{"format": "long"}'
./03-concepts -debug-mode=1 -debug-collect-counts=1 -debug-collect-interval=5s -plugin-config='{"format": "other"}'
./03-concepts -debug-mode=1 -debug-collect-counts=1 -debug-collect-interval=5s -plugin-config='{}'
./03-concepts -debug-mode=1 -debug-collect-counts=1 -debug-collect-interval=5s
```

gives output:
```bash
Gathered metrics (length=5):
example.date.day 8 {map[]}
example.date.month 9 {map[]}
example.time.hour 20 {map[]}
example.time.minute 34 {map[]}
example.time.second 21 {map[]}
```

#### State

At times, you will need to remember values and objects between consecutive collections, ie:
- credentials,
- objects representing any client created during Load phase (like: TPC, REST, GRPC etc),
- cashe,
- post-processed configuration (see: [State and configuration](https://github.com/librato/snap-plugin-lib-go/blob/master/v2/tutorial/03-concepts/README.md#state-and-configuration)),
- custom statistics.

In that case, `ctx.Store()` and `ctx.Load()` come in handy, allowing to store and read objects for a given task (context).

> You shouldn't use `simpleCollector` struct members to store state, it's not task-aware.

Let's add a new metric (`/example/count/running`) which provide plugin runtime information (or being more precise, load time of a particular task)

In order to enable plugin running duration calculation, we have to save current time in `Load()` method:
```go
func (s simpleCollector) Load(ctx plugin.Context) error {
	ctx.Store("startTime", time.Now())
	return nil
}
```

We modify `Collect()` method to use that variable every iteration:
```go

func (s simpleCollector) Collect(ctx plugin.CollectContext) error {
	// ...

	// Count metrics
	startTime, _ := ctx.Load("startTime")
	runningDuration := int(time.Now().Sub(startTime.(time.Time)).Seconds())
	_ = ctx.AddMetric("/example/count/running", runningDuration)

	return nil;
}
```

Execution command:
```
./03-concepts -debug-mode=1 -debug-collect-counts=3 -debug-collect-interval=5s
```

Output:
```
Gathered metrics (length=6):
example.date.day 9 {map[]}
example.date.month 9 {map[]}
example.time.hour 6 {map[]}
example.time.minute 53 {map[]}
example.time.second 14 {map[]}
example.count.running 0 {map[]}

Gathered metrics (length=6):
example.date.day 9 {map[]}
example.date.month 9 {map[]}
example.time.hour 6 {map[]}
example.time.minute 53 {map[]}
example.time.second 19 {map[]}
example.count.running 5 {map[]}

Gathered metrics (length=6):
example.date.day 9 {map[]}
example.date.month 9 {map[]}
example.time.hour 6 {map[]}
example.time.minute 53 {map[]}
example.time.second 24 {map[]}
example.count.running 10 {map[]}
```

#### State and Configuration

Additionally state can be used to optimize processing configuration values. 
In the [previous section](https://github.com/librato/snap-plugin-lib-go/blob/master/v2/tutorial/03-concepts/README.md#configuration) "format" option was read during each collection.
Alternatively, we could read it only once during `Load()` and store in context. 

Example:
```go
func (s simpleCollector) format(ctx plugin.Context) string {
    fm, _ := ctx.Load("configFormat")
    return fm.(string)
}
  
func (s simpleCollector) Load(ctx plugin.Context) error {
    fm, _ := ctx.Config("format")
    if fm == "short" {
        ctx.Store("configFormat", "short")
    } else {
        ctx.Store("configFormat", "long")
    }
  
    ctx.Store("startTime", time.Now())
    return nil
 }
```

This approach will be used also in [Chapter 9](/v2/tutorial/09-config/README.md).

----

* [Table of contents](/v2/README.md)
- Previous Chapter: [Testing](/v2/tutorial/02-testing/README.md)
- Next Chapter: [Metrics - filters, definition, tags](/v2/tutorial/04-metrics/README.md)
