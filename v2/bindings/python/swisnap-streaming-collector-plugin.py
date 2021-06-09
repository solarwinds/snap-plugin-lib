import time

from swisnap_plugin_lib_py import BaseStreamingCollector, LOGLEVEL_INFO
from swisnap_plugin_lib_py.runner import start_streaming_collector


class ExampleStreamingCollectorPlugin(BaseStreamingCollector):
    def streaming_collect(self, ctx):
        ctx.log(LOGLEVEL_INFO, "Streaming collect requested", {"name": self._name})

        metric_value = 0

        while True:
            metric_value += 2
            ctx.add_metric("/example/group1/metric1", metric_value)
            time.sleep(2)

    def load(self, ctx):
        ctx.log(LOGLEVEL_INFO, "Plugin is being loaded", {"name": self._name})

    def unload(self, ctx):
        ctx.log(LOGLEVEL_INFO, "Plugin is being unloaded", {"name": self._name})


if __name__ == "__main__":
    collector = ExampleStreamingCollectorPlugin("collector-example", "0.0.1")
    start_streaming_collector(collector)
