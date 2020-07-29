# Snap Plugin Library

Snap Plugin Library helps developers writing plugins (like collectors and publishers) that are able to work in Snap environment:
* collector is a small application which role is to gather metrics from a monitored system (like CPU usage, db metrics, etc.)
* publisher is a small application which role is to send metrics to a specific backend 

The library natively supports Go language, although writing Python plugins is also supported via Cgo integration. There is a plan for supporting other programming languages as well.

Currenlty active and maintained version is v2.

## Example

Writing a simple plugin code in Go is very straightforward:

```go
package main

import (
    "github.com/librato/snap-plugin-lib-go/v2/plugin"
    "github.com/librato/snap-plugin-lib-go/v2/runner"
)

type myCollector struct {}

func (c *myCollector) Collect(ctx plugin.CollectContext) error {
    _ = ctx.AddMetric("/example/static/value", 34)
    return nil
}

func main() {
    runner.StartCollector(&myCollector{}, "example-collector", "v1.0.0")
}
```

The same functionality in Python:

```python
from snap_plugin_lib_py import BasePlugin, start_collector


class ExamplePlugin(BasePlugin):

    def collect(self, ctx):
        ctx.add_metric("/example/metric/value", 45)


if __name__ == '__main__':
    start_collector(ExamplePlugin("example", "0.0.1"))
```

More complicated examples can be found in:
* (Go) [Examples folder](examples/v2)
* (Go) [Tutorial](v2/tutorial/09-config/collector)
* (Python) [Example](v2/plugin-lib/snap-plugin-lib-example.py)

## Documentation



* [Documentation of Ver. 1](/v1/README.md)
* [Documentation of Ver. 2](/v2/README.md)