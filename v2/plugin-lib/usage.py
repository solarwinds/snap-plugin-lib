from ctypes import *

lib_file = "mylib.dll"
my_fun = CDLL(lib_file)

class GoNamespaceEl(Structure):
    _fields_ = [
        ("name", c_char_p),
        ("value", c_char_p),
        ("description", c_char_p),
    ]

class GoNamespace(Structure):
    _fields_ = [
        ("length", c_longlong),
        ("elements", POINTER(GoNamespaceEl))
    ]

class GoString(Structure):
    _fields_ = [
        ("p", c_char_p),
        ("n", c_longlong)
    ]

@CFUNCTYPE(None, c_char_p)
def collect(ctxId):
    my_fun.ctx_add_metric_int(ctxId, b"/python/example/metric1", 10)
    my_fun.ctx_add_metric_int(ctxId, b"/python/example/metric2", 20)
    my_fun.ctx_add_metric_int(ctxId, b"/python/example/metric3", 40)


class cCollector(Structure):
    _fields_ = [
        ("collect_callback", CFUNCTYPE(None, c_char_p))
    ]

my_fun.StartCollector(collect, b"python-collector", b"0.0.1")
