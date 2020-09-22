package envparse

// Env maps a string key to a list of values.
// see net/url/url.go
type Env map[string][]string

// Set sets the key to value. It replaces any existing
// values.
func (e Env) Set(key, value string) {
	e[key] = []string{value}
}

// Add adds the value to key. It appends to any existing
// values associated with key.
func (e Env) Add(key, value string) {
	e[key] = append(e[key], value)
}

// Del deletes the values associated with key.
func (e Env) Del(key string) {
	delete(e, key)
}

// Parse parses raw key=value
// pairs into Env(some sort of url.Values)
func Parse(raw string) (Env, error) {
	e := make(Env)

	cb := func(key, value string) {
		e.Set(key, value)
	}

	err := ParseRaw(raw, cb)
	if err != nil {
		return nil, err
	}

	return e, nil
}
