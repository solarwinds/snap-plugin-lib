import os.path
from collections import defaultdict
from ctypes import CDLL, c_char_p, c_void_p, c_longlong, POINTER, CFUNCTYPE
import platform
from itertools import count

from .convertions import string_to_bytes, dict_to_tags, CError, to_value_t
from .exceptions import throw_exception_if_error, throw_exception_if_null

# Dependent library
PLUGIN_LIB_EXTENSION = ".dll" if (platform.system() == "Windows") else ".so"
PLUGIN_LIB_FILE = "snap-plugin-lib%s" % PLUGIN_LIB_EXTENSION
PLUGIN_LIB_OBJ = CDLL(os.path.join(os.path.dirname(__file__), PLUGIN_LIB_FILE))

# Used to store object for a given context. Access example: storedObjectMap[ctx_id][key]
storedObjectMap = defaultdict(dict)

# Reference to user-defined collector
global collector_py

###############################################################################
# C functions metadata

PLUGIN_LIB_OBJ.define_example_config.restype = POINTER(CError)
PLUGIN_LIB_OBJ.ctx_add_metric.restype = POINTER(CError)
PLUGIN_LIB_OBJ.ctx_config.restype = c_char_p
PLUGIN_LIB_OBJ.ctx_raw_config.restype = c_char_p
PLUGIN_LIB_OBJ.ctx_is_done.restype = c_longlong
PLUGIN_LIB_OBJ.ctx_add_warning.restype = c_void_p
PLUGIN_LIB_OBJ.ctx_log.restype = c_void_p
PLUGIN_LIB_OBJ.ctx_dismiss_all_modifiers.restype = c_void_p
PLUGIN_LIB_OBJ.ctx_requested_metrics.restype = POINTER(c_char_p)


###############################################################################
# Python wrappers to context object - will call C functions and performing some conversions.
# Load, store are exceptions since it's safer to keep Python references on Python side

class DefineContext:
    @staticmethod
    def define_tasks_per_instance_limit(limit):
        PLUGIN_LIB_OBJ.define_tasks_per_instance_limit(limit)

    @staticmethod
    def define_instances_limit(limit):
        PLUGIN_LIB_OBJ.define_instances_limit(limit)

    @staticmethod
    def define_metric(namespace, unit, is_default, description):
        PLUGIN_LIB_OBJ.define_metric(string_to_bytes(namespace),
                                     string_to_bytes(unit),
                                     int(is_default),
                                     string_to_bytes(description))

    @staticmethod
    def define_group(name, description):
        PLUGIN_LIB_OBJ.define_group(string_to_bytes(name),
                                    string_to_bytes(description))

    @staticmethod
    @throw_exception_if_error
    def define_example_config(config):
        return PLUGIN_LIB_OBJ.define_example_config(string_to_bytes(config))


class Context:
    def __init__(self, ctx_id):
        self.__ctx_id = ctx_id

    @throw_exception_if_null("object with given key doesn't exists")
    def config(self, key: str):
        return PLUGIN_LIB_OBJ.ctx_config(self._ctx_id(),
                                         string_to_bytes(key)).decode(encoding='utf-8')

    def raw_config(self):
        return PLUGIN_LIB_OBJ.ctx_raw_config(self._ctx_id()).decode(encoding='utf-8')

    def requested_metrics(self):
        req_mts_c = PLUGIN_LIB_OBJ.ctx_requested_metrics(self._ctx_id())
        req_mts = []
        for i in count(0):
            if req_mts_c[i] is None:
                break
            req_mts.append(req_mts_c[i].decode(encoding='utf-8'))

        return req_mts

    def store(self, key, obj):
        storedObjectMap[self._ctx_id()][key] = obj

    def load(self, key):
        return storedObjectMap[self._ctx_id()][key]

    def log(self, level, message, fields):
        return PLUGIN_LIB_OBJ.ctx_log(self._ctx_id(),
                                      level,
                                      string_to_bytes(message),
                                      dict_to_tags(fields),
                                      len(fields))

    def add_warning(self, message):
        return PLUGIN_LIB_OBJ.ctx_add_warning(self._ctx_id(),
                                              string_to_bytes(message))

    def is_done(self):
        return bool(PLUGIN_LIB_OBJ.ctx_is_done(self._ctx_id()))

    def dismiss_all_modifiers(self):
        PLUGIN_LIB_OBJ.ctx_dismiss_all_modifiers(self._ctx_id())

    def _ctx_id(self):
        return self.__ctx_id


class CollectContext(Context):
    @throw_exception_if_error
    def add_metric(self, namespace, value):
        return PLUGIN_LIB_OBJ.ctx_add_metric(self._ctx_id(),
                                             string_to_bytes(namespace),
                                             to_value_t(value))

    def should_process(self, namespace):
        return bool(PLUGIN_LIB_OBJ.ctx_should_process(self._ctx_id(),
                                                      string_to_bytes(namespace)))


###############################################################################
# Callback related functions (called from C library)

@CFUNCTYPE(None)
def define_handler():
    collector_py.define_plugin(DefineContext())


@CFUNCTYPE(None, c_char_p)
def collect_handler(ctx_id):
    collector_py.collect(CollectContext(ctx_id))


@CFUNCTYPE(None, c_char_p)
def load_handler(ctx_id):
    collector_py.load(Context(ctx_id))


@CFUNCTYPE(None, c_char_p)
def unload_handler(ctx_id):
    collector_py.unload(Context(ctx_id))


###############################################################################
# Collector setup

def start_c_collector(collector):
    global collector_py

    name = collector.name()
    version = collector.version()
    collector_py = collector

    PLUGIN_LIB_OBJ.start_collector(collect_handler,
                                   load_handler,
                                   unload_handler,
                                   define_handler,
                                   string_to_bytes(name), string_to_bytes(version))
