from .cbridge import CollectContext, DefineContext, Context


class BasePlugin:
    def __init__(self, name, version):
        self._name = name
        self._version = version

    def define_plugin(self, ctx: DefineContext):
        pass

    def collect(self, ctx: CollectContext):
        pass

    def load(self, ctx: Context):
        pass

    def unload(self, ctx: Context):
        pass

    def name(self):
        return self._name

    def version(self):
        return self._version
