from datetime import datetime
from .convertions import unpack_value_t, cmap_to_dict, ctimewithns_to_time


class NamespaceElement:
    def __init__(self, value, name, description):
        self.value = value
        self.name = name
        self.description = description

    @classmethod
    def unpack_from_ne_struct(cls, nm_element_struct):
        return cls(
            nm_element_struct.value.decode(encoding="utf-8"),
            nm_element_struct.name.decode(encoding="utf-8"),
            nm_element_struct.description.decode(encoding="utf-8"),
        )

    def is_dynamic(self) -> bool:
        return bool(not self.name)


class Namespace:
    def __init__(self, namespace_elements, length, string):
        self.length = length
        self.namespace_elements = namespace_elements
        self.string = string

    @classmethod
    def unpack_from_nm_struct(cls, namespace_struct):
        _length = namespace_struct.length
        _str = namespace_struct.string.decode(encoding="utf-8")
        elements = namespace_struct.elements
        _ne_arr = []
        for i in range(_length):
            _el = NamespaceElement.unpack_from_ne_struct(elements[i])
            _ne_arr.append(_el)
        return cls(_ne_arr, _length, _str)

    def __repr__(self):
        return self.string

    def __iter__(self):
        for element in self.namespace_elements:
            yield element
        
    def __len__(self)->int:
        return len(self.namespace_elements)

    def __getitem__(self, index: int)->NamespaceElement:
        return self.namespace_elements[index]

class Metric:
    def __init__(
        self,
        namespace="",
        description="",
        value="",
        value_type=None,
        timestamp="",
        tags="",
    ):
        self.namespace = namespace
        self.description = description

        self.value = value
        self.unit = value_type
        self.timestamp = timestamp
        self.tags = tags

    def _tags_to_str(self) -> str:
        all_tags = ""
        if self.tags:
            for k, v in self.tags.items():
                _tags = ":".join([str(k), str(v)])
                all_tags = " ".join([all_tags, _tags])
        return all_tags

    def __repr__(self) -> str:
        _repr = "{} {} {} {} {}".format(
            self.namespace,
            self.unit,
            self.value,
            self.description,
            datetime.utcfromtimestamp(self.timestamp),
        )
        tags = self._tags_to_str()
        _repr = " ".join([_repr, tags])
        return _repr

    @classmethod
    def unpack_from_metric_struct(cls, mt_struct):
        _namespace = Namespace.unpack_from_nm_struct(mt_struct.namespace.contents)
        _desc = mt_struct.description.decode(encoding="utf-8")
        _value, _unit = unpack_value_t(mt_struct.value)
        _time = ctimewithns_to_time(mt_struct.timestamp)
        _tags = cmap_to_dict(mt_struct.tags)
        return cls(_namespace, _desc, _value, _unit, _time, _tags,)
