from ctypes import *

lib_file = "mylib.dll"
lib_obj = CDLL(lib_file)

# class GoNamespaceEl(Structure):
#     _fields_ = [
#         ("name", c_char_p),
#         ("value", c_char_p),
#         ("description", c_char_p),
#     ]

# class GoNamespace(Structure):
#     _fields_ = [
#         ("length", c_longlong),
#         ("elements", POINTER(GoNamespaceEl))
#     ]

# class GoString(Structure):
#     _fields_ = [
#         ("p", c_char_p),
#         ("n", c_longlong)
#     ]

class Tags(Structure):
    _fields_ = [
        ("key", c_char_p),
        ("value", c_char_p)
    ]

@CFUNCTYPE(None)
def define():
    print("** python *** Define plugin\n")

    lib_obj.define_example_config('{"ip": "127.0.0.1", "port": 5434}')
    lib_obj.define_tasks_per_instance_limit(4)
    lib_obj.define_instances_limit(3)

    lib_obj.define_group(b"dyn", "Dynamic element from python")

    lib_obj.define_metric(b"/python/group1/metric1", b"C", 1, "1st metric")
    lib_obj.define_metric(b"/python/group1/metric2", b"C", 1, "2nd metric")
    lib_obj.define_metric(b"/python/group1/metric3", b"C", 0, "3rd metric")
    lib_obj.define_metric(b"/python/group2/[dyn]/metric4", b"C", 0, "1st dynamic metric")
    lib_obj.define_metric(b"/python/group2/[dyn]/metric5", b"C", 0, "2nd dynamic metric")


@CFUNCTYPE(None, c_char_p)
def collect(ctxId):
    print("** python *** Collect called\n")

    lib_obj.ctx_add_metric(ctxId, b"/python/group1/metric1", 10)
    lib_obj.ctx_add_metric(ctxId, b"/python/group1/metric2", 20)
    lib_obj.ctx_add_metric(ctxId, b"/python/group1/metric3", 40)
    lib_obj.ctx_add_metric(ctxId, b"/python/group2/dyn1/metric4", 40)
    lib_obj.ctx_add_metric(ctxId, b"/python/group2/dyn15/metric4", 11)
    lib_obj.ctx_add_metric(ctxId, b"/python/group2/dyn24/metric5", 34)

    tags = (Tags * 2)()
    tags[0].key = b"tag1"
    tags[0].value = b"value1"
    tags[1].key = b"tag2"
    tags[1].value = b"value2"
    lib_obj.ctx_add_metric_with_tags(ctxId, b"/python/group1/metric1", 14, tags, 2)


@CFUNCTYPE(None, c_char_p)
def load(ctxId):
    print("** python *** Load called\n")

    lib_obj.ctx_config.restype = c_char_p
    res = lib_obj.ctx_config(ctxId, b"configip")
    print(res)


@CFUNCTYPE(None, c_char_p)
def unload(ctxId):
    print("** python *** Unload called\n")


lib_obj.StartCollector(collect, load, unload, define, b"python-collector", b"0.0.1")
