import time

from snap_plugin_lib_py import BasePublisher, start_publisher, LOGLEVEL_INFO
from snap_plugin_lib_py.exceptions import PluginLibException
import pprint


class ExamplePublisherPlugin(BasePublisher):
    def define_plugin(self, ctx):
        pass

    def publish(self, ctx):
        ctx.log(LOGLEVEL_INFO, "Plugin is publishing", {"name": self._name})
        pprint.pprint("Number of received metrics: {}".format(ctx.count()))
        pprint.pprint(ctx.list_all_metrics())

    def load(self, ctx):
        ctx.log(LOGLEVEL_INFO, "Plugin is being loaded", {"name": self._name})

    def unload(self, ctx):
        ctx.log(LOGLEVEL_INFO, "Plugin is being unloaded", {"name": self._name})


if __name__ == "__main__":
    publisher = ExamplePublisherPlugin("publisher-example", "0.0.1")
    start_publisher(publisher)
