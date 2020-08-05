from ctypes import (
    Structure,
    Union,
    c_char_p,
    c_longlong,
    c_ulonglong,
    c_double,
    c_int,
    POINTER,
)

min_int = -9223372036854775808
max_int = 9223372036854775807
max_uint = 18446744073709551615

_, TYPE_INT64, TYPE_UINT64, TYPE_DOUBLE, TYPE_BOOL = range(5)
(
    _,
    LOGLEVEL_PANIC,
    LOGLEVEL_FATAL,
    LOGLEVEL_ERROR,
    LOGLEVEL_WARN,
    LOGLEVEL_INFO,
    LOGLEVEL_DEBUG,
    LOGLEVEL_TRACE,
) = range(8)


class MapElement(Structure):
    _fields_ = [("key", c_char_p), ("value", c_char_p)]


class Map(Structure):
    _fields_ = [("elements", POINTER(MapElement)), ("length", c_int)]


class TimeWithNs(Structure):
    _fields_ = [("sec", c_int), ("nsec", c_int)]


class Modifiers(Structure):
    _fields_ = [
        ("tags_to_add", POINTER(Map)),
        ("tags_to_remove", POINTER(Map)),
        ("timestamp", POINTER(TimeWithNs)),
        ("description", c_char_p),
        ("unit", c_char_p),
    ]


class CError(Structure):
    _fields_ = [("msg", c_char_p)]


class ValueUnion(Union):
    _fields_ = [
        ("v_int64", c_longlong),
        ("v_uint64", c_ulonglong),
        ("v_double", c_double),
        ("v_bool", c_int),
    ]


class CValue(Structure):
    _fields_ = [("value", ValueUnion), ("v_type", c_int)]


class CNamespaceElement(Structure):
    _fields_ = [
        ("name", c_char_p),
        ("value", c_char_p),
        ("description", c_char_p),
    ]


class CNamespace(Structure):
    _fields_ = [
        ("elements", POINTER(CNamespaceElement)),
        ("length", c_int),
        ("string", c_char_p),
    ]


class CMetricStruct(Structure):
    _fields_ = [
        ("namespace", POINTER(CNamespace)),
        ("description", c_char_p),
        ("value", POINTER(CValue)),
        ("timestamp", POINTER(TimeWithNs)),
        ("tags", POINTER(Map)),
    ]
