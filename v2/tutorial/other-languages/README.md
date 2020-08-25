# Writing plugins in Python and C#

## Overview

The library provides bindings to enable developers to write custom collectors (and publishers) in languages other that Go. 
Currently supported languages are listed in the table below.

Language | Location                                   | Collector API | Publisher API
---------|------------------------------------------- |---------------|--------------
Python   | [/v2/bindings/python](/v2/bindings/python) | Yes           | Yes
C#       | [/v2/bindings/csharp](/v2/bindings/csharp) | Yes           | No

The bindings have API as similar as possible to one provided by Go. 
Thereby, the majority of the tutorial for Go language should be easily translated to the language of your choice.
Differences (very often related to language syntax) are presented in the further chapters.
Each language binding contains at least one example, based on which you can learn how to properly use the API.

## Minimal plugin definitions

### Python

```python
from swisnap_plugin_lib_py import BaseCollector, start_collector


class ExampleCollectorPlugin(BaseCollector):
    def collect(self, ctx):
        # Collect metrics 
        pass


if __name__ == "__main__":
    collector = ExampleCollectorPlugin("collector-example", "0.0.1")
    start_collector(collector)

```

### C#

```csharp
using System;
using SnapPluginLib;

namespace CollectorExample
{
    public class CollectorExample : CollectorPluginBase
    {
        public CollectorExample(string name, Version version) : base(name, version)
        {
        }

        public override void Collect(ICollectContext ctx)
        {
            // Collect metrics 
        }

        static void Main(string[] args)
        {
            Runner.StartCollector(new CollectorExample("collector-example", new Version(0, 0, 1)));
        }
    }
}
```

## API Support

### Plugin definitions

Go                                    | Python | C#
--------------------------------------|--------|-----
``def.DefineTasksPerInstanceLimit()`` | Yes    | Yes
``def.DefineInstancesLimit()``        | Yes    | Yes
``def.DefineMetric()``                | Yes    | Yes
``def.DefineGroup())``                | Yes    | Yes
``def.DefineExampleConfig()``         | Yes    | Yes


### Context

Go                    | Python  | C#
----------------------|---------|--------
``ctx.ConfigValue()`` | Yes     | Yes
``ctx.ConfigKeys()``  | Yes     | Yes
``ctx.RawConfig()``   | Yes     | Yes
``ctx.Store()``       | Yes     | Yes
``ctx.Load()``        | Yes [(1)](/v2/tutorial/other-languages#1) | Yes [(2)](/v2/tutorial/other-languages#2)
``ctx.LoadTo()``      | No      | No 
``ctx.AddWarning()``  | Yes     | Yes
``ctx.IsDone()``      | Yes     | Yes
``ctx.Done()``        | No      | No
``ctx.Logger()``      | Yes [(3)](/v2/tutorial/other-languages#3) | Yes [(3)](/v2/tutorial/other-languages#3)

#### **(1)**
In Python ``Load()`` throws an exception when the object with a given key wasn't loaded.

#### **(2)**
In C# ``Load()`` is a template method. 
The developer provides a concrete type of object he/she expects.
In case the object wasn't found, an exception is thrown.

Example:
```csharp
ctx.Load<Dictionary<string, int>>("stored_object");
```

#### **(3)** 
In Go language ``Logger()`` returns reference to ``Logger`` object.
Thereby setting log level and structured information is provided by ``Logger`` API.
In Python and C# ``Logger()`` method contains that information as parameters.

Go:
```go
log.WithFields(logrus.Fields{"language": "Go"}).Info("Log message")
```

Python:
```python
ctx.log(LOGLEVEL_INFO, "Log message from Python", {"language": "python"})
```

C#:
```csharp
ctx.Log(LogLevel.Info, "Log message from C#", new Dictionary<string, string>
{
    {"language", "c#"}, 
});
```

### Collector Context (superset of Context)

Go                            | Python    | C#
------------------------------|-----------|---------
``ctx.AddMetric()``           | Yes [(4)](/v2/tutorial/other-languages#4) | Yes [(5)](/v2/tutorial/other-languages#5)
``ctx.AlwaysApply()``         | Yes [(6)](/v2/tutorial/other-languages#6) | Yes [(6)](/v2/tutorial/other-languages#6)
``ctx.DismissAllModifiers()`` | Yes       | Yes
``ctx.ShouldProcess()``       | Yes       | Yes
``ctx.RequestedMetrics()``    | Yes       | Yes

#### **(4)** 
In Python modifiers are provided via named arguments.

```python
ctx.add_metric("/example/group/metric1", 30, tags={"os": "windows"}, description="custom description", unit="custom unit")
```

#### **(5)** 
In C# modifiers are provided similarly to Go. Example:

```csharp
ctx.AddMetric("/example/group/metric1", 12.4,
    Modifiers.Tags(new Dictionary<string, string>
    {
        {"os", "windows"}
    }),
    Modifiers.Description("custom description")
);
```

#### **(6)** 
In Python and C# ``AlwaysApply()`` doesn't return an object (used to dismiss only given modification).

## Manual compilation of CGo dependency

In order to manually build library required by bindings you need to execute the following command in ``v2/bindings`` (requires ``gcc`` installed on the system):

On Windows:
```bash
go build --buildmode=c-shared -o swisnap-plugin-lib.dll
```

On Linux:
```bash
go build --buildmode=c-shared -o swisnap-plugin-lib.so
```


----

* [Table of contents](/v2/README.md)
- Previous Chapter: [Handle configuration](/v2/tutorial/09-config/README.md)
- Next Chapter: [FAQ](/v2/tutorial/faq/README.md)