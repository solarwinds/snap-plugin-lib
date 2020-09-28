# FAQ (Frequently Asked Questions)

##### Do I have to implement all 4 methods (`Collect`, `Load`, `Unload`, `DefinePlugin`) in order to enable plugin to work in snap environment?

No. In simple cases only `Collect` method is required.

However, implementing all of the methods gives the following benefits:
- `Collect` takes care of gathering measurements
- `Load` (in pair with `Unload`) is used for processing configuration, creating shared objects etc. and enables the plugin to handle more than one task simultaneously.
- `DefinePlugin` helps in plugin maintenance by: providing metadata of what collector may produce.

----

##### Is there a template for a minimal collector?

You can take a look at collector's code from [Chapter 1](/v2/tutorial/01-simple/README.md) (see: [main.go](/v2/tutorial/01-simple/main.go))

----

##### Should I handle error value from `ctx.AddMetric()`?

Generally, no. 
Notice that error may be returned in the following situations:
- you are trying to add metric not aligned with metrics set defined in `DefinePlugin` (developer error)
- you are trying to add metric with namespace containing disallowed characters (developer error) 
- you are trying to add metric which is filtered out (due to using `-plugin-filter` flag or defining subset of requested metrics) - actually not a logical error 

If you find yourself in a situation when you can't find why metric is not present in a result set you can check the value returned from `AddMetric` in debug mode (or temporarily printing it out).

----

##### How can I test if my plugin is working correctly?

There are several ways to achieve that at different levels:
- you can write unit tests - library provides `mock.Context` if you want to test `Collect` method (take a look at `./collector/09-config/collector/*_test.go` to see how it can be achieved)
- you can run a plugin in [Debug mode](/v2/tutorial/02-testing#debug-mode)
- you can run a plugin with [Snap-mock](/v2/tutorial/02-testing#running-plugin-with-snap-mock)
- you can run a plugin in the Snap environment

----

* [Table of contents](/v2/tutorial/README.md)
- Previous Chapter: [Writing plugins in Python and C#](/v2/tutorial/other-languages/README.md)
