package pl_test

import (
	"testing"

	"github.com/lesomnus/pl"
	"github.com/stretchr/testify/require"
)

func TestPropsGet(t *testing.T) {
	tcs := []struct {
		desc     string
		props    pl.Props
		key      pl.RefKey
		expected any
	}{
		{
			desc:     "get string from map",
			props:    pl.NewProps(pl.M{"name": "Kujo Jotaro"}),
			key:      pl.K("name"),
			expected: "Kujo Jotaro",
		},
		{
			desc:     "get int from map",
			props:    pl.NewProps(pl.M{"age": 28}),
			key:      pl.K("age"),
			expected: 28,
		},
		{
			desc:     "get string from array",
			props:    pl.NewProps(pl.A{"Dio", 122}),
			key:      pl.K(0),
			expected: "Dio",
		},
		{
			desc:     "get int from array",
			props:    pl.NewProps(pl.A{"Dio", 122}),
			key:      pl.K(1),
			expected: 122,
		},
	}
	for _, tc := range tcs {
		t.Run(tc.desc, func(t *testing.T) {
			require := require.New(t)
			v, ok := tc.props.Get(tc.key)
			require.True(ok)
			require.Equal(tc.expected, v)
		})
	}

	t.Run("fails if", func(t *testing.T) {
		tcs := []struct {
			desc  string
			props pl.Props
			key   pl.RefKey
		}{
			{
				desc:  "access by int for map",
				props: pl.NewProps(pl.M{"name": "Kujo Jotaro"}),
				key:   pl.K(42),
			},
			{
				desc:  "access by string for array",
				props: pl.NewProps(pl.A{"Dio", 122}),
				key:   pl.K("age"),
			},
			{
				desc:  "access an array out of bound",
				props: pl.NewProps(pl.A{"Dio", 122}),
				key:   pl.K(2),
			},
			{
				desc:  "access by empty key",
				props: pl.NewProps(pl.A{"Dio", 122}),
				key:   pl.RefKey{},
			},
		}
		for _, tc := range tcs {
			t.Run(tc.desc, func(t *testing.T) {
				require := require.New(t)
				_, ok := tc.props.Get(tc.key)
				require.False(ok)
			})
		}
	})
}

func TestPropsSet(t *testing.T) {
	tcs := []struct {
		desc   string
		before pl.Props
		after  pl.Props
		key    pl.RefKey
		value  any
	}{
		{
			desc:   "set value to map",
			before: pl.NewProps(pl.M{"a": 42}),
			after:  pl.NewProps(pl.M{"a": 36}),
			key:    pl.K("a"),
			value:  36,
		},
		{
			desc:   "set value to arr",
			before: pl.NewProps(pl.A{"a", "b", "c"}),
			after:  pl.NewProps(pl.A{"a", "z", "c"}),
			key:    pl.K(1),
			value:  "z",
		},
	}
	for _, tc := range tcs {
		t.Run(tc.desc, func(t *testing.T) {
			require := require.New(t)

			ok := tc.before.Set(tc.key, tc.value)
			require.True(ok)
			require.Equal(tc.after, tc.before)
		})
	}

	t.Run("fails if", func(t *testing.T) {
		tcs := []struct {
			desc  string
			props pl.Props
			key   pl.RefKey
			value any
		}{
			{
				desc:  "access map as array",
				props: pl.NewProps(pl.M{}),
				key:   pl.K(0),
				value: nil,
			},
			{
				desc:  "access array as map",
				props: pl.NewProps(pl.A{}),
				key:   pl.K("a"),
				value: nil,
			},
			{
				desc:  "access array out of bound",
				props: pl.NewProps(pl.A{}),
				key:   pl.K(1),
				value: nil,
			},
			{
				desc:  "access with empty key",
				props: pl.NewProps(pl.A{}),
				key:   pl.RefKey{},
				value: nil,
			},
		}
		for _, tc := range tcs {
			t.Run(tc.desc, func(t *testing.T) {
				require := require.New(t)

				ok := tc.props.Set(tc.key, tc.value)
				require.False(ok)
			})
		}
	})
}

func TestPropsNext(t *testing.T) {
	props := pl.NewProps(pl.M{
		"map": pl.M{
			"a": 42,
		},
		"arr": pl.A{"b", 3.14, "c"},
		"str": "foo",
	})

	t.Run("map", func(t *testing.T) {
		require := require.New(t)

		sub, ok := props.Next(pl.K("map"))
		require.True(ok)

		v, ok := sub.Get(pl.K("a"))
		require.True(ok)
		require.Equal(42, v)
	})

	t.Run("array", func(t *testing.T) {
		require := require.New(t)

		sub, ok := props.Next(pl.K("arr"))
		require.True(ok)

		v, ok := sub.Get(pl.K(0))
		require.True(ok)
		require.Equal("b", v)
	})

	t.Run("non-prop field should be false", func(t *testing.T) {
		require := require.New(t)

		_, ok := props.Next(pl.K("str"))
		require.False(ok)
	})

	t.Run("non-exists field should be false", func(t *testing.T) {
		require := require.New(t)

		_, ok := props.Next(pl.K("not exists"))
		require.False(ok)
	})
}

func TestPropsResolve(t *testing.T) {
	tcs := []struct {
		desc     string
		props    pl.Props
		ref      []pl.RefKey
		expected any
	}{
		{
			desc:     "access nested maps",
			props:    pl.NewProps(pl.M{"a": pl.M{"b": pl.M{"c": 42}}}),
			ref:      must(pl.NewRef("a", "b", "c")),
			expected: 42,
		},
		{
			desc:     "access nested arrays",
			props:    pl.NewProps(pl.A{"a", pl.A{pl.A{1, 2, "bar"}, 3.14}}),
			ref:      must(pl.NewRef(1, 0, 2)),
			expected: "bar",
		},
		{
			desc:     "access nested maps and arrays",
			props:    pl.NewProps(pl.A{"a", pl.M{"a": pl.M{"b": pl.A{1, "baz", 21}}}}),
			ref:      must(pl.NewRef(1, "a", "b", 2)),
			expected: 21,
		},
	}
	for _, tc := range tcs {
		t.Run(tc.desc, func(t *testing.T) {
			require := require.New(t)

			v, err := tc.props.Resolve(tc.ref)
			require.NoError(err)
			require.Equal(tc.expected, v)
		})
	}

	t.Run("fails if access non-exists key", func(t *testing.T) {
		require := require.New(t)

		props := pl.NewProps(pl.M{"a": pl.A{1, "2", 3}})
		_, err := props.Resolve(must(pl.NewRef("a", 42)))
		require.Error(err)
		require.ErrorContains(err, "$.a")
		require.ErrorContains(err, "[42]")
	})
}
