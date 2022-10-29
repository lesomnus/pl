package pl

import (
	"fmt"
	"strings"
)

type Props struct {
	vs any
}

type M map[string]any

type A []any

func NewProps[T M | A](vs T) Props {
	return Props{vs}
}

func (p *Props) Get(key RefKey) (any, bool) {
	if key.Name != nil {
		vs, ok := p.vs.(M)
		if !ok {
			return nil, false
		}

		if v, ok := vs[*key.Name]; !ok {
			return nil, false
		} else {
			return v, true
		}
	} else if key.Index != nil {
		vs, ok := p.vs.(A)
		if !ok {
			return nil, false
		}
		if len(vs) <= *key.Index {
			return nil, false
		}

		return vs[*key.Index], true
	} else {
		// Invalid key
		return nil, false
	}
}

func (p *Props) Set(key RefKey, v any) bool {
	if key.Name != nil {
		vs, ok := p.vs.(M)
		if !ok {
			return false
		}

		vs[*key.Name] = v
		return true
	} else if key.Index != nil {
		vs, ok := p.vs.(A)
		if !ok {
			return false
		}
		if len(vs) <= *key.Index {
			return false
		}

		vs[*key.Index] = v
		return true
	} else {
		// Invalid key
		return false
	}
}

func (p *Props) Next(key RefKey) (Props, bool) {
	next, ok := p.Get(key)
	if !ok {
		return Props{}, false
	}

	switch next.(type) {
	case M:
	case A:

	default:
		return Props{}, false
	}

	return Props{next}, true
}

func (p *Props) Resolve(keys []RefKey) (any, error) {
	var (
		vs any = p.vs
		ok bool
	)
	for i, key := range keys {
		props := Props{vs}
		vs, ok = props.Get(key)
		if ok {
			continue
		}

		paths := make([]string, i)
		for j, k := range keys[:i] {
			paths[j] = k.String()
		}

		return nil, fmt.Errorf("key %s not exists at $%s", key.String(), strings.Join(paths, ""))
	}

	return vs, nil
}
