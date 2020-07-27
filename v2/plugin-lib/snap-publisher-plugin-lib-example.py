import time

from snap_plugin_lib_py import BasePublisher, start_publisher, LOGLEVEL_INFO
from snap_plugin_lib_py.exceptions import PluginLibException



class ExamplePublisherPlugin(BasePublisher):
    def define_plugin(self, ctx):
        print('1')

    def publish(self, ctx):
        ctx.log(LOGLEVEL_INFO, "Plugin is publishing", {
            "name": "py-example"
        })

    def load(self, ctx):
        print('3')
        ctx.log(LOGLEVEL_INFO, "Plugin is being loaded", {
            "name": "py-example"
        })

    def unload(self, ctx):

        ctx.log(LOGLEVEL_INFO, "Plugin is being unloaded", {
            "name": "py-example"
        })
        print('4')

if __name__ == '__main__':
    publisher = ExamplePublisherPlugin("example", "0.0.1")
    start_publisher(publisher)
