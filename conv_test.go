package pl_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/lesomnus/pl"
	"github.com/stretchr/testify/require"
)

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

		v, err := convs.Convert(reflect.TypeOf(int64(0)), reflect.TypeOf(int32(0)), reflect.ValueOf(int32(42)))
		require.NoError(err)
		require.Equal(v, int64(42))
	})

	t.Run("returns error if conversion is not defined", func(t *testing.T) {
		require := require.New(t)

		_, err := convs.ConvertTo(reflect.TypeOf(float64(0)), float32(0))
		require.Error(err)
		require.ErrorContains(err, "for float32")

		_, err = convs.ConvertTo(reflect.TypeOf(float64(0)), int32(0))
		require.Error(err)
		require.ErrorContains(err, "to float64 from int32")

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
