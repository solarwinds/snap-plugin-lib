# FAQ (Frequently Asked Questions)

##### Do I have to implement all 4 methods (`Collect`, `Load`, `Unload`, `DefinePlugin`) in order to enable plugin to work in snap environment?

No. In simple cases, it's required to implement only `Collect` method.

However, implementing all of the methods gives the following benefits:
- `Collect` takes care about gathering measurements
- `Load` (in pair with `Unload`) is used for processing configuration, creating shared objects etc. and enables the plugin to handle more than one task simultaneously.
- `DefinePlugin` helps in plugin maintenance by: providing metadata of what collector may produce.

----

##### Is there a template for a minimal collector?

You can take a look at collector's code from [Chapter 1](/tutorial/01-simple/README.md) (see: [main.go](/tutorial/01-simple/main.go))

----

##### Should I handle error value from `ctx.AddMetric()` and `ctx.AddMetricWithTag()`?

Generally, no. 
Notice that error may be returned in following situations:
- you are trying to add metric not aligned with metrics set defined in `DefinePlugin` (developer error)
- you are trying to add metric with namespace containing disallowed characters (developer error) 
- you are trying to add metric which is filtered out (due to using `-plugin-filter` flag or defining subset of requested metrics) - actually not a logical error 

If you find yourself in a situation when you can't find why metric is not present in a result set you can check the value returned from `AddMetric` in debug mode (or temporarily printing it out).

----

##### How I can test if my plugin is working correctly?

There are several ways to achieve that at different levels:
- you can write unit tests - library provides `mock.Context` if you want to test `Collect` method (take a loot at `./collector/09-config/collector/*_test.go` to see how it can be done)
- you can run a plugin in [Debug mode](https://github.com/librato/snap-plugin-lib-go/tree/ao-12231-tutorial_ch89/tutorial/02-testing#debug-mode)
- you can run a plugin with [Snap-mock](https://github.com/librato/snap-plugin-lib-go/tree/ao-12231-tutorial_ch89/tutorial/02-testing#running-plugin-with-snap-mock)
- you can run a plugin in the Snap environment

----

* [Table of contents](/tutorial/README.md)
- Previous Chapter: [Handle configuration](/tutorial/09-config/README.md)
