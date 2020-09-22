package envparse

import (
	"fmt"
	"strings"
)

// Env maps a string key to a list of values.
// see net/url/url.go
type Env map[string]string

// Set sets the key to value.
func (e Env) Set(key, value string) {
	e[key] = value
}

// Del deletes the values associated with key.
func (e Env) Del(key string) {
	delete(e, key)
}

// String returns string representation
// of Env
func (e Env) String() string {
	var b strings.Builder

	left := len(e)

	for k, v := range e {
		b.WriteString(k)
		b.WriteByte('=')
		b.WriteString(fmt.Sprintf("%q", v))

		left--

		// got more elements to emit
		if left != 0 {
			b.WriteByte(' ')
		}
	}

	return b.String()
}

// Parse parses raw key=value
// pairs into Env(some sort of url.Values)
func Parse(raw string) (Env, error) {
	if len(raw) == 0 {
		return nil, nil
	}

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
