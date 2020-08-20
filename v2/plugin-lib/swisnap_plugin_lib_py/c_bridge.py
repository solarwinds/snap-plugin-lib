from collections import defaultdict
from ctypes import c_char_p, c_void_p, c_longlong, POINTER, CFUNCTYPE, pointer, cast
from .convertions import (
    string_to_bytes,
    dict_to_cmap,
    to_value_t,
    cstrarray_to_list,
    time_to_ctimewithns,
)
from .metric import Metric
from .snap_ctypes import CMetricStruct, CError, Modifiers
from .dynamic_lib import PLUGIN_LIB_OBJ
from .exceptions import throw_exception_if_error

# Used to store object for a given context. Access example: storedObjectMap[ctx_id][key]
storedObjectMap = defaultdict(dict)

# Reference to user-defined collector
global plugin_py
###############################################################################
# C functions metadata

PLUGIN_LIB_OBJ.ctx_list_all_metrics.restype = POINTER(POINTER(CMetricStruct))

PLUGIN_LIB_OBJ.ctx_add_metric.restype = POINTER(CError)
PLUGIN_LIB_OBJ.ctx_always_apply.restype = POINTER(CError)
PLUGIN_LIB_OBJ.ctx_dismiss_all_modifiers.restype = c_void_p
PLUGIN_LIB_OBJ.ctx_should_process.restype = c_longlong
PLUGIN_LIB_OBJ.ctx_requested_metrics.restype = POINTER(c_char_p)

PLUGIN_LIB_OBJ.ctx_config.restype = c_void_p  # -> string
PLUGIN_LIB_OBJ.ctx_config_keys.restype = POINTER(c_char_p)
PLUGIN_LIB_OBJ.ctx_raw_config.restype = c_void_p  # -> string
PLUGIN_LIB_OBJ.ctx_add_warning.restype = c_void_p
PLUGIN_LIB_OBJ.ctx_is_done.restype = c_longlong
PLUGIN_LIB_OBJ.ctx_log.restype = c_void_p

PLUGIN_LIB_OBJ.define_example_config.restype = POINTER(CError)


###############################################################################
# Python wrappers to context object - will call C functions and performing some conversions.
# Load, store are exceptions since it's safer to keep Python references on Python side


class DefineCommonContext:
    @staticmethod
    def define_tasks_per_instance_limit(limit):
        PLUGIN_LIB_OBJ.define_tasks_per_instance_limit(limit)

    @staticmethod
    def define_instances_limit(limit):
        PLUGIN_LIB_OBJ.define_instances_limit(limit)


class DefineCollectorContext(DefineCommonContext):
    @staticmethod
    def define_metric(namespace, unit, is_default, description):
        PLUGIN_LIB_OBJ.define_metric(
            string_to_bytes(namespace),
            string_to_bytes(unit),
            int(is_default),
            string_to_bytes(description),
        )

    @staticmethod
    def define_group(name, description):
        PLUGIN_LIB_OBJ.define_group(string_to_bytes(name), string_to_bytes(description))

    @staticmethod
    @throw_exception_if_error
    def define_example_config(config):
        return PLUGIN_LIB_OBJ.define_example_config(string_to_bytes(config))


class Context:
    def __init__(self, ctx_id):
        self.__ctx_id = ctx_id

    # @throw_exception_if_null("object with given key doesn't exist")
    def config(self, key: str):
        ret_ptr = PLUGIN_LIB_OBJ.ctx_config(self._ctx_id(), string_to_bytes(key))

        ret_char_ptr = cast(ret_ptr, c_char_p)
        ret_str = ret_char_ptr.value.decode(encoding="utf-8")
        PLUGIN_LIB_OBJ.dealloc_charp(ret_char_ptr)

        return ret_str

    def config_keys(self):
        config_list_c = PLUGIN_LIB_OBJ.ctx_config_keys(self._ctx_id())

        ret_list = cstrarray_to_list(config_list_c)
        PLUGIN_LIB_OBJ.dealloc_str_array(config_list_c)

        return ret_list

    def raw_config(self):
        ret_ptr = PLUGIN_LIB_OBJ.ctx_raw_config(self._ctx_id())

        ret_char_ptr = cast(ret_ptr, c_char_p)
        ret_str = ret_char_ptr.value.decode(encoding="utf-8")
        PLUGIN_LIB_OBJ.dealloc_charp(ret_char_ptr)

        return ret_str

    def store(self, key, obj):
        storedObjectMap[self._ctx_id()][key] = obj

    def load(self, key):
        return storedObjectMap[self._ctx_id()][key]

    def add_warning(self, message):
        return PLUGIN_LIB_OBJ.ctx_add_warning(self._ctx_id(), string_to_bytes(message))

    def is_done(self):
        return bool(PLUGIN_LIB_OBJ.ctx_is_done(self._ctx_id()))

    def log(self, level, message, fields):
        return PLUGIN_LIB_OBJ.ctx_log(
            self._ctx_id(),
            level,
            string_to_bytes(message),
            dict_to_cmap(fields),
            len(fields),
        )

    def _ctx_id(self):
        return self.__ctx_id


class PublishContext(Context):
    def list_all_metrics(self) -> []:
        _mts_ptr = PLUGIN_LIB_OBJ.ctx_list_all_metrics(self._ctx_id())

        _mt_list = self._c_mt_array_to_list(_mts_ptr, self.count())
        PLUGIN_LIB_OBJ.dealloc_metric_array(_mts_ptr, self.count())
        return _mt_list

    def count(self) -> int:
        return PLUGIN_LIB_OBJ.ctx_count(self._ctx_id())

    def _c_mt_array_to_list(self, mt_arr_ptr, mt_arr_size: int) -> []:
        """Converts C **metric_t to list of python managed objects of Metric class"""
        result_list = []
        for i in range(mt_arr_size):
            if mt_arr_ptr[i] is None:
                break
            mt = Metric.unpack_from_metric_struct(mt_arr_ptr[i].contents)
            result_list.append(mt)
        return result_list


class CollectContext(Context):
    @throw_exception_if_error
    def add_metric(
        self,
        namespace,
        value,
        *,
        tags=None,
        timestamp=None,
        description=None,
        unit=None
    ):
        return PLUGIN_LIB_OBJ.ctx_add_metric(
            self._ctx_id(),
            string_to_bytes(namespace),
            to_value_t(value),
            self.__create_modifiers(tags, None, timestamp, description, unit),
        )

    def always_apply(
        self,
        namespace,
        *,
        tags_to_add=None,
        tags_to_remove=None,
        timestamp=None,
        description=None,
        unit=None
    ):
        return PLUGIN_LIB_OBJ.ctx_always_apply(
            self._ctx_id(),
            string_to_bytes(namespace),
            self.__create_modifiers(
                tags_to_add, tags_to_remove, timestamp, description, unit
            ),
        )

    def dismiss_all_modifiers(self):
        PLUGIN_LIB_OBJ.ctx_dismiss_all_modifiers(self._ctx_id())

    def should_process(self, namespace):
        return bool(
            PLUGIN_LIB_OBJ.ctx_should_process(
                self._ctx_id(), string_to_bytes(namespace)
            )
        )

    def requested_metrics(self):
        req_mts_c = PLUGIN_LIB_OBJ.ctx_requested_metrics(self._ctx_id())
        return cstrarray_to_list(req_mts_c)

    @staticmethod
    def __create_modifiers(tags_to_add, tags_to_remove, timestamp, description, unit):
        modifiers = Modifiers()
        modifiers.tags_to_add = (
            dict_to_cmap(tags_to_add) if tags_to_add is not None else None
        )
        modifiers.tags_to_remove = (
            dict_to_cmap(tags_to_remove) if tags_to_remove is not None else None
        )
        modifiers.description = (
            c_char_p(string_to_bytes(description)) if description is not None else None
        )
        modifiers.unit = c_char_p(string_to_bytes(unit)) if unit is not None else None
        modifiers.timestamp = (
            time_to_ctimewithns(timestamp) if timestamp is not None else None
        )
        return pointer(modifiers)


###############################################################################
# Callback related functions (called from C library)


@CFUNCTYPE(None)
def define_publisher_handler():
    plugin_py.define_plugin(DefineCommonContext())


@CFUNCTYPE(None)
def define_collector_handler():
    plugin_py.define_plugin(DefineCollectorContext())


@CFUNCTYPE(None, c_char_p)
def collector_handler(ctx_id):
    plugin_py.collect(CollectContext(ctx_id))


@CFUNCTYPE(None, c_char_p)
def publisher_handler(ctx_id):
    plugin_py.publish(PublishContext(ctx_id))


@CFUNCTYPE(None, c_char_p)
def load_handler(ctx_id):
    plugin_py.load(Context(ctx_id))


@CFUNCTYPE(None, c_char_p)
def unload_handler(ctx_id):
    plugin_py.unload(Context(ctx_id))


###############################################################################
# Collector setup


def start_c_collector(collector):
    global plugin_py

    name = collector.name()
    version = collector.version()
    plugin_py = collector

    PLUGIN_LIB_OBJ.start_collector(
        collector_handler,
        load_handler,
        unload_handler,
        define_collector_handler,
        string_to_bytes(name),
        string_to_bytes(version),
    )


def start_c_publisher(publisher):
    global plugin_py

    name = publisher.name()
    version = publisher.version()
    plugin_py = publisher

    PLUGIN_LIB_OBJ.start_publisher(
        publisher_handler,
        load_handler,
        unload_handler,
        define_publisher_handler,
        string_to_bytes(name),
        string_to_bytes(version),
    )
