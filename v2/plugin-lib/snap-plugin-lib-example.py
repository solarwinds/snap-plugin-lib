from snap_plugin_lib_py import *


class ExamplePlugin(BasePlugin):
    def define_plugin(self, ctx):
        pass

    def collect(self, ctx):
        pass

    def load(self, ctx):
        pass

    def unload(self, ctx):
        pass


if __name__ == '__main__':
    collector = ExamplePlugin("example", "0.0.1")
    start_collector(collector)
