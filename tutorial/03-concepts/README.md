# Basic concepts

In [Chapter 1](/tutorial/01-simple/README.md) and [Chapter 2](/tutorial/02-testing/README.md) you've learned how to write and validate a simple collector. 
Now we will introduce some more advanced concepts that version 2 of plugin-lib-go introduced to help developer build their plugins easier and faster.

## Tasks

In typical situation plugins are controlled and managed by snap.
In snap v3 (along with plugin-lib-go v2) user configures one or several tasks that will be requested from a single or several plugins.
A single task contains information about configuration and requested metrics.
When several tasks are requested from the same plugin, by default only one instance of plugin's binary will be run by snap. 
Plugin-lib-go provides facilities for maintaining different tasks.

### Context 

You might have noticed, that when we defined plugin algorithm, we had to provide implementation of `Collect` method:  

```go
func (s simpleCollector) Collect(ctx plugin.Context) error {
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

When a task is defined for a plugin, snap will send a `Load()` request to plugin containing:
- task identifier - unique value maintained by snap
- JSON-like object with configuration fields
- list of metrics that user wants to gather - we can request only subset of "measurements"

When a task is no longer needed snap may send an `Unload()` request. 
As for `Load()`, we can provide some custom code which will be executed when the task is finished.

Let's introduce empty custom implementation of those methods:
```go
func (s simpleCollector) Load(ctx plugin.Context) error {
	return nil
}

func (s simpleCollector) Unload(ctx plugin.Context) error {
	return nil
}
```

> Keep in mind, however, this doesn't mean the plugin is no longer needed as there might be some other tasks still running that rely on this plugin.

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
- by accessing method `ctx.Config` which implements simplified access to configuration values
- by accessing method `ctx.RawConfig` which returns JSON object (bytes) that needs to be manually unmarshaled.

In this chapter we will use first method. Second one will be presented in [Chapter 9](/tutorial/09-config/README.md) **TODO**linktoanchor**

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
`ctx.Config` returns two values: string associated with field, and bool flag indicating that field was present in configuration.
If `format` field had value `short` we will return it, otherwise default `long` is returned in other situations.

> Notice, that `ctx.Config()` always return string even if value in JSON was provided in other format like int or bool.
> If you want strict type control use `ctx.RawConfig()`.

Modified version of `Collect` method:

```go
func (s simpleCollector) Collect(ctx plugin.Context) error {
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

Now we can test it from a command line (requesting short format):
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

Analogically, long format requested by any of following command:
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

Sometimes you need to remember some values and objects between consecutive collections, ie:
- credentials,
- objects representing any client created during Load phase (like: TPC, REST, GRPC etc),
- cashes,
- post-processed configuration (see: [State and configuration](https://github.com/librato/snap-plugin-lib-go/tree/ao-12231-tutorial/tutorial/03-concepts#state-and-configuration)),
- custom statistics.

In that case, `ctx.Store()` and `ctx.Load()` come in handy, allowing to store and read objects for a given task (context).

> You shouldn't use `simpleCollector` struct members to store state, since it's not task-aware.

Let's add a new metric (`/example/count/running`) which provide information how long plugin is running (or being more precisely: for how long task is being loaded)

In order to enable plugin running duration calculation, we have to save current time in `Load()` method:
```go
func (s simpleCollector) Load(ctx plugin.Context) error {
	ctx.Store("startTime", time.Now())
	return nil
}
```

We modify `Collect()` method to use that variable every iteration:
```go

func (s simpleCollector) Collect(ctx plugin.Context) error {
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

State may be also used to optimize processing configuration values. 
In [previous section](https://github.com/librato/snap-plugin-lib-go/tree/ao-12231-tutorial/tutorial/03-concepts#configuration) "format" option was read during each collection.
We could read it only once during `Load()` and store in context. 

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

This approach will be used also in [Chapter 9](/tutorial/09-config/README.md).

----

* [Table of contents](/tutorial/README.md)
- Previous Chapter: [Testing](/tutorial/02-testing/README.md)
- Next Chapter: [Metrics - filters, definition, tags](/tutorial/04-metrics/README.md)