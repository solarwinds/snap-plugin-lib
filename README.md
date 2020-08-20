# Snap Plugin Library

Snap Plugin Library helps developers writing plugins (like collectors and publishers) that are able to work in Snap environment:
* a collector is a small application whose role is to gather metrics from a monitored system (like CPU usage, db metrics, etc.)
* a publisher is a small application whose role is to send metrics to a specific backend 

The library natively supports Go language, although writing Python plugins is also supported via Cgo integration. There is a plan for supporting other programming languages as well.

The currently active and maintained version is v2.

## Development setup

For a simple development setup you don't need any dependencies outside this repository.
Please refer to  [Testing](/v2/tutorial/02-testing/README.md) for more information.

For a complete development setup please refer to [AppOptics Knowledge Base](https://documentation.solarwinds.com/en/Success_Center/appoptics/Content/kb/host_infrastructure/host_agent.htm)

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
from swisnap_plugin_lib_py import BaseCollector, start_collector


class ExamplePlugin(BaseCollector):

    def collect(self, ctx):
        ctx.add_metric("/example/metric/value", 45)


if __name__ == '__main__':
    start_collector(ExamplePlugin("example", "0.0.1"))
```

More complicated examples can be found in:
* (Go) [Examples folder](examples/v2)
* (Go) [Tutorial](v2/tutorial/09-config/collector)
* (Python) [Collector Example](v2/plugin-lib/swisnap-collector-plugin-lib-example.py)
* (Python) [Publisher Example](v2/plugin-lib/swisnap-publisher-plugin-lib-example.py)

## Documentation

Plugin Library contains a comprehensive tutorial describing how to write custom plugins.

Covered topics:
- [Overview](/v2/README.md)
- Simple example - Date/Time collector:
  * [Introduction](/v2/tutorial/01-simple/README.md)
  * [Testing](/v2/tutorial/02-testing/README.md)
  * [Configuration and state](/v2/tutorial/03-concepts/README.md)
  * [Metrics - filters, definition, tags](/v2/tutorial/04-metrics/README.md)
  * [Useful tools](/v2/tutorial/05-tools/README.md)
- Advanced example - System collector:
  * [Overview](/v2/tutorial/06-overview/README.md)
  * [Gathering data (Proxy)](/v2/tutorial/07-proxy/README.md)
  * [Collector](/v2/tutorial/08-collector/README.md)
  * [Handle configuration](/v2/tutorial/09-config/README.md)
- [FAQ](/v2/tutorial/faq/README.md)

## Contributing Guide

### Issue Reporting Guidelines

* Always check if the problem wasn't already reported by other developers. 
* Please fill in [Issue submission form](https://github.com/librato/snap-plugin-lib-go/issues/new).

### Pull Request Guidelines

* General:
    * Make sure that the library is compiling without errors,
    * Make sure that *ALL* unit tests pass, 
    * Try to add the accompanying test case,
    * It's OK to have multiple small commits as you work on the PR - they will be squashed before merging,
    * Provide meaningful commit messages.

* For issues:
    * Add issue number to your PR title for a better release log.

* For a new development:
    * Provide a convincing reason to add this feature.

## Version 1

The repository contains a legacy version of Plugin Library (v1) which is no longer supported either by SolarWinds or Intel (which is originally forked from). New defects raised for v1 won't be fixed by maintainers.

Please use v2 for new development.

Links:
* [Documentation of Ver. 1](/v1/README.md)
* [Community plugins list](https://github.com/intelsdi-x/snap/blob/master/docs/PLUGIN_CATALOG.md) (no longer maintained)
