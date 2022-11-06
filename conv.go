package pl

import (
	"errors"
	"reflect"
	"strconv"
)

var ErrNotFound = errors.New("not found")

var string_t = reflect.TypeOf("")

type ConvMap map[reflect.Type]map[reflect.Type](func(v reflect.Value) (any, error))

func NewConvMap() ConvMap {
	rst := make(ConvMap)
	for in, convs := range defaultConversions {
		tgt := make(map[reflect.Type]func(v reflect.Value) (any, error))

		for out, conv := range convs {
			tgt[out] = conv
		}

		rst[in] = tgt
	}

	return rst
}

func (m ConvMap) Set(from reflect.Type, to reflect.Type, conv func(v reflect.Value) (any, error)) {
	tgt, ok := m[from]
	if !ok {
		tgt = make(map[reflect.Type]func(v reflect.Value) (any, error))
		m[from] = tgt
	}

	tgt[to] = conv
}

func (m ConvMap) MergeWith(other ConvMap) {
	for from, convs := range other {
		tgt, ok := m[from]
		if !ok {
			tgt = make(map[reflect.Type]func(v reflect.Value) (any, error))
			m[from] = tgt
		}

		for to, conv := range convs {
			tgt[to] = conv
		}
	}
}

func (m ConvMap) Convert(out reflect.Type, in reflect.Type, v reflect.Value) (any, error) {
	convs, ok := m[in]
	if !ok {
		return nil, ErrNotFound
	}

	conv, ok := convs[out]
	if !ok {
		return nil, ErrNotFound
	}

	return conv(v)
}

func (m ConvMap) ConvertTo(out reflect.Type, in any) (any, error) {
	return m.Convert(out, reflect.TypeOf(in), reflect.ValueOf(in))
}

var defaultConversions = ConvMap{
	reflect.TypeOf(int(0)): map[reflect.Type]func(v reflect.Value) (any, error){
		reflect.TypeOf(string("")): func(v reflect.Value) (any, error) { return strconv.Itoa(int(v.Int())), nil },
		reflect.TypeOf(float64(0)): func(v reflect.Value) (any, error) { return float64(v.Int()), nil },
	},
}
