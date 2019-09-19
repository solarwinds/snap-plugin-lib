# Testing

## Compile plugin

To build executable file simply go to folder containing written code (see: [Previous Chapter](/tutorial/01-simple)) and execute command:
```
go build
```
Binary name is depending on folder in which `main` file is located. 

To force own custom name of binary, use:
```
go build -o binaryName
```

 Instead of compiling your own code, you can reuse tutorial example(s):
```bash
cd $GOPATH/src/github.com/librato/snap-plugin-lib-go/tutorial/02-testing
go build
```

> Note:
> Further commands will use binary name based on examples in tutorial folders (`01-simple`, `02-testing`, etc)
> For windows you probably should replace it with `01-simple.exe`, `02-testing.exe` etc. 

## Execution

### Smoke-test (No-argument provided)

The simplest way of validating plugin is run binary without any arguments.
```bash
./02-testing
```
In valid scenario collector should print out metadata information and quit after some time with error message:
```
{"GRPCVersion":"2.0.0","Plugin":{"Name":"example","Version":"1.0.0"},"GRPC":{"IP":"127.0.0.1","Port":56302},"Profiling":{"Enabled":false,"Location":""},"Stats":{"Enabled":false,"IP":"","Port":0}}
time="2019-09-03T15:12:10+02:00" level=warning msg="Ping timeout occurred" layer=lib max=3 missed=1 module=plugin-rpc
time="2019-09-03T15:12:13+02:00" level=warning msg="Ping timeout occurred" layer=lib max=3 missed=2 module=plugin-rpc
time="2019-09-03T15:12:16+02:00" level=warning msg="Ping timeout occurred" layer=lib max=3 missed=3 module=plugin-rpc
time="2019-09-03T15:12:16+02:00" level=error msg="Major error occurred - plugin will be shut down" error="ping message missed 3 times (timeout: 3s)" layer=lib module=plugin-rpc
```

Observed scenario is valid. Executing without any arguments means that plugin expects to be controlled by snap daemon.
Since snap is not running, plugin quits.

### Validating metric

Framework offers 3 way of validating metric collections:
1) running plugin in debug-mode
2) running plugin with snap-mock (minimal version of snap daemon)
3) running plugin with snap daemon (production)

This chapter will cover methods 1 and 2. 

#### Command-line arguments

You can execute plugin binary with some additional arguments. To list them you can simply call:
```
./02-testing -h
```

From the perspective of plugin developer the most important options are:

|Flag                     | Description                                                                     |
|-------------------------|---------------------------------------------------------------------------------|
| -debug-mode             | Run plugin in debug mode (no snap daemon required)                              |
| -debug-collect-counts   | Number of collect requests executed in debug mode (0 - infinitely) (default 1)  |
| -debug-collect-interval | Interval between consecutive collect requests (default 5s)                      |
| -log-level              |  Minimal level of logged messages (you should use either `debug` or `trace`)    | 

> Other useful flags, like: `-plugin-config`, `-plugin-filter` and `*stats*` related will be discussed later. (see: [Stats](https://github.com/librato/snap-plugin-lib-go/tree/ao-12231-tutorial/tutorial/05-tools#stats-server))
  
#### Debug-mode

To execute collection, run binary in debug-mode.

```bash
./02-testing -debug-mode=1
```

It should generate output, similar to:
```
Gathered metrics (length=5):
example.date.day 3 {map[]}
example.date.month September {map[]}
example.time.hour 15 {map[]}
example.time.minute 58 {map[]}
example.time.second 53 {map[]}
```

Values of metrics depends of course on the current date and time and will be different on your testing environment, but the metric names should be the same.

You can request several collects, using other flags:
```bash
./02-testing -debug-mode=1 -debug-collect-counts=3 -debug-collect-interval=5s
```

Output will contain 3 sets of metrics calculated every 5 seconds:
```
Gathered metrics (length=5):
example.date.day 3 {map[]}
example.date.month September {map[]}
example.time.hour 16 {map[]}
example.time.minute 5 {map[]}
example.time.second 34 {map[]}

Gathered metrics (length=5):
example.date.day 3 {map[]}
example.date.month September {map[]}
example.time.hour 16 {map[]}
example.time.minute 5 {map[]}
example.time.second 39 {map[]}

Gathered metrics (length=5):
example.date.day 3 {map[]}
example.date.month September {map[]}
example.time.hour 16 {map[]}
example.time.minute 5 {map[]}
example.time.second 44 {map[]}
```

#### Running plugin with snap-mock

Debug mode should be sufficient in majority of cases, although there will be scenarios when running with snap-mock enables additional testing capability.

To work correctly plugin should be run with several flags as follows:
```bash
./02-testing -grpc-ping-max-missed=0 -grpc-port=50123
```

By providing `-grpc-ping-max-missed=0` plugin won't exit after not receiving 3 pings from snap (or its equivalent, like snap-mock).

Now, in other console, you should locate snap-mock, compile it and execute:

```bash
cd $GOPATH/src/github.com/librato/snap-plugin-lib-go/v2/snap-mock
go build
./snap-mock -plugin-port=50123 -max-collect-requests=3 -collect-interval=5s -send-kill=1
```

> Make sure that `-grpc-port` (provided for plugin) and `-plugin-port` (provided for snap-mock) are the same.

> `-send-kill` flag causes that both: plugin and snap-mock will complete execution after working.

Output of snap-mock:
```

Received 5 metric(s)
 example.date.day v_int64:3  [map[]]
 example.date.month v_int64:9  [map[]]
 example.time.hour v_int64:16  [map[]]
 example.time.minute v_int64:39  [map[]]
 example.time.second v_int64:30  [map[]]

Received 5 metric(s)
 example.date.day v_int64:3  [map[]]
 example.date.month v_int64:9  [map[]]
 example.time.hour v_int64:16  [map[]]
 example.time.minute v_int64:39  [map[]]
 example.time.second v_int64:35  [map[]]

Received 5 metric(s)
 example.date.day v_int64:3  [map[]]
 example.date.month v_int64:9  [map[]]
 example.time.hour v_int64:16  [map[]]
 example.time.minute v_int64:39  [map[]]
 example.time.second v_int64:40  [map[]]
```

Output of plugin:
```
{"GRPCVersion":"2.0.0","Plugin":{"Name":"example","Version":"1.0.0"},"GRPC":{"IP":"127.0.0.1","Port":50123},"Profiling":{"Enabled":false,"Location":""},"Stats":{"Enabled":false,"IP":"","Port":0}}
time="2019-09-03T16:39:30+02:00" level=trace msg="GRPC Ping() received" layer=lib module=plugin-rpc
time="2019-09-03T16:39:30+02:00" level=trace msg="GRPC Load() received" layer=lib module=plugin-rpc
time="2019-09-03T16:39:30+02:00" level=trace msg="GRPC Collect() received" layer=lib module=plugin-rpc
time="2019-09-03T16:39:30+02:00" level=debug msg="Collect completed" elapsed=0s layer=lib module=plugin-proxy
time="2019-09-03T16:39:30+02:00" level=debug msg="metrics chunk has been sent to snap" layer=lib len=5 module=plugin-rpc
time="2019-09-03T16:39:32+02:00" level=trace msg="GRPC Ping() received" layer=lib module=plugin-rpc
time="2019-09-03T16:39:34+02:00" level=trace msg="GRPC Ping() received" layer=lib module=plugin-rpc
time="2019-09-03T16:39:35+02:00" level=trace msg="GRPC Collect() received" layer=lib module=plugin-rpc
time="2019-09-03T16:39:35+02:00" level=debug msg="Collect completed" elapsed=0s layer=lib module=plugin-proxy
time="2019-09-03T16:39:35+02:00" level=debug msg="metrics chunk has been sent to snap" layer=lib len=5 module=plugin-rpc
time="2019-09-03T16:39:36+02:00" level=trace msg="GRPC Ping() received" layer=lib module=plugin-rpc
time="2019-09-03T16:39:38+02:00" level=trace msg="GRPC Ping() received" layer=lib module=plugin-rpc
time="2019-09-03T16:39:40+02:00" level=trace msg="GRPC Ping() received" layer=lib module=plugin-rpc
time="2019-09-03T16:39:40+02:00" level=trace msg="GRPC Collect() received" layer=lib module=plugin-rpc
time="2019-09-03T16:39:40+02:00" level=debug msg="Collect completed" elapsed=0s layer=lib module=plugin-proxy
time="2019-09-03T16:39:40+02:00" level=debug msg="metrics chunk has been sent to snap" layer=lib len=5 module=plugin-rpc
time="2019-09-03T16:39:41+02:00" level=trace msg="GRPC Unload() received" layer=lib module=plugin-rpc
time="2019-09-03T16:39:41+02:00" level=trace msg="GRPC Kill() received" layer=lib module=plugin-rpc
time="2019-09-03T16:39:41+02:00" level=debug msg="GRPC server stopped gracefully" layer=lib module=plugin-rpc
```

What we can see from debug logs:
- GRPC `Collect()` was called 3 times, every 5 seconds 
- Plugin receives GRPC Ping periodically (`"GRPC Ping() received"`)
- GRPC `Kill()` is called at the end causing plugin to quit

There are also other entries:
- first line is metadata in JSON format used by snap. 
- GRPC `Load()` and GRPC `Unload()` calls will be covered later. In short, collection can be requested by several independent tasks and `Load()` allows to handle independent configuration for a single task. 
- `metrics chunk has been sent to snap` means that metrics are sent in portions to snap (generally to avoid some internal limits of GRPC when plugin want to send large portions of metrics)

As you can see we achieved almost the adequate testing results as in the debug mode (3 request every 5 seconds).
The difference is that with current approach we trigger collection via GRPC API (the same way that snap does), so it's the same way, which plugin would be triggered in production environment. 
Debug-mode calls defined methods internally (without utilizing GRPC communication), but if want to validate collection logic it suffice.
Snap-mock will be useful to observe how plugin reacts with different tasks (several configurations requested at the same time).

----

* [Table of contents](/tutorial/README.md)
- Previous Chapter: [Introduction - Simple plugin](/tutorial/01-simple/README.md)
- Next Chapter: [Basic concepts - Configuration and state](/tutorial/03-concepts/README.md)