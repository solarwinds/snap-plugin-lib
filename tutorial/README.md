# Tutorial

## What's collector

Collector is a small application gathering information about any observed system. 
You can collect CPU utilization, database statistics, message broker queue's sizes or whatever you want in a form of metric.

## What's metric

Metric holds information about single measurement (it's usually pair of string and number)
**TODO**More

### Snap deamon vs snap collectors

In production environment plugins are controlled by snap:
1) Snap reads and forward configuration (credentials, requested metrics etc.) to plugin
2) Periodically (usally 60s) snap requests collection of metrics for different plugins 
3) Each plugin responds with own set of metrics. 

## Intro

This tutorial will teach you how to write custom collector plugin in Go language, able to work with snap deamon. 

We will start from a very simple example - building a minimal plugin and testing that it's working correctly. 
After you obtain basics, we will teach you how to write real, useful collector (gathering sysytem metrics) utilizing advanced features that plugin-go V2 has to offer:
- defining supported metrics
- build-in filtering mechanism
- defining custom configuration 
- store state between independent collection requests
- etc.

Each part of tutorial contain complete golang source code, which you can run and modify dependending on your needs.

## Version 2

Only version 2 of plugin-lib-go is covered in this tutorial. Examples and plugin catalog related to version 1 can be found in [Writing a Plugin](https://github.com/librato/snap-plugin-lib-go/tree/ao-12231-tutorial#writing-a-plugin) section.

# Content 

- Simple date/time collector:
  * [Introduction - Simple plugin](/tutorial/01-simple/README.md)
  * [Testing](/tutorial/02-testing/README.md)
  * [Configuration and state](/tutorial/03-concepts/README.md)
  * [Metrics - filters, definition, tags](/tutorial/04-metrics/README.md)
  * [Useful tools](/tutorial/05-tools/README.md)
- Process collector:
  * [Advanced plugin - overview](/tutorial/06-overview/README.md)
  * [Dynamic metrics](/tutorial/07-dynmetrics/README.md)
  * [Other features](/tutorial/08-collector/README.md)
- [FAQ](/tutorial/faq/README.md)