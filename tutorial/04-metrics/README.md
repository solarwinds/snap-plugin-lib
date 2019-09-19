# Metrics

## Filtering metrics

So far we have written collector providing 6 metrics:
```
/example/date/day
/example/date/month
/example/time/hour
/example/time/minute
/example/time/second
/example/count/running
```

Every time collector is running, all 6 metrics are being gathered.
That might not always be the case.
We might want to collect only subsets of all possible measurements.
Good news - all filtering is done internally by plugin so there is no need to change any line of code.

> In production environment requesting all metrics may cause huge impact on system if plugin requires advanced processing.
> You should always benchmark you plugin in different situations.

To filter metrics we provide additional parameter during execution (in production it would be entries in task file **TODO**link**):
```bash
./04-metrics -debug-mode=1 -debug-collect-counts=3 -debug-collect-interval=5s -plugin-filter="/example/date/*"
```

Now the output is restricted only to date metrics:
```
Gathered metrics (length=2):
example.date.day 9 {map[]}
example.date.month 9 {map[]}

Gathered metrics (length=2):
example.date.day 9 {map[]}
example.date.month 9 {map[]}

Gathered metrics (length=2):
example.date.day 9 {map[]}
example.date.month 9 {map[]}
```

Other examples of filters:
- `/example/*/day` - return only day metrics
- `/example/count/running` - return only running metric
- **TODO** regexp - command line issue fix

You can combine filters using `;`, ie.
```
$ ./04-metrics -debug-mode=1 -debug-collect-counts=1 -debug-collect-interval=5s -plugin-filter="/example/date/*;/example/count/running"
```

Output:
```
Gathered metrics (length=3):
example.date.day 9 {map[]}
example.date.month 9 {map[]}
example.count.running 0 {map[]}
```

## Defining metrics 

Plugin creator can add some useful metadata, like list of supported metrics.
To do so, simply define a new method:

```
func (s simpleCollector) PluginDefinition(def plugin.CollectorDefinition) error {
	def.DefineMetric("/example/date/day", "", true, "Current day")
	def.DefineMetric("/example/date/month", "", true, "Current month")
	def.DefineMetric("/example/time/hour", "h", true, "Current hour")
	def.DefineMetric("/example/time/minute", "m", true, "Current minute")
	def.DefineMetric("/example/time/second", "s", true, "Current second")
	def.DefineMetric("/example/count/running", "s", false, "Time since task was loaded")

	return nil
}
```

`DefineMetric()` has four arguments:
- metric name,
- unit (used i.e. by AppOptics front end),
- indication if metric is default (see: [Example task](https://github.com/librato/snap-plugin-lib-go/tree/ao-12231-tutorial/tutorial/05-tools#printing-example-task-file))
- metric description (used by AppOptics front-end)

There are two major advantages of providing this list:
1. user can obtain accurate default task (see: [Example task](https://github.com/librato/snap-plugin-lib-go/tree/ao-12231-tutorial/tutorial/05-tools#printing-example-task-file))
2. additional validation when metrics are calculated - user can't add metric which was not defined

Try to add following code at the end of `Collect()` metric:

```go
func (s simpleCollector) Collect(ctx plugin.Context) error {
    // ...

	err := ctx.AddMetric("/example/other/undefined", 10)
	fmt.Printf("%v\n", err)

	return nil
}
```

When executing the code (by requesting metrics) 
```bash
./04-metrics -debug-mode=1 -debug-collect-counts=1 -debug-collect-interval=5s
```

you will have the information that new metric could not have been added (due to not matching definition)
```
couldn't match metric with plugin definition: /example/other/undefined
Gathered metrics (length=6):
example.date.day 9 {map[]}
example.date.month 9 {map[]}
example.time.hour 7 {map[]}
example.time.minute 54 {map[]}
example.time.second 25 {map[]}
example.count.running 0 {map[]}
```

## Tags

Metric can hold additional information apart from its name and value. 
Tags are pairs of strings associated with metric.

We can add information about weekday (ie. Monday) to existing metric `/example/date/day` by calling:
```go
_ = ctx.AddMetricWithTags("/example/date/day", t.Day(), map[string]string{"weekday": t.Weekday().String()})
```
instead of:
```go
_ = ctx.AddMetric("/example/date/day", t.Day()
```

After requesting metrics:
```bash
$ ./04-metrics -debug-mode=1 -debug-collect-counts=1 -debug-collect-interval=5s
```

We got additional information:
```
Gathered metrics (length=6):
example.date.day 9 {map[weekday:Monday]}
example.date.month 9 {map[]}
example.time.hour 8 {map[]}
example.time.minute 17 {map[]}
example.time.second 2 {map[]}
example.count.running 0 {map[]}
```

----

* [Table of contents](/tutorial/README.md)
- Previous Chapter: [Basic concepts - Configuration and state](/tutorial/03-concepts/README.md)
- Next Chapter: [Useful tools](/tutorial/05-tools/README.md)

