from ctypes import *

lib_file = "mylib.dll"
lib_obj = CDLL(lib_file)


class Tags(Structure):
    _fields_ = [
        ("key", c_char_p),
        ("value", c_char_p)
    ]


class CError(Structure):
    _fields_ = [
        ("msg", c_char_p)
    ]


@CFUNCTYPE(None)
def define():
    print("** python *** Define plugin\n")

    lib_obj.define_example_config('{"ip": "127.0.0.1", "port": 5434}')
    lib_obj.define_tasks_per_instance_limit(4)
    lib_obj.define_instances_limit(3)

    lib_obj.define_group(b"dyn", "Dynamic element from python")

    lib_obj.define_metric(b"/python/group1/metric1", b"C", 1, "1st metric")
    lib_obj.define_metric(b"/python/group1/metric2", b"C", 1, "2nd metric")
    lib_obj.define_metric(b"/python/group1/metric3", b"C", 0, "3rd metric")
    lib_obj.define_metric(b"/python/group2/[dyn]/metric4", b"C", 0, "1st dynamic metric")
    lib_obj.define_metric(b"/python/group2/[dyn]/metric5", b"C", 0, "2nd dynamic metric")


@CFUNCTYPE(None, c_char_p)
def collect(ctxId):
    print("** python *** Collect called\n")

    lib_obj.ctx_add_metric.restype = POINTER(CError)
    lib_obj.ctx_add_metric(ctxId, b"/python/group1/metric1", 10)
    lib_obj.ctx_add_metric(ctxId, b"/python/group1/metric2", 20)
    lib_obj.ctx_add_metric(ctxId, b"/python/group1/metric3", 40)
    lib_obj.ctx_add_metric(ctxId, b"/python/group2/dyn1/metric4", 40)
    lib_obj.ctx_add_metric(ctxId, b"/python/group2/dyn15/metric4", 11)
    res1 = lib_obj.ctx_add_metric(ctxId, b"/python/group2/dyn24/metric5", 34)

    res = lib_obj.ctx_add_metric(ctxId, b"/python/group1/metricWRONG", 10)

    tags = (Tags * 2)()
    tags[0].key = b"tag1"
    tags[0].value = b"value1"
    tags[1].key = b"tag2"
    tags[1].value = b"value2"
    lib_obj.ctx_add_metric_with_tags(ctxId, b"/python/group1/metric1", 14, tags, 2)


@CFUNCTYPE(None, c_char_p)
def load(ctxId):
    print("** python *** Load called\n")

    lib_obj.ctx_config.restype = c_char_p
    res = lib_obj.ctx_config(ctxId, b"configip")
    print(res)


@CFUNCTYPE(None, c_char_p)
def unload(ctxId):
    print("** python *** Unload called\n")


###############################################################################

plugin_lib_filename = "mylib.dll"


class BasePlugin():
    def __init__(self, name, version):
        self._name = name
        self._version = version
        self.plugin_lib = CDLL(lib_file)

    @staticmethod
    @CFUNCTYPE(None, c_char_p)
    def __collect_handler(ctx_id):
        pass

    @staticmethod
    @CFUNCTYPE(None, c_char_p)
    def __load_handler(ctx_id):
        pass

    @staticmethod
    @CFUNCTYPE(None, c_char_p)
    def __unload_handler(ctx_id):
        pass

    @staticmethod
    @CFUNCTYPE(None)
    def __define_plugin_handler():
        pass

    def start(self):
        self.plugin_lib.start_collector(BasePlugin.__collect_handler,
                                        BasePlugin.__load_handler,
                                        BasePlugin.__unload_handler,
                                        BasePlugin.__define_plugin_handler,
                                        bytes(self._version, 'utf-8'), bytes(self._name, 'utf-8'))


###############################################################################

if __name__ == '__main__':
    lib_obj.start_collector(collect, load, unload, define, b"python-collector", b"0.0.1")
