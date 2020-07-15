import math
from ctypes import Structure, Union, c_char_p, c_longlong, c_ulonglong, c_double, c_int, POINTER, pointer
from itertools import count

from .exceptions import PluginLibException

min_int = -9223372036854775807
max_int = 9223372036854775807
max_uint = 18446744073709551615

_, TYPE_INT64, TYPE_UINT64, TYPE_DOUBLE, TYPE_BOOL = range(5)
_, LOGLEVEL_PANIC, LOGLEVEL_FATAL, LOGLEVEL_ERROR, LOGLEVEL_WARN, \
    LOGLEVEL_INFO, LOGLEVEL_DEBUG, LOGLEVEL_TRACE = range(8)


class MapElement(Structure):
    _fields_ = [
        ("key", c_char_p),
        ("value", c_char_p)
    ]


class Map(Structure):
    _fields_ = [
        ("elements", POINTER(MapElement)),
        ("length", c_int)
    ]


class TimeWithNs(Structure):
    _fields_ = [
        ("sec", c_int),
        ("nsec", c_int)
    ]


class Modifiers(Structure):
    _fields_ = [
        ("tags_to_add", POINTER(Map)),
        ("tags_to_remove", POINTER(Map)),
        ("timestamp", POINTER(TimeWithNs)),
        ("description", POINTER(c_char_p)),
        ("unit", POINTER(c_char_p))
    ]


class CError(Structure):
    _fields_ = [
        ("msg", c_char_p)
    ]


class ValueUnion(Union):
    _fields_ = [
        ("v_int64", c_longlong),
        ("v_uint64", c_ulonglong),
        ("v_double", c_double),
        ("v_bool", c_int),
    ]


class CValue(Structure):
    _fields_ = [
        ("value", ValueUnion),
        ("v_type", c_int)
    ]


def string_to_bytes(s):
    """
    Converts string to bytes if necessary.
    Allow to use string type in Python code and covert it to required char *
    (bytes) when calling C Api
    """

    if isinstance(s, type("")):
        return bytes(s, 'utf-8')
    elif isinstance(s, type(b"")):
        return s
    else:
        raise Exception("Invalid type, expected string or bytes")


def dict_to_cmap(d):
    """Converts python dictionary to C map pointer"""

    cmap = Map()
    cmap.elements = (MapElement * len(d))()
    cmap.length = len(d)

    for i, (k, v) in enumerate(d.items()):
        cmap.elements[i].key = string_to_bytes(k)
        cmap.elements[i].value = string_to_bytes(v)

    return pointer(cmap)


def cstrarray_to_list(arr):
    """Converts C **char to Python list"""
    result_list = []
    for i in count(0):
        if arr[i] is None:
            break
        result_list.append(arr[i].decode(encoding='utf-8'))

    return result_list


def time_to_ctimewithns(timestamp):
    sec = int(math.floor(timestamp))
    nsec = int(math.floor(timestamp - sec) * 1e9)
    return pointer(TimeWithNs(sec, nsec))


def to_value_t(v):
    val_ptr = (CValue * 1)()
    val = val_ptr[0]

    if isinstance(v, bool):
        val.value.v_bool = c_int(v)
        val.v_type = TYPE_BOOL
    elif isinstance(v, int):
        if min_int <= v <= max_int:
            val.value.v_int64 = c_longlong(v)
            val.v_type = TYPE_INT64
        else:
            if v <= max_uint:
                val.value.v_uint64 = c_ulonglong(v)
                val.v_type = TYPE_UINT64
            else:
                val.value.v_double = c_double(v)
                val.v_type = TYPE_DOUBLE
    elif isinstance(v, float):
        val.value.v_double = c_double(v)
        val.v_type = TYPE_DOUBLE
    else:
        raise PluginLibException("invalid metric value type")

    return val_ptr
