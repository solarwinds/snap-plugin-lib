# FAQ (Frequently Asked Questions)

##### Do I have to implement all 4 methods in order (`Collect`, `Load`, `Unload`, `DifinePlugin`) to enable plugin working in snap environment.

Simple plugin may require to implement only `Collect` method.

For a "production" plugin implementing all 4 methods are beneficial:
- `Collect` takes care about gathering measurement
- `Load` (in pair with `Unload`) is used for processing configuration, creating shared objects etc.
- `DefinePlugin` helps in maintenance of plugin: providing metadata of what collector may produce.

----

##### Is there a template for a minimal collector

You can take a look at collector's code from [Chapter 1](/tutorial/01-simple/README.md) (see: [main.go](/tutorial/01-simple/main.go))

----

##### Should I have to handle error value from ctx.AddMetric() and ctx.AddMetricWithTag() ?

Generally, no. 
Notice that error may be return in following situations.
- you are trying to add metric not aligned with set defined in `DefinePlugin` (developer error)
- you are trying to add metric which namespace contains disallowed characters (developer error) 
- you are trying to add metric which is filtered (due to using `-plugin-filter` flag or defining subset of requested metrics) - actually not a logical error 

If you find yourself in situation when you can't find why metric is not present in a result set you can check the value returned from `AddMetric` in debug mode (or temporary printing out return value).

----

##### How I can test that my plugin is working correctly?

There are several way to achieve that at different level:
- you can write unit tests - library provides `mock.Context` if you want to test `Collect` method (take a loot at `./collector/09-config/collector/*_test.go` to see how it can be done)
- you can run a plugin in [Debug mode](https://github.com/librato/snap-plugin-lib-go/tree/ao-12231-tutorial_ch89/tutorial/02-testing#debug-mode)
- you can run a plugin with [Snap-mock](https://github.com/librato/snap-plugin-lib-go/tree/ao-12231-tutorial_ch89/tutorial/02-testing#running-plugin-with-snap-mock)
- you can run a plugin in snap environment

----

* [Table of contents](/tutorial/README.md)
- Previous Chapter: [Handle configuration](/tutorial/09-config/README.md)
