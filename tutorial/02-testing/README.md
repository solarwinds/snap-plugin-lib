# Testing

## Compile plugin

To compile plugin simply go to folder containing written code in [Previous Chapter](/tutorial/01-simple) and execute command:
```
go build
```
which result in building executable file named as folder containing go code, or
```
go build -o binaryName
```
to produce binary with custom name.

Instead of compiling you own code, you can reuse tutorial example:
```bash
cd $GOPATH/src/github.com/librato/snap-plugin-lib-go/tutorial/02-testing
go build
``` 

> Note:
> Further commands will use binary name based on examples in tutorial folder (`01-simple`, `02-testing`, etc)
> For windows you should replace it with `01-simple.exe`, `02-testing.exe` etc. 

## Execution

### Smoke-test (No-argument provided)

The simplest way of validating plugin is run binary without any arguments.
```bash
./02-testing
```
In valid scenario collector should print out metadata information and after some time quit execution with error message:
```
{"GRPCVersion":"2.0.0","Plugin":{"Name":"example","Version":"1.0.0"},"GRPC":{"IP":"127.0.0.1","Port":56302},"Profiling":{"Enabled":false,"Location":""},"Stats":{"Enabled":false,"IP":"","Port":0}}
time="2019-09-03T15:12:10+02:00" level=warning msg="Ping timeout occurred" layer=lib max=3 missed=1 module=plugin-rpc
time="2019-09-03T15:12:13+02:00" level=warning msg="Ping timeout occurred" layer=lib max=3 missed=2 module=plugin-rpc
time="2019-09-03T15:12:16+02:00" level=warning msg="Ping timeout occurred" layer=lib max=3 missed=3 module=plugin-rpc
time="2019-09-03T15:12:16+02:00" level=error msg="Major error occurred - plugin will be shut down" error="ping message missed 3 times (timeout: 3s)" layer=lib module=plugin-rpc
```

Observed scenario is valid. Executing without any arguments means that plugin expects to be controller by snap daemon.
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

> Other useful flags, like: `-plugin-config`, `-plugin-filter` and `*stats*` related will be discussed later. **TODO**link
  
#### Debug-mode

To execute collection run binary in debug-mode.

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

Output will contains 3 sets of metrics calculated every 5 seconds:
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
