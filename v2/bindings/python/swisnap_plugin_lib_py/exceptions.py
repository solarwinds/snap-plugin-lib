from .dynamic_lib import PLUGIN_LIB_OBJ


class PluginLibException(Exception):
    pass


# Decorator - used to throw exception when function return not-null error
def throw_exception_if_error(func):
    def func_wrapper(*args, **kwargs):
        err = func(*args, **kwargs)
        if err.contents.msg is None:
            return None

        err_msg = str(err.contents.msg.decode(encoding="utf-8"))
        if err_msg is not None:
            PLUGIN_LIB_OBJ.dealloc_error(err)
            raise PluginLibException(err_msg)

    return func_wrapper


# Decorator - used to throw exception when function return NULL
def throw_exception_if_null(exception_msg):
    def func_wrapper(func):
        def inner_wrapper(*args, **kwargs):
            ret_value = func(*args, **kwargs)
            if ret_value is None:
                raise PluginLibException(exception_msg)

            return ret_value

        return inner_wrapper

    return func_wrapper
