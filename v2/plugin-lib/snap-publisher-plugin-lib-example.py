import time

from snap_plugin_lib_py import BasePublisher, start_publisher, LOGLEVEL_INFO
from snap_plugin_lib_py.exceptions import PluginLibException
import pprint


class ExamplePublisherPlugin(BasePublisher):
    def define_plugin(self, ctx):
        print(' define: 1')

    def publish(self, ctx):
        print(' publish: 2 count:')
        print(ctx.count())
        ctx.log(LOGLEVEL_INFO, "Plugin is publishing", {
            "name": "py-example"
        })
        print("test")
        pprint.pprint(ctx.list_all_metrics())

    def load(self, ctx):
        print('load: 3')
        ctx.log(LOGLEVEL_INFO, "Plugin is being loaded", {
            "name": "py-example"
        })

    def unload(self, ctx):

        ctx.log(LOGLEVEL_INFO, "Plugin is being unloaded", {
            "name": "py-example"
        })
        print('unload: 4')

if __name__ == '__main__':
    publisher = ExamplePublisherPlugin("example", "0.0.1")
    start_publisher(publisher)
