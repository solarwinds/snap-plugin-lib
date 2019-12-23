class BasePlugin:
    def __init__(self, name, version):
        self._name = name
        self._version = version

    def define_plugin(self, ctx):
        pass

    def collect(self, ctx):
        pass

    def load(self, ctx):
        pass

    def unload(self, ctx):
        pass

    def name(self):
        return self._name

    def version(self):
        return self._version