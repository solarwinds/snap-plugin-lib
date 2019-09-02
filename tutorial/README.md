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

* [Introduction - Simple plugin](/tutorial/01-simple/README.md)
* [Testing](/tutorial/02/README.md)
* [Basic concepts]((/tutorial/03/README.md))
* [Advanced plugin - overview](/tutorial/04/README.md)
* [Defining plugin metadata (metrics and configuration)](/tutorial/05/README.md)
* [Other features](/tutorial/06/README.md)
