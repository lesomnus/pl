package pl

import (
	"fmt"
	"reflect"
	"strings"
)

type Ref []RefKey

func (r Ref) String() string {
	paths := make([]string, len(r))
	for i, k := range r {
		paths[i] = k.String()
	}

	return strings.Join(paths, "")
}

func Resolve(data any, ref Ref) (any, error) {
	cursor := reflect.ValueOf(data)
	for i, key := range ref {
		t := cursor.Type()
		for {
			switch t.Kind() {
			case reflect.Pointer:
				fallthrough
			case reflect.Interface:
				cursor = cursor.Elem()
				t = cursor.Type()
				continue
			}

			break
		}

		if key.Name != nil {
			switch t.Kind() {
			case reflect.Map:
				if t.Key().Kind() != reflect.String {
					return nil, fmt.Errorf("%s is a map but key type is not a string", ref[:i].String())
				}

				cursor = cursor.MapIndex(reflect.ValueOf(*key.Name))
				if !cursor.IsValid() {
					return nil, fmt.Errorf("$%s has no key %s", ref[:i].String(), *key.Name)
				}

			case reflect.Struct:
				cursor = cursor.FieldByName(*key.Name)
				if !cursor.IsValid() {
					return nil, fmt.Errorf("$%s has no field %s", ref[:i].String(), *key.Name)
				}

			default:
				return nil, fmt.Errorf("$%s is not an object but %s", ref[:i].String(), t.String())
			}
		} else if key.Index != nil {
			switch t.Kind() {
			case reflect.Map:
				var index any
				switch t.Key().Kind() {
				case reflect.Int:
					index = *key.Index
				case reflect.Int16:
					index = int16(*key.Index)
				case reflect.Int32:
					index = int32(*key.Index)
				case reflect.Int64:
					index = int64(*key.Index)
				case reflect.Uint:
					index = uint(*key.Index)
				case reflect.Uint16:
					index = uint16(*key.Index)
				case reflect.Uint32:
					index = uint32(*key.Index)
				case reflect.Uint64:
					index = uint64(*key.Index)

				default:
					return nil, fmt.Errorf("%s is a map but key type is not an integer", ref[:i].String())
				}

				cursor = cursor.MapIndex(reflect.ValueOf(index))
				if !cursor.IsValid() {
					return nil, fmt.Errorf("$%s has no key %d", ref[:i].String(), *key.Index)
				}

				continue

			case reflect.Array:
			case reflect.Slice:

			default:
				return nil, fmt.Errorf("$%s is not a list but %s", ref[:i].String(), t.String())
			}

			l := cursor.Len()
			if l <= *key.Index {
				return nil, fmt.Errorf("$%s: out of range", ref[:i].String())
			}

			cursor = cursor.Index(*key.Index)
		} else {
			return nil, fmt.Errorf("invalid key at %d", i)
		}
	}

	return cursor.Interface(), nil
}
