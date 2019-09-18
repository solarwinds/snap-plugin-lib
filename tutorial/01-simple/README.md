# Introduction - Simple plugin

## Code

The simplest plugin gathering only 1 metric (every time with the same value) can be written as follows.

```go
package main

import (
    "github.com/librato/snap-plugin-lib-go/v2/plugin"
    "github.com/librato/snap-plugin-lib-go/v2/runner"
)

type simpleCollector struct{}

func (s simpleCollector) Collect(ctx plugin.Context) error {
    _ = ctx.AddMetric("/example/metric1", 10)
    return nil;
}

func main() {
    runner.StartCollector(&simpleCollector{}, "example", "1.0.0")
}
```

### Headers

```go
package main
```

This instruction will tell go compiler that should build executable file from this package.

```go
import (
    "github.com/librato/snap-plugin-lib-go/v2/plugin"
    "github.com/librato/snap-plugin-lib-go/v2/runner"
)
```

Import section enumerates required dependencies:
- `/v2/plugin` - contains interfaces definition (ie. `Collector`) which we should implement
- `/v2/runner` - contains implementation of `StartCollector()` which is used in main function

### Collector code

```go
type simpleCollector struct{}

func (s simpleCollector) Collect(ctx plugin.Context) error {
    _ = ctx.AddMetric("/example/metric1", 10)
    return nil;
}
```

Above code shows the simplest implementation of collector.
Method `Collect()` needs to be always provided - it's the "heart" of plugin, saying what needs to be done.
We are gathering one metric `/example/metric1` containing value `10`, by calling `ctx.AddMetric(metricName, value)`.
In real application those values would vary in time (depending on the state of observed system).

Following code is a little more complicated.
It gather current date and time, producing 5 metrics associated with current day, month, hour, minute and second of time.

> Although its practicality is doubtful, it will be suffient to show different set of plugin-lib v2 features.
> If you want to learn straightaway how to write useful collector, visit [Chapter 6](/tutorial/06-overview/README.md) of this tutorial. 

```go
type simpleCollector struct{}

func (s simpleCollector) Collect(ctx plugin.Context) error {
    // Collect data
    t := time.Now()

    // Convert data to metric form
    _ = ctx.AddMetric("/example/date/day", t.Day())
    _ = ctx.AddMetric("/example/date/month", t.Month())
    _ = ctx.AddMetric("/example/time/hour", t.Hour())
    _ = ctx.AddMetric("/example/time/minute", t.Minute())
    _ = ctx.AddMetric("/example/time/second", t.Second())

    return nil;
}
```

At the beginning we are simulating data collection by obtaining current system time. 
In real plugins those code would be replaced with some REST request, file reading, SQL query etc.

When we have access to `time` object, we create 5 metrics - each of which is associated with a separated value.
Also we have introduced metric groups: date and time (second position in metric name).
Metric form will be described later in details but in short metrics should contain at least two strings separated by "/". 
Usually the first one is plugin name and the last one is metric purpose.
Intermediate strings serves as groups (collector sub-functions) and simplifies filtering.

### Runner

```go
func main() {
    runner.StartCollector(&simpleCollector{}, "example", "1.0.0")
}
```

main() function will usually have the same implementation (with different parameters depending of plugin).
Runner takes care about establishing valid communication between snap daemon and plugin.
User can focus only on the collector implementation details.
