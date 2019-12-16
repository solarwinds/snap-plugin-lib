from ctypes import *

lib_file = "mylib.dll"
my_fun = CDLL(lib_file)

class GoNamespace(Structure):
    _fields_ = [
        ("name", c_char_p),
        ("value", c_char_p),
        ("description", c_char_p)
    ]

class GoString(Structure):
    _fields_ = [
        ("p", c_char_p),
        ("n", c_longlong)
    ]

@CFUNCTYPE(None, GoNamespace)
def fun_callback(ns):
    print(ns.name)

# my_fun.Clear.argtypes = [Structure]
# my_fun.Clear(GoString(b"task-345", 8))

# my_fun.ListMetrics(GoString(b"task-123", 9), fun_callback)

my_fun.ListMetrics(fun_callback)