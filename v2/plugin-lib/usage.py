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

@CFUNCTYPE(None)
def define():
    print("** python *** Define plugin\n")


@CFUNCTYPE(None, c_char_p)
def collect(ctxId):
    print("** python *** Collect called\n")

    my_fun.ctx_add_metric_int(ctxId, b"/python/example/metric1", 10)
    my_fun.ctx_add_metric_int(ctxId, b"/python/example/metric2", 20)
    my_fun.ctx_add_metric_int(ctxId, b"/python/example/metric3", 40)

@CFUNCTYPE(None, c_char_p)
def load(ctxId):
    print("** python *** Load called\n")

@CFUNCTYPE(None, c_char_p)
def unload(ctxId):
    print("** python *** Unload called\n")

my_fun.StartCollector(collect, load, unload, define, b"python-collector", b"0.0.1")
