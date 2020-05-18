package cmd

// Defines string keys map
type StringKeyMap map[string]interface{}

// ContainsKey gets all map's keys
func (m *StringKeyMap) Keys() []string {
	keys := make([]string, 0, len(*m))
	for k := range *m {
		keys = append(keys, k)
	}
	return keys
}

// ContainsKey gets whether a key is presented in map
func (m *StringKeyMap) ContainsKey(key string) bool {
	_, ok := (*m)[key]
	return ok
}
