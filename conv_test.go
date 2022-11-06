package pl_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/lesomnus/pl"
	"github.com/stretchr/testify/require"
)

func TestConvMapMergeWith(t *testing.T) {
	require := require.New(t)

	string_t := reflect.TypeOf("")
	int16_t := reflect.TypeOf(int16(0))
	int32_t := reflect.TypeOf(int32(0))
	int64_t := reflect.TypeOf(int64(0))

	lhs := make(pl.ConvMap)
	rhs := make(pl.ConvMap)

	lhs.Set(int64_t, int16_t, func(v reflect.Value) (any, error) { return 6416, nil })
	lhs.Set(int64_t, int32_t, func(v reflect.Value) (any, error) { return 523, nil })

	rhs.Set(int64_t, int32_t, func(v reflect.Value) (any, error) { return 6432, nil })
	rhs.Set(string_t, int64_t, func(v reflect.Value) (any, error) { return 42, nil })

	lhs.MergeWith(rhs)

	v, err := lhs[int64_t][int16_t](reflect.ValueOf(0))
	require.NoError(err)
	require.Equal(v, 6416)

	v, err = lhs[int64_t][int32_t](reflect.ValueOf(0))
	require.NoError(err)
	require.Equal(v, 6432)

	v, err = lhs[string_t][int64_t](reflect.ValueOf(0))
	require.NoError(err)
	require.Equal(v, 42)
}

func TestConvMapConvert(t *testing.T) {
	int32_t := reflect.TypeOf(int32(0))
	int64_t := reflect.TypeOf(int64(0))

	convs := pl.ConvMap{
		int32_t: map[reflect.Type]func(v reflect.Value) (any, error){
			int64_t: func(v reflect.Value) (any, error) {
				return v.Int(), nil
			},
		},
	}

	t.Run("convert with type instance", func(t *testing.T) {
		require := require.New(t)

		v, err := convs.Convert(int64_t, int32_t, reflect.ValueOf(int32(42)))
		require.NoError(err)
		require.Equal(v, int64(42))
	})

	t.Run("convert with temporal type instance", func(t *testing.T) {
		require := require.New(t)

		v, err := convs.Convert(int64_t, int32_t, reflect.ValueOf(int32(42)))
		require.NoError(err)
		require.Equal(v, int64(42))
	})

	t.Run("fails if conversion is not defined", func(t *testing.T) {
		require := require.New(t)

		_, err := convs.ConvertTo(reflect.TypeOf(float64(0)), float32(0))
		require.ErrorIs(err, pl.ErrNotFound)

		_, err = convs.ConvertTo(reflect.TypeOf(float64(0)), int32(0))
		require.ErrorIs(err, pl.ErrNotFound)

	})
}

func TestConvMapDefaultConversions(t *testing.T) {
	convs := pl.NewConvMap()

	tcs := []struct {
		in  any
		out any
	}{
		// from int
		{
			in:  42,
			out: "42",
		},
		{
			in:  42,
			out: 42.0,
		},
	}
	for _, tc := range tcs {
		t.Run(fmt.Sprintf("from %s to %s", reflect.TypeOf(tc.in).Name(), reflect.TypeOf(tc.out).Name()), func(t *testing.T) {
			require := require.New(t)

			v, err := convs.ConvertTo(reflect.TypeOf(tc.out), tc.in)
			require.NoError(err)
			require.Equal(tc.out, v)
		})
	}
}
