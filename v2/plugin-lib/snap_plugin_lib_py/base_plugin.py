from abc import ABC
from .c_bridge import Context, CollectContext, DefineContext, PublishContext


class BasePlugin(ABC):
    def __init__(self, name, version):
        self._name = name
        self._version = version

    def define_plugin(self, ctx: DefineContext):
        pass

    def load(self, ctx: Context):
        pass

    def unload(self, ctx: Context):
        pass

    def name(self):
        return self._name

    def version(self):
        return self._version


class BaseCollector(BasePlugin):
    def collect(self, ctx: CollectContext):
        pass

class BasePublisher(BasePlugin):
    def publish(self, ctx: PublishContext):
        pass
