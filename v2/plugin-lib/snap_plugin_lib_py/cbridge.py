from ctypes import *
import os.path
from .exceptions import PluginLibException, throw_exception_if_error

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


###############################################################################
# C function metadata

plugin_lib_obj.ctx_add_metric.restype = POINTER(CError)
plugin_lib_obj.ctx_add_metric_with_tags.restype = POINTER(CError)
plugin_lib_obj.ctx_apply_tags_by_path.restype = POINTER(CError)
plugin_lib_obj.ctx_apply_tags_by_regexp.restype = POINTER(CError)
plugin_lib_obj.ctx_apply_tags_by_regexp.restype = c_longlong


###############################################################################

class DefineContext:
    def define_tasks_per_instance_limit(self, limit):
        plugin_lib_obj.define_tasks_per_instance_limit(limit)

    def define_instances_limit(self, limit):
        plugin_lib_obj.define_instances_limit(limit)

    def define_metric(self, namespace, unit, is_default, description):
        plugin_lib_obj.define_metric(string_to_bytes(namespace),
                                     string_to_bytes(unit),
                                     int(is_default),
                                     string_to_bytes(description))

    def define_group(self, name, description):
        plugin_lib_obj.define_group(string_to_bytes(name),
                                    string_to_bytes(description))

    def define_global_tags(self, selector, tags):
        plugin_lib_obj.define_global_tags(string_to_bytes(selector),
                                          dict_to_tags(tags),
                                          len(tags))

    def define_example_config(self, config):
        pass


class Context:
    def __init__(self, ctx_id):
        self._ctx_id = ctx_id

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

    def ctx_id(self):
        return self._ctx_id


class CollectContext(Context):
    @throw_exception_if_error
    def add_metric(self, namespace, value):
        return plugin_lib_obj.ctx_add_metric(self.ctx_id(),
                                             string_to_bytes(namespace),
                                             c_longlong(value))

    @throw_exception_if_error
    def add_metric_with_tags(self, namespace, value, tags):
        return plugin_lib_obj.ctx_add_metric_with_tags(self.ctx_id(),
                                                       string_to_bytes(namespace),
                                                       c_longlong(value),
                                                       dict_to_tags(tags),
                                                       len(tags))

    @throw_exception_if_error
    def apply_tags_by_path(self, namespace, tags):
        return plugin_lib_obj.ctx_apply_tags_by_path(self.ctx_id(),
                                                     string_to_bytes(namespace),
                                                     dict_to_tags(tags),
                                                     len(tags))

    @throw_exception_if_error
    def apply_tags_by_regexp(self, selector, tags):
        return plugin_lib_obj.ctx_apply_tags_by_regexp(self.ctx_id(),
                                                       string_to_bytes(selector),
                                                       dict_to_tags(tags),
                                                       len(tags))

    def should_process(self, namespace):
        return bool(plugin_lib_obj.ctx_should_process(self.ctx_id(),
                                                      string_to_bytes(namespace)))


###############################################################################
# Callback related functions (called from C library)

@CFUNCTYPE(None)
def define_handler():
    print("** cpython *** Define plugin\n")
    collector_py.define_plugin(DefineContext())


@CFUNCTYPE(None, c_char_p)
def collect_handler(ctxId):
    collector_py.collect(CollectContext(ctxId))
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


def dict_to_tags(d):
    tags = (Tags * len(d))()

    for i, (k, v) in enumerate(d.items()):
        tags[i].key = string_to_bytes(k)
        tags[i].value = string_to_bytes(v)

    return tags
