from ctypes import Structure, c_char_p


class Tags(Structure):
    _fields_ = [
        ("key", c_char_p),
        ("value", c_char_p)
    ]


class CError(Structure):
    _fields_ = [
        ("msg", c_char_p)
    ]


# Convert string to bytes if necessary.
# Allow to use string type in Python code and covert it to required char *
# (bytes) when calling C Api
def string_to_bytes(s):
    if isinstance(s, type("")):
        return bytes(s, 'utf-8')
    elif isinstance(s, type(b"")):
        return s
    else:
        raise Exception("Invalid type, expected string or bytes")


# Converting python dictionary to array of objects
def dict_to_tags(d):
    tags = (Tags * len(d))()

    for i, (k, v) in enumerate(d.items()):
        tags[i].key = string_to_bytes(k)
        tags[i].value = string_to_bytes(v)

    return tags
