class PluginLibException(Exception):
    pass


# Decorator - used to throw exception when function return not-null string
def throw_exception_if_error(func):
    def func_wrapper(*args, **kwargs):
        err = func(*args, **kwargs)
        if err.contents.msg is None:
            return None

        err_msg = str(err.contents.msg)
        if err_msg is not None:
            raise PluginLibException(err_msg)

    return func_wrapper


# Decorator - used to throw exception when function return NULL
def throw_exception_if_null(exception_msg):
    def func_wrapper(func):
        def inner_wrapper(*args, **kwargs):
            err = func(*args, **kwargs)
            print(err)
            if err is None:
                raise PluginLibException(exception_msg)

        return inner_wrapper

    return func_wrapper