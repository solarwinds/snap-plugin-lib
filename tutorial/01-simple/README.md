# Introduction - Simple plugin

## Code

The simplest plugin gathering only 1 metric (every time with the same value) can written as follows.

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

Import section lists required dependencies:
- /v2/plugin - contains interfaces definition (ie. Collector) which we can implement
- /v2/runner - contains implementation of StartCollector() which is used in main function

### Collector code

```go
type simpleCollector struct{}

func (s simpleCollector) Collect(ctx plugin.Context) error {
    _ = ctx.AddMetric("/example/metric1", 10)
    return nil;
}
```

Here we defined collector code. 
In short it's an object of any type implementing Collect() method, which is the "heart" of plugin - saying what needs to be done. 
In our simple case we are gathering one metric "/example/metric1" containing value 10, by calling `ctx.AddMetric(metricName, value)` but in real case those values would vary in time depending on the state of observed system.

Following code will present date and time collector, which produce 5 metrics associated with day, month, hour, minute and second of time.

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

When we have access to time object, we create 5 metrics - each of which is associated with a separated value associated with time.
Also we have introducted metric groups: date and time (second position in metric name).
Metric form will be described later in details but in short metrics can contain at least two strings separated by "/". 
Usually the first one is plugin name and the last one is metric purpose.
Intermediate strings serves as groups (collector subfunctions) and simplifies filtering.

### Runner

```go
func main() {
    runner.StartCollector(&simpleCollector{}, "example", "1.0.0")
}
```

main() function will usually have the same implementation (with different parameters depending of plugin).
Runner takes care about establishing valid communication between snap daemon and plugin.
