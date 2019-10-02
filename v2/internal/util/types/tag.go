package types

type Tags map[string]string

func (m Tags) ContainsKey(key string) bool {
	_, ok := m[key]
	return ok
}

func (m Tags) ContainsValue(value string) bool {
	for _, v := range m {
		if v == value {
			return true
		}
	}

	return false
}

func (m Tags) Contains(key string, value string) bool {
	for k, v := range m {
		if k == key && v == value {
			return true
		}
	}

	return false
}
