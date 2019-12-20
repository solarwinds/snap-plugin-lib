from usage import BasePlugin


class ExamplePlugin(BasePlugin):
    def define_plugin(self):
        pass

    def collect(self, ctx):
        pass

    def load(self, ctx):
        pass

    def unload(self, ctx):
        pass


if __name__ == '__main__':
    plugin = ExamplePlugin("example-plugin", "0.0.1")
    plugin.start()
