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

@CFUNCTYPE(None, POINTER(GoNamespace))
def fun_callback(ns):
    for i in range(4):
        print(ns[0].length)
        print(ns[0].elements[i].name)
        print(ns[0].elements[i].value)
        print(ns[0].elements[i].description)

# my_fun.Clear.argtypes = [Structure]
# my_fun.Clear(GoString(b"task-345", 8))

# my_fun.ListMetrics(GoString(b"task-123", 9), fun_callback)

##extern void StartCollector(cCollectorT* p0, char* p1, char* p2);

@CFUNCTYPE(None, c_void_p)
def collect(rawCtx):
    my_fun.ctx_add_metric(rawCtx, b"/python/example/metric")

class cCollector(Structure):
    _fields_ = [
        ("collect_callback", CFUNCTYPE(None, c_void_p))
    ]

# cc = cCollector()
# cc.collect_callback = collect
# 
# my_fun.StartCollector(cc, "python-collector", "0.0.1")

my_fun.StartCollector(collect, b"python-collector", b"0.0.1")
