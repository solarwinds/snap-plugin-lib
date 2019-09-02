# Tutorial

## Intro

This tutorial will teach you how to write custom plugin in Go language working with snap deamon. 

We will start from very simple example, when you will learn how to build a minimal plugin and test that it's working correctly. 
After you obtain basics, we will teach you how to write real, useful plugin (gathering sysytem metrics) which utilizes advanced features that plugin-go V2 has to offer:
- defining supported metrics
- build-in filtering mechanism
- defining custom configuration 
- store state between independent collection requests
- optimazing gathering metrics
- etc.

Each part of tutorial contain complete golang source code.

## Version 2

Only version 2 of plugin-lib-go is covered in this tutorial. Examples and plugin catalog related to version 1 can be found in [Writing a Plugin](https://github.com/librato/snap-plugin-lib-go/tree/ao-12231-tutorial#writing-a-plugin) section.

# Content 

* [Introduction - Simple plugin](/tutorial/01/README.md)
* [Testing](/tutorial/02/README.md)
* [Basic concepts]((/tutorial/03/README.md))
* [Advanced plugin - overview](/tutorial/04/README.md)
* [Defining plugin metadata (metrics and configuration)](/tutorial/05/README.md)
* [Other features](/tutorial/06/README.md)
