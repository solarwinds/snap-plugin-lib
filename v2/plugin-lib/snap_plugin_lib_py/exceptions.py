class PluginLibException(Exception):
    pass


def throw_exception_if_error(func):
    def wrapper(*args, **kwargs):
        err = func(*args, **kwargs)
        if err.contents.msg is None:
            return None

        err_msg = str(err.contents.msg)
        if err_msg is not None:
            raise PluginLibException(err_msg)

    return wrapper