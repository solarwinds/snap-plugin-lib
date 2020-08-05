import math
from ctypes import (
    c_longlong,
    c_ulonglong,
    c_double,
    c_int,
    pointer,
)
from itertools import count
from .snap_ctypes import Map, MapElement, TimeWithNs, CValue, TYPE_INT64, TYPE_UINT64, TYPE_DOUBLE, TYPE_BOOL, max_int, max_uint, min_int

from .exceptions import PluginLibException


def string_to_bytes(s):
    """
    Converts string to bytes if necessary.
    Allow to use string type in Python code and covert it to required char *
    (bytes) when calling C Api
    """

    if isinstance(s, type("")):
        return bytes(s, "utf-8")
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


def cmap_to_dict(cmap_ptr):
    """Converts C map pointer to python dict"""
    map_len = cmap_ptr.contents.length
    _map = dict()
    if map_len != 0:
        for i in range(map_len):
            el = cmap_ptr.contents.elements[i]
            _map[el.key.decode(encoding="utf-8")] = el.value.decode(encoding="utf-8")
    return _map


def cstrarray_to_list(arr):
    """Converts C **char to Python list"""
    result_list = []
    for i in count(0):
        if arr[i] is None:
            break
        result_list.append(arr[i].decode(encoding="utf-8"))

    return result_list


def time_to_ctimewithns(timestamp):
    sec = int(math.floor(timestamp))
    nsec = int(math.floor(timestamp - sec) * 1e9)
    return pointer(TimeWithNs(sec, nsec))


def ctimewithns_to_time(ctime_ptr):
    sec = ctime_ptr.contents.sec
    nsec = ctime_ptr.contents.nsec / 1e9
    return sec + nsec


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


def unpack_value_t(val_ptr):
    v_type = val_ptr.contents.v_type
    if v_type == TYPE_DOUBLE:
        value = val_ptr.contents.value.v_double
        unit = float
    elif v_type == TYPE_BOOL:
        value = val_ptr.contents.value.v_bool
        unit = bool
    elif v_type == TYPE_INT64:
        value = val_ptr.contents.value.v_int64
        unit = int
    elif v_type == TYPE_UINT64:
        value = val_ptr.contents.value.v_uint64
        unit = int
    return (value, unit)
