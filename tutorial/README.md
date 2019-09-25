# Tutorial

## What's collector

Collector is a small application gathering information about observed system. 
You can collect CPU utilization, database statistics, message broker queue's sizes or whatever you want in a form of metric.

## What's metric

Metric holds information about single measurement. It consists of:
- name: (ie. `/system/total/cpu`)
- associated value (ie. `10`)
- unit (ie. `%`, `s.`)
- description 
- tags (additional textual data)
- time of measurements

### Snap daemon vs snap collectors

When used in production, plugins are controlled by snap daemon.
Simplified algorithm may be described as follows:
1) Snap reads and forward configuration (credentials, requested metrics etc.) to plugin
2) Periodically (typically 60s) snap requests collection of metrics for different plugins 
3) Each plugin responds with its own set of metrics

## Intro

This tutorial will teach you how to write a custom collector plugin in Go language, which in turn can be used in your production environment. 

We will start with a very simple example - building a minimal plugin and test that it's working correctly. 
After you obtain the basics, we will instruct you how to write real, useful collector (gathering system metrics) utilizing advanced features that plugin-go V2 has to offer:
- defining supported metrics
- build-in filtering mechanism
- defining custom configuration 
- store state between independent collection requests
- etc.

Each part of the tutorial contain complete golang source code, which you can run and modify depending on your needs.

## Version 2

Only version 2 of plugin-lib-go is covered in this tutorial. Examples and plugin catalog related to version 1 can be found in the [Writing a Plugin](https://github.com/librato/snap-plugin-lib-go/tree/ao-12231-tutorial#writing-a-plugin) section.

# Content 

- Simple example - Date/Time collector:
  * [Introduction](/tutorial/01-simple/README.md)
  * [Testing](/tutorial/02-testing/README.md)
  * [Configuration and state](/tutorial/03-concepts/README.md)
  * [Metrics - filters, definition, tags](/tutorial/04-metrics/README.md)
  * [Useful tools](/tutorial/05-tools/README.md)
- Advanced example - System collector:
  * [Overview](/tutorial/06-overview/README.md)
  * [Gathering data (Proxy)](/tutorial/07-proxy/README.md)
  * [Collector](/tutorial/08-collector/README.md)
  * [Handle configuration](/tutorial/09-config/README.md)
- [FAQ](/tutorial/faq/README.md)
