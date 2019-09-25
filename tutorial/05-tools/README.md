# Useful tools

## Example configuration

Plugin creator can provide example configuration to simplify task-file creation.
Currently supported format is YAML (due to ability to add comments).

Example:
```go
 func (s simpleCollector) PluginDefinition(def plugin.CollectorDefinition) error {
 	cfg := `
 format: short # format of hour (short 0-12, long 0-24)
 options:
   - zone: UTC # time zone
 `
 	_ = def.DefineExampleConfig(cfg)
    
    // ...
}
```

## Printing example task-file

User can print/create default task based on metadata provided by the plugin creator.

```bash
./05-tools -print-example-task
```

Output:
```yaml
version: 2
schedule:
    type: simple
    interval: 60s
plugins:
  - name: examplecollector
    metrics:
      - /example/date/day
      - /example/date/month
      - /example/time/hour
      - /example/time/minute
      - /example/time/second
    config:
        format: short # format of hour (short 0-12, long 0-24)
        options:
          - zone: UTC # time zone
    tags:
        /example:
            plugin_tag: tag
    publish:
        config:
            period: 60
            floor_second: 60
```

## Stats server

When plugin is controlled by snap-mock, user can gather several statistics:
- options used to start plugin,
- number of tasks currently loaded,
- number of collect request,
- processing time of last collect request,
- total number of metrics
- etc.

To enable this feature simply run plugin with the following arguments:
```bash
./05-tools -grpc-port=50123 -log-level=debug -enable-stats -enable-stats-server -stats-port=8080
```

To simulate two independent tasks, call snap-mock twice:
```bash
./snap-mock -plugin-port=50123 -max-collect-requests=12 -collect-interval=5s -task-id=1 &
./snap-mock -plugin-port=50123 -max-collect-requests=12 -collect-interval=5s -task-id=2 &
```

To see output of stats server, open a browser on address http://127.0.0.1:8080/stats.

Example output:
```json
{
    "Plugin info": {
        "Name": "example",
        "Version": "1.0.0",
        "Command-line options": "-grpc-port=50123 -log-level=debug -enable-stats -enable-stats-server -stats-port=8080",
        "Options": {
            "PluginIp": "127.0.0.1",
            "GrpcPort": 50123,
            "GrpcPingTimeout": 3000000000,
            "GrpcPingMaxMissed": 3,
            "LogLevel": "debug",
            "EnableProfiling": false,
            "EnableStats": true,
            "EnableStatsServer": true,
            "StatsPort": 8080
        },
        "Started": {
            "Time": "Sep  9 09:32:48.848641",
            "Ago": "47.632944s"
        }
    },
    "Tasks summary": {
        "Counters": {
            "Currently active tasks": 2,
            "Total active tasks": 2,
            "Total collect requests": 16
        },
        "Processing times": {
            "Total": "5.6107ms",
            "Average": "350.668µs",
            "Maximum": "1.0125ms"
        }
    },
    "Task details": {
        "1": {
            "Configuration": {},
            "Requested metrics (filters)": [],
            "Counters": {
                "Collect requests": 9,
                "Total metrics": 54,
                "Average metrics / Collect": 6
            },
            "Loaded": {
                "Time": "Sep  9 09:32:54.126860",
                "Ago": "42.3549103s"
            },
            "Processing times": {
                "Total": "3.3691ms",
                "Average": "374.344µs",
                "Maximum": "1.0125ms"
            },
            "Last measurement": {
                "Occurred": {
                    "Time": "Sep  9 09:33:34.650584",
                    "Ago": "1.8311868s"
                },
                "Collected metrics": 6
            }
        },
        "2": {
            "Configuration": {},
            "Requested metrics (filters)": [],
            "Counters": {
                "Collect requests": 7,
                "Total metrics": 42,
                "Average metrics / Collect": 6
            },
            "Loaded": {
                "Time": "Sep  9 09:33:02.914434",
                "Ago": "33.567337s"
            },
            "Processing times": {
                "Total": "2.2416ms",
                "Average": "320.228µs",
                "Maximum": "815.2µs"
            },
            "Last measurement": {
                "Occurred": {
                    "Time": "Sep  9 09:33:33.433020",
                    "Ago": "3.0487507s"
                },
                "Collected metrics": 6
            }
        }
    }
}
```

## Profiling

When plugin is controlled by snap-mock user can run profiling server in the background by executing:
```bash
./05-tools -grpc-port=50123 -log-level=debug -enable-profiling -pprof-port=8081
```

To access pprof Web-GUI browse http://127.0.0.1:8081/debug/pprof/

----

* [Table of contents](/tutorial/README.md)
- Previous Chapter: [Metrics - filters, definition, tags](/tutorial/04-metrics/README.md)
- Next Chapter: [Advanced Plugin - Introduction](/tutorial/06-overview/README.md)
