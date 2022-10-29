package pl

import (
	"fmt"
)

func NewPl(fs ...*Fn) *Pl {
	return &Pl{Funcs: fs}
}

func NewFn(name string, args ...interface{}) (*Fn, error) {
	args_, err := NewArgs(args...)
	if err != nil {
		return nil, err
	}

	return &Fn{Name: name, Args: args_}, nil
}

func NewArgs(args ...interface{}) ([]*Arg, error) {
	rst := make([]*Arg, len(args))

	for i, arg := range args {
		rst[i] = &Arg{}

		switch v := arg.(type) {
		case string:
			rst[i].String = &v
		case float64:
			rst[i].Float = &v
		case int:
			rst[i].Int = &v
		case []RefKey:
			rst[i].Ref = v
		case *Pl:
			rst[i].Nested = v

		default:
			return nil, fmt.Errorf("invalid type of argument at %d", i)
		}
	}

	return rst, nil
}

func NewRef(entries ...interface{}) ([]RefKey, error) {
	rst := make([]RefKey, len(entries))

	for i, entry := range entries {
		switch v := entry.(type) {
		case string:
			rst[i].Name = &v
		case int:
			rst[i].Index = &v

		default:
			return nil, fmt.Errorf("invalid type of argument at %d", i)
		}
	}

	return rst, nil
}

func NewRefKey[T string | int](key T) RefKey {
	switch v := interface{}(key).(type) {
	case string:
		return RefKey{Name: &v}
	case int:
		return RefKey{Index: &v}
	}

	panic(key)
}

func K[T string | int](key T) RefKey {
	return NewRefKey(key)
}
