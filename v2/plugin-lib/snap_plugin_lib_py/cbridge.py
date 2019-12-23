from ctypes import *
import os.path

plugin_lib_file = "snap-plugin-lib.dll"
plugin_lib_obj = CDLL(os.path.join(os.path.dirname(__file__), plugin_lib_file))

global collector_py


###############################################################################

class Tags(Structure):
    _fields_ = [
        ("key", c_char_p),
        ("value", c_char_p)
    ]


class CError(Structure):
    _fields_ = [
        ("msg", c_char_p)
    ]


class DefineContext:
    def define_tasks_per_instance_limit(self, limit):
        pass

    def define_instances_limit(self, limit):
        pass

    def define_metric(self, namespace, unit, is_default, description):
        pass

    def define_group(self, name, description):
        pass

    def define_global_tags(self, selector, tags):
        pass

    def define_example_config(self, config):
        pass


class Context:
    def config(self, key):
        pass

    def config_keys(self):
        pass

    def raw_config(self):
        pass

    def store(self, obj):
        pass

    def load(self, obj):
        pass


class CollectContext:
    def add_metric(self, namespace, value):
        pass

    def add_metric_with_tags(self, namespace, value, tags):
        pass

    def apply_tags_by_path(self, namespace, tags):
        pass

    def apply_tags_by_regexp(self, selector, tags):
        pass

    def should_process(self, namespace):
        pass


###############################################################################

@CFUNCTYPE(None)
def define_handler():
    global collector_py

    print("** cpython *** Define plugin\n")
    collector_py.define_plugin(DefineContext())


@CFUNCTYPE(None, c_char_p)
def collect_handler(ctxId):
    print("** cpython *** Collect called\n")


@CFUNCTYPE(None, c_char_p)
def load_handler(ctxId):
    print("** cpython *** Load called\n")


@CFUNCTYPE(None, c_char_p)
def unload_handler(ctxId):
    print("** cpython *** Unload called\n")


###############################################################################

# todo: only 1 instance should be run # check this
def start_c_collector(collector):
    global collector_py

    name = collector.name()
    version = collector.version()
    collector_py = collector

    plugin_lib_obj.start_collector(collect_handler,
                                   load_handler,
                                   unload_handler,
                                   define_handler,
                                   string_to_bytes(name), string_to_bytes(version))


###############################################################################


def string_to_bytes(s):
    if isinstance(s, type("")):
        return bytes(s, 'utf-8')
    elif isinstance(s, type(b"")):
        return s
    else:
        raise Exception("Invalid type, expected string or bytes")
